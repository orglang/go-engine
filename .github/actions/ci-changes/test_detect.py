import importlib.util
import os
import subprocess
import tempfile
import unittest
from pathlib import Path

MODULE_PATH = Path(__file__).with_name("detect.py")
SPEC = importlib.util.spec_from_file_location("ci_detect", MODULE_PATH)
ci_detect = importlib.util.module_from_spec(SPEC)
assert SPEC.loader is not None
SPEC.loader.exec_module(ci_detect)


class ClassifyPathsTest(unittest.TestCase):
    def test_application_dependencies_are_not_sources(self):
        changed = ci_detect.classify_paths(["go.mod", "go.sum"])

        self.assertTrue(changed["app-deps"])
        self.assertFalse(changed["app-sources"])
        self.assertFalse(changed["tool-deps"])

    def test_tool_dependencies(self):
        changed = ci_detect.classify_paths(
            ["tsp/package.json", "tsp/package-lock.json"]
        )

        self.assertTrue(changed["tool-deps"])
        self.assertFalse(changed["app-deps"])
        self.assertFalse(changed["app-sources"])

    def test_specs_and_generated_openapi(self):
        changed = ci_detect.classify_paths(
            ["tsp/main.tsp", "tsp/models/pool/comp.tsp", "openapi/openapi.yaml"]
        )

        self.assertTrue(changed["app-specs"])
        self.assertFalse(changed["tool-specs"])

    def test_spec_compiler_configuration_is_tool_specs(self):
        changed = ci_detect.classify_paths(["tsp/tspconfig.yaml"])

        self.assertTrue(changed["tool-specs"])
        self.assertFalse(changed["app-specs"])

    def test_go_and_sql_sources(self):
        changed = ci_detect.classify_paths(
            ["internal/app/service.go", "db/postgres/owner/000001_initial.up.sql"]
        )

        self.assertTrue(changed["app-sources"])
        self.assertFalse(changed["app-deps"])

    def test_binary_entrypoint_is_classified_as_app_binary_input(self):
        changed = ci_detect.classify_paths(["cmd/engine/main.go"])

        self.assertTrue(changed["app-binaries"])
        self.assertTrue(changed["app-sources"])

    def test_docker_and_compose_files_are_distribution_inputs(self):
        changed = ci_detect.classify_paths(
            ["db/postgres/migrations.Dockerfile", "db/postgres/compose.yaml"]
        )

        self.assertTrue(changed["app-distros"])
        self.assertFalse(changed["app-binaries"])

    def test_workflow_change_invalidates_all_groups(self):
        changed = ci_detect.classify_paths(
            [".github/workflows/sources-verify.yaml"]
        )

        self.assertTrue(all(changed.values()))

    def test_task_definition_change_invalidates_all_groups(self):
        changed = ci_detect.classify_paths(["db/taskfile.yaml"])

        self.assertTrue(all(changed.values()))

    def test_mise_toolchain_change_invalidates_all_groups(self):
        changed = ci_detect.classify_paths(["mise.sources-dev.lock"])

        self.assertTrue(all(changed.values()))

    def test_multiple_changed_paths_accumulate_groups(self):
        changed = ci_detect.classify_paths(
            ["tsp/main.tsp", "internal/app/service.go", "go.mod"]
        )

        self.assertTrue(changed["app-specs"])
        self.assertTrue(changed["app-sources"])
        self.assertTrue(changed["app-deps"])

    def test_unrelated_documentation_does_not_select_checks(self):
        changed = ci_detect.classify_paths(["README.md", "docs/adr/design.md"])

        self.assertFalse(any(changed.values()))


class GitDiffTest(unittest.TestCase):
    def test_rename_is_reported_as_old_and_new_paths(self):
        with tempfile.TemporaryDirectory() as directory:
            repo = Path(directory)
            subprocess.run(["git", "init", "-q", str(repo)], check=True)
            subprocess.run(
                ["git", "-C", str(repo), "config", "user.email", "ci@example.test"],
                check=True,
            )
            subprocess.run(
                ["git", "-C", str(repo), "config", "user.name", "CI Test"],
                check=True,
            )

            old_path = repo / "README.md"
            old_path.write_text("content\\n", encoding="utf-8")
            subprocess.run(["git", "-C", str(repo), "add", "."], check=True)
            subprocess.run(
                ["git", "-C", str(repo), "commit", "-qm", "base"], check=True
            )
            base_sha = subprocess.check_output(
                ["git", "-C", str(repo), "rev-parse", "HEAD"], text=True
            ).strip()

            new_path = repo / "internal" / "service.go"
            new_path.parent.mkdir(parents=True)
            old_path.rename(new_path)
            subprocess.run(["git", "-C", str(repo), "add", "-A"], check=True)
            subprocess.run(
                ["git", "-C", str(repo), "commit", "-qm", "rename"], check=True
            )
            head_sha = subprocess.check_output(
                ["git", "-C", str(repo), "rev-parse", "HEAD"], text=True
            ).strip()

            previous_directory = Path.cwd()
            try:
                os.chdir(repo)
                paths = ci_detect.changed_paths(base_sha, head_sha)
            finally:
                os.chdir(previous_directory)

            self.assertEqual(set(paths), {"README.md", "internal/service.go"})
            changed = ci_detect.classify_paths(paths)
            self.assertTrue(changed["app-sources"])


class EventBehaviorTest(unittest.TestCase):
    def test_non_pull_request_events_select_all_groups(self):
        for event_name in ("push", "merge_group", "workflow_dispatch"):
            with self.subTest(event_name=event_name):
                changed = ci_detect.groups_for_event(event_name)
                self.assertTrue(all(changed.values()))

    def test_pull_request_requires_both_commit_shas(self):
        with self.assertRaises(ValueError):
            ci_detect.groups_for_event("pull_request")


if __name__ == "__main__":
    unittest.main()
