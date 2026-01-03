package typedef

type DefSpecVP struct {
	TypeNS string `form:"ns" json:"ns"`
	TypeSN string `form:"name" json:"name"`
}

type DefRefVP struct {
	DefID string `form:"id" json:"id" param:"id"`
	DefRN int64  `form:"def_rn" json:"def_rn"`
	Title string `form:"name" json:"title"`
}

type DefSnapVP struct {
	DefID  string     `json:"def_id"`
	DefRN  int64      `json:"def_rn"`
	Title  string     `json:"title"`
	TypeTS TermSpecME `json:"type_ts"`
}
