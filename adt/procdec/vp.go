package procdec

type DecSpecVP struct {
	ProcNS string `form:"ns" json:"ns"`
	ProcSN string `form:"sn" json:"sn"`
	Title  string `form:"name" json:"title"`
}

type DecRefVP struct {
	DecID string `form:"dec_id" json:"dec_id" param:"id"`
	DecRN int64  `form:"dec_rn" json:"dec_rn"`
}

type DecSnapVP struct {
	DecID string `json:"dec_id"`
	DecRN int64  `json:"dec_rn"`
	Title string `json:"title"`
}
