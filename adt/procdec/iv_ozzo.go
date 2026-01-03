package procdec

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"orglang/orglang/adt/identity"
	"orglang/orglang/adt/qualsym"
	"orglang/orglang/adt/termctx"
)

func (dto DecSpecME) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.X, validation.Required),
		validation.Field(&dto.ProcQN, qualsym.Required...),
		validation.Field(&dto.Ys, termctx.Optional...),
	)
}

func (dto IdentME) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.DecID, identity.Required...),
	)
}

func (dto DecRefME) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.DecID, identity.Required...),
	)
}

func (dto DecSpecVP) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.ProcNS, qualsym.Required...),
		validation.Field(&dto.ProcSN, qualsym.Required...),
	)
}

func (dto DecRefVP) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.DecID, identity.Required...),
	)
}
