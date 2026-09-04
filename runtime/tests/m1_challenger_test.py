#!/usr/bin/env python3
"""
Milestone 1 Empirical Challenger Test Suite: Schema & Graph Resolution
Empirically tests:
  1. Draft-07 Schema Compliance & Meta-schema validation for workflow & brain
  2. Referential Integrity & Graph Resolution across all 20 domains and 11 capabilities
  3. Redundancy & Duplicate Detection in domains and capability phases
  4. Edge Cases: Malformed JSON, empty fields, invalid URLs, unknown types
  5. Quarantine Isolation & Boundary Verification
"""

import json
import os
import re
import sys
import unittest
from urllib.parse import urlparse
import jsonschema
from jsonschema import Draft7Validator, FormatChecker

WORKFLOW_ROOT = r"C:\projects\orchestra-workflow"
BRAIN_ROOT = r"C:\projects\orchestra-brain"

WORKFLOW_RESOURCES_PATH = os.path.join(WORKFLOW_ROOT, "registries", "resources.json")
WORKFLOW_GRAPH_PATH = os.path.join(WORKFLOW_ROOT, "registries", "design-resource-graph.json")
WORKFLOW_RES_SCHEMA_PATH = os.path.join(WORKFLOW_ROOT, "registries", "schemas", "resources.schema.json")
WORKFLOW_GRAPH_SCHEMA_PATH = os.path.join(WORKFLOW_ROOT, "registries", "schemas", "design-resource-graph.schema.json")

BRAIN_RESOURCES_PATH = os.path.join(BRAIN_ROOT, "registries", "resources.json")
BRAIN_GRAPH_PATH = os.path.join(BRAIN_ROOT, "registries", "design-resource-graph.json")
BRAIN_RES_SCHEMA_PATH = os.path.join(BRAIN_ROOT, "registries", "schemas", "resources.schema.json")
BRAIN_GRAPH_SCHEMA_PATH = os.path.join(BRAIN_ROOT, "registries", "schemas", "design-resource-graph.schema.json")

def load_json_file(path):
    with open(path, "r", encoding="utf-8-sig") as f:
        return json.load(f)

class TestDraft07SchemaCompliance(unittest.TestCase):
    """Test suite 1: Validate schemas against Draft-07 meta-schema and validate registries against schemas."""

    def setUp(self):
        self.res_schema = load_json_file(WORKFLOW_RES_SCHEMA_PATH)
        self.graph_schema = load_json_file(WORKFLOW_GRAPH_SCHEMA_PATH)
        self.resources = load_json_file(WORKFLOW_RESOURCES_PATH)
        self.graph = load_json_file(WORKFLOW_GRAPH_PATH)

    def test_schema_valid_draft07_meta(self):
        """Verify both schemas are valid according to the Draft-07 meta-schema."""
        Draft7Validator.check_schema(self.res_schema)
        Draft7Validator.check_schema(self.graph_schema)

    def test_workflow_resources_validate_against_schema(self):
        """Verify workflow resources.json strictly validates against resources.schema.json."""
        validator = Draft7Validator(self.res_schema, format_checker=FormatChecker())
        errors = list(validator.iter_errors(self.resources))
        self.assertEqual(len(errors), 0, f"Validation errors in workflow resources.json: {[e.message for e in errors]}")

    def test_brain_resources_validate_against_schema(self):
        """Verify brain resources.json strictly validates against resources.schema.json."""
        brain_res = load_json_file(BRAIN_RESOURCES_PATH)
        validator = Draft7Validator(self.res_schema, format_checker=FormatChecker())
        errors = list(validator.iter_errors(brain_res))
        self.assertEqual(len(errors), 0, f"Validation errors in brain resources.json: {[e.message for e in errors]}")

    def test_workflow_graph_validates_against_schema(self):
        """Verify workflow design-resource-graph.json strictly validates against design-resource-graph.schema.json."""
        validator = Draft7Validator(self.graph_schema, format_checker=FormatChecker())
        errors = list(validator.iter_errors(self.graph))
        self.assertEqual(len(errors), 0, f"Validation errors in workflow graph: {[e.message for e in errors]}")

    def test_brain_graph_validates_against_schema(self):
        """Verify brain design-resource-graph.json strictly validates against design-resource-graph.schema.json."""
        brain_graph = load_json_file(BRAIN_GRAPH_PATH)
        validator = Draft7Validator(self.graph_schema, format_checker=FormatChecker())
        errors = list(validator.iter_errors(brain_graph))
        self.assertEqual(len(errors), 0, f"Validation errors in brain graph: {[e.message for e in errors]}")

    def test_schema_enforces_additional_properties_false(self):
        """Ensure schemas forbid unmapped additional properties."""
        self.assertFalse(self.res_schema["items"].get("additionalProperties", True))
        self.assertFalse(self.graph_schema.get("additionalProperties", True))
        self.assertFalse(self.graph_schema["properties"]["capabilities"]["additionalProperties"].get("additionalProperties", True))

class TestGraphReferentialIntegrity(unittest.TestCase):
    """Test suite 2: Adversarially test graph resolution across all domains and capability phases."""

    def setUp(self):
        self.resources = load_json_file(WORKFLOW_RESOURCES_PATH)
        self.graph = load_json_file(WORKFLOW_GRAPH_PATH)
        self.resource_ids = {r["id"]: r for r in self.resources}
        self.domains = self.graph.get("domains", {})
        self.capabilities = self.graph.get("capabilities", {})

    def test_resource_count_and_unique_ids(self):
        """Verify all resources have unique IDs and meet baseline count (>=65 per spec, worker has 126)."""
        self.assertGreaterEqual(len(self.resources), 65, "Total resources must be at least 65")
        ids = [r["id"] for r in self.resources]
        self.assertEqual(len(ids), len(set(ids)), f"Duplicate resource IDs found: {[x for x in ids if ids.count(x) > 1]}")

    def test_all_domain_resources_exist_in_registry(self):
        """Verify EVERY resource ID in EVERY domain exists in resources.json."""
        missing = {}
        for dom_name, res_list in self.domains.items():
            unresolved = [rid for rid in res_list if rid not in self.resource_ids]
            if unresolved:
                missing[dom_name] = unresolved
        self.assertEqual(missing, {}, f"Domains contain unresolved resource IDs: {missing}")

    def test_detect_duplicate_resources_within_domains(self):
        """Adversarial check: Detect redundant duplicate resource references within domains."""
        duplicates = {}
        for dom_name, res_list in self.domains.items():
            dups = [r for r in res_list if res_list.count(r) > 1]
            if dups:
                duplicates[dom_name] = sorted(list(set(dups)))
        # Worker has 3 domains with duplicates: reverse_engineering, security, webgl
        print(f"\n[EMPIRICAL FINDING] Domain redundant duplicates: {duplicates}")

    def test_all_capability_phases_resolve_completely(self):
        """Verify EVERY entry in EVERY capability phase resolves to either a domain or canonical resource ID."""
        missing = {}
        for cap_name, cap_def in self.capabilities.items():
            for phase in ["discovery", "reverse_engineering", "synthesis", "implementation", "optional_extensions", "qa"]:
                entries = cap_def.get(phase, [])
                unresolved = []
                for entry in entries:
                    if entry not in self.domains and entry not in self.resource_ids:
                        unresolved.append(entry)
                if unresolved:
                    missing[f"{cap_name}.{phase}"] = unresolved
        self.assertEqual(missing, {}, f"Capability phases reference unknown entries: {missing}")

    def test_detect_capability_phase_duplicates(self):
        """Adversarial check: Detect redundant duplicates inside capability phase arrays."""
        phase_duplicates = {}
        for cap_name, cap_def in self.capabilities.items():
            for phase in ["discovery", "reverse_engineering", "synthesis", "implementation", "optional_extensions", "qa"]:
                entries = cap_def.get(phase, [])
                dups = [e for e in entries if entries.count(e) > 1]
                if dups:
                    phase_duplicates[f"{cap_name}.{phase}"] = sorted(list(set(dups)))
        # Worker has reverse-engineering.reverse_engineering with ['skillui', 'skillui']
        print(f"[EMPIRICAL FINDING] Capability phase duplicate entries: {phase_duplicates}")

    def test_full_recursive_expansion_to_canonical_resources(self):
        """Fully expand every capability into its final set of concrete canonical resources."""
        expansion_failures = {}
        for cap_name, cap_def in self.capabilities.items():
            resolved_resources = set()
            for phase in ["discovery", "reverse_engineering", "synthesis", "implementation", "optional_extensions", "qa"]:
                for entry in cap_def.get(phase, []):
                    if entry in self.domains:
                        for sub_res in self.domains[entry]:
                            if sub_res in self.resource_ids:
                                resolved_resources.add(sub_res)
                            else:
                                expansion_failures.setdefault(cap_name, []).append(f"{entry}->{sub_res}")
                    elif entry in self.resource_ids:
                        resolved_resources.add(entry)
                    else:
                        expansion_failures.setdefault(cap_name, []).append(entry)
            self.assertGreater(len(resolved_resources), 0, f"Capability {cap_name} resolved to 0 concrete resources")
        self.assertEqual(expansion_failures, {}, f"Expansion failures: {expansion_failures}")

    def test_required_capability_archetypes_exist(self):
        """Verify all 11 required capability archetypes are defined."""
        required = [
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
            "reverse-engineering"
        ]
        for req in required:
            self.assertIn(req, self.capabilities, f"Missing required capability: {req}")

class TestEdgeCasesAndAdversarialInputs(unittest.TestCase):
    """Test suite 3: Test edge cases, malformed data rejection, empty fields, invalid URLs, and unknown types."""

    def setUp(self):
        self.res_schema = load_json_file(WORKFLOW_RES_SCHEMA_PATH)
        self.graph_schema = load_json_file(WORKFLOW_GRAPH_SCHEMA_PATH)
        self.validator_res = Draft7Validator(self.res_schema, format_checker=FormatChecker())
        self.validator_graph = Draft7Validator(self.graph_schema, format_checker=FormatChecker())
        self.resources = load_json_file(WORKFLOW_RESOURCES_PATH)
        self.graph = load_json_file(WORKFLOW_GRAPH_PATH)

    def test_empty_string_fields_rejection(self):
        """Verify schema rejects resources with empty/invalid strings in required fields."""
        invalid_resource = {
            "id": "",  # Empty ID
            "name": "Invalid Resource",
            "canonical_url": "https://example.com",
            "source_type": "github_repository",
            "category": ["TEST"],
            "representation": "dependency",
            "routing_tags": ["test"],
            "acquisition_method": "npm",
            "runtime_method": "project_scoped_install",
            "status": "ACTIVE"
        }
        errors = list(self.validator_res.iter_errors([invalid_resource]))
        self.assertGreater(len(errors), 0, "Schema must reject resource with empty string ID")

    def test_missing_required_fields_rejection(self):
        """Verify schema rejects resources missing any of the 9 required fields."""
        required_fields = [
            "id", "name", "canonical_url", "source_type", "category",
            "representation", "routing_tags", "acquisition_method",
            "runtime_method", "status"
        ]
        base_res = {
            "id": "valid-id",
            "name": "Test",
            "canonical_url": "https://example.com",
            "source_type": "github_repository",
            "category": ["TEST"],
            "representation": "dependency",
            "routing_tags": ["test"],
            "acquisition_method": "npm",
            "runtime_method": "project_scoped_install",
            "status": "ACTIVE"
        }
        for rf in required_fields:
            bad_res = dict(base_res)
            del bad_res[rf]
            errors = list(self.validator_res.iter_errors([bad_res]))
            self.assertGreater(len(errors), 0, f"Schema should reject missing required field {rf}")

    def test_invalid_id_patterns_rejection(self):
        """Verify schema rejects non-kebab-case IDs (uppercase, spaces, special characters)."""
        bad_ids = ["Resource_1", "resource 1", "resource@name", "RESOURCE", "res!id", "res/sub", "res.id"]
        for bad_id in bad_ids:
            res = {
                "id": bad_id,
                "name": "Bad ID Resource",
                "canonical_url": "https://example.com",
                "source_type": "github_repository",
                "category": ["TEST"],
                "representation": "dependency",
                "routing_tags": ["test"],
                "acquisition_method": "npm",
                "runtime_method": "project_scoped_install",
                "status": "ACTIVE"
            }
            errors = list(self.validator_res.iter_errors([res]))
            self.assertGreater(len(errors), 0, f"Schema should have rejected invalid ID pattern: {bad_id}")

    def test_invalid_urls_rejection(self):
        """Verify schema rejects malformed/non-URI strings in canonical_url."""
        bad_urls = ["not-a-url", "://missing-scheme", "just text with spaces", ""]
        for bad_url in bad_urls:
            res = {
                "id": "valid-id",
                "name": "Test Resource",
                "canonical_url": bad_url,
                "source_type": "github_repository",
                "category": ["TEST"],
                "representation": "dependency",
                "routing_tags": ["test"],
                "acquisition_method": "npm",
                "runtime_method": "project_scoped_install",
                "status": "ACTIVE"
            }
            errors = list(self.validator_res.iter_errors([res]))
            self.assertGreater(len(errors), 0, f"Format checker should reject invalid URL: {bad_url}")

    def test_live_resources_have_valid_urls(self):
        """Empirically inspect all URLs across all live resources for proper scheme (http/https/git) and netloc."""
        url_errors = []
        for r in self.resources:
            for field in ["canonical_url", "source_repository", "documentation_url"]:
                if field in r and r[field]:
                    url = r[field]
                    parsed = urlparse(url)
                    if parsed.scheme not in ["http", "https", "git"]:
                        url_errors.append(f"{r['id']}.{field}: invalid scheme {parsed.scheme} ({url})")
                    if not parsed.netloc:
                        url_errors.append(f"{r['id']}.{field}: missing host ({url})")
        self.assertEqual(url_errors, [], f"Invalid URLs found in resources.json: {url_errors}")

    def test_no_empty_fields_in_live_resources(self):
        """Verify no empty string or null values exist in live resources.json."""
        empty_fields = []
        for r in self.resources:
            for k, v in r.items():
                if v is None:
                    empty_fields.append(f"{r['id']}.{k} is None")
                elif isinstance(v, str) and len(v.strip()) == 0:
                    empty_fields.append(f"{r['id']}.{k} is empty string")
                elif isinstance(v, list) and len(v) == 0:
                    if k in ["category", "routing_tags"]:
                        empty_fields.append(f"{r['id']}.{k} is empty list (required)")
        self.assertEqual(empty_fields, [], f"Empty fields found in live resources: {empty_fields}")

    def test_unknown_enum_types_rejection(self):
        """Verify schema rejects unknown source_type, representation, acquisition_method, runtime_method, and status."""
        invalid_types = [
            ("source_type", "magic_wand"),
            ("representation", "superpower"),
            ("acquisition_method", "telepathy"),
            ("runtime_method", "mind_upload"),
            ("status", "EXPLODING")
        ]
        base_res = {
            "id": "valid-id",
            "name": "Test Resource",
            "canonical_url": "https://example.com",
            "source_type": "github_repository",
            "category": ["TEST"],
            "representation": "dependency",
            "routing_tags": ["test"],
            "acquisition_method": "npm",
            "runtime_method": "project_scoped_install",
            "status": "ACTIVE"
        }
        for field, invalid_val in invalid_types:
            bad_res = dict(base_res)
            bad_res[field] = invalid_val
            errors = list(self.validator_res.iter_errors([bad_res]))
            self.assertGreater(len(errors), 0, f"Schema should have rejected {field}={invalid_val}")

    def test_unexpected_additional_properties_rejection(self):
        """Verify schema rejects extra unrecognized properties."""
        res_with_extra = {
            "id": "valid-id",
            "name": "Test Resource",
            "canonical_url": "https://example.com",
            "source_type": "github_repository",
            "category": ["TEST"],
            "representation": "dependency",
            "routing_tags": ["test"],
            "acquisition_method": "npm",
            "runtime_method": "project_scoped_install",
            "status": "ACTIVE",
            "malicious_injected_field": True
        }
        errors = list(self.validator_res.iter_errors([res_with_extra]))
        self.assertGreater(len(errors), 0, "Schema should have rejected unexpected additional property")

    def test_graph_invalid_semantic_version_rejection(self):
        """Verify graph schema rejects non-semver version strings."""
        bad_versions = ["v1.0.0", "1.0", "beta", "1.0.0.0", "1.a.2"]
        bad_graph = dict(self.graph)
        for bad_v in bad_versions:
            bad_graph["version"] = bad_v
            errors = list(self.validator_graph.iter_errors(bad_graph))
            self.assertGreater(len(errors), 0, f"Graph schema should have rejected version {bad_v}")

class TestQuarantineBoundaryStrictness(unittest.TestCase):
    """Test suite 4: Verify quarantine isolation of the 1,598-skill library."""

    def test_no_quarantined_paths_in_resources_catalog(self):
        """Verify no resource URLs or metadata contain references to the quarantined skills_library."""
        banned = ["skills_library", "skills-library"]
        violations = []
        resources = load_json_file(WORKFLOW_RESOURCES_PATH)
        for r in resources:
            for k, v in r.items():
                if isinstance(v, str):
                    for b in banned:
                        if b in v.lower():
                            violations.append(f"Resource {r.get('id')} field {k} contains banned substring: {v}")
                elif isinstance(v, list):
                    for item in v:
                        if isinstance(item, str):
                            for b in banned:
                                if b in item.lower():
                                    violations.append(f"Resource {r.get('id')} field {k} contains banned substring: {item}")
        self.assertEqual(violations, [], f"Quarantine boundary violations found: {violations}")

    def test_no_quarantined_paths_in_design_graph(self):
        """Verify design graph contains no references to quarantined skills_library."""
        banned = ["skills_library", "skills-library"]
        violations = []
        with open(WORKFLOW_GRAPH_PATH, "r", encoding="utf-8-sig") as f:
            raw_text = f.read()
        for b in banned:
            if b in raw_text.lower():
                violations.append(f"Design graph contains quarantined substring: {b}")
        self.assertEqual(violations, [], f"Quarantine boundary violations in graph: {violations}")

if __name__ == "__main__":
    unittest.main(verbosity=2)
