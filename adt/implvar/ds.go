package implvar

type VarRecDS struct {
	ImplID string `db:"impl_id"`
	ImplRN int64  `db:"impl_rn"`
	CommID string `db:"comm_id"`
	ChnlID string `db:"chnl_id"`
	ChnlPH string `db:"chnl_ph"`
	ExpVK  int64  `db:"exp_vk"`
	ChnlBS int8   `db:"side"`
}
