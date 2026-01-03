package typedef

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"orglang/orglang/adt/identity"
	"orglang/orglang/adt/qualsym"
	"orglang/orglang/adt/revnum"
)

func (dto DefSpecME) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.TypeQN, qualsym.Required...),
		validation.Field(&dto.TypeES, validation.Required),
	)
}

func (dto IdentME) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.DefID, identity.Required...),
	)
}

func (dto DefRefME) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.DefID, identity.Required...),
		validation.Field(&dto.DefRN, revnum.Optional...),
	)
}

func (dto DefSnapME) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.DefID, identity.Required...),
		validation.Field(&dto.DefRN, revnum.Optional...),
		validation.Field(&dto.TypeES, validation.Required),
	)
}

func (dto DefSpecVP) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.TypeNS, qualsym.Required...),
		validation.Field(&dto.TypeSN, qualsym.Required...),
	)
}

func (dto DefRefVP) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.DefID, identity.Required...),
		validation.Field(&dto.DefRN, revnum.Optional...),
		validation.Field(&dto.Title, qualsym.Required...),
	)
}
