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
		validation.Field(&dto.TypeTS, validation.Required),
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
		validation.Field(&dto.TypeTS, validation.Required),
	)
}

func (dto TermSpecME) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.K, kindRequired...),
		validation.Field(&dto.Link, validation.Required.When(dto.K == LinkKind), validation.Skip.When(dto.K != LinkKind)),
		validation.Field(&dto.Tensor, validation.Required.When(dto.K == TensorKind), validation.Skip.When(dto.K != TensorKind)),
		validation.Field(&dto.Lolli, validation.Required.When(dto.K == LolliKind), validation.Skip.When(dto.K != LolliKind)),
		validation.Field(&dto.Plus, validation.Required.When(dto.K == PlusKind), validation.Skip.When(dto.K != PlusKind)),
		validation.Field(&dto.With, validation.Required.When(dto.K == WithKind), validation.Skip.When(dto.K != WithKind)),
	)
}

func (dto LinkSpecME) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.TypeQN, qualsym.Required...),
	)
}

func (dto ProdSpecME) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.ValTS, validation.Required),
		validation.Field(&dto.ContTS, validation.Required),
	)
}

func (dto SumSpecME) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.Choices,
			validation.Required,
			validation.Length(1, 10),
			validation.Each(validation.Required),
		),
	)
}

func (dto ChoiceSpecME) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.Label, qualsym.Required...),
		validation.Field(&dto.ContTS, validation.Required),
	)
}

func (dto TermRefME) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.TermID, identity.Required...),
		validation.Field(&dto.K, kindRequired...),
	)
}

var kindRequired = []validation.Rule{
	validation.Required,
	validation.In(OneKind, LinkKind, TensorKind, LolliKind, PlusKind, WithKind),
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
