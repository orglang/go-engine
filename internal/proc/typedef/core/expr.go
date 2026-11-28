package core

import (
	"fmt"
	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/polarity"
	"orglang/go-engine/internal/adt/uniqsym"
)

type ExprSpec interface {
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
	Val  ExprSpec // val to send
	Cont ExprSpec // cont
}

func (TensorSpec) spec() {}

type LolliSpec struct {
	Val  ExprSpec // val to receive
	Cont ExprSpec // cont
}

func (LolliSpec) spec() {}

// aka Internal Choice
type PlusSpec struct {
	Choices map[uniqsym.ADT]ExprSpec // conts
}

func (PlusSpec) spec() {}

// aka External Choice
type WithSpec struct {
	Choices map[uniqsym.ADT]ExprSpec // conts
}

func (WithSpec) spec() {}

type UpSpec struct {
	Cont ExprSpec // cont
}

func (UpSpec) spec() {}

type DownSpec struct {
	Cont ExprSpec // cont
}

func (DownSpec) spec() {}

type ExprRef interface {
	contkey.ContentKeyable
}

type OneRef struct {
	ExprCK contkey.ADT
}

func (r OneRef) Key() contkey.ADT { return r.ExprCK }

type LinkRef struct {
	ExprCK contkey.ADT
}

func (r LinkRef) Key() contkey.ADT { return r.ExprCK }

type PlusRef struct {
	ExprCK contkey.ADT
}

func (r PlusRef) Key() contkey.ADT { return r.ExprCK }

type WithRef struct {
	ExprCK contkey.ADT
}

func (r WithRef) Key() contkey.ADT { return r.ExprCK }

type TensorRef struct {
	ExprCK contkey.ADT
}

func (r TensorRef) Key() contkey.ADT { return r.ExprCK }

type LolliRef struct {
	ExprCK contkey.ADT
}

func (r LolliRef) Key() contkey.ADT { return r.ExprCK }

type UpRef struct {
	ExprCK contkey.ADT
}

func (r UpRef) Ident() contkey.ADT { return r.ExprCK }

type DownRef struct {
	ExprCK contkey.ADT
}

func (r DownRef) Ident() contkey.ADT { return r.ExprCK }

// aka Stype
type ExprRec interface {
	polarity.Polarizable
	contkey.ContentKeyable
}

type ProdRec interface {
	Next() contkey.ADT
}

type SumRec interface {
	Next(uniqsym.ADT) contkey.ADT
}

type OneRec struct {
	ExprCK contkey.ADT
}

func (OneRec) spec() {}

func (r OneRec) Key() contkey.ADT { return r.ExprCK }

func (OneRec) Pol() polarity.ADT { return polarity.Pos }

// aka TpName
type LinkRec struct {
	ExprCK contkey.ADT
	TypeQN uniqsym.ADT
}

func (LinkRec) spec() {}

func (r LinkRec) Key() contkey.ADT { return r.ExprCK }

func (LinkRec) Pol() polarity.ADT { return polarity.Zero }

// aka Internal Choice
type PlusRec struct {
	ExprCK  contkey.ADT
	Choices map[uniqsym.ADT]ExprRec
}

func (PlusRec) spec() {}

func (r PlusRec) Key() contkey.ADT { return r.ExprCK }

func (r PlusRec) Next(l uniqsym.ADT) contkey.ADT { return r.Choices[l].Key() }

func (PlusRec) Pol() polarity.ADT { return polarity.Pos }

// aka External Choice
type WithRec struct {
	ExprCK  contkey.ADT
	Choices map[uniqsym.ADT]ExprRec
}

func (WithRec) spec() {}

func (r WithRec) Key() contkey.ADT { return r.ExprCK }

func (r WithRec) Next(l uniqsym.ADT) contkey.ADT { return r.Choices[l].Key() }

func (WithRec) Pol() polarity.ADT { return polarity.Neg }

type TensorRec struct {
	ExprCK contkey.ADT
	Val    ExprRec
	Cont   ExprRec
}

func (TensorRec) spec() {}

func (r TensorRec) Key() contkey.ADT { return r.ExprCK }

func (r TensorRec) Next() contkey.ADT { return r.Cont.Key() }

func (TensorRec) Pol() polarity.ADT { return polarity.Pos }

type LolliRec struct {
	ExprCK contkey.ADT
	Val    ExprRec
	Cont   ExprRec
}

func (LolliRec) spec() {}

func (r LolliRec) Key() contkey.ADT { return r.ExprCK }

func (r LolliRec) Next() contkey.ADT { return r.Cont.Key() }

func (LolliRec) Pol() polarity.ADT { return polarity.Neg }

type UpRec struct {
	ExprCK contkey.ADT
	Cont   ExprRec
}

func (UpRec) spec() {}

func (r UpRec) Key() contkey.ADT { return r.ExprCK }

func (UpRec) Pol() polarity.ADT { return polarity.Zero }

type DownRec struct {
	ExprCK contkey.ADT
	Cont   ExprRec
}

func (DownRec) spec() {}

func (r DownRec) Key() contkey.ADT { return r.ExprCK }

func (DownRec) Pol() polarity.ADT { return polarity.Zero }

// aka eqtp
func CheckSpec(got, want ExprSpec) error {
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
func CheckRec(got, want ExprRec) error {
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

func ErrSpecTypeUnexpected(got ExprSpec) error {
	return fmt.Errorf("spec type unexpected: %T", got)
}

func ErrRefTypeUnexpected(got ExprRef) error {
	return fmt.Errorf("ref type unexpected: %T", got)
}

func ErrRecTypeUnexpected(got ExprRec) error {
	return fmt.Errorf("rec type unexpected: %T", got)
}

func ErrSpecTypeMismatch(got, want ExprSpec) error {
	return fmt.Errorf("spec type mismatch: want %T, got %T", want, got)
}

func ErrSnapTypeMismatch(got, want ExprRec) error {
	return fmt.Errorf("root type mismatch: want %T, got %T", want, got)
}
