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

func TestLayerAFailClosedWithoutURL(t *testing.T) {
	r := NewPlaywrightRunner()
	rep, err := r.Run(context.Background(), LayerARequest{OutputDir: t.TempDir(), Viewports: DefaultViewports(false)})
	if err != nil {
		t.Fatal(err)
	}
	if rep.Passed {
		t.Fatal("Layer A must not pass without a running app")
	}
	if rep.EnvironmentError == "" {
		t.Fatal("expected environment error")
	}
}

func TestLayerARejectsMockScreenshot(t *testing.T) {
	dir := t.TempDir()
	mockPath := filepath.Join(dir, "desktop.png")
	if err := os.WriteFile(mockPath, []byte("MOCK_SCREENSHOT_desktop_1440x900"), 0644); err != nil {
		t.Fatal(err)
	}
	rep := &LayerAReport{Owner: CanonicalOwner, Passed: true, Captures: []Capture{{
		Viewport: Viewport{Name: "desktop", Width: 1440, Height: 900},
		Route:    "/", ScreenshotPath: mockPath, Passed: true,
	}}}
	if err := rejectInventedEvidence(rep); err == nil {
		t.Fatal("mock screenshots must not be accepted")
	}
	if AcceptableScreenshot([]byte("MOCK_SCREENSHOT_x")) {
		t.Fatal("mock bytes must not be acceptable")
	}
	if !AcceptableScreenshot(TinyPNG) {
		t.Fatal("real PNG must be acceptable")
	}
}

func TestLayerAResponsiveFixtureEvidence(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		b, _ := os.ReadFile(filepath.Join("testdata", "fixture.html"))
		_, _ = w.Write(b)
	}))
	defer srv.Close()

	dir := t.TempDir()
	rep, err := ScriptedRunner{OverflowMobile: true}.Run(context.Background(), LayerARequest{
		BaseURL: srv.URL, OutputDir: dir, Routes: []string{"/"}, Viewports: DefaultViewports(false), Surface: "home",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rep.Captures) != 2 {
		t.Fatalf("desktop+mobile required, got %d", len(rep.Captures))
	}
	seen := map[string]bool{}
	for _, c := range rep.Captures {
		seen[c.Viewport.Name] = true
		if c.ScreenshotPath == "" {
			t.Fatal("screenshot path required")
		}
		if c.Viewport.Name == ViewportMobile && !c.Overflow {
			t.Fatal("mobile overflow must be reported independently of desktop")
		}
		if c.Viewport.Name == ViewportDesktop && c.Overflow {
			t.Fatal("desktop must not inherit mobile overflow")
		}
	}
	if !seen[ViewportDesktop] || !seen[ViewportMobile] {
		t.Fatalf("missing viewport: %v", seen)
	}
	if rep.Passed {
		t.Fatal("mobile overflow must fail Layer A")
	}
}

func TestLayerALivePlaywright(t *testing.T) {
	requireLivePlaywright(t)
	html := fixtureHTML(t)
	srv := serveHTML(t, html)
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	rep, err := NewPlaywrightRunner().Run(ctx, LayerARequest{
		BaseURL: srv.URL, OutputDir: t.TempDir(), Routes: []string{"/"},
		Viewports: DefaultViewports(false), RequiredSelectors: []string{"[data-qa=hero]"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if rep.PlaywrightMissing {
		t.Skipf("Playwright package missing at runtime: %s", rep.EnvironmentError)
	}
	if len(rep.Captures) < 2 {
		t.Fatalf("expected desktop+mobile captures, got %d (%s)", len(rep.Captures), rep.EnvironmentError)
	}
}
