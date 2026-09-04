#!/usr/bin/env python3
"""
Milestone 2 Empirical Challenger Test Suite: Research Coordinator & Token Weighting Adversary
Audited and executed by challenger_m2_2.

Empirically tests:
  1. Attempted Quarantine Evasion:
     - Rejection of paths containing 'skills_library' (varied slashes, case variations)
     - Rejection of 8.3 short name aliases ('SKILLS~1', 'CURATE~1')
     - Resolution and rejection of NTFS directory junctions pointing to quarantined targets
     - Filtering of quarantined candidate IDs/URLs during candidate discovery
  2. Offline Resilience:
     - Zero network latency simulation (<30ms, non-blocking execution)
     - Structural integrity of CuratedSourceFixtures (awwwards, jiro-design, cari-institute, awesome-design-md, godly-design, refero-design)
     - Rich visual design profiles (6+ color palette tokens, AAA contrast, dark matte base #0B0E14, typography triad, spring dynamics)
  3. Diversity Enforcement:
     - Rejection of single-source suggestions for high-visual tasks
     - Multi-source pull (>=2 distinct sources) across diverse archetype families
     - Validation of the 4 Archetype Families (AestheticBenchmark, LayoutComposition, MovementTaxonomy, SpecialistEcho)
     - Enforcement of the 2-source gate in reference-log.md generation
  4. Token Weight Scaling:
     - Validation across all 10 Design Archetypes:
       premium-website, 3d-portfolio, operator-hud, b2b-portal, academic-reader,
       micro-interactions, physics-canvas, saas-dashboard, mobile-app, security-audit
     - Bound verification (~1,500 baseline to >8,000+ for high-visual composite archetypes)
     - Resource deduplication & idempotency
     - Monotonic complexity scaling across task tiers
"""

import json
import os
import re
import subprocess
import sys
import tempfile
import time
import unittest

WORKFLOW_ROOT = r"C:\projects\orchestra-workflow"
RUNTIME_ROOT = os.path.join(WORKFLOW_ROOT, "runtime")
REGISTRIES_ROOT = os.path.join(WORKFLOW_ROOT, "registries")
RESOURCES_JSON = os.path.join(REGISTRIES_ROOT, "resources.json")
GRAPH_JSON = os.path.join(REGISTRIES_ROOT, "design-resource-graph.json")


def load_json(path):
    with open(path, "r", encoding="utf-8-sig") as f:
        return json.load(f)


class TestM2ResearchAndTokenAdversarial(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.resources = load_json(RESOURCES_JSON)
        cls.graph = load_json(GRAPH_JSON)
        cls.res_by_id = {r["id"]: r for r in cls.resources}

    # =========================================================================
    # 1. QUARANTINE BOUNDARY & EVASION ATTEMPTS
    # =========================================================================

    def test_quarantine_path_patterns_in_registries(self):
        """Verify zero canonical resources reference quarantined skills_library or 8.3 aliases."""
        banned = ["skills_library", "skills-library", "skills~1", "curate~1"]
        for res in self.resources:
            res_id = res.get("id", "")
            for pattern in banned:
                self.assertNotIn(pattern, res_id.lower(), f"Quarantined pattern in resource ID: {res_id}")
                canon = res.get("canonical_url", "")
                self.assertNotIn(pattern, canon.lower(), f"Quarantined pattern in canonical_url of {res_id}")
                repo = res.get("source_repository", "")
                self.assertNotIn(pattern, repo.lower(), f"Quarantined pattern in source_repository of {res_id}")
                doc = res.get("documentation_url", "")
                self.assertNotIn(pattern, doc.lower(), f"Quarantined pattern in documentation_url of {res_id}")

    def test_quarantine_go_test_execution(self):
        """Execute Go adversarial quarantine evasion tests."""
        cmd = ["go", "test", "-v", "-count=1", "./internal/research", "-run", "TestAdversarial_Quarantine"]
        proc = subprocess.run(cmd, cwd=RUNTIME_ROOT, capture_output=True, text=True)
        self.assertEqual(proc.returncode, 0, f"Go quarantine tests failed:\nSTDOUT:\n{proc.stdout}\nSTDERR:\n{proc.stderr}")
        self.assertIn("PASS: TestAdversarial_Quarantine_DirectBannedPaths", proc.stdout)
        self.assertIn("PASS: TestAdversarial_Quarantine_8dot3Aliases", proc.stdout)
        self.assertIn("PASS: TestAdversarial_Quarantine_NTFSJunctionResolution", proc.stdout)
        self.assertIn("PASS: TestAdversarial_Quarantine_CandidateInjection", proc.stdout)

    # =========================================================================
    # 2. OFFLINE RESILIENCE & CURATED FIXTURES
    # =========================================================================

    def test_curated_fixtures_defined_and_valid(self):
        """Verify the 6 curated fallback fixtures exist in resources.json and have complete metadata."""
        fixtures = ["awwwards", "jiro-design", "cari-institute", "awesome-design-md", "godly-design", "refero-design"]
        for f_id in fixtures:
            self.assertIn(f_id, self.res_by_id, f"Mandatory offline fixture missing from canonical registry: {f_id}")
            res = self.res_by_id[f_id]
            self.assertEqual(res.get("status"), "ACTIVE", f"Fixture {f_id} must have ACTIVE status")
            self.assertTrue(res.get("canonical_url", "").startswith("https://"), f"Fixture {f_id} URL invalid")
            self.assertGreater(len(res.get("routing_tags", [])), 0, f"Fixture {f_id} has empty routing_tags")

    def test_offline_resilience_go_test_execution(self):
        """Execute Go offline resilience adversarial tests."""
        cmd = ["go", "test", "-v", "-count=1", "./internal/research", "-run", "TestAdversarial_OfflineResilience"]
        proc = subprocess.run(cmd, cwd=RUNTIME_ROOT, capture_output=True, text=True)
        self.assertEqual(proc.returncode, 0, f"Go offline resilience tests failed:\nSTDOUT:\n{proc.stdout}\nSTDERR:\n{proc.stderr}")
        self.assertIn("PASS: TestAdversarial_OfflineResilience_ZeroNetworkSimulation", proc.stdout)
        self.assertIn("PASS: TestAdversarial_OfflineResilience_CuratedFixturesIntegrity", proc.stdout)
        self.assertIn("PASS: TestAdversarial_OfflineResilience_RichVisualProfiles", proc.stdout)

    # =========================================================================
    # 3. DIVERSITY ENFORCEMENT & ARCHETYPE FAMILIES
    # =========================================================================

    def test_diversity_enforcement_go_test_execution(self):
        """Execute Go diversity enforcement adversarial tests."""
        cmd = ["go", "test", "-v", "-count=1", "./internal/research", "-run", "TestAdversarial_Diversity"]
        proc = subprocess.run(cmd, cwd=RUNTIME_ROOT, capture_output=True, text=True)
        self.assertEqual(proc.returncode, 0, f"Go diversity enforcement tests failed:\nSTDOUT:\n{proc.stdout}\nSTDERR:\n{proc.stderr}")
        self.assertIn("PASS: TestAdversarial_Diversity_SingleSourceSuggestionRejection", proc.stdout)
        self.assertIn("PASS: TestAdversarial_Diversity_CrossFamilyDistribution", proc.stdout)
        self.assertIn("PASS: TestAdversarial_Diversity_MinSourceViolationGate", proc.stdout)

    # =========================================================================
    # 4. TOKEN WEIGHT SCALING ACROSS ALL 10 ARCHETYPES
    # =========================================================================

    def test_router_token_scaling_go_test_execution(self):
        """Execute Go router adversarial token scaling and bounds tests."""
        cmd = ["go", "test", "-v", "-count=1", "./internal/router", "-run", "TestAdversarial_Router"]
        proc = subprocess.run(cmd, cwd=RUNTIME_ROOT, capture_output=True, text=True)
        self.assertEqual(proc.returncode, 0, f"Go router adversarial tests failed:\nSTDOUT:\n{proc.stdout}\nSTDERR:\n{proc.stderr}")
        self.assertIn("PASS: TestAdversarial_Router_All10DesignArchetypesTokenScaling", proc.stdout)
        self.assertIn("PASS: TestAdversarial_Router_TokenDeduplication", proc.stdout)
        self.assertIn("PASS: TestAdversarial_Router_MonotonicComplexityScaling", proc.stdout)
        self.assertIn("PASS: TestAdversarial_Router_NilCatalogAndGraphFallback", proc.stdout)

    def test_all_10_archetypes_present_in_graph(self):
        """Verify all 10 design archetypes exist in design-resource-graph.json with defined capabilities."""
        archetypes = [
            "premium-website",
            "3d-portfolio",
            "operator-hud",
            "b2b-portal",
            "academic-reader",
            "micro-interactions",
            "physics-canvas",
            "saas-dashboard",
            "mobile-app",
            "security-audit",
        ]
        caps = self.graph.get("capabilities", {})
        for arch in archetypes:
            self.assertIn(arch, caps, f"Mandatory capability missing from graph: {arch}")
            cap = caps[arch]
            self.assertTrue(cap.get("name"), f"Capability {arch} missing name")
            self.assertTrue(cap.get("primary_archetype"), f"Capability {arch} missing primary_archetype")
            self.assertGreater(len(cap.get("trigger_tags", [])), 0, f"Capability {arch} missing trigger_tags")
            self.assertGreater(len(cap.get("anti_patterns", [])), 0, f"Capability {arch} missing anti_patterns")


if __name__ == "__main__":
    unittest.main(verbosity=2)
