package designdomain

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestDomainsStayOnTheClassifierParent(t *testing.T) {
	file := loadFile(t)
	cases := []struct {
		name       string
		capability string
		text       string
		wantDomain string
		wantParent string
	}{
		{"editorial", "premium-website", "An editorial magazine longform for a publisher", "editorial", "premium-website"},
		{"award", "premium-website", "An award-winning agency site", "award-agency", "premium-website"},
		{"shop", "premium-website", "A storefront with checkout", "ecommerce", "premium-website"},
		{"plain premium", "premium-website", "A landing page for a coffee roastery", "premium-websites", "premium-website"},
		{"spatial", "3d-portfolio", "A 3D WebGL portfolio with R3F camera orbits", "immersive-3d", "3d-portfolio"},
		{"shader", "3d-portfolio", "A custom fragment shader for the hero", "shader-experience", "3d-portfolio"},
		{"dashboard", "saas-dashboard", "A scheduling dashboard with attendance charts", "data-dashboard", "saas-dashboard"},
		{"saas", "saas-dashboard", "A multi-tenant SaaS with billing", "saas-application", "saas-dashboard"},
		{"hud", "operator-hud", "A mission control cockpit for telemetry", "mission-control", "operator-hud"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Annotate(file, tc.capability, tc.text)
			if !got.Found {
				t.Fatal("expected a domain")
			}
			if got.Domain.ID != tc.wantDomain || got.Domain.Parent != tc.wantParent {
				t.Fatalf("domain %s parent %s", got.Domain.ID, got.Domain.Parent)
			}
			if tc.capability == "premium-website" || tc.capability == "3d-portfolio" || tc.capability == "saas-dashboard" {
				if len(got.Gates) == 0 || got.Gates[0] != "responsive" {
					t.Fatalf("responsive gate missing: %v", got.Gates)
				}
			}
		})
	}
}

func TestThreeParentsDoNotCollapse(t *testing.T) {
	file := loadFile(t)
	restaurant := Annotate(file, "premium-website", "Build a distinctive restaurant website")
	school := Annotate(file, "saas-dashboard", "Build a scheduling dashboard for a school with attendance charts")
	spatial := Annotate(file, "3d-portfolio", "Build a 3D WebGL portfolio with R3F")
	if restaurant.Domain.Parent == school.Domain.Parent || school.Domain.Parent == spatial.Domain.Parent || restaurant.Domain.Parent == spatial.Domain.Parent {
		t.Fatal("restaurant, school, and 3D domains collapsed onto one parent")
	}
}

func loadFile(t *testing.T) File {
	t.Helper()
	_, here, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("caller")
	}
	root := filepath.Join(filepath.Dir(here), "..", "..", "..")
	file, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	return file
}
