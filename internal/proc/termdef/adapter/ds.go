package adapter

type ExpRec struct {
	K     expKind   `json:"k"`
	Close *closeRec `json:"close,omitempty"`
	Wait  *waitRec  `json:"wait,omitempty"`
	Send  *sendRec  `json:"send,omitempty"`
	Recv  *recvRec  `json:"recv,omitempty"`
	Lab   *labRec   `json:"lab,omitempty"`
	Case  *caseRec  `json:"case,omitempty"`
	Fwd   *fwdRec   `json:"fwd,omitempty"`
}

type ExpSpec struct {
	K     expKind    `json:"k"`
	Close *closeSpec `json:"close,omitempty"`
	Wait  *waitSpec  `json:"wait,omitempty"`
	Send  *sendSpec  `json:"send,omitempty"`
	Recv  *recvSpec  `json:"recv,omitempty"`
	Lab   *labSpec   `json:"lab,omitempty"`
	Case  *caseSpec  `json:"case,omitempty"`
	Fwd   *fwdSpec   `json:"fwd,omitempty"`
}

type expKind int

const (
	unkKind expKind = iota
	closeExp
	waitExp
	sendExp
	recvExp
	labExp
	caseExp
	linkExp
	spawnExp
	fwdExp
)

type closeSpec struct {
	X string `json:"x"`
}
type closeRec struct {
	X string `json:"x"`
}

type waitSpec struct {
	X      string  `json:"x"`
	ContES ExpSpec `json:"cont"`
}
type waitRec struct {
	X      string  `json:"x"`
	ContES ExpSpec `json:"cont"`
}

type sendSpec struct {
	X string `json:"x"`
	Y string `json:"y"`
}
type sendRec struct {
	X string `json:"x"`
	A string `json:"a"`
	B string `json:"b"`
}

type recvSpec struct {
	X      string  `json:"x"`
	Y      string  `json:"y"`
	ContES ExpSpec `json:"cont"`
}
type recvRec struct {
	X      string  `json:"x"`
	A      string  `json:"a"`
	Y      string  `json:"y"`
	ContES ExpSpec `json:"cont"`
}

type labSpec struct {
	X     string `json:"x"`
	Label string `json:"lab"`
}
type labRec struct {
	X     string `json:"x"`
	Label string `json:"lab"`
}

type caseSpec struct {
	X        string       `json:"x"`
	Branches []branchSpec `json:"brs"`
}
type caseRec struct {
	X        string      `json:"x"`
	Branches []branchRec `json:"brs"`
}

type branchSpec struct {
	Label  string  `json:"lab"`
	ContES ExpSpec `json:"cont"`
}
type branchRec struct {
	Label  string  `json:"lab"`
	ContES ExpSpec `json:"cont"`
}

type fwdSpec struct {
	X string `json:"x"`
	Y string `json:"y"`
}
type fwdRec struct {
	X string `json:"x"`
	B string `json:"b"`
}
