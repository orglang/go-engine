package ws

import (
	"slices"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (dto exchangeCS) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.Protocol, validation.Required),
		validation.Field(&dto.Server, validation.Required),
	)
}

func (dto protocolCS) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.Modes, validation.Required, validation.Each(validation.In(httpProto))),
		validation.Field(&dto.HTTP, validation.Required.When(slices.Contains(dto.Modes, httpProto))),
	)
}

const (
	MinPort = 80
	MaxPort = 20000
)

func (dto httpCS) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.Port, validation.Required, validation.Min(MinPort), validation.Max(MaxPort)),
	)
}

func (dto serverCS) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.Mode, validation.Required, validation.In(echoServer)),
		// validation.Field(&dto.Echo, validation.Required.When(dto.Mode == echoMode)),
	)
}

func (dto echoCS) Validate() error {
	return validation.ValidateStruct(&dto)
}
