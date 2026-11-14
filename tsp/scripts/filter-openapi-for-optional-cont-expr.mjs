import { readFile, writeFile } from "node:fs/promises";
import { parseDocument, stringify } from "yaml";

const [input, output] = process.argv.slice(2);
if (!input || !output) {
  throw new Error("usage: filter-openapi-for-optional-cont-expr.mjs <input.yaml> <output.yaml>");
}

const document = parseDocument(await readFile(input, "utf8"));
if (document.errors.length > 0) {
  throw new AggregateError(document.errors, "failed to parse OpenAPI document");
}

const paths = document.get("paths");
if (!paths) throw new Error("OpenAPI document has no paths");

const retained = new Set([
  "/api/v1/pools/executions/spawns",
  "/api/v1/pools/executions/steps",
]);
for (const key of [...paths.keys()]) {
  if (!retained.has(key)) paths.delete(key);
}
if (paths.size !== retained.size) {
  throw new Error(`expected to retain ${retained.size} Pool.StepSpec operations, found ${paths.size}`);
}

await writeFile(output, stringify(document));
console.log("retained only Pool.StepSpec spawn/step operations for the cont_expr recursion experiment");
