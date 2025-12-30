package procdecl

type BndSpecME struct {
	ChnlPH string `json:"chnl_ph"`
	TypeQN string `json:"type_qn"`
}

type SigSpecME struct {
	X     BndSpecME   `json:"x"`
	SigQN string      `json:"sig_qn"`
	Ys    []BndSpecME `json:"ys"`
}

type IdentME struct {
	SigID string `json:"id" param:"id"`
}

type SigRefME struct {
	SigID string `json:"id" param:"id"`
	Title string `json:"title"`
	SigRN int64  `json:"rev"`
}

type SigSnapME struct {
	X     BndSpecME   `json:"x"`
	SigID string      `json:"sig_id"`
	Ys    []BndSpecME `json:"ys"`
	Title string      `json:"title"`
	SigRN int64       `json:"sig_rn"`
}
