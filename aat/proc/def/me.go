package def

type SemKind string

const (
	Msg = SemKind("msg")
	Svc = SemKind("svc")
)

type RootME struct {
	ID string  `json:"id"`
	K  SemKind `json:"kind"`
}

type TermKind string

const (
	Close = TermKind("close")
	Wait  = TermKind("wait")
	Send  = TermKind("send")
	Recv  = TermKind("recv")
	Lab   = TermKind("lab")
	Case  = TermKind("case")
	Call  = TermKind("call")
	Link  = TermKind("link")
	Spawn = TermKind("spawn")
	Fwd   = TermKind("fwd")
)

type TermSpecME struct {
	K     TermKind     `json:"kind"`
	Close *CloseSpecME `json:"close,omitempty"`
	Wait  *WaitSpecME  `json:"wait,omitempty"`
	Send  *SendSpecME  `json:"send,omitempty"`
	Recv  *RecvSpecME  `json:"recv,omitempty"`
	Lab   *LabSpecME   `json:"lab,omitempty"`
	Case  *CaseSpecME  `json:"case,omitempty"`
	Spawn *SpawnSpecME `json:"spawn,omitempty"`
	Fwd   *FwdSpecME   `json:"fwd,omitempty"`
	Call  *CallSpecME  `json:"call,omitempty"`
}

type CloseSpecME struct {
	X string `json:"x"`
}

type WaitSpecME struct {
	X    string     `json:"x"`
	Cont TermSpecME `json:"cont"`
}

type SendSpecME struct {
	X string `json:"x"`
	Y string `json:"y"`
}

type RecvSpecME struct {
	X    string     `json:"x"`
	Y    string     `json:"y"`
	Cont TermSpecME `json:"cont"`
}

type LabSpecME struct {
	X     string `json:"x"`
	Label string `json:"label"`
}

type CaseSpecME struct {
	X   string         `json:"x"`
	Brs []BranchSpecME `json:"branches"`
}

type BranchSpecME struct {
	Label string     `json:"label"`
	Cont  TermSpecME `json:"cont"`
}

type CallSpecME struct {
	X     string   `json:"x"`
	SigPH string   `json:"sig_ph"`
	Ys    []string `json:"ys"`
}

type SpawnSpecME struct {
	X     string      `json:"x"`
	SigID string      `json:"sig_id"`
	Ys    []string    `json:"ys"`
	Cont  *TermSpecME `json:"cont"`
}

type FwdSpecME struct {
	X string `json:"x"`
	Y string `json:"y"`
}
