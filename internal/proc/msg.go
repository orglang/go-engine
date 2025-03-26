package proc

type RefMsg struct {
	ProcID string `json:"proc_id" param:"id"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
var (
	MsgToRef   func(RefMsg) (Ref, error)
	MsgFromRef func(Ref) RefMsg
)
