package termexp

type Repo interface {
}

type ExpSpecDS struct {
	K       expKind       `json:"k"`
	Acquire *grantSpecDS  `json:"acquire,omitempty"`
	Accept  *grantSpecDS  `json:"accept,omitempty"`
	Hire    *coopSpecDS   `json:"hire,omitempty"`
	Apply   *coopSpecDS   `json:"apply,omitempty"`
	Release *revokeSpecDS `json:"release,omitempty"`
	Detach  *revokeSpecDS `json:"detach,omitempty"`
}

type grantSpecDS struct {
	CommChnlPH string    `json:"ph"`
	ContExp    ExpSpecDS `json:"exp"`
}

type coopSpecDS struct {
	CommChnlPH string    `json:"ph"`
	ProcTermQN string    `json:"qn"`
	ContExp    ExpSpecDS `json:"exp"`
}

type revokeSpecDS struct {
	CommChnlPH string `json:"ph"`
}

type ExpRecDS struct {
	K       expKind      `json:"k"`
	Acquire *grantRecDS  `json:"acquire,omitempty"`
	Accept  *grantRecDS  `json:"accept,omitempty"`
	Hire    *coopRecDS   `json:"hire,omitempty"`
	Apply   *coopRecDS   `json:"apply,omitempty"`
	Release *revokeRecDS `json:"release,omitempty"`
	Detach  *revokeRecDS `json:"detach,omitempty"`
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

type grantRecDS struct {
	ContChnlPH string    `json:"ph"`
	ContExp    ExpSpecDS `json:"exp"`
}

type coopRecDS struct {
	ContChnlPH string    `json:"ph"`
	ProcTermQN string    `json:"qn"`
	ContExp    ExpSpecDS `json:"exp"`
}

type revokeRecDS struct {
	ContChnlPH string `json:"ph"`
}
