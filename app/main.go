//go:build !goverter

package main

import (
	"go.uber.org/fx"

	"orglang/orglang/lib/lf"
	"orglang/orglang/lib/sd"
	"orglang/orglang/lib/ws"

	"orglang/orglang/adt/poolxec"
	"orglang/orglang/adt/procdec"
	"orglang/orglang/adt/procdef"
	"orglang/orglang/adt/procxec"
	"orglang/orglang/adt/syndec"
	"orglang/orglang/adt/typedef"

	"orglang/orglang/app/web"
)

func main() {
	fx.New(
		// lib
		ws.Module,
		sd.Module,
		lf.Module,
		// adt
		syndec.Module,
		poolxec.Module,
		typedef.Module,
		procdef.Module,
		procdec.Module,
		procxec.Module,
		// app
		web.Module,
	).Run()
}
