package procdef

type semKindME string

const (
	Msg = "msg"
	Svc = "svc"
)

type DefRecME struct {
	ID string    `json:"id"`
	K  semKindME `json:"kind"`
}
