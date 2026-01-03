package procdef

type semKindME string

const (
	Msg = semKindME("msg")
	Svc = semKindME("svc")
)

type DefRecME struct {
	ID string    `json:"id"`
	K  semKindME `json:"kind"`
}

type termKindME string

const (
	Close = termKindME("close")
	Wait  = termKindME("wait")
	Send  = termKindME("send")
	Recv  = termKindME("recv")
	Lab   = termKindME("lab")
	Case  = termKindME("case")
	Call  = termKindME("call")
	Link  = termKindME("link")
	Spawn = termKindME("spawn")
	Fwd   = termKindME("fwd")
)

type TermSpecME struct {
	K     termKindME   `json:"kind"`
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
