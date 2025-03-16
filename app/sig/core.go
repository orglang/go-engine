package sig

import (
	"context"
	"fmt"
	"log/slog"

	"smecalculus/rolevod/lib/data"
	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/ph"
	"smecalculus/rolevod/lib/rev"
	"smecalculus/rolevod/lib/sym"

	"smecalculus/rolevod/internal/alias"
	"smecalculus/rolevod/internal/chnl"

	"smecalculus/rolevod/app/role"
)

type ID = id.ADT
type QN = sym.ADT
type Name = string

type Spec struct {
	QN sym.ADT
	Ys []chnl.Spec // vals
	X  chnl.Spec   // via
}

type Ref struct {
	SigID id.ADT
	Title string
	Rev   rev.ADT
}

type Snap struct {
	SigID id.ADT
	Title string
	Ys    []chnl.Spec
	X     chnl.Spec
	Rev   rev.ADT
}

// aka ExpDec or ExpDecDef without expression
type Root struct {
	SigID id.ADT
	Title string
	Ys    []chnl.Spec
	X     chnl.Spec
	Rev   rev.ADT
}

// aka ChanTp
type EP struct {
	ChnlPH ph.ADT
	RoleQN sym.ADT
}

type API interface {
	Incept(QN) (Ref, error)
	Create(Spec) (Root, error)
	Retrieve(id.ADT) (Root, error)
	RetreiveRefs() ([]Ref, error)
}

type service struct {
	sigs     Repo
	aliases  alias.Repo
	operator data.Operator
	log      *slog.Logger
}

func newService(sigs Repo, aliases alias.Repo, operator data.Operator, l *slog.Logger) *service {
	name := slog.String("name", "sigService")
	return &service{sigs, aliases, operator, l.With(name)}
}

// for compilation purposes
func newAPI() API {
	return &service{}
}

func (s *service) Incept(fqn sym.ADT) (_ Ref, err error) {
	ctx := context.Background()
	fqnAttr := slog.Any("fqn", fqn)
	s.log.Debug("inception started", fqnAttr)
	newAlias := alias.Root{Sym: fqn, ID: id.New(), Rev: rev.Initial()}
	newRoot := Root{SigID: newAlias.ID, Rev: newAlias.Rev, Title: newAlias.Sym.Name()}
	s.operator.Explicit(ctx, func(ds data.Source) error {
		err = s.aliases.Insert(ds, newAlias)
		if err != nil {
			return err
		}
		err = s.sigs.Insert(ds, newRoot)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		s.log.Error("inception failed", fqnAttr)
		return Ref{}, err
	}
	s.log.Debug("inception succeeded", fqnAttr, slog.Any("id", newRoot.SigID))
	return ConvertRootToRef(newRoot), nil
}

func (s *service) Create(spec Spec) (_ Root, err error) {
	ctx := context.Background()
	fqnAttr := slog.Any("fqn", spec.QN)
	s.log.Debug("creation started", fqnAttr, slog.Any("spec", spec))
	root := Root{
		SigID: id.New(),
		Title: spec.QN.Name(),
		Ys:    spec.Ys,
		X:     spec.X,
		Rev:   rev.Initial(),
	}
	s.operator.Explicit(ctx, func(ds data.Source) error {
		err = s.sigs.Insert(ds, root)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		s.log.Error("creation failed", fqnAttr)
		return Root{}, err
	}
	s.log.Debug("creation succeeded", fqnAttr, slog.Any("id", root.SigID))
	return root, nil
}

func (s *service) Retrieve(rid ID) (root Root, err error) {
	ctx := context.Background()
	s.operator.Implicit(ctx, func(ds data.Source) error {
		root, err = s.sigs.SelectByID(ds, rid)
		return err
	})
	if err != nil {
		s.log.Error("retrieval failed", slog.Any("id", rid))
		return Root{}, err
	}
	return root, nil
}

func (s *service) RetreiveRefs() (refs []Ref, err error) {
	ctx := context.Background()
	s.operator.Implicit(ctx, func(ds data.Source) error {
		refs, err = s.sigs.SelectAll(ds)
		return err
	})
	if err != nil {
		s.log.Error("retrieval failed")
		return nil, err
	}
	return refs, nil
}

type Repo interface {
	Insert(data.Source, Root) error
	SelectAll(data.Source) ([]Ref, error)
	SelectByID(data.Source, ID) (Root, error)
	SelectByIDs(data.Source, []ID) ([]Root, error)
	SelectEnv(data.Source, []ID) (map[ID]Root, error)
}

func CollectEnv(sigs []Root) []role.QN {
	roleQNs := []role.QN{}
	for _, sig := range sigs {
		roleQNs = append(roleQNs, sig.X.RoleQN)
		for _, ce := range sig.Ys {
			roleQNs = append(roleQNs, ce.RoleQN)
		}
	}
	return roleQNs
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
var (
	ConvertRootToRef func(Root) Ref
)

func ErrRootMissingInEnv(rid ID) error {
	return fmt.Errorf("root missing in env: %v", rid)
}
