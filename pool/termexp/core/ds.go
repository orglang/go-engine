package core

type Repo interface {
}

type ExpSpecDS struct {
	K       ExpKind       `json:"k"`
	Acquire *GrantSpecDS  `json:"acquire,omitempty"`
	Accept  *GrantSpecDS  `json:"accept,omitempty"`
	Hire    *CoopSpecDS   `json:"hire,omitempty"`
	Apply   *CoopSpecDS   `json:"apply,omitempty"`
	Release *RevokeSpecDS `json:"release,omitempty"`
	Detach  *RevokeSpecDS `json:"detach,omitempty"`
}

type GrantSpecDS struct {
	CommChnlPH string    `json:"ph"`
	ContExp    ExpSpecDS `json:"exp"`
}

type CoopSpecDS struct {
	CommChnlPH string    `json:"ph"`
	ProcTermQN string    `json:"qn"`
	ContExp    ExpSpecDS `json:"exp"`
}

type RevokeSpecDS struct {
	CommChnlPH string `json:"ph"`
}

type ExpRecDS struct {
	K       ExpKind      `json:"k"`
	Acquire *GrantRecDS  `json:"acquire,omitempty"`
	Accept  *GrantRecDS  `json:"accept,omitempty"`
	Hire    *CoopRecDS   `json:"hire,omitempty"`
	Apply   *CoopRecDS   `json:"apply,omitempty"`
	Release *RevokeRecDS `json:"release,omitempty"`
	Detach  *RevokeRecDS `json:"detach,omitempty"`
}

type ExpKind int

const (
	UnkKind ExpKind = iota
	AcquireKind
	AcceptKind
	HireKind
	ApplyKind
	ReleaseKind
	DetachKind
)

type GrantRecDS struct {
	ContChnlPH string    `json:"ph"`
	ContExp    ExpSpecDS `json:"exp"`
}

type CoopRecDS struct {
	ContChnlPH string    `json:"ph"`
	ProcTermQN string    `json:"qn"`
	ContExp    ExpSpecDS `json:"exp"`
}

type RevokeRecDS struct {
	ContChnlPH string `json:"ph"`
}
