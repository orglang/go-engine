# Explicit application lifecycle

The `app.App` type owns one `samber/do` injector and one application lifecycle coordinator. It does not replace `do`'s own service shutdown behavior: `App.Stop(ctx)` first stops registered application hooks, then calls `Injector.ShutdownWithContext(ctx)` for instantiated services that implement one of `do`'s shutdown interfaces.

## Startup and shutdown order

- Constructors register `app.Hook` values through the `app.Lifecycle` dependency.
- The composition root resolves all required lazy services before calling `App.Start(ctx)`.
- Startup hooks run in registration order.
- If startup fails, hooks that completed startup are rolled back in reverse order. The failed hook is not considered started.
- Normal shutdown runs hooks in reverse order, then shuts down the DI container. Both shutdown phases are attempted and their errors are joined.
- Hooks cannot be appended after startup begins. `Stop` is idempotent; repeated `Start` is an error.

## Constructor example

```go
func newHTTPServer(i do.Injector) (*HTTPServer, error) {
    lifecycle := do.MustInvoke[app.Lifecycle](i)
    server := &HTTPServer{/* ... */}

    if err := lifecycle.Append(app.Hook{
        OnStart: server.Start,
        OnStop:  server.Stop,
    }); err != nil {
        return nil, err
    }

    return server, nil
}
```

The constructor should register the hook only; it should not start the server. The composition root is responsible for resolving the server and other required entrypoint services before starting the app.

Resources may alternatively implement the `do.ShutdownerWithContextAndError` interface and rely on container shutdown. Avoid registering the same close operation both as an application hook and as a `do` shutdown method.

## Scope

This is lifecycle infrastructure only. The production Fx composition root and existing Fx modules remain unchanged. Migration of `db`, `ws`, and `wp` providers is a subsequent step so their startup readiness, shutdown semantics, and ownership can be changed and tested individually.
