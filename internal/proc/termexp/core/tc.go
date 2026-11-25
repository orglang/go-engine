package core

const (
	CloseKind int16 = iota + 1
	WaitKind
	SendKind
	RecvKind
	LabKind
	CaseKind
	LinkKind
	FwdKind
	CallKind
	SpawnKind
	AcquireKind
	AcceptKind
	DetachKind
	ReleaseKind
)
