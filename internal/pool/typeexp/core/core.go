package core

import (
 "fmt"
 "orglang/go-engine/internal/adt/contkey"
 typedef "orglang/go-engine/internal/pool/typedef/core"
)

type ExpSpec = typedef.ExpSpec
type OneSpec = typedef.OneSpec
type LinkSpec = typedef.LinkSpec
type PlusSpec = typedef.PlusSpec
type WithSpec = typedef.WithSpec
type UpSpec = typedef.UpSpec
type DownSpec = typedef.DownSpec
type ExpRef = typedef.ExpRef
type OneRef = typedef.OneRef
type LinkRef = typedef.LinkRef
type PlusRef = typedef.PlusRef
type WithRef = typedef.WithRef
type ExpRec = typedef.ExpRec
type ProdRec = typedef.ProdRec
type SumRec = typedef.SumRec
type OneRec = typedef.OneRec
type LinkRec = typedef.LinkRec
type PlusRec = typedef.PlusRec
type WithRec = typedef.WithRec
type UpRec = typedef.UpRec
type DownRec = typedef.DownRec

func ErrSpecTypeUnexpected(got ExpSpec) error { return fmt.Errorf("spec type unexpected: %T", got) }
func ErrRefTypeUnexpected(got ExpRef) error { return fmt.Errorf("ref type unexpected: %T", got) }
func ErrRecTypeUnexpected(got ExpRec) error { return fmt.Errorf("rec type unexpected: %T", got) }
func ErrDoesNotExist(want contkey.ADT) error { return fmt.Errorf("root doesn't exist: %v", want) }
