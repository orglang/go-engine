package poolexp

type Repo interface {
}

type ExpSpecDS struct {
	K       expKind      `json:"k"`
	Acquire *shiftSpecDS `json:"acquire,omitempty"`
	Accept  *shiftSpecDS `json:"accept,omitempty"`
	Hire    *laborSpecDS `json:"hire,omitempty"`
	Apply   *laborSpecDS `json:"apply,omitempty"`
}

type shiftSpecDS struct {
	CommChnlPH string    `json:"ph"`
	ContExp    ExpSpecDS `json:"exp"`
}

type laborSpecDS struct {
	CommChnlPH string
	ProcDescQN string
}

type ExpRecDS struct {
	K       expKind     `json:"k"`
	Acquire *shiftRecDS `json:"acquire,omitempty"`
	Accept  *shiftRecDS `json:"accept,omitempty"`
}

type expKind int

const (
	unkKind expKind = iota
	acquireKind
	acceptKind
	hireKind
	applyKind
)

type shiftRecDS struct {
	ContChnlID string    `json:"id"`
	ContExp    ExpSpecDS `json:"exp"`
}
