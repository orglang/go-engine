# Temporary Go client generator comparison

This directory is a disposable prototype for issue
[orglang/go-workspace#114](https://github.com/orglang/go-workspace/issues/114).

Both generators consume the same committed OpenAPI 3.1 document:

- `ogen v1.24.0`: client-only generation, configured in `ogen.yml`.
- `oapi-codegen v2.8.2`: generated models and HTTP client.

Run the prototype through the repository's existing Generate artifacts
workflow (or locally with `mise run generate` on this branch).

Compare the outputs on:

1. Whether the complete contract generates without unsupported-schema errors.
2. Generated API ergonomics: operation names, request/response types, optional
   and nullable values, and discriminated expression unions.
3. Runtime dependencies and generated code size.
4. Error/response handling, customization points, and maintainability.
5. Reproducibility and ease of pinning the generator version.

No generated package here is intended for production use. All files in this
directory and the temporary task/workflow wiring will be removed after the
comparison is recorded.
