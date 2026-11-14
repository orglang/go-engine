import { readFile, writeFile } from "node:fs/promises";
import { parseDocument, stringify } from "yaml";

const [input, output] = process.argv.slice(2);
if (!input || !output) throw new Error("usage: filter-openapi-for-optional-cont-expr.mjs <input.yaml> <output.yaml>");
const document = parseDocument(await readFile(input, "utf8"));
if (document.errors.length > 0) throw new AggregateError(document.errors, "failed to parse OpenAPI document");
const paths = document.get("paths");
if (!paths || !Array.isArray(paths.items)) throw new Error("OpenAPI document has no paths map");
const retained = new Set(["/api/v1/pools/executions/spawns", "/api/v1/pools/executions/steps"]);
for (const pair of [...paths.items]) if (!retained.has(pair.key.value)) paths.delete(pair.key);
if (paths.items.length !== retained.size) throw new Error(`expected to retain ${retained.size} Pool.StepSpec operations, found ${paths.items.length}`);
await writeFile(output, stringify(document));
console.log("retained only Pool.StepSpec spawn/step operations for the cont_expr recursion experiment");
