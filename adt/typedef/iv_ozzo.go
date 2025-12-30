package typedef

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"orglang/orglang/adt/identity"
	"orglang/orglang/adt/qualsym"
	"orglang/orglang/adt/revnum"
)

func (dto TypeSpecME) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.TypeQN, qualsym.Required...),
		validation.Field(&dto.TypeTS, validation.Required),
	)
}

func (dto IdentME) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.ID, identity.Required...),
	)
}

func (dto TypeRefME) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.TypeID, identity.Required...),
		validation.Field(&dto.TypeRN, revnum.Optional...),
	)
}

func (dto TypeSnapME) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.TypeID, identity.Required...),
		validation.Field(&dto.TypeRN, revnum.Optional...),
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
		validation.Field(&dto.QN, qualsym.Required...),
	)
}

func (dto ProdSpecME) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.Value, validation.Required),
		validation.Field(&dto.Cont, validation.Required),
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
		validation.Field(&dto.Cont, validation.Required),
	)
}

func (dto TermRefME) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.ID, identity.Required...),
		validation.Field(&dto.K, kindRequired...),
	)
}

var kindRequired = []validation.Rule{
	validation.Required,
	validation.In(OneKind, LinkKind, TensorKind, LolliKind, PlusKind, WithKind),
}

func (dto TypeSpecVP) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.NS, qualsym.Required...),
		validation.Field(&dto.Name, qualsym.Required...),
	)
}

func (dto TypeRefVP) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.RoleID, identity.Required...),
		validation.Field(&dto.RoleRN, revnum.Optional...),
		validation.Field(&dto.Title, qualsym.Required...),
	)
}
