package polarity

type ADT int8

const (
	Zero = ADT(0)
	// provider -> client
	Pos = ADT(+1)
	// client -> provider
	Neg = ADT(-1)
)

type Polarizable interface {
	Pol() ADT
}
