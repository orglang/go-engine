package adapter

type ExpKind int

const (
	UnkKind ExpKind = iota
	OneKind
	LinkKind
	TensorKind
	LolliKind
	PlusKind
	WithKind
	UpKind
	DownKind
)

type ExpRef struct {
	ExpCK int64   `db:"exp_ck" json:"exp_ck"`
	K     ExpKind `db:"kind" json:"kind"`
}

type ExpRec struct {
	ExpCK  int64
	States []State
}

type State struct {
	ExpCK    int64   `db:"exp_ck"`
	SupExpCK int64   `db:"sup_exp_ck"`
	K        ExpKind `db:"kind"`
	Spec     ExpSpec `db:"spec" fieldopt:"noexpand"`
}

type ExpSpec struct {
	Link   string `json:"link,omitempty"`
	Tensor *Prod  `json:"tensor,omitempty"`
	Lolli  *Prod  `json:"lolli,omitempty"`
	Plus   *Sum   `json:"plus,omitempty"`
	With   *Sum   `json:"with,omitempty"`
	Up     *Shift `json:"up,omitempty"`
	Down   *Shift `json:"down,omitempty"`
}

type Prod struct {
	ValExpCK  int64 `json:"on"`
	ContExpCK int64 `json:"to"`
}

type Sum struct {
	ProcQNs   []string `json:"on"`
	ContExpCK int64    `json:"to"`
}

type Shift struct {
	ContExpCK int64 `json:"to"`
}
