package adapter

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var qualifiedSymbolPattern = regexp.MustCompile(`^[0-9A-Za-z_-]+(\.[0-9A-Za-z_-]+)*$`)

func (dto DefSpecVP) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.TypeQN,
			validation.Required,
			validation.Length(1, 512),
			validation.Match(qualifiedSymbolPattern),
		),
	)
}
