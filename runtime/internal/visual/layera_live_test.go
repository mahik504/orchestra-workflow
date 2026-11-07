package visual

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func requireLivePlaywright(t *testing.T) {
	t.Helper()
	probe := ProbePlaywright()
	if !probe.Available {
		t.Skipf("Playwright library or Chromium is not installed: %s", probe.Error)
	}
	t.Logf("playwright library=%s chromium=%s node=%s", probe.Library, probe.Chromium, probe.Node)
}

func fixtureHTML(t *testing.T) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", "fixture.html"))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func serveHTML(t *testing.T, html []byte) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(html)
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestLayerASmokePlaywright(t *testing.T) {
	requireLivePlaywright(t)
	srv := serveHTML(t, fixtureHTML(t))
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	dir := t.TempDir()
	start := time.Now()
	rep, err := NewPlaywrightRunner().Run(ctx, LayerARequest{
		BaseURL: srv.URL, OutputDir: dir, Routes: []string{"/"},
		Viewports: DefaultViewports(false), RequiredSelectors: []string{"[data-qa=hero]"},
		Surface: "home",
	})
	elapsed := time.Since(start)
	if err != nil {
		t.Fatal(err)
	}
	if rep.PlaywrightMissing {
		t.Skipf("Playwright missing at runtime: %s", rep.EnvironmentError)
	}
	if rep.EnvironmentError != "" && len(rep.Captures) == 0 {
		t.Fatalf("BLOCKED: chromium did not capture: %s", rep.EnvironmentError)
	}
	if !rep.BrowserClosed {
		t.Fatal("browser must close")
	}
	if rep.BrowserPID > 0 && processAlive(rep.BrowserPID) {
		t.Fatalf("orphaned chromium pid %d", rep.BrowserPID)
	}
	if len(rep.Captures) != 2 {
		t.Fatalf("desktop+mobile required, got %d (%s)", len(rep.Captures), rep.EnvironmentError)
	}
	if rep.BrowserName != "chromium" {
		t.Fatalf("browser_name=%s", rep.BrowserName)
	}
	if rep.BrowserVersion == "" {
		t.Fatal("browser version missing")
	}
	t.Logf("browser=%s %s startup_ms=%d screenshot_ms=%d wall_ms=%d chromium=%s",
		rep.BrowserName, rep.BrowserVersion, rep.StartupMS, rep.ScreenshotMS, elapsed.Milliseconds(), rep.ChromiumPath)
	for _, c := range rep.Captures {
		b, err := os.ReadFile(c.ScreenshotPath)
		if err != nil {
			t.Fatal(err)
		}
		if IsMockScreenshot(b) {
			t.Fatalf("MOCK_SCREENSHOT cannot pass: %s", c.ScreenshotPath)
		}
		if !IsPNG(b) || !c.PNGMagic {
			t.Fatalf("PNG magic missing for %s size=%d", c.Viewport.Name, len(b))
		}
		if len(b) <= len(TinyPNG) {
			t.Fatalf("%s screenshot too small to be a real capture: %d bytes", c.Viewport.Name, len(b))
		}
		if c.ByteSize != len(b) {
			t.Fatalf("byte_size %d != file %d", c.ByteSize, len(b))
		}
		if c.URL == "" {
			t.Fatalf("capture URL missing for %s", c.Viewport.Name)
		}
		t.Logf("capture viewport=%s %dx%d bytes=%d png=%v overflow=%v path=%s url=%s",
			c.Viewport.Name, c.Viewport.Width, c.Viewport.Height, len(b), c.PNGMagic, c.Overflow, c.ScreenshotPath, c.URL)
	}
}

func TestLayerAResponsiveRealBrowser(t *testing.T) {
	requireLivePlaywright(t)
	srv := serveHTML(t, fixtureHTML(t))
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	rep, err := NewPlaywrightRunner().Run(ctx, LayerARequest{
		BaseURL: srv.URL, OutputDir: t.TempDir(), Routes: []string{"/"},
		Viewports: DefaultViewports(false), Surface: "home",
	})
	if err != nil {
		t.Fatal(err)
	}
	if rep.PlaywrightMissing {
		t.Skipf("Playwright missing at runtime: %s", rep.EnvironmentError)
	}
	if len(rep.Captures) < 2 {
		t.Fatalf("expected at least two captures, got %d (%s)", len(rep.Captures), rep.EnvironmentError)
	}
	var desk, mob *Capture
	for i := range rep.Captures {
		c := &rep.Captures[i]
		switch c.Viewport.Name {
		case ViewportDesktop:
			desk = c
		case ViewportMobile:
			mob = c
		}
	}
	if desk == nil || mob == nil {
		t.Fatal("desktop and mobile captures required")
	}
	if desk.Viewport.Width != 1440 || mob.Viewport.Width != 390 {
		t.Fatalf("viewports desktop=%d mobile=%d", desk.Viewport.Width, mob.Viewport.Width)
	}
	if desk.Overflow {
		t.Fatal("desktop must not inherit mobile overflow")
	}
	if !mob.Overflow {
		t.Fatal("mobile overflow must be a Layer A finding")
	}
	if desk.Passed == mob.Passed {
		t.Fatal("desktop pass must not equal mobile pass on this fixture")
	}
	if mob.URL == "" || mob.ScreenshotPath == "" {
		t.Fatal("mobile finding needs viewport + URL + screenshot path")
	}
	var sawOverflow bool
	for _, f := range rep.Findings {
		if f.Type == "overflow" && f.Viewport == ViewportMobile && f.EvidencePath != "" {
			sawOverflow = true
		}
	}
	if !sawOverflow {
		t.Fatalf("expected mobile overflow finding with screenshot path: %+v", rep.Findings)
	}
	if rep.Passed {
		t.Fatal("mobile overflow must fail Layer A")
	}
}
