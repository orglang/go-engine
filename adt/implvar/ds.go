package implvar

type VarRecDS struct {
	ImplID string `db:"impl_id"`
	ImplRN int64  `db:"impl_rn"`
	ChnlID string `db:"chnl_id"`
	ChnlPH string `db:"chnl_ph"`
	ChnlBS uint8  `db:"chnl_bs"`
	ExpVK  int64  `db:"exp_vk"`
}
