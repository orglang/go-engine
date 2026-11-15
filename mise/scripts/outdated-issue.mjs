// Keep one open issue listing mise tools whose newest release is not the
// version pinned by the project's mise lockfiles. Run weekly by
// .github/workflows/mise-tools.yaml; requires mise, gh, and GH_TOKEN.

import { execFileSync } from "node:child_process";

const title = "Update the mise tools";

const environments = [
  { name: "default", args: [] },
  { name: "specs", args: ["-E", "specs"] },
  { name: "sources", args: ["-E", "sources"] },
];

function output(command, args) {
  return execFileSync(command, args, { encoding: "utf8" }).trim();
}

function isOutdated(value) {
  return (
    typeof value === "object" &&
    value !== null &&
    typeof value.name === "string" &&
    typeof value.current === "string" &&
    typeof value.latest === "string"
  );
}

function outdatedTools() {
  return environments.flatMap((environment) => {
    const args = [...environment.args, "outdated", "--local", "--json"];
    const parsed = JSON.parse(output("mise", args));

    if (typeof parsed !== "object" || parsed === null || Array.isArray(parsed)) {
      throw new Error(`mise ${args.join(" ")} returned an unexpected JSON value`);
    }

    return Object.entries(parsed)
      .toSorted(([a], [b]) => a.localeCompare(b))
      .map(([name, tool]) => ({ name, ...tool }))
      .filter(isOutdated)
      .map((tool) => ({ ...tool, environment: environment.name }));
  });
}

function openIssue() {
  const issues = JSON.parse(
    output("gh", [
      "issue", "list",
      "--state", "open",
      "--search", `"${title}" in:title`,
      "--json", "number,title",
    ]),
  );

  if (!Array.isArray(issues)) return undefined;
  const found = issues.find((issue) => issue.title === title);
  return found?.number;
}

function bodyOf(tools) {
  const rows = tools.map((tool) => {
    const newest = tool.release_url
      ? `[${tool.latest}](${tool.release_url})`
      : tool.latest;
    const prefix = tool.environment === "default"
      ? "mise upgrade"
      : `mise -E ${tool.environment} upgrade`;
    return `| ${tool.environment} | ${tool.name} | ${tool.current} | ${newest} | \`${prefix} ${tool.name}\` |`;
  });

  return `These tools are declared in \`mise.toml\` and pinned by the corresponding mise lockfiles.
The list is refreshed weekly. Update a tool with the command in the last column, run the
relevant project checks, and commit the resulting lockfile changes.

| Environment | Tool | Locked | Newest | Upgrade command |
|---|---|---|---|---|
${rows.join("\n")}
`;
}

const tools = outdatedTools();
const number = openIssue();

if (tools.length === 0) {
  if (number !== undefined) {
    execFileSync(
      "gh",
      ["issue", "close", String(number), "--comment", "All mise tools are up to date."],
      { stdio: "inherit" },
    );
  }
} else if (number === undefined) {
  execFileSync(
    "gh",
    ["issue", "create", "--title", title, "--body", bodyOf(tools)],
    { stdio: "inherit" },
  );
} else {
  execFileSync(
    "gh",
    ["issue", "edit", String(number), "--body", bodyOf(tools)],
    { stdio: "inherit" },
  );
}
