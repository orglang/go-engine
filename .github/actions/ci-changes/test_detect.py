import importlib.util
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

    def test_renamed_paths_can_be_evaluated_independently(self):
        changed = ci_detect.classify_paths(
            ["tsp/main.tsp", "internal/app/service.go"]
        )

        self.assertTrue(changed["app-specs"])
        self.assertTrue(changed["app-sources"])

    def test_unrelated_documentation_does_not_select_checks(self):
        changed = ci_detect.classify_paths(["README.md", "docs/adr/design.md"])

        self.assertFalse(any(changed.values()))


if __name__ == "__main__":
    unittest.main()
