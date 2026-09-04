package acquisition

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/user/orchestra-v3/internal/research"
	"github.com/user/orchestra-v3/internal/resources"
)

// WebAdapter implements AcquisitionAdapter for web references and HTTP documentation
type WebAdapter struct {
	httpClient     *http.Client
	offlineMode    bool
	maxPayloadSize int64
	allowPrivateIP bool
}

// NewWebAdapter constructs an initialized WebAdapter
func NewWebAdapter(offlineMode bool) *WebAdapter {
	adapter := &WebAdapter{
		offlineMode:    offlineMode,
		maxPayloadSize: 10 * 1024 * 1024, // 10MB cap
		allowPrivateIP: false,
	}
	transport := &http.Transport{
		MaxIdleConns:        50,
		IdleConnTimeout:     60 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	adapter.httpClient = &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("stopped after 10 redirects")
			}
			// Block redirects to private, loopback, link-local, unspecified, or metadata IPs
			h := req.URL.Hostname()
			if strings.EqualFold(h, "localhost") {
				return fmt.Errorf("%w: redirect to localhost blocked", ErrSSRFDetected)
			}
			if ip := net.ParseIP(h); ip != nil {
				if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsUnspecified() || ip.String() == "169.254.169.254" {
					return fmt.Errorf("%w: redirect to blocked IP %s", ErrSSRFDetected, ip.String())
				}
			}
			ips, err := net.LookupIP(h)
			if err == nil {
				for _, ip := range ips {
					if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsUnspecified() || ip.String() == "169.254.169.254" {
						return fmt.Errorf("%w: redirect resolved to blocked IP %s", ErrSSRFDetected, ip.String())
					}
				}
			}
			return nil
		},
	}
	return adapter
}

func (w *WebAdapter) Name() string {
	return "web"
}

func (w *WebAdapter) CanHandle(method string) bool {
	return strings.EqualFold(method, "web_fetch")
}

// SetAllowPrivateIP configures whether loopback/private IPs are permitted (used for httptest unit tests)
func (w *WebAdapter) SetAllowPrivateIP(allow bool) {
	w.allowPrivateIP = allow
}

// SetOfflineMode toggles offline fallback mode
func (w *WebAdapter) SetOfflineMode(offline bool) {
	w.offlineMode = offline
}

// validateTargetURL inspects scheme, quarantine boundaries, and performs SSRF detection
func (w *WebAdapter) validateTargetURL(urlStr string) (*url.URL, error) {
	if err := resources.CheckQuarantineBoundary(urlStr); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQuarantineViolation, err)
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("malformed URL: %w", err)
	}

	scheme := strings.ToLower(u.Scheme)
	if scheme != "http" && scheme != "https" {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedURLScheme, u.Scheme)
	}

	hostname := u.Hostname()
	if !w.allowPrivateIP {
		if strings.EqualFold(hostname, "localhost") {
			return nil, fmt.Errorf("%w: localhost address blocked", ErrSSRFDetected)
		}
		if parsedIP := net.ParseIP(hostname); parsedIP != nil {
			if parsedIP.IsLoopback() || parsedIP.IsPrivate() || parsedIP.IsLinkLocalUnicast() || parsedIP.IsLinkLocalMulticast() || parsedIP.IsUnspecified() || parsedIP.String() == "169.254.169.254" {
				return nil, fmt.Errorf("%w: resolved to blocked IP %s", ErrSSRFDetected, parsedIP.String())
			}
		}
		ips, err := net.LookupIP(hostname)
		if err == nil {
			for _, ip := range ips {
				if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsUnspecified() {
					return nil, fmt.Errorf("%w: resolved to blocked IP %s", ErrSSRFDetected, ip.String())
				}
				// Metadata IP check (AWS, GCP, Azure 169.254.169.254)
				if ip.String() == "169.254.169.254" {
					return nil, fmt.Errorf("%w: cloud metadata IP blocked", ErrSSRFDetected)
				}
			}
		}
	}

	return u, nil
}

// getOfflineFixture attempts to retrieve curated reference findings for offline resilience
func (w *WebAdapter) getOfflineFixture(resID string) ([]byte, bool) {
	fixture, exists := research.CuratedSourceFixtures[resID]
	if exists {
		data, err := json.MarshalIndent(fixture, "", "  ")
		if err == nil {
			return data, true
		}
	}
	if resID != "" {
		fallback := map[string]any{
			"source_id": resID,
			"status":    "offline_fixture",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		}
		data, err := json.MarshalIndent(fallback, "", "  ")
		if err == nil {
			return data, true
		}
	}
	return nil, false
}

// Acquire fetches the web reference, computes SHA256, and stores it or falls back to fixtures
func (w *WebAdapter) Acquire(ctx context.Context, res *resources.Resource, dest string) (*AcquisitionResult, error) {
	start := time.Now()

	if res == nil || res.ID == "" {
		return nil, ErrResourceNotFound
	}
	if !w.CanHandle(res.AcquisitionMethod) {
		return nil, fmt.Errorf("%w: resource '%s' requires acquisition method '%s', not web_fetch", ErrResourceNotAllowed, res.ID, res.AcquisitionMethod)
	}

	targetURL := res.CanonicalURL
	if targetURL == "" {
		targetURL = res.DocumentationURL
	}
	if targetURL == "" {
		targetURL = res.DocURL
	}

	// 1. Offline Mode Fast-Path
	if w.offlineMode {
		if fixtureData, ok := w.getOfflineFixture(res.ID); ok {
			h := sha256.Sum256(fixtureData)
			shaHex := hex.EncodeToString(h[:])
			return &AcquisitionResult{
				ResourceID:        res.ID,
				AdapterName:       w.Name(),
				AcquisitionMethod: "web_fetch",
				SourceURL:         targetURL,
				VersionOrSHA:      shaHex[:16],
				SHA256Hash:        shaHex,
				InstalledPath:     "",
				ResolvedTarget:    targetURL,
				AlreadyInstalled:  false,
				Ephemeral:         true,
				Duration:          time.Since(start),
				Output:            string(fixtureData),
				Metadata: map[string]string{
					"source":  "offline_fixture",
					"offline": "true",
				},
			}, nil
		}
	}

	// 2. Validate URL and SSRF Safety
	validURL, err := w.validateTargetURL(targetURL)
	if err != nil {
		return nil, err
	}

	// 3. HTTP GET Request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, validURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("User-Agent", "Orchestra-Acquisition-Engine/3.0 (Security Guarded)")

	resp, err := w.httpClient.Do(req)
	if err != nil {
		// Do NOT mask SSRF errors with offline fixtures!
		if errors.Is(err, ErrSSRFDetected) {
			return nil, err
		}
		// Network failure fallback to curated fixture
		if fixtureData, ok := w.getOfflineFixture(res.ID); ok {
			h := sha256.Sum256(fixtureData)
			shaHex := hex.EncodeToString(h[:])
			return &AcquisitionResult{
				ResourceID:        res.ID,
				AdapterName:       w.Name(),
				AcquisitionMethod: "web_fetch",
				SourceURL:         targetURL,
				VersionOrSHA:      shaHex[:16],
				SHA256Hash:        shaHex,
				InstalledPath:     "",
				ResolvedTarget:    targetURL,
				AlreadyInstalled:  false,
				Ephemeral:         true,
				Duration:          time.Since(start),
				Output:            string(fixtureData),
				Metadata: map[string]string{
					"source": "network_fallback_fixture",
					"error":  err.Error(),
				},
			}, nil
		}
		return nil, fmt.Errorf("failed fetching URL %s: %w", targetURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP GET %s returned status %d", targetURL, resp.StatusCode)
	}

	// 4. Enforce 10MB Payload Cap
	limitReader := io.LimitReader(resp.Body, w.maxPayloadSize+1)
	payload, err := io.ReadAll(limitReader)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %w", err)
	}
	if int64(len(payload)) > w.maxPayloadSize {
		return nil, fmt.Errorf("%w: received > %d bytes", ErrPayloadTooLarge, w.maxPayloadSize)
	}

	// 5. Compute SHA256 Checksum
	h := sha256.Sum256(payload)
	shaHex := hex.EncodeToString(h[:])

	// 6. Write to destination if specified
	installedPath := ""
	if dest != "" {
		if err := resources.CheckQuarantineBoundary(dest); err != nil {
			return nil, fmt.Errorf("%w: destination '%s' violates quarantine: %w", ErrQuarantinedDestination, dest, err)
		}
		// If dest is a directory, save as <res.ID>.json
		targetFile := dest
		if fi, statErr := os.Stat(dest); statErr == nil && fi.IsDir() {
			targetFile = filepath.Join(dest, res.ID+".json")
		}
		if writeErr := os.WriteFile(targetFile, payload, 0644); writeErr == nil {
			installedPath = targetFile
		}
	}

	return &AcquisitionResult{
		ResourceID:        res.ID,
		AdapterName:       w.Name(),
		AcquisitionMethod: "web_fetch",
		SourceURL:         targetURL,
		VersionOrSHA:      shaHex[:16],
		SHA256Hash:        shaHex,
		InstalledPath:     installedPath,
		ResolvedTarget:    targetURL,
		AlreadyInstalled:  false,
		Ephemeral:         installedPath == "",
		Duration:          time.Since(start),
		Output:            string(payload),
		Metadata: map[string]string{
			"status": "fetched",
			"bytes":  fmt.Sprintf("%d", len(payload)),
		},
	}, nil
}
