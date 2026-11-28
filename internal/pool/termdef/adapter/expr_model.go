package adapter

const (
	AcquireKind int16 = iota + 1
	AcceptKind
	HireKind
	ApplyKind
	ReleaseKind
	DetachKind
	FireKind
	QuitKind
	SpawnKind
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
	CommChnlPH  string   `json:"comm_chnl_ph,omitempty"`
	ProcTermQN  string   `json:"proc_term_qn,omitempty"`
	ProcExecQNs []string `json:"proc_exec_qns,omitempty"`
	ContExprCK  int64    `json:"cont_expr_ck,omitempty"`
}
