package exec

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"orglang/orglang/avt/id"
)

func (dto PoolSpecME) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.SupID, id.Optional...),
	)
}

func (dto StepSpecME) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.PoolID, id.Required...),
		validation.Field(&dto.ProcID, id.Required...),
		validation.Field(&dto.Term, validation.Required),
	)
}
