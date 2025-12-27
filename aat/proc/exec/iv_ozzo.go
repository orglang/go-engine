package exec

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"orglang/orglang/avt/id"
)

func (dto SpecME) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.ProcID, id.Required...),
		validation.Field(&dto.PoolID, id.Required...),
		validation.Field(&dto.Term, validation.Required),
	)
}
