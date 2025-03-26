package pool

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/step"
)

type SpecMsg struct {
	Title  string   `json:"title"`
	SupID  string   `json:"sup_id"`
	DepIDs []string `json:"dep_ids"`
}

func (dto SpecMsg) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.SupID, id.Optional...),
	)
}

type IdentMsg struct {
	PoolID string `json:"id" param:"id"`
}

type RefMsg struct {
	PoolID string `json:"id" param:"id"`
	Title  string `json:"title"`
}

type SnapMsg struct {
	PoolID string   `json:"id"`
	Title  string   `json:"title"`
	Subs   []RefMsg `json:"subs"`
}

type RootMsg struct {
	PoolID string `json:"id"`
	ProcID string `json:"proc_id"`
	Title  string `json:"title"`
	SupID  string `json:"sup_id"`
	Rev    int64  `json:"rev"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
var (
	MsgToSpec    func(SpecMsg) (Spec, error)
	MsgFromSpec  func(Spec) SpecMsg
	MsgToRoot    func(RootMsg) (Root, error)
	MsgFromRoot  func(Root) RootMsg
	MsgFromRoots func([]Root) []RootMsg
	MsgToSnap    func(SnapMsg) (SubSnap, error)
	MsgFromSnap  func(SubSnap) SnapMsg
)

type TranSpecMsg struct {
	PoolID string       `json:"pool_id"`
	ProcID string       `json:"proc_id"`
	Term   step.TermMsg `json:"term"`
}

func (dto TranSpecMsg) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.PoolID, id.Required...),
		validation.Field(&dto.ProcID, id.Required...),
		validation.Field(&dto.Term, validation.Required),
	)
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
// goverter:extend smecalculus/rolevod/internal/step:Msg.*
var (
	MsgFromTranSpec func(TranSpec) TranSpecMsg
	MsgToTranSpec   func(TranSpecMsg) (TranSpec, error)
)
