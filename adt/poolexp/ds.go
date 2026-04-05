package poolexp

type Repo interface {
}

type ExpRecDS struct {
	K       expKind     `json:"k"`
	Acquire *shiftRecDS `json:"acquire,omitempty"`
	Accept  *shiftRecDS `json:"accept,omitempty"`
}

type expKind int

const (
	unkExp expKind = iota
	acquireExp
	acceptExp
)

type shiftRecDS struct {
	ContChnlID string `json:"cont"`
}
