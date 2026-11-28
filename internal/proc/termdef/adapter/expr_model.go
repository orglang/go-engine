package adapter

const (
	CloseKind int16 = iota + 1
	WaitKind
	SendKind
	RecvKind
	LabKind
	CaseKind
	LinkKind
	FwdKind
	CallKind
	SpawnKind
	AcquireKind
	AcceptKind
	DetachKind
	ReleaseKind
)

type ExprData struct {
	ExprCK int64       `db:"expr_ck"`
	States []ExprState `db:"states"`
}
type ExprState struct {
	ExprCK    int64         `db:"expr_ck"`
	SupExprCK int64         `db:"sup_expr_ck"`
	K         int16         `db:"kind"`
	Spec      ExprStateSpec `db:"spec"`
}
type ExprStateSpec struct {
	ContChnlPH  string           `json:"cont_chnl_ph,omitempty"`
	CommChnlPH  string           `json:"comm_chnl_ph,omitempty"`
	ValChnlPH   string           `json:"val_chnl_ph,omitempty"`
	ValLabQN    string           `json:"val_lab_qn,omitempty"`
	NewChnlPH   string           `json:"new_chnl_ph,omitempty"`
	ContExprCK  int64            `json:"cont_expr_ck,omitempty"`
	ContExprCKs map[string]int64 `json:"cont_expr_cks,omitempty"`
	ProcTermQN  string           `json:"proc_term_qn,omitempty"`
	X           string           `json:"x,omitempty"`
	Ys          []string         `json:"ys,omitempty"`
	ValChnlPHs  []string         `json:"val_chnl_phs,omitempty"`
	NewChnlPHs  []string         `json:"new_chnl_phs,omitempty"`
}
