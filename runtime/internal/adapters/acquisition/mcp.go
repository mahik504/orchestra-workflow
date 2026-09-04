package acquisition

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/user/orchestra-v3/internal/resources"
	"github.com/user/orchestra-v3/internal/runner"
)

// BannedMCPServers defines servers blocked by security or duplication policy
var BannedMCPServers = map[string]string{
	"higgsfield":        "paid extra studio; Stitch remains 2D path",
	"magicui":           "always-on kit glow; rejected by taste policy",
	"21st":              "always-on kit glow; rejected by taste policy",
	"agent-reach":       "extra MCP writing into all AI tools; redundant",
	"code-review-graph": "unnecessary file scanning overhead",
	"open-design":       "second design studio vs Stitch; DESIGN.md is canonical",
}

// MCPConfig models .orchestra/mcp.json or mcp_config.json
type MCPConfig struct {
	MCPServers map[string]*MCPServerConfig `json:"mcpServers"`
}

// MCPServerConfig represents configuration for an individual MCP server
type MCPServerConfig struct {
	Command     string            `json:"command,omitempty"`
	Args        []string          `json:"args,omitempty"`
	Env         map[string]string `json:"env,omitempty"`
	URL         string            `json:"url,omitempty"`
	Transport   string            `json:"transport,omitempty"` // "stdio" | "sse"
	Disabled    bool              `json:"disabled,omitempty"`
	AutoApprove []string          `json:"autoApprove,omitempty"`
	Comment     string            `json:"comment,omitempty"`
}

// MCPAdapter implements AcquisitionAdapter for Model Context Protocol servers
type MCPAdapter struct {
	runner     runner.CommandRunner
	httpClient *http.Client
	timeout    time.Duration
}

// NewMCPAdapter creates an initialized MCPAdapter
func NewMCPAdapter(r runner.CommandRunner) *MCPAdapter {
	if r == nil {
		r = runner.NewOSCommandRunner()
	}
	return &MCPAdapter{
		runner: r,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		timeout: 5 * time.Second,
	}
}

func (m *MCPAdapter) Name() string {
	return "mcp"
}

func (m *MCPAdapter) CanHandle(method string) bool {
	lower := strings.ToLower(strings.TrimSpace(method))
	return lower == "mcp" || lower == "mcp_install" || lower == "mcp_connection"
}

var serverNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// ValidateConfig executes policy, sanitization, and anti-global install validation on an MCP server config
func (m *MCPAdapter) ValidateConfig(name string, cfg *MCPServerConfig) error {
	normName := strings.ToLower(strings.TrimSpace(name))
	if reason, banned := BannedMCPServers[normName]; banned {
		return fmt.Errorf("%w: server %q is banned (%s)", ErrMCPRejected, name, reason)
	}

	if !serverNameRegex.MatchString(name) {
		return fmt.Errorf("invalid server name %q: must contain only alphanumeric characters, underscores, or hyphens", name)
	}

	if cfg == nil {
		return fmt.Errorf("MCP server configuration cannot be nil")
	}

	// 1. Programmatic Anti-Global Installation Block
	if cfg.Command != "" {
		lowerCmd := strings.ToLower(cfg.Command)
		if lowerCmd == "-g" || lowerCmd == "--global" {
			return ErrGlobalInstallBlocked
		}
	}
	for _, arg := range cfg.Args {
		lower := strings.ToLower(strings.TrimSpace(arg))
		if lower == "-g" || lower == "--global" || lower == "global" || strings.HasPrefix(lower, "--global=") || strings.HasPrefix(lower, "-g=") {
			return fmt.Errorf("%w: argument %q attempts global installation", ErrGlobalInstallBlocked, arg)
		}
		if strings.ContainsAny(arg, "|;&$`\n\r><") {
			return fmt.Errorf("%w: argument %q contains prohibited shell metacharacters", ErrCommandInjectionRisk, arg)
		}
	}

	// 2. Command Injection & Metacharacter Check
	if cfg.Command != "" {
		if strings.ContainsAny(cfg.Command, "|;&$`\n\r><") {
			return fmt.Errorf("%w: command %q contains prohibited shell characters", ErrCommandInjectionRisk, cfg.Command)
		}
		// Disallow raw shell interpreters
		lowerCmd := strings.ToLower(cfg.Command)
		if lowerCmd == "cmd.exe" || lowerCmd == "powershell.exe" || lowerCmd == "bash" || lowerCmd == "sh" {
			return fmt.Errorf("%w: executing via raw shell %q is prohibited", ErrCommandInjectionRisk, cfg.Command)
		}
	}

	// 3. Remote URL SSRF Check for SSE / HTTP endpoints
	if cfg.URL != "" {
		u, err := url.Parse(cfg.URL)
		if err != nil {
			return fmt.Errorf("invalid MCP server URL: %w", err)
		}
		scheme := strings.ToLower(u.Scheme)
		if scheme != "http" && scheme != "https" {
			return fmt.Errorf("%w: %s", ErrUnsupportedURLScheme, u.Scheme)
		}
		if strings.EqualFold(u.Hostname(), "169.254.169.254") {
			return fmt.Errorf("%w: cloud metadata endpoint blocked", ErrSSRFDetected)
		}
	}

	return nil
}

// CheckReachability verifies that the server responds to JSON-RPC ping or HTTP health check
func (m *MCPAdapter) CheckReachability(ctx context.Context, cfg *MCPServerConfig) error {
	if cfg.URL != "" && (cfg.Transport == "sse" || strings.HasPrefix(cfg.URL, "http")) {
		req, err := http.NewRequestWithContext(ctx, http.MethodHead, cfg.URL, nil)
		if err != nil {
			return err
		}
		resp, err := m.httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("remote MCP endpoint unreachable: %w", err)
		}
		defer resp.Body.Close()
		return nil
	}

	if cfg.Command != "" {
		// Verify binary presence
		_, err := m.runner.LookPath(cfg.Command)
		if err != nil {
			return fmt.Errorf("%w: %s (%v)", ErrMCPBinaryNotFound, cfg.Command, err)
		}

		// Perform JSON-RPC 2.0 handshake if executable is available
		execCtx, cancel := context.WithTimeout(ctx, m.timeout)
		defer cancel()

		cmd := exec.CommandContext(execCtx, cfg.Command, cfg.Args...)
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return fmt.Errorf("failed to open stdin pipe: %w", err)
		}
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return fmt.Errorf("failed to open stdout pipe: %w", err)
		}

		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to start MCP process: %w", err)
		}

		initReq := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"orchestra-kernel","version":"3.0.0"}}}` + "\n"
		_, _ = stdin.Write([]byte(initReq))

		scanner := bufio.NewScanner(stdout)
		readDone := make(chan bool, 1)
		go func() {
			if scanner.Scan() {
				readDone <- true
			} else {
				readDone <- false
			}
		}()

		select {
		case ok := <-readDone:
			_ = stdin.Close()
			_ = cmd.Process.Kill()
			if !ok {
				return fmt.Errorf("empty response during MCP JSON-RPC handshake")
			}
			return nil
		case <-execCtx.Done():
			_ = stdin.Close()
			if cmd.Process != nil {
				_ = cmd.Process.Kill()
			}
			return fmt.Errorf("MCP JSON-RPC handshake timed out")
		}
	}

	return nil
}

// Acquire validates and registers the MCP server in the target workspace
func (m *MCPAdapter) Acquire(ctx context.Context, res *resources.Resource, dest string) (*AcquisitionResult, error) {
	start := time.Now()

	if res == nil || res.ID == "" {
		return nil, ErrResourceNotFound
	}
	if !m.CanHandle(res.AcquisitionMethod) {
		return nil, fmt.Errorf("%w: resource '%s' requires acquisition method '%s', not mcp", ErrResourceNotAllowed, res.ID, res.AcquisitionMethod)
	}

	// 1. Build default MCPServerConfig based on Resource specification
	cfg := &MCPServerConfig{
		Comment: fmt.Sprintf("Managed by Orchestra: %s", res.Name),
	}

	if res.CanonicalURL != "" && strings.HasPrefix(res.CanonicalURL, "http") {
		cfg.URL = res.CanonicalURL
		cfg.Transport = "sse"
	} else {
		// CLI-based MCP server (e.g. npx -y @playwright/mcp or npx -y @upstash/context7-mcp)
		cfg.Command = "npx"
		cfg.Args = []string{"-y", res.ID}
		cfg.Transport = "stdio"
	}

	// 2. Validate configuration against policies
	if err := m.ValidateConfig(res.ID, cfg); err != nil {
		return nil, err
	}

	// 3. Quarantine Boundary Check on destination
	if dest != "" {
		if err := resources.CheckQuarantineBoundary(dest); err != nil {
			return nil, fmt.Errorf("%w: destination '%s' violates quarantine: %v", ErrQuarantinedDestination, dest, err)
		}
	}

	// 4. Compute Configuration Hash
	cfgBytes, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to serialize MCP config: %w", err)
	}
	h := sha256.Sum256(cfgBytes)
	cfgHash := hex.EncodeToString(h[:])

	// 5. Write or update .orchestra/mcp.json in workspace
	installedPath := ""
	if dest != "" {
		orchDir := filepath.Join(dest, ".orchestra")
		_ = os.MkdirAll(orchDir, 0755)
		mcpPath := filepath.Join(orchDir, "mcp.json")

		var currentConfig MCPConfig
		if existingData, err := os.ReadFile(mcpPath); err == nil {
			_ = json.Unmarshal(existingData, &currentConfig)
		}
		if currentConfig.MCPServers == nil {
			currentConfig.MCPServers = make(map[string]*MCPServerConfig)
		}

		currentConfig.MCPServers[res.ID] = cfg
		updatedData, err := json.MarshalIndent(currentConfig, "", "  ")
		if err == nil {
			_ = os.WriteFile(mcpPath, updatedData, 0644)
			installedPath = mcpPath
		}
	}

	return &AcquisitionResult{
		ResourceID:        res.ID,
		AdapterName:       m.Name(),
		AcquisitionMethod: "mcp_install",
		SourceURL:         res.CanonicalURL,
		PackageName:       res.ID,
		VersionOrSHA:      cfgHash[:16],
		SHA256Hash:        cfgHash,
		InstalledPath:     installedPath,
		ResolvedTarget:    res.CanonicalURL,
		AlreadyInstalled:  false,
		Ephemeral:         false,
		Duration:          time.Since(start),
		Metadata: map[string]string{
			"transport": cfg.Transport,
			"status":    "configured",
		},
	}, nil
}
