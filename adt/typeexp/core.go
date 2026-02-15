package typeexp

import (
	"fmt"

	"orglang/go-engine/adt/polarity"
	"orglang/go-engine/adt/revnum"
	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/adt/valkey"
)

type ExpSpec interface {
	spec()
}

type OneSpec struct{}

func (OneSpec) spec() {}

// aka TpName
type LinkSpec struct {
	TypeQN uniqsym.ADT
}

func (LinkSpec) spec() {}

type TensorSpec struct {
	Val  ExpSpec // val to send
	Cont ExpSpec // cont
}

func (TensorSpec) spec() {}

type LolliSpec struct {
	Val  ExpSpec // val to receive
	Cont ExpSpec // cont
}

func (LolliSpec) spec() {}

// aka Internal Choice
type PlusSpec struct {
	Choices map[uniqsym.ADT]ExpSpec // conts
}

func (PlusSpec) spec() {}

// aka External Choice
type WithSpec struct {
	Choices map[uniqsym.ADT]ExpSpec // conts
}

func (WithSpec) spec() {}

type UpSpec struct {
	Cont ExpSpec // cont
}

func (UpSpec) spec() {}

type DownSpec struct {
	Cont ExpSpec // cont
}

func (DownSpec) spec() {}

type ExpRef interface {
	valkey.Keyable
}

type OneRef struct {
	ExpVK valkey.ADT
}

func (r OneRef) Key() valkey.ADT { return r.ExpVK }

type LinkRef struct {
	ExpVK valkey.ADT
}

func (r LinkRef) Key() valkey.ADT { return r.ExpVK }

type PlusRef struct {
	ExpVK valkey.ADT
}

func (r PlusRef) Key() valkey.ADT { return r.ExpVK }

type WithRef struct {
	ExpVK valkey.ADT
}

func (r WithRef) Key() valkey.ADT { return r.ExpVK }

type TensorRef struct {
	ExpVK valkey.ADT
}

func (r TensorRef) Key() valkey.ADT { return r.ExpVK }

type LolliRef struct {
	ExpVK valkey.ADT
}

func (r LolliRef) Key() valkey.ADT { return r.ExpVK }

type UpRef struct {
	ExpVK valkey.ADT
}

func (r UpRef) Ident() valkey.ADT { return r.ExpVK }

type DownRef struct {
	ExpVK valkey.ADT
}

func (r DownRef) Ident() valkey.ADT { return r.ExpVK }

// aka Stype
type ExpRec interface {
	polarity.Polarizable
	valkey.Keyable
}

type ProdRec interface {
	Next() valkey.ADT
}

type SumRec interface {
	Next(uniqsym.ADT) valkey.ADT
}

type OneRec struct {
	ExpVK valkey.ADT
}

func (OneRec) spec() {}

func (r OneRec) Key() valkey.ADT { return r.ExpVK }

func (OneRec) Pol() polarity.ADT { return polarity.Pos }

// aka TpName
type LinkRec struct {
	ExpVK  valkey.ADT
	TypeQN uniqsym.ADT
}

func (LinkRec) spec() {}

func (r LinkRec) Key() valkey.ADT { return r.ExpVK }

func (LinkRec) Pol() polarity.ADT { return polarity.Zero }

// aka Internal Choice
type PlusRec struct {
	ExpVK   valkey.ADT
	Choices map[uniqsym.ADT]ExpRec
}

func (PlusRec) spec() {}

func (r PlusRec) Key() valkey.ADT { return r.ExpVK }

func (r PlusRec) Next(l uniqsym.ADT) valkey.ADT { return r.Choices[l].Key() }

func (PlusRec) Pol() polarity.ADT { return polarity.Pos }

// aka External Choice
type WithRec struct {
	ExpVK   valkey.ADT
	Choices map[uniqsym.ADT]ExpRec
}

func (WithRec) spec() {}

func (r WithRec) Key() valkey.ADT { return r.ExpVK }

func (r WithRec) Next(l uniqsym.ADT) valkey.ADT { return r.Choices[l].Key() }

func (WithRec) Pol() polarity.ADT { return polarity.Neg }

type TensorRec struct {
	ExpVK valkey.ADT
	Val   ExpRec
	Cont  ExpRec
}

func (TensorRec) spec() {}

func (r TensorRec) Key() valkey.ADT { return r.ExpVK }

func (r TensorRec) Next() valkey.ADT { return r.Cont.Key() }

func (TensorRec) Pol() polarity.ADT { return polarity.Pos }

type LolliRec struct {
	ExpVK valkey.ADT
	Val   ExpRec
	Cont  ExpRec
}

func (LolliRec) spec() {}

func (r LolliRec) Key() valkey.ADT { return r.ExpVK }

func (r LolliRec) Next() valkey.ADT { return r.Cont.Key() }

func (LolliRec) Pol() polarity.ADT { return polarity.Neg }

type UpRec struct {
	ExpVK valkey.ADT
	Cont  ExpRec
}

func (UpRec) spec() {}

func (r UpRec) Key() valkey.ADT { return r.ExpVK }

func (UpRec) Pol() polarity.ADT { return polarity.Zero }

type DownRec struct {
	ExpVK valkey.ADT
	Cont  ExpRec
}

func (DownRec) spec() {}

func (r DownRec) Key() valkey.ADT { return r.ExpVK }

func (DownRec) Pol() polarity.ADT { return polarity.Zero }

type Context struct {
	Assets map[uniqsym.ADT]ExpRec
	Liabs  map[uniqsym.ADT]ExpRec
}

func ErrSymMissingInEnv(want uniqsym.ADT) error {
	return fmt.Errorf("qn missing in env: %v", want)
}

func errConcurrentModification(got revnum.ADT, want revnum.ADT) error {
	return fmt.Errorf("entity concurrent modification: want revision %v, got revision %v", want, got)
}

func errOptimisticUpdate(got revnum.ADT) error {
	return fmt.Errorf("entity concurrent modification: got revision %v", got)
}

func CheckRef(got, want valkey.ADT) error {
	if got != want {
		return fmt.Errorf("type mismatch: want %+v, got %+v", want, got)
	}
	return nil
}

// aka eqtp
func CheckSpec(got, want ExpSpec) error {
	switch wantSpec := want.(type) {
	case OneSpec:
		_, ok := got.(OneSpec)
		if !ok {
			return ErrSpecTypeMismatch(got, want)
		}
		return nil
	case TensorSpec:
		gotSpec, ok := got.(TensorSpec)
		if !ok {
			return ErrSpecTypeMismatch(got, want)
		}
		err := CheckSpec(gotSpec.Val, wantSpec.Val)
		if err != nil {
			return err
		}
		return CheckSpec(gotSpec.Cont, wantSpec.Cont)
	case LolliSpec:
		gotSpec, ok := got.(LolliSpec)
		if !ok {
			return ErrSpecTypeMismatch(got, want)
		}
		err := CheckSpec(gotSpec.Val, wantSpec.Val)
		if err != nil {
			return err
		}
		return CheckSpec(gotSpec.Cont, wantSpec.Cont)
	case PlusSpec:
		gotSpec, ok := got.(PlusSpec)
		if !ok {
			return ErrSpecTypeMismatch(got, want)
		}
		if len(gotSpec.Choices) != len(wantSpec.Choices) {
			return fmt.Errorf("choices mismatch: want %v items, got %v items", len(wantSpec.Choices), len(gotSpec.Choices))
		}
		for wantLab, wantCont := range wantSpec.Choices {
			gotCont, ok := gotSpec.Choices[wantLab]
			if !ok {
				return fmt.Errorf("label mismatch: want %v, got nothing", wantLab)
			}
			err := CheckSpec(gotCont, wantCont)
			if err != nil {
				return err
			}
		}
		return nil
	case WithSpec:
		gotSpec, ok := got.(WithSpec)
		if !ok {
			return ErrSpecTypeMismatch(got, want)
		}
		if len(gotSpec.Choices) != len(wantSpec.Choices) {
			return fmt.Errorf("choices mismatch: want %v items, got %v items", len(wantSpec.Choices), len(gotSpec.Choices))
		}
		for wantLab, wantCont := range wantSpec.Choices {
			gotCont, ok := gotSpec.Choices[wantLab]
			if !ok {
				return fmt.Errorf("label mismatch: want %v, got nothing", wantLab)
			}
			err := CheckSpec(gotCont, wantCont)
			if err != nil {
				return err
			}
		}
		return nil
	default:
		panic(ErrSpecTypeUnexpected(want))
	}
}

// aka eqtp
func CheckRec(got, want ExpRec) error {
	switch wantRec := want.(type) {
	case OneRec:
		_, ok := got.(OneRec)
		if !ok {
			return ErrSnapTypeMismatch(got, want)
		}
		return nil
	case TensorRec:
		gotRec, ok := got.(TensorRec)
		if !ok {
			return ErrSnapTypeMismatch(got, want)
		}
		err := CheckRec(gotRec.Val, wantRec.Val)
		if err != nil {
			return err
		}
		return CheckRec(gotRec.Cont, wantRec.Cont)
	case LolliRec:
		gotRec, ok := got.(LolliRec)
		if !ok {
			return ErrSnapTypeMismatch(got, want)
		}
		err := CheckRec(gotRec.Val, wantRec.Val)
		if err != nil {
			return err
		}
		return CheckRec(gotRec.Cont, wantRec.Cont)
	case PlusRec:
		gotRec, ok := got.(PlusRec)
		if !ok {
			return ErrSnapTypeMismatch(got, want)
		}
		if len(gotRec.Choices) != len(wantRec.Choices) {
			return fmt.Errorf("choices mismatch: want %v items, got %v items", len(wantRec.Choices), len(gotRec.Choices))
		}
		for wantLab, wantCont := range wantRec.Choices {
			gotCont, ok := gotRec.Choices[wantLab]
			if !ok {
				return fmt.Errorf("label mismatch: want %v, got nothing", wantLab)
			}
			err := CheckRec(gotCont, wantCont)
			if err != nil {
				return err
			}
		}
		return nil
	case WithRec:
		gotRec, ok := got.(WithRec)
		if !ok {
			return ErrSnapTypeMismatch(got, want)
		}
		if len(gotRec.Choices) != len(wantRec.Choices) {
			return fmt.Errorf("choices mismatch: want %v items, got %v items", len(wantRec.Choices), len(gotRec.Choices))
		}
		for wantLab, wantChoice := range wantRec.Choices {
			gotChoice, ok := gotRec.Choices[wantLab]
			if !ok {
				return fmt.Errorf("label mismatch: want %v, got nothing", wantLab)
			}
			err := CheckRec(gotChoice, wantChoice)
			if err != nil {
				return err
			}
		}
		return nil
	default:
		panic(ErrRecTypeUnexpected(want))
	}
}

func ErrSpecTypeUnexpected(got ExpSpec) error {
	return fmt.Errorf("spec type unexpected: %T", got)
}

func ErrRefTypeUnexpected(got ExpRef) error {
	return fmt.Errorf("ref type unexpected: %T", got)
}

func ErrRecTypeUnexpected(got ExpRec) error {
	return fmt.Errorf("rec type unexpected: %T", got)
}

func ErrDoesNotExist(want valkey.ADT) error {
	return fmt.Errorf("root doesn't exist: %v", want)
}

func ErrMissingInEnv(want valkey.ADT) error {
	return fmt.Errorf("root missing in env: %v", want)
}

func ErrMissingInCfg(want valkey.ADT) error {
	return fmt.Errorf("root missing in cfg: %v", want)
}

func ErrMissingInCtx(want uniqsym.ADT) error {
	return fmt.Errorf("root missing in ctx: %v", want)
}

func ErrSpecTypeMismatch(got, want ExpSpec) error {
	return fmt.Errorf("spec type mismatch: want %T, got %T", want, got)
}

func ErrSnapTypeMismatch(got, want ExpRec) error {
	return fmt.Errorf("root type mismatch: want %T, got %T", want, got)
}

func ErrPolarityUnexpected(got ExpRec) error {
	return fmt.Errorf("root polarity unexpected: %v", got.Pol())
}

func ErrPolarityMismatch(a, b ExpRec) error {
	return fmt.Errorf("root polarity mismatch: %v != %v", a.Pol(), b.Pol())
}
