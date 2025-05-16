//go:build !goverter

package main

import (
	"go.uber.org/fx"

	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/data"
	"smecalculus/rolevod/lib/msg"

	"smecalculus/rolevod/internal/alias"

	poolexec "smecalculus/rolevod/app/pool/exec"
	procdec "smecalculus/rolevod/app/proc/dec"
	procdef "smecalculus/rolevod/app/proc/def"
	procexec "smecalculus/rolevod/app/proc/exec"
	typedef "smecalculus/rolevod/app/type/def"

	"smecalculus/rolevod/app/web"
)

func main() {
	fx.New(
		// lib
		core.Module,
		data.Module,
		msg.Module,
		alias.Module,
		// app
		procdef.Module,
		poolexec.Module,
		typedef.Module,
		procexec.Module,
		procdec.Module,
		web.Module,
	).Run()
}
