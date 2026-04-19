package termexp

type Repo interface {
}

type ExpSpecDS struct {
	K       expKind      `json:"k"`
	Acquire *upSpecDS    `json:"acquire,omitempty"`
	Accept  *upSpecDS    `json:"accept,omitempty"`
	Hire    *laborSpecDS `json:"hire,omitempty"`
	Apply   *laborSpecDS `json:"apply,omitempty"`
	Release *downSpecDS  `json:"release,omitempty"`
	Detach  *downSpecDS  `json:"detach,omitempty"`
}

type upSpecDS struct {
	CommChnlPH string    `json:"ph"`
	ContExp    ExpSpecDS `json:"exp"`
}

type laborSpecDS struct {
	CommChnlPH string    `json:"ph"`
	ProcTermQN string    `json:"qn"`
	ContExp    ExpSpecDS `json:"exp"`
}

type downSpecDS struct {
	CommChnlPH string `json:"ph"`
}

type ExpRecDS struct {
	K       expKind     `json:"k"`
	Acquire *upRecDS    `json:"acquire,omitempty"`
	Accept  *upRecDS    `json:"accept,omitempty"`
	Hire    *laborRecDS `json:"hire,omitempty"`
	Apply   *laborRecDS `json:"apply,omitempty"`
	Release *downRecDS  `json:"release,omitempty"`
	Detach  *downRecDS  `json:"detach,omitempty"`
}

type expKind int

const (
	unkKind expKind = iota
	acquireKind
	acceptKind
	hireKind
	applyKind
	releaseKind
	detachKind
)

type upRecDS struct {
	ContChnlID string    `json:"id"`
	ContExp    ExpSpecDS `json:"exp"`
}

type laborRecDS struct {
	ContChnlID string    `json:"id"`
	ProcTermQN string    `json:"qn"`
	ContExp    ExpSpecDS `json:"exp"`
}

type downRecDS struct {
	ContChnlID string `json:"id"`
}
