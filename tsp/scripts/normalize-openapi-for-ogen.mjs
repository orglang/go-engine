import { readFile, writeFile } from "node:fs/promises";
import { parseDocument, stringify } from "yaml";

const [input, output] = process.argv.slice(2);
if (!input || !output) {
  throw new Error("usage: normalize-openapi-for-ogen.mjs <input.yaml> <output.yaml>");
}

const document = parseDocument(await readFile(input, "utf8"));
if (document.errors.length > 0) {
  throw new AggregateError(document.errors, "failed to parse generated OpenAPI document");
}

const schemas = document.get("components")?.get("schemas");
if (!schemas) throw new Error("OpenAPI document has no components.schemas");

let normalized = 0;
for (const pair of schemas.items) {
  const schema = pair.value;
  if (schema.get("oneOf") && schema.get("discriminator") && schema.has("type")) {
    schema.delete("type");
    normalized++;
  }
}

// ogen v1.24.0 currently creates a conflicting struct and sum type name
// for discriminator unions declared with both type: object and oneOf.
// The variants already constrain these schemas to objects, so omitting the
// redundant type is equivalent for these four union schemas.
if (normalized !== 4) {
  throw new Error(`expected to normalize 4 discriminator unions, found ${normalized}`);
}

await writeFile(output, stringify(document));
console.log(`normalized ${normalized} discriminator unions for ogen`);
