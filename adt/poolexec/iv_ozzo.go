package poolexec

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"orglang/orglang/adt/identity"
)

func (dto PoolSpecME) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.SupID, identity.Optional...),
	)
}

func (dto StepSpecME) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.PoolID, identity.Required...),
		validation.Field(&dto.ProcID, identity.Required...),
		validation.Field(&dto.Term, validation.Required),
	)
}
