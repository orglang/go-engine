package ep

import (
	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/ph"
	"smecalculus/rolevod/lib/sym"
)

type Spec struct {
	ChnlPH ph.ADT
	RoleQN sym.ADT
}

type Impl struct {
	ChnlPH  ph.ADT
	ChnlID  id.ADT
	StateID id.ADT
}
