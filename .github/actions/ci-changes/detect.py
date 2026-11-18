#!/usr/bin/env python3
"""Classify changed repository paths for CI job selection."""

from __future__ import annotations

import fnmatch
import os
import subprocess
import sys
from collections.abc import Iterable

GROUPS = (
    "app-deps",
    "tool-deps",
    "app-specs",
    "tool-specs",
    "app-sources",
    "tool-sources",
    "app-binaries",
    "tool-binaries",
    "app-distros",
    "tool-distros",
)

# These inputs define CI orchestration or shared task/toolchain behavior. Their
# impact is intentionally broad until a narrower dependency graph is proven.
BROAD_PATTERNS = (
    ".github/workflows/**",
    ".github/actions/**",
    "taskfile.yaml",
    "**/taskfile.yaml",
    "mise*.toml",
    "mise*.lock",
    ".miserc.toml",
)

PATH_PATTERNS = {
    "app-deps": (
        "go.mod",
        "go.sum",
        "**/go.mod",
        "**/go.sum",
    ),
    "tool-deps": (
        "tsp/package.json",
        "tsp/package-lock.json",
    ),
    "app-specs": (
        "tsp/**/*.tsp",
        "tsp/*.tsp",
        "openapi/openapi.yaml",
    ),
    "tool-specs": (
        "tsp/tspconfig.yaml",
    ),
    "app-sources": (
        "*.go",
        "**/*.go",
        "*.sql",
        "**/*.sql",
    ),
    "tool-sources": (
        "scripts/*.py",
        "scripts/**/*.py",
        "scripts/*.sh",
        "scripts/**/*.sh",
        "*.sh",
        "**/*.sh",
    ),
    "app-binaries": (
        "cmd/engine/**",
    ),
    "tool-binaries": (),
    "app-distros": (
        "Dockerfile",
        "Dockerfile.*",
        "**/Dockerfile",
        "**/Dockerfile.*",
        "*.Dockerfile",
        "**/*.Dockerfile",
        ".dockerignore",
        "compose.yaml",
        "compose.yml",
        "*.compose.yaml",
        "*.compose.yml",
        "**/compose.yaml",
        "**/compose.yml",
        "**/*.compose.yaml",
        "**/*.compose.yml",
    ),
    "tool-distros": (
        "mise.distros.toml",
        "mise.distros.lock",
    ),
}


def matches(path: str, patterns: Iterable[str]) -> bool:
    return any(fnmatch.fnmatchcase(path, pattern) for pattern in patterns)


def classify_paths(paths: Iterable[str]) -> dict[str, bool]:
    """Return true for every group directly touched by the supplied paths."""
    changed = {group: False for group in GROUPS}

    for raw_path in paths:
        path = raw_path.strip().replace("\\", "/").removeprefix("./")
        if not path:
            continue

        if matches(path, BROAD_PATTERNS):
            changed = {group: True for group in GROUPS}
            continue

        for group, patterns in PATH_PATTERNS.items():
            if matches(path, patterns):
                changed[group] = True

    return changed


def changed_paths(base_sha: str, head_sha: str) -> list[str]:
    if not base_sha or not head_sha:
        raise ValueError("base_sha and head_sha are required for pull_request events")

    # Disable rename detection so both sides of a move are evaluated as paths.
    result = subprocess.run(
        [
            "git",
            "diff",
            "--name-only",
            "--no-renames",
            f"{base_sha}...{head_sha}",
        ],
        check=True,
        capture_output=True,
        text=True,
    )
    return result.stdout.splitlines()


def groups_for_event(
    event_name: str, base_sha: str = "", head_sha: str = ""
) -> dict[str, bool]:
    if event_name == "pull_request":
        return classify_paths(changed_paths(base_sha, head_sha))

    # Do not accidentally filter merge queue, push, or manual runs.
    return {group: True for group in GROUPS}


def main() -> int:
    groups = groups_for_event(
        os.environ.get("EVENT_NAME", ""),
        os.environ.get("BASE_SHA", ""),
        os.environ.get("HEAD_SHA", ""),
    )

    output_path = os.environ.get("GITHUB_OUTPUT")
    if not output_path:
        print("GITHUB_OUTPUT is not set", file=sys.stderr)
        return 1

    with open(output_path, "a", encoding="utf-8") as output:
        for group in GROUPS:
            output.write(f"{group}={str(groups[group]).lower()}\\n")

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
