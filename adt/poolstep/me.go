package poolstep

type Exch interface {
	SendSpec(StepSpec) error
	SendRec(StepRec) error
}
