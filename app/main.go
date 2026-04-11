package main

import (
	"go.uber.org/fx"

	"orglang/go-engine/lib/db"
	"orglang/go-engine/lib/kv"
	"orglang/go-engine/lib/lf"
	"orglang/go-engine/lib/wp"
	"orglang/go-engine/lib/ws"

	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/descsem"
	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/poolcomm"
	"orglang/go-engine/adt/pooldec"
	"orglang/go-engine/adt/poolexec"
	"orglang/go-engine/adt/poolimpl"
	"orglang/go-engine/adt/poolstep"
	"orglang/go-engine/adt/poolvar"
	"orglang/go-engine/adt/procconn"
	"orglang/go-engine/adt/procdec"
	"orglang/go-engine/adt/procdef"
	"orglang/go-engine/adt/procexec"
	"orglang/go-engine/adt/typedef"
	"orglang/go-engine/adt/typeexp"
	"orglang/go-engine/adt/xactdef"
	"orglang/go-engine/adt/xactexp"

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
		commsem.Module,
		descsem.Module,
		implsem.Module,
		xactdef.Module,
		xactexp.Module,
		poolcomm.Module,
		pooldec.Module,
		poolexec.Module,
		poolimpl.Module,
		poolstep.Module,
		poolvar.Module,
		typedef.Module,
		typeexp.Module,
		procconn.Module,
		procdef.Module,
		procdec.Module,
		procexec.Module,
		// app
		web.Module,
	).Run()
}
