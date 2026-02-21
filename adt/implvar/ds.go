package implvar

type VarSpecDS struct {
	ChnlPH string `json:"chnl_ph"`
	TypeQN string `json:"type_qn"`
}

type VarRecDS struct {
	ImplID string `db:"exec_id"`
	ImplRN int64  `db:"exec_rn"`
	ChnlBS uint8  `db:"chnl_bs"`
	ChnlPH string `db:"chnl_ph"`
	ChnlID string `db:"chnl_id"`
	ExpVK  int64  `db:"exp_vk"`
}
