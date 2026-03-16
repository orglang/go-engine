package implvar

type VarRecDS struct {
	ImplID string `db:"impl_id"`
	CommID string `db:"comm_id"`
	ChnlID string `db:"chnl_id"`
	ChnlPH string `db:"chnl_ph"`
	ChnlBS int8   `db:"chnl_bs"`
	ExpVK  int64  `db:"exp_vk"`
}
