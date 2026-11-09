package main

import (
	"go.uber.org/fx"

	"orglang/go-engine/internal/lib/db"
	"orglang/go-engine/internal/lib/kv"
	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/wp"
	"orglang/go-engine/internal/lib/ws"

	poolconn "orglang/go-engine/internal/pool/conn/adapter"
	poolconfexec "orglang/go-engine/internal/pool/compexec/adapter"
	compvar "orglang/go-engine/internal/pool/compvar/adapter"
	pooltermdef "orglang/go-engine/internal/pool/termdef/adapter"
	pooltypedef "orglang/go-engine/internal/pool/typedef/adapter"
	pooltypeexp "orglang/go-engine/internal/pool/typeexp/adapter"
	procconnadapter "orglang/go-engine/internal/proc/conn/adapter"

	proccompexecadapter "orglang/go-engine/internal/proc/compexec/adapter"
	termdec "orglang/go-engine/internal/proc/termdec/adapter"
	termdef "orglang/go-engine/internal/proc/termdef/adapter"
	typedef "orglang/go-engine/internal/proc/typedef/adapter"
	typeexp "orglang/go-engine/internal/proc/typeexp/adapter"

	"orglang/go-engine/internal/app/web"
)

func main() {
	fx.New(
		// lib
		db.Module,
		kv.Module,
		lf.Module,
		wp.Module,
		ws.Module,
		// adt
		pooltypedef.Module,
		pooltypeexp.Module,
		poolconn.Module,
		pooltermdef.Module,
		poolconfexec.Module,
		compvar.Module,
		typedef.Module,
		typeexp.Module,

		procconnadapter.Module,
		termdef.Module,
		termdec.Module,
		proccompexecadapter.Module,
		// app
		web.Module,
	).Run()
}
