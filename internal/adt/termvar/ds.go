package termvar

type VarRecDS struct {
	TypeID string `json:"id"`
	ChnlPH string `json:"ph"`
	// ExpCK keeps the legacy JSON name for wire compatibility.
	ExpCK  int64  `json:"vk"`
}
