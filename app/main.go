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
	"orglang/go-engine/pool/compvar"
	pooltermdef "orglang/go-engine/pool/termdef"
	pooltypedef "orglang/go-engine/pool/typedef"
	pooltypeexp "orglang/go-engine/pool/typeexp"
	proccommexch "orglang/go-engine/proc/commexch"
	"orglang/go-engine/proc/compexec"
	"orglang/go-engine/proc/termdec"
	"orglang/go-engine/proc/termdef"
	"orglang/go-engine/proc/typedef"
	"orglang/go-engine/proc/typeexp"

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
		proccommexch.Module,
		termdef.Module,
		termdec.Module,
		compexec.Module,
		// app
		web.Module,
	).Run()
}
