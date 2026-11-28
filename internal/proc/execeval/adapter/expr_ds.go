package adapter

type ExprKind int

const (
	NonExpr ExprKind = iota
	CloseExpr
	WaitExpr
	SendExpr
	RecvExpr
	LabExpr
	CaseExpr
	LinkExpr
	SpawnExpr
	FwdExpr
)

type ExprSpec struct {
	K     ExprKind    `json:"k"`
	Close *CloseSpec `json:"close,omitempty"`
	Wait  *WaitSpec  `json:"wait,omitempty"`
	Send  *SendSpec  `json:"send,omitempty"`
	Recv  *RecvSpec  `json:"recv,omitempty"`
	Lab   *LabSpec   `json:"lab,omitempty"`
	Case  *CaseSpec  `json:"case,omitempty"`
	Fwd   *FwdSpec   `json:"fwd,omitempty"`
}

type ExprRec struct {
	K     ExprKind   `json:"k"`
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
	ContES ExprSpec `json:"cont"`
}
type WaitRec struct {
	X      string  `json:"x"`
	ContES ExprSpec `json:"cont"`
}

type SendSpec struct {
	X string `json:"x"`
	Y string `json:"y"`
}
type SendRec struct {
	X        string `json:"x"`
	A        string `json:"a"`
	B        string `json:"b"`
	CommID   string `json:"comm_id"`
	ValExprCK int64  `json:"val_expr_ck"`
}

type RecvSpec struct {
	X      string  `json:"x"`
	Y      string  `json:"y"`
	ContES ExprSpec `json:"cont"`
}
type RecvRec struct {
	X      string  `json:"x"`
	A      string  `json:"a"`
	Y      string  `json:"y"`
	ContES ExprSpec `json:"cont"`
}

type LabSpec struct {
	X      string  `json:"x"`
	Label  string  `json:"lab"`
	ContES ExprSpec `json:"cont"`
}
type LabRec struct {
	X     string `json:"x"`
	A     string `json:"a"`
	Label string `json:"lab"`
	// ValExprCK keeps the legacy JSON property for wire compatibility.
	ValExprCK int64 `json:"val_expr_vk"`
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
	ContES ExprSpec `json:"cont"`
}
type BranchRec struct {
	Label  string  `json:"lab"`
	ContES ExprSpec `json:"cont"`
}

type FwdSpec struct {
	X string `json:"x"`
	Y string `json:"y"`
}
type FwdRec struct {
	X string `json:"x"`
	B string `json:"b"`
}
