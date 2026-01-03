package typeexp

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"orglang/orglang/adt/identity"
	"orglang/orglang/adt/qualsym"
)

func (dto ExpSpecME) Validate() error {
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
		validation.Field(&dto.ValES, validation.Required),
		validation.Field(&dto.ContES, validation.Required),
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
		validation.Field(&dto.ContES, validation.Required),
	)
}

func (dto ExpRefME) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.ExpID, identity.Required...),
		validation.Field(&dto.K, kindRequired...),
	)
}

var kindRequired = []validation.Rule{
	validation.Required,
	validation.In(OneKind, LinkKind, TensorKind, LolliKind, PlusKind, WithKind),
}
