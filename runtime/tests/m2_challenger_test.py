#!/usr/bin/env python3
"""
M2 Challenger Adversarial Stress Test Suite for Orchestra V3 Design-First Go Engine
Tests CLI integration, state transitions, gate halting, quarantine defense,
artifact generation, and determinism.
"""

import json
import os
import shutil
import subprocess
import tempfile
import unittest

ORCHESTRA_CMD = ["go", "run", "cmd/orchestra/main.go"]
RUNTIME_DIR = r"C:\projects\orchestra-workflow\runtime"


class TestM2CLIAdversarialExecution(unittest.TestCase):
    """Adversarial stress-testing of the CLI interface and 8-stage pipeline."""

    def setUp(self):
        self.test_dir = tempfile.mkdtemp(prefix="orch_m2_challenger_")

    def tearDown(self):
        shutil.rmtree(self.test_dir, ignore_errors=True)

    def run_orchestra(self, subcmd, *args, expect_exit=0):
        cmd = ORCHESTRA_CMD + [subcmd] + list(args)
        proc = subprocess.run(
            cmd,
            cwd=RUNTIME_DIR,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            timeout=30,
        )
        if expect_exit is not None and proc.returncode != expect_exit:
            self.fail(
                f"Command '{' '.join(cmd)}' failed with exit code {proc.returncode} (expected {expect_exit}).\n"
                f"Stdout: {proc.stdout}\nStderr: {proc.stderr}"
            )
        return proc

    def test_cli_doctor_diagnostic(self):
        """Verify doctor diagnostic inspects all components and enforces quarantine."""
        proc = self.run_orchestra("doctor", expect_exit=0)
        self.assertIn("=== Orchestra V3 Environment & System Doctor ===", proc.stdout)
        self.assertIn("Go Runtime:", proc.stdout)
        self.assertIn("canonical resources", proc.stdout)
        self.assertIn("Quarantine:      [PASS]", proc.stdout)

    def test_cli_human_gate_halting_without_auto_approve(self):
        """Verify high-visual task exits with human gate halt (os.Exit(2) wrapped by go run as exit code 1 with stderr 'exit status 2')."""
        proc = self.run_orchestra(
            "run",
            "--task",
            "Build award-winning 3D spatial portfolio",
            "--workdir",
            self.test_dir,
            expect_exit=1,
        )
        self.assertIn("[GATE HALTED]", proc.stdout)
        self.assertIn("Human approval required", proc.stdout)
        self.assertIn("exit status 2", proc.stderr)

        # Confirm downstream artifacts were NOT created
        design_md = os.path.join(self.test_dir, ".orchestra", "DESIGN.md")
        self.assertFalse(os.path.exists(design_md), "DESIGN.md created prematurely")

    def test_cli_run_with_auto_approve_produces_all_artifacts(self):
        """Verify full 8-stage execution with --auto-approve produces all required artifacts."""
        proc = self.run_orchestra(
            "run",
            "--task",
            "Build award-winning 3D spatial portfolio showcase",
            "--workdir",
            self.test_dir,
            "--auto-approve",
            expect_exit=0,
        )
        self.assertIn("Status:        SUCCESS", proc.stdout)
        self.assertIn("DESIGN.md", proc.stdout)
        self.assertIn("Reference Log", proc.stdout)
        self.assertIn("Handoff State", proc.stdout)

        # Inspect generated DESIGN.md
        design_md_path = os.path.join(self.test_dir, ".orchestra", "DESIGN.md")
        self.assertTrue(os.path.exists(design_md_path), "DESIGN.md missing on disk")
        with open(design_md_path, "r", encoding="utf-8") as f:
            content = f.read()
            self.assertIn("# DESIGN.md — Design System Contract", content)
            self.assertIn("--bg-base:", content)
            self.assertIn("Display Font", content)

        # Inspect reference-log.md
        ref_log_path = os.path.join(self.test_dir, ".orchestra", "reference-log.md")
        self.assertTrue(os.path.exists(ref_log_path), "reference-log.md missing")
        with open(ref_log_path, "r", encoding="utf-8") as f:
            ref_content = f.read()
            self.assertIn("# Visual Research Reference Log", ref_content)

        # Inspect state.json
        state_path = os.path.join(self.test_dir, ".orchestra", "handoff", "state.json")
        self.assertTrue(os.path.exists(state_path), "state.json missing")
        with open(state_path, "r", encoding="utf-8") as f:
            state = json.load(f)
            self.assertEqual(state["version"], 3)
            self.assertIn("completed_steps", state)
            self.assertIn("Implement", state["completed_steps"])

    def test_cli_quarantine_rejection_at_entry(self):
        """Verify attempt to execute inside quarantined skills_library is rejected."""
        banned = r"C:\Users\mockuser\.gemini\config\skills_library"
        proc = self.run_orchestra(
            "run",
            "--task",
            "Build component",
            "--workdir",
            banned,
            "--auto-approve",
            expect_exit=1,
        )
        self.assertIn("quarantine", proc.stdout.lower() + proc.stderr.lower())

    def test_cli_plan_dry_run_json_output(self):
        """Verify plan command outputs valid JSON with dynamic token calculation."""
        proc = self.run_orchestra(
            "plan",
            "--task",
            "Build 3D interactive WebGL showcase",
            "--output",
            "json",
            expect_exit=0,
        )
        plan = json.loads(proc.stdout)
        self.assertIn("estimated_token_cost", plan)
        self.assertGreater(plan["estimated_token_cost"], 1500)
        self.assertTrue(plan.get("requires_human_gate"))
        self.assertIn("primary_archetype", plan)

    def test_cli_verify_multi_viewport(self):
        """Verify verify command tests desktop, tablet, and mobile viewports."""
        proc = self.run_orchestra(
            "verify",
            "--viewports",
            "desktop,tablet,mobile",
            "--output-dir",
            os.path.join(self.test_dir, "qa"),
            expect_exit=0,
        )
        self.assertIn("=== Orchestra V3 Multi-Viewport Visual QA Runner ===", proc.stdout)
        self.assertIn("Viewport desktop  [PASS]", proc.stdout)
        self.assertIn("Viewport tablet   [PASS]", proc.stdout)
        self.assertIn("Viewport mobile   [PASS]", proc.stdout)
        self.assertIn("[PASS] All viewports verified successfully. Zero horizontal overflow.", proc.stdout)


if __name__ == "__main__":
    unittest.main(verbosity=2)
