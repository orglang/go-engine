package procbind

type BindSpecDS struct {
	ChnlPH string `json:"chnl_ph"`
	TypeQN string `json:"type_qn"`
}

type BindRecDS struct {
	ExecID  string
	ExecRN  int64
	ChnlPH  string
	ChnlID  string
	StateID string
}
