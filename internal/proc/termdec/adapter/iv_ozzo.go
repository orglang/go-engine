package adapter

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var qualifiedSymbolPattern = regexp.MustCompile(`^[0-9A-Za-z_-]+(\.[0-9A-Za-z_-]+)*$`)

func (dto DecSpecVP) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.TermQN,
			validation.Required,
			validation.Length(1, 512),
			validation.Match(qualifiedSymbolPattern),
		),
	)
}
