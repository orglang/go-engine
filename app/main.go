package main

import (
	"go.uber.org/fx"

	"orglang/go-engine/lib/db"
	"orglang/go-engine/lib/kv"
	"orglang/go-engine/lib/lf"
	"orglang/go-engine/lib/wp"
	"orglang/go-engine/lib/ws"

	"orglang/go-engine/pool/commexch"
	"orglang/go-engine/pool/commturn"
	poolconfexec "orglang/go-engine/pool/compexec"
	compvar "orglang/go-engine/pool/compvar/adapter"
	pooltermdef "orglang/go-engine/pool/termdef"
	pooltypedef "orglang/go-engine/pool/typedef/adapter"
	pooltypeexp "orglang/go-engine/pool/typeexp/adapter"
	proccommexchadapter "orglang/go-engine/proc/commexch/adapter"

	"orglang/go-engine/proc/compexec"
	"orglang/go-engine/proc/termdec"
	termdef "orglang/go-engine/proc/termdef/adapter"
	typedef "orglang/go-engine/proc/typedef/adapter"
	typeexp "orglang/go-engine/proc/typeexp/adapter"

	"orglang/go-engine/app/web"
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
		commexch.Module,
		pooltermdef.Module,
		poolconfexec.Module,
		commturn.Module,
		compvar.Module,
		typedef.Module,
		typeexp.Module,

		proccommexchadapter.Module,
		termdef.Module,
		termdec.Module,
		compexec.Module,
		// app
		web.Module,
	).Run()
}
