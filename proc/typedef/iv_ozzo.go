package typedef

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/orglang/go-sdk/adt/uniqsym"
)

func (dto DefSpecVP) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.TypeQN, uniqsym.Required...),
	)
}
