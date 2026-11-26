package core

import (
 "fmt"
 "orglang/go-engine/internal/adt/contkey"
 "orglang/go-engine/internal/adt/seqnum"
 "orglang/go-engine/internal/adt/uniqsym"
 typedef "orglang/go-engine/internal/proc/typedef/core"
)

type Context = typedef.Context
type ExpSpec = typedef.ExpSpec
type OneSpec = typedef.OneSpec
type LinkSpec = typedef.LinkSpec
type TensorSpec = typedef.TensorSpec
type LolliSpec = typedef.LolliSpec
type PlusSpec = typedef.PlusSpec
type WithSpec = typedef.WithSpec
type UpSpec = typedef.UpSpec
type DownSpec = typedef.DownSpec
type ExpRef = typedef.ExpRef
type OneRef = typedef.OneRef
type LinkRef = typedef.LinkRef
type PlusRef = typedef.PlusRef
type WithRef = typedef.WithRef
type TensorRef = typedef.TensorRef
type LolliRef = typedef.LolliRef
type UpRef = typedef.UpRef
type DownRef = typedef.DownRef
type ExpRec = typedef.ExpRec
type ProdRec = typedef.ProdRec
type SumRec = typedef.SumRec
type OneRec = typedef.OneRec
type LinkRec = typedef.LinkRec
type PlusRec = typedef.PlusRec
type WithRec = typedef.WithRec
type TensorRec = typedef.TensorRec
type LolliRec = typedef.LolliRec
type UpRec = typedef.UpRec
type DownRec = typedef.DownRec

func CheckSpec(got, want ExpSpec) error { return typedef.CheckSpec(got, want) }
func CheckRec(got, want ExpRec) error { return typedef.CheckRec(got, want) }
func CheckRef(got, want contkey.ADT) error { return typedef.CheckRef(got, want) }

func ErrSymMissingInEnv(want uniqsym.ADT) error {
	return fmt.Errorf("qn missing in env: %v", want)
}

func errConcurrentModification(got seqnum.ADT, want seqnum.ADT) error {
	return fmt.Errorf("entity concurrent modification: want revision %v, got revision %v", want, got)
}

func CheckRef(got, want contkey.ADT) error {
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

func ErrDoesNotExist(want contkey.ADT) error {
	return fmt.Errorf("root doesn't exist: %v", want)
}

func ErrMissingInEnv(want contkey.ADT) error {
	return fmt.Errorf("root missing in env: %v", want)
}

func ErrMissingInCfg(want contkey.ADT) error {
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

