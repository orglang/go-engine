package adapter

type ExpKind int

const (
	NonExp ExpKind = iota
	CloseExp
	WaitExp
	SendExp
	RecvExp
	LabExp
	CaseExp
	LinkExp
	SpawnExp
	FwdExp
)

type ExpSpec struct {
	K     ExpKind    `json:"k"`
	Close *CloseSpec `json:"close,omitempty"`
	Wait  *WaitSpec  `json:"wait,omitempty"`
	Send  *SendSpec  `json:"send,omitempty"`
	Recv  *RecvSpec  `json:"recv,omitempty"`
	Lab   *LabSpec   `json:"lab,omitempty"`
	Case  *CaseSpec  `json:"case,omitempty"`
	Fwd   *FwdSpec   `json:"fwd,omitempty"`
}

type ExpRec struct {
	K     ExpKind   `json:"k"`
	Close *CloseRec `json:"close,omitempty"`
	Wait  *WaitRec  `json:"wait,omitempty"`
	Send  *SendRec  `json:"send,omitempty"`
	Recv  *RecvRec  `json:"recv,omitempty"`
	Lab   *LabRec   `json:"lab,omitempty"`
	Case  *CaseRec  `json:"case,omitempty"`
	Fwd   *FwdRec   `json:"fwd,omitempty"`
}

type CloseSpec struct {
	X string `json:"x"`
}
type CloseRec struct {
	X string `json:"x"`
}

type WaitSpec struct {
	X      string  `json:"x"`
	ContES ExpSpec `json:"cont"`
}
type WaitRec struct {
	X      string  `json:"x"`
	ContES ExpSpec `json:"cont"`
}

type SendSpec struct {
	X string `json:"x"`
	Y string `json:"y"`
}
type SendRec struct {
	X        string `json:"x"`
	A        string `json:"a"`
	B        string `json:"b"`
	ConnID   string `json:"conn_id"`
	ValExpCK int64  `json:"val_exp_ck"`
}

type RecvSpec struct {
	X      string  `json:"x"`
	Y      string  `json:"y"`
	ContES ExpSpec `json:"cont"`
}
type RecvRec struct {
	X      string  `json:"x"`
	A      string  `json:"a"`
	Y      string  `json:"y"`
	ContES ExpSpec `json:"cont"`
}

type LabSpec struct {
	X      string  `json:"x"`
	Label  string  `json:"lab"`
	ContES ExpSpec `json:"cont"`
}
type LabRec struct {
	X     string `json:"x"`
	A     string `json:"a"`
	Label string `json:"lab"`
	// ValExpCK keeps the legacy JSON property for wire compatibility.
	ValExpCK int64 `json:"val_exp_vk"`
}

type CaseSpec struct {
	X        string       `json:"x"`
	Branches []BranchSpec `json:"brs"`
}
type CaseRec struct {
	X        string      `json:"x"`
	A        string      `json:"a"`
	Branches []BranchRec `json:"brs"`
}

type BranchSpec struct {
	Label  string  `json:"lab"`
	ContES ExpSpec `json:"cont"`
}
type BranchRec struct {
	Label  string  `json:"lab"`
	ContES ExpSpec `json:"cont"`
}

type FwdSpec struct {
	X string `json:"x"`
	Y string `json:"y"`
}
type FwdRec struct {
	X string `json:"x"`
	B string `json:"b"`
}

type DescData struct {
	ExpCK  int64       `db:"exp_ck"`
	States []DescState `db:"states"`
}

type DescState struct {
	ExpCK    int64     `db:"exp_ck"`
	SupExpCK int64     `db:"sup_exp_ck"`
	K        int16     `db:"kind"`
	Spec     StateSpec `db:"spec"`
}

type StateSpec struct {
	ContChnlPH string           `json:"cont_chnl_ph,omitempty"`
	ConnChnlPH string           `json:"conn_chnl_ph,omitempty"`
	ValChnlPH  string           `json:"val_chnl_ph,omitempty"`
	ValLabQN   string           `json:"val_lab_qn,omitempty"`
	NewChnlPH  string           `json:"new_chnl_ph,omitempty"`
	ContExpCK  int64            `json:"cont_exp_ck,omitempty"`
	ContExpCKs map[string]int64 `json:"cont_exp_cks,omitempty"`
	ProcTermQN string           `json:"proc_term_qn,omitempty"`
	X          string           `json:"x,omitempty"`
	Ys         []string         `json:"ys,omitempty"`
	ValChnlPHs []string         `json:"val_chnl_phs,omitempty"`
	NewChnlPHs []string         `json:"new_chnl_phs,omitempty"`
}
