package procexp

type expKindME string

const (
	Close = expKindME("close")
	Wait  = expKindME("wait")
	Send  = expKindME("send")
	Recv  = expKindME("recv")
	Lab   = expKindME("lab")
	Case  = expKindME("case")
	Call  = expKindME("call")
	Link  = expKindME("link")
	Spawn = expKindME("spawn")
	Fwd   = expKindME("fwd")
)

type ExpSpecME struct {
	K     expKindME    `json:"kind"`
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
	X      string    `json:"x"`
	ContES ExpSpecME `json:"cont"`
}

type SendSpecME struct {
	X string `json:"x"`
	Y string `json:"y"`
}

type RecvSpecME struct {
	X      string    `json:"x"`
	Y      string    `json:"y"`
	ContES ExpSpecME `json:"cont"`
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
	Label  string    `json:"label"`
	ContES ExpSpecME `json:"cont"`
}

type CallSpecME struct {
	X      string   `json:"x"`
	ProcQN string   `json:"proc_qn"` // раньше был SigPH
	Ys     []string `json:"ys"`
}

type SpawnSpecME struct {
	X      string     `json:"x"`
	DecID  string     `json:"dec_id"`
	Ys     []string   `json:"ys"`
	ContES *ExpSpecME `json:"cont"`
}

type FwdSpecME struct {
	X string `json:"x"`
	Y string `json:"y"`
}
