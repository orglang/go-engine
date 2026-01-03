package typeexp

type ExpSpecME struct {
	K      termKindME  `json:"kind"`
	Link   *LinkSpecME `json:"link,omitempty"`
	Tensor *ProdSpecME `json:"tensor,omitempty"`
	Lolli  *ProdSpecME `json:"lolli,omitempty"`
	Plus   *SumSpecME  `json:"plus,omitempty"`
	With   *SumSpecME  `json:"with,omitempty"`
}

type LinkSpecME struct {
	TypeQN string `json:"type_qn"`
}

type ProdSpecME struct {
	ValES  ExpSpecME `json:"value"`
	ContES ExpSpecME `json:"cont"`
}

type SumSpecME struct {
	Choices []ChoiceSpecME `json:"choices"`
}

type ChoiceSpecME struct {
	Label  string    `json:"label"`
	ContES ExpSpecME `json:"cont"`
}

type ExpRefME struct {
	ExpID string     `json:"exp_id" param:"id"`
	K     termKindME `json:"kind"`
}

type termKindME string

const (
	OneKind    = termKindME("one")
	LinkKind   = termKindME("link")
	TensorKind = termKindME("tensor")
	LolliKind  = termKindME("lolli")
	PlusKind   = termKindME("plus")
	WithKind   = termKindME("with")
)
