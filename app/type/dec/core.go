package dec

import (
	"context"
	"fmt"
	"log/slog"

	"smecalculus/rolevod/lib/data"
	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/rn"
	"smecalculus/rolevod/lib/sym"

	"smecalculus/rolevod/app/type/def"

	"smecalculus/rolevod/internal/alias"
)

type TypeSpec struct {
	TypeNS sym.ADT
	TypeSN sym.ADT
	TypeTS def.TermSpec
}

type TypeRef struct {
	TypeID id.ADT
	Title  string
	TypeRN rn.ADT
}

// aka TpDef
type TypeRec struct {
	TypeID id.ADT
	Title  string
	TermID id.ADT
	TypeRN rn.ADT
}

type TypeSnap struct {
	TypeID id.ADT
	Title  string
	TypeQN sym.ADT
	TypeTS def.TermSpec
	TypeRN rn.ADT
}

type API interface {
	Incept(sym.ADT) (TypeRef, error)
	Create(TypeSpec) (TypeSnap, error)
	Modify(TypeSnap) (TypeSnap, error)
	Retrieve(id.ADT) (TypeSnap, error)
	retrieveSnap(TypeRec) (TypeSnap, error)
	RetreiveRefs() ([]TypeRef, error)
}

type service struct {
	roles    Repo
	states   def.TermRepo
	aliases  alias.Repo
	operator data.Operator
	log      *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return &service{}
}

func newService(
	roles Repo,
	states def.TermRepo,
	aliases alias.Repo,
	operator data.Operator,
	l *slog.Logger,
) *service {
	return &service{roles, states, aliases, operator, l}
}

func (s *service) Incept(qn sym.ADT) (_ TypeRef, err error) {
	ctx := context.Background()
	qnAttr := slog.Any("roleQN", qn)
	s.log.Debug("inception started", qnAttr)
	newAlias := alias.Root{QN: qn, ID: id.New(), RN: rn.Initial()}
	newType := TypeRec{TypeID: newAlias.ID, TypeRN: newAlias.RN, Title: newAlias.QN.SN()}
	s.operator.Explicit(ctx, func(ds data.Source) error {
		err = s.aliases.Insert(ds, newAlias)
		if err != nil {
			return err
		}
		err = s.roles.InsertType(ds, newType)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		s.log.Error("inception failed", qnAttr)
		return TypeRef{}, err
	}
	s.log.Debug("inception succeeded", qnAttr, slog.Any("roleID", newType.TypeID))
	return ConvertRecToRef(newType), nil
}

func (s *service) Create(spec TypeSpec) (_ TypeSnap, err error) {
	ctx := context.Background()
	qnAttr := slog.Any("typeQN", spec.TypeSN)
	s.log.Debug("creation started", qnAttr, slog.Any("spec", spec))
	newAlias := alias.Root{QN: spec.TypeSN, ID: id.New(), RN: rn.Initial()}
	newTerm := def.ConvertSpecToRec(spec.TypeTS)
	newType := TypeRec{
		TypeID: newAlias.ID,
		TypeRN: newAlias.RN,
		Title:  newAlias.QN.SN(),
		TermID: newTerm.Ident(),
	}
	s.operator.Explicit(ctx, func(ds data.Source) error {
		err = s.aliases.Insert(ds, newAlias)
		if err != nil {
			return err
		}
		err = s.states.InsertTerm(ds, newTerm)
		if err != nil {
			return err
		}
		err = s.roles.InsertType(ds, newType)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		s.log.Error("creation failed", qnAttr)
		return TypeSnap{}, err
	}
	s.log.Debug("creation succeeded", qnAttr, slog.Any("typeID", newType.TypeID))
	return TypeSnap{
		TypeID: newType.TypeID,
		TypeRN: newType.TypeRN,
		Title:  newType.Title,
		TypeQN: newAlias.QN,
		TypeTS: def.ConvertRecToSpec(newTerm),
	}, nil
}

func (s *service) Modify(snap TypeSnap) (_ TypeSnap, err error) {
	ctx := context.Background()
	idAttr := slog.Any("typeID", snap.TypeID)
	s.log.Debug("modification started", idAttr)
	var rec TypeRec
	s.operator.Implicit(ctx, func(ds data.Source) error {
		rec, err = s.roles.SelectByID(ds, snap.TypeID)
		return err
	})
	if err != nil {
		s.log.Error("modification failed", idAttr)
		return TypeSnap{}, err
	}
	if snap.TypeRN != rec.TypeRN {
		s.log.Error("modification failed", idAttr)
		return TypeSnap{}, errConcurrentModification(snap.TypeRN, rec.TypeRN)
	} else {
		snap.TypeRN = rn.Next(snap.TypeRN)
	}
	curSnap, err := s.retrieveSnap(rec)
	if err != nil {
		s.log.Error("modification failed", idAttr)
		return TypeSnap{}, err
	}
	s.operator.Explicit(ctx, func(ds data.Source) error {
		if def.CheckSpec(snap.TypeTS, curSnap.TypeTS) != nil {
			newTerm := def.ConvertSpecToRec(snap.TypeTS)
			err = s.states.InsertTerm(ds, newTerm)
			if err != nil {
				return err
			}
			rec.TermID = newTerm.Ident()
			rec.TypeRN = snap.TypeRN
		}
		if rec.TypeRN == snap.TypeRN {
			err = s.roles.Update(ds, rec)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		s.log.Error("modification failed", idAttr)
		return TypeSnap{}, err
	}
	s.log.Debug("modification succeeded", idAttr)
	return snap, nil
}

func (s *service) Retrieve(recID id.ADT) (_ TypeSnap, err error) {
	ctx := context.Background()
	var root TypeRec
	s.operator.Implicit(ctx, func(ds data.Source) error {
		root, err = s.roles.SelectByID(ds, recID)
		return err
	})
	if err != nil {
		s.log.Error("retrieval failed", slog.Any("roleID", recID))
		return TypeSnap{}, err
	}
	return s.retrieveSnap(root)
}

func (s *service) RetrieveImpl(recID id.ADT) (enty TypeRec, err error) {
	ctx := context.Background()
	s.operator.Implicit(ctx, func(ds data.Source) error {
		enty, err = s.roles.SelectByID(ds, recID)
		return err
	})
	if err != nil {
		s.log.Error("retrieval failed", slog.Any("roleID", recID))
		return TypeRec{}, err
	}
	return enty, nil
}

func (s *service) retrieveSnap(typeRec TypeRec) (_ TypeSnap, err error) {
	ctx := context.Background()
	var termRec def.TermRec
	s.operator.Implicit(ctx, func(ds data.Source) error {
		termRec, err = s.states.SelectTermByID(ds, typeRec.TermID)
		return err
	})
	if err != nil {
		s.log.Error("retrieval failed", slog.Any("roleID", typeRec.TypeID))
		return TypeSnap{}, err
	}
	return TypeSnap{
		TypeID: typeRec.TypeID,
		TypeRN: typeRec.TypeRN,
		Title:  typeRec.Title,
		TypeTS: def.ConvertRecToSpec(termRec),
	}, nil
}

func (s *service) RetreiveRefs() (refs []TypeRef, err error) {
	ctx := context.Background()
	s.operator.Implicit(ctx, func(ds data.Source) error {
		refs, err = s.roles.SelectRefs(ds)
		return err
	})
	if err != nil {
		s.log.Error("retrieval failed")
		return nil, err
	}
	return refs, nil
}

func CollectEnv(roots []TypeRec) []id.ADT {
	termIDs := []id.ADT{}
	for _, r := range roots {
		termIDs = append(termIDs, r.TermID)
	}
	return termIDs
}

type Repo interface {
	InsertType(data.Source, TypeRec) error
	Update(data.Source, TypeRec) error
	SelectRefs(data.Source) ([]TypeRef, error)
	SelectByID(data.Source, id.ADT) (TypeRec, error)
	SelectByIDs(data.Source, []id.ADT) ([]TypeRec, error)
	SelectByQN(data.Source, sym.ADT) (TypeRec, error)
	SelectByQNs(data.Source, []sym.ADT) ([]TypeRec, error)
	SelectEnv(data.Source, []sym.ADT) (map[sym.ADT]TypeRec, error)
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
// goverter:extend smecalculus/rolevod/app/type/def:Convert.*
var (
	ConvertRecToRef  func(TypeRec) TypeRef
	ConvertSnapToRef func(TypeSnap) TypeRef
)

func ErrMissingInEnv(want sym.ADT) error {
	return fmt.Errorf("root missing in env: %v", want)
}

func errConcurrentModification(got rn.ADT, want rn.ADT) error {
	return fmt.Errorf("entity concurrent modification: want revision %v, got revision %v", want, got)
}

func errOptimisticUpdate(got rn.ADT) error {
	return fmt.Errorf("entity concurrent modification: got revision %v", got)
}
