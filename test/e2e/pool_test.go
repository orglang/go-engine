package e2e

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/fx"

	"orglang/go-engine/lib/e2e"

	"github.com/orglang/go-sdk/lib/rc"

	"github.com/orglang/go-sdk/adt/descvar"
	"github.com/orglang/go-sdk/adt/implsem"
	"github.com/orglang/go-sdk/adt/implvar"
	"github.com/orglang/go-sdk/adt/pooldec"
	"github.com/orglang/go-sdk/adt/poolexec"
	"github.com/orglang/go-sdk/adt/poolexp"
	"github.com/orglang/go-sdk/adt/poolstep"
	"github.com/orglang/go-sdk/adt/procdec"
	"github.com/orglang/go-sdk/adt/procexp"
	"github.com/orglang/go-sdk/adt/procstep"
	"github.com/orglang/go-sdk/adt/typedef"
	"github.com/orglang/go-sdk/adt/typeexp"
	"github.com/orglang/go-sdk/adt/xactdef"
	"github.com/orglang/go-sdk/adt/xactexp"
)

func TestPool(t *testing.T) {
	s := suite{}
	s.beforeAll(t)
	t.Run("WaitClose", s.waitClose)
	// t.Run("RecvSend", s.recvSend)
	// t.Run("CaseLab", s.caseLab)
	// t.Run("Call", s.call)
	// t.Run("Fwd", s.fwd)
}

type suite struct {
	fx.In
	PoolDecAPI  e2e.PoolDecAPI
	PoolExecAPI e2e.PoolExecAPI
	XactDefAPI  e2e.XactDefAPI
	ProcDecAPI  e2e.ProcDecAPI
	ProcExecAPI e2e.ProcExecAPI
	TypeDefAPI  e2e.TypeDefAPI
	DB          *sql.DB `optional:"true"`
}

func (s *suite) beforeAll(t *testing.T) {
	db, err := sql.Open("pgx", "postgres://orglang:orglang@localhost:5432/orglang")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		err := db.Close()
		if err != nil {
			t.Logf("stop failed: %s", err)
		}
	})
	app := fx.New(rc.Module, e2e.Module, fx.Populate(s))
	s.DB = db
	t.Cleanup(func() {
		err := app.Stop(t.Context())
		if err != nil {
			t.Logf("stop failed: %s", err)
		}
	})
}

func (s *suite) beforeEach(t *testing.T) {
	tables := []string{
		"desc_sems", "desc_binds",
		"impl_sems", "impl_binds",
		"xact_defs", "xact_exps",
		"pool_decs", "pool_execs", "pool_vars",
		"type_defs", "type_exps",
		"proc_decs", "proc_vars", "proc_steps",
	}
	for _, table := range tables {
		_, err := s.DB.Exec(fmt.Sprintf("truncate table %v", table))
		if err != nil {
			t.Fatal(err)
		}
	}
}

func (s *suite) waitClose(t *testing.T) {
	s.beforeEach(t)
	// given
	closerProcQN := "closer-proc-qn"
	waiterProcQN := "waiter-proc-qn"
	// and
	withXactQN := "with-xact-qn"
	_, err := s.XactDefAPI.Create(xactdef.DefSpec{
		XactQN: withXactQN,
		XactES: xactexp.ExpSpec{
			K: xactexp.With,
			With: &xactexp.LaborSpec{
				Choices: []xactexp.ChoiceSpec{
					// пул заявляет способность работать closerProcQN
					{ProcQN: closerProcQN, ContES: xactexp.ExpSpec{
						K:    xactexp.Link,
						Link: &xactexp.LinkSpec{XactQN: withXactQN},
					}},
					// пул заявляет способность работать waiterProcQN
					{ProcQN: waiterProcQN, ContES: xactexp.ExpSpec{
						K:    xactexp.Link,
						Link: &xactexp.LinkSpec{XactQN: withXactQN},
					}},
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	poolDescQN := "pool-desc-qn"
	poolImplQN := "pool-impl-qn"
	poolClientPH := "pool-client-ph"
	poolProviderPH := "pool-provider-ph"
	_, err = s.PoolDecAPI.Create(pooldec.DecSpec{
		DescQN: poolDescQN,
		ProviderVS: descvar.VarSpec{
			ChnlPH: poolProviderPH,
			DescQN: withXactQN,
		},
		ClientVSes: []descvar.VarSpec{{
			ChnlPH: poolClientPH,
			DescQN: withXactQN,
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	myPoolExec, err := s.PoolExecAPI.Create(poolexec.ExecSpec{
		DescQN: poolDescQN,
		ProviderVS: implvar.VarSpec{
			ChnlPH: poolProviderPH,
			ImplQN: poolImplQN,
		},
		ClientVSes: []implvar.VarSpec{{
			ChnlPH: poolClientPH,
			ImplQN: poolImplQN,
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	err = s.PoolExecAPI.Take(poolstep.StepSpec{
		ImplRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Hire, // пул делает предложение поработать в качестве closerProcQN
			Hire: &poolexp.HireSpec{
				CommChnlPH: poolClientPH, // пул выступает в качестве клиента самого себя
				ProcDescQN: closerProcQN,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	err = s.PoolExecAPI.Take(poolstep.StepSpec{
		ImplRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Apply, // пул принимает предложение поработать в качестве closerProcQN
			Apply: &poolexp.ApplySpec{
				CommChnlPH: poolProviderPH, // пул выступает в качестве провайдера самому себе
				ProcDescQN: closerProcQN,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	oneTypeQN := "one-type-qn"
	_, err = s.TypeDefAPI.Create(typedef.DefSpec{
		TypeQN: oneTypeQN,
		TypeES: typeexp.ExpSpec{K: typeexp.One},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	closerProcDec, err := s.ProcDecAPI.Create(procdec.DecSpec{
		DescQN: closerProcQN,
		ProviderVS: descvar.VarSpec{
			ChnlPH: "closer-provider-ph",
			DescQN: oneTypeQN,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	closerProcExec, err := s.PoolExecAPI.Spawn(poolstep.StepSpec{
		ImplRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Spawn,
			Spawn: &poolexp.SpawnSpec{
				ProcDescRef: closerProcDec.DescRef,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	err = s.PoolExecAPI.Take(poolstep.StepSpec{
		ImplRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Hire,
			Hire: &poolexp.HireSpec{
				CommChnlPH: poolClientPH,
				ProcDescQN: waiterProcQN,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	err = s.PoolExecAPI.Take(poolstep.StepSpec{
		ImplRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Apply,
			Apply: &poolexp.ApplySpec{
				CommChnlPH: poolProviderPH,
				ProcDescQN: waiterProcQN,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	waiterProcDec, err := s.ProcDecAPI.Create(procdec.DecSpec{
		DescQN:     waiterProcQN,
		ProviderVS: descvar.VarSpec{ChnlPH: "waiter-provider-ph", DescQN: oneTypeQN},
		ClientVSes: []descvar.VarSpec{
			{ChnlPH: "waiter-closer-ph", DescQN: oneTypeQN},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	waiterProcExec, err := s.PoolExecAPI.Spawn(poolstep.StepSpec{
		ImplRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Spawn,
			Spawn: &poolexp.SpawnSpec{
				ProcDescRef:  waiterProcDec.DescRef,
				ProcImplRefs: []implsem.SemRef{closerProcExec},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// when
	err = s.ProcExecAPI.Take(procstep.StepSpec{
		ExecRef: closerProcExec,
		ProcES: procexp.ExpSpec{
			K: procexp.Close,
			Close: &procexp.CloseSpec{
				CommChnlPH: closerProcDec.ProviderVR.ChnlPH,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	err = s.ProcExecAPI.Take(procstep.StepSpec{
		ExecRef: waiterProcExec,
		ProcES: procexp.ExpSpec{
			K: procexp.Wait,
			Wait: &procexp.WaitSpec{
				CommChnlPH: waiterProcDec.ClientVRs[0].ChnlPH,
				ContES: procexp.ExpSpec{
					K: procexp.Close,
					Close: &procexp.CloseSpec{
						CommChnlPH: waiterProcDec.ProviderVR.ChnlPH,
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// then
	// TODO добавить проверку
}

func (s *suite) recvSend(t *testing.T) {
	s.beforeEach(t)
	// given
	senderProcQN := "sender-proc-qn"
	messageProcQN := "message-proc-qn"
	receiverProcQN := "receiver-proc-qn"
	// and
	withXactQN := "with-xact-qn"
	_, err := s.XactDefAPI.Create(xactdef.DefSpec{
		XactQN: withXactQN,
		XactES: xactexp.ExpSpec{
			K: xactexp.With,
			With: &xactexp.LaborSpec{
				Choices: []xactexp.ChoiceSpec{
					{ProcQN: senderProcQN, ContES: xactexp.ExpSpec{
						K:    xactexp.Link,
						Link: &xactexp.LinkSpec{XactQN: withXactQN},
					}},
					{ProcQN: receiverProcQN, ContES: xactexp.ExpSpec{
						K:    xactexp.Link,
						Link: &xactexp.LinkSpec{XactQN: withXactQN},
					}},
					{ProcQN: messageProcQN, ContES: xactexp.ExpSpec{
						K:    xactexp.Link,
						Link: &xactexp.LinkSpec{XactQN: withXactQN},
					}},
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	myPoolQN := "my-pool-qn"
	poolProviderPH := "pool-provider-ph"
	poolClientPH := "pool-client-ph"
	_, err = s.PoolDecAPI.Create(pooldec.DecSpec{
		DescQN: myPoolQN,
		ProviderVS: descvar.VarSpec{
			ChnlPH: poolProviderPH,
			DescQN: withXactQN,
		},
		ClientVSes: []descvar.VarSpec{{
			ChnlPH: poolClientPH,
			DescQN: withXactQN,
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	myPoolExec, err := s.PoolExecAPI.Create(poolexec.ExecSpec{
		DescQN: myPoolQN,
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	lolliTypeQN := "lolli-type-qn"
	_, err = s.TypeDefAPI.Create(typedef.DefSpec{
		TypeQN: lolliTypeQN,
		TypeES: typeexp.ExpSpec{
			K: typeexp.Lolli,
			Lolli: &typeexp.ProdSpec{
				Val:  typeexp.ExpSpec{K: typeexp.One},
				Cont: typeexp.ExpSpec{K: typeexp.One},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	oneTypeQN := "one-type-qn"
	_, err = s.TypeDefAPI.Create(typedef.DefSpec{
		TypeQN: oneTypeQN,
		TypeES: typeexp.ExpSpec{K: typeexp.One},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	receiverProcDec, err := s.ProcDecAPI.Create(procdec.DecSpec{
		DescQN: receiverProcQN,
		ProviderVS: descvar.VarSpec{
			ChnlPH: "receiver-provider-ph",
			DescQN: lolliTypeQN,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	messageProcDec, err := s.ProcDecAPI.Create(procdec.DecSpec{
		DescQN: messageProcQN,
		ProviderVS: descvar.VarSpec{
			ChnlPH: "message-provider-ph",
			DescQN: oneTypeQN,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	senderProcDec, err := s.ProcDecAPI.Create(procdec.DecSpec{
		DescQN: senderProcQN,
		ProviderVS: descvar.VarSpec{
			ChnlPH: "sender-provider-ph",
			DescQN: oneTypeQN,
		},
		ClientVSes: []descvar.VarSpec{
			{ChnlPH: "sender-receiver-ph", DescQN: lolliTypeQN},
			{ChnlPH: "sender-message-ph", DescQN: oneTypeQN},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	err = s.PoolExecAPI.Take(poolstep.StepSpec{
		ImplRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Hire,
			Hire: &poolexp.HireSpec{
				CommChnlPH: poolClientPH,
				ProcDescQN: receiverProcQN,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	receiverProcExec, err := s.PoolExecAPI.Spawn(poolstep.StepSpec{
		ImplRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Spawn,
			Spawn: &poolexp.SpawnSpec{
				ProcDescRef: receiverProcDec.DescRef,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	err = s.PoolExecAPI.Take(poolstep.StepSpec{
		ImplRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Hire,
			Hire: &poolexp.HireSpec{
				CommChnlPH: poolClientPH,
				ProcDescQN: messageProcQN,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	messageProcExec, err := s.PoolExecAPI.Spawn(poolstep.StepSpec{
		ImplRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Spawn,
			Spawn: &poolexp.SpawnSpec{
				ProcDescRef: messageProcDec.DescRef,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	err = s.PoolExecAPI.Take(poolstep.StepSpec{
		ImplRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Hire,
			Hire: &poolexp.HireSpec{
				CommChnlPH: poolClientPH,
				ProcDescQN: senderProcQN,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	senderProcExec, err := s.PoolExecAPI.Spawn(poolstep.StepSpec{
		ImplRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Spawn,
			Spawn: &poolexp.SpawnSpec{
				ProcDescRef:  senderProcDec.DescRef,
				ProcImplRefs: []implsem.SemRef{receiverProcExec, messageProcExec},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// when
	err = s.ProcExecAPI.Take(procstep.StepSpec{
		ExecRef: receiverProcExec,
		ProcES: procexp.ExpSpec{
			K: procexp.Recv,
			Recv: &procexp.RecvSpec{
				CommChnlPH: receiverProcDec.ProviderVR.ChnlPH,
				BindChnlPH: "receiver-message-ph",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	err = s.ProcExecAPI.Take(procstep.StepSpec{
		ExecRef: senderProcExec,
		ProcES: procexp.ExpSpec{
			K: procexp.Send,
			Send: &procexp.SendSpec{
				CommChnlPH: senderProcDec.ClientVRs[0].ChnlPH,
				ValChnlPH:  senderProcDec.ClientVRs[1].ChnlPH,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// then
	// TODO добавить проверку
}

func (s *suite) caseLab(t *testing.T) {
	s.beforeEach(t)
	// given
	myPoolQN := "my-pool-qn"
	myPoolExec, err := s.PoolExecAPI.Create(poolexec.ExecSpec{
		DescQN: myPoolQN,
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	patternQN := "pattern-qn"
	// and
	withType, err := s.TypeDefAPI.Create(typedef.DefSpec{
		TypeQN: "with-type-qn",
		TypeES: typeexp.ExpSpec{
			With: &typeexp.SumSpec{
				Choices: []typeexp.ChoiceSpec{
					{LabQN: patternQN, Cont: typeexp.ExpSpec{K: typeexp.One}},
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	oneType, err := s.TypeDefAPI.Create(typedef.DefSpec{
		TypeQN: "one-type-qn",
		TypeES: typeexp.ExpSpec{K: typeexp.One},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	followerProcQN := "follower-proc-qn"
	followerProcDec := procdec.DecSpec{
		DescQN: followerProcQN,
		ProviderVS: descvar.VarSpec{
			ChnlPH: "follower-provider-ph",
			DescQN: withType.DefSpec.TypeQN,
		},
	}
	followerProcDec2, err := s.ProcDecAPI.Create(followerProcDec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	deciderProcQN := "decider-proc-qn"
	deciderProcDec, err := s.ProcDecAPI.Create(procdec.DecSpec{
		DescQN:     deciderProcQN,
		ClientVSes: []descvar.VarSpec{followerProcDec.ProviderVS},
		ProviderVS: descvar.VarSpec{
			ChnlPH: "decider-provider-ph",
			DescQN: oneType.DefSpec.TypeQN,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	poolFollowerPH := "pool-follower-ph"
	err = s.PoolExecAPI.Take(poolstep.StepSpec{
		ImplRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Hire,
			Hire: &poolexp.HireSpec{
				ProcDescQN: followerProcQN,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	followerProcExec, err := s.PoolExecAPI.Spawn(poolstep.StepSpec{
		ImplRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Spawn,
			Spawn: &poolexp.SpawnSpec{
				ProcDescRef: followerProcDec2.DescRef,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	err = s.PoolExecAPI.Take(poolstep.StepSpec{
		ImplRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Hire,
			Hire: &poolexp.HireSpec{
				ProcDescQN: deciderProcQN,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	deciderProcExec, err := s.PoolExecAPI.Spawn(poolstep.StepSpec{
		ImplRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Spawn,
			Spawn: &poolexp.SpawnSpec{
				ProcDescRef:  deciderProcDec.DescRef,
				ProcImplRefs: []implsem.SemRef{followerProcExec},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// when
	err = s.ProcExecAPI.Take(procstep.StepSpec{
		ExecRef: followerProcExec,
		ProcES: procexp.ExpSpec{
			K: procexp.Case,
			Case: &procexp.CaseSpec{
				CommChnlPH: poolFollowerPH,
				ContBSes: []procexp.BranchSpec{
					{PatternQN: patternQN, ContES: procexp.ExpSpec{
						K: procexp.Close,
						Close: &procexp.CloseSpec{
							CommChnlPH: poolFollowerPH,
						},
					},
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	err = s.ProcExecAPI.Take(procstep.StepSpec{
		ExecRef: deciderProcExec,
		ProcES: procexp.ExpSpec{
			K: procexp.Lab,
			Lab: &procexp.LabSpec{
				CommChnlPH: poolFollowerPH,
				PatternQN:  patternQN,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// then
	// TODO добавить проверку
}

func (s *suite) call(t *testing.T) {
	s.beforeEach(t)
	// given
	myPoolQN := "my-pool-qn"
	myPoolExec, err := s.PoolExecAPI.Create(poolexec.ExecSpec{
		DescQN: myPoolQN,
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	oneType, err := s.TypeDefAPI.Create(
		typedef.DefSpec{
			TypeQN: "one-type-qn",
			TypeES: typeexp.ExpSpec{K: typeexp.One},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	// and
	injecteeProcQN := "injectee-proc-qn"
	injecteeProcDec, err := s.ProcDecAPI.Create(procdec.DecSpec{
		DescQN: injecteeProcQN,
		ProviderVS: descvar.VarSpec{
			ChnlPH: "injectee-provider-ph",
			DescQN: oneType.DefSpec.TypeQN,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	err = s.PoolExecAPI.Take(poolstep.StepSpec{
		ImplRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Hire,
			Hire: &poolexp.HireSpec{
				ProcDescQN: injecteeProcQN,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	injecteeProcExec, err := s.PoolExecAPI.Spawn(poolstep.StepSpec{
		ImplRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Spawn,
			Spawn: &poolexp.SpawnSpec{
				ProcDescRef: injecteeProcDec.DescRef,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	callerProcQN := "caller-proc-qn"
	callerProcDec, err := s.ProcDecAPI.Create(procdec.DecSpec{
		DescQN: callerProcQN,
		ProviderVS: descvar.VarSpec{
			ChnlPH: "caller-provider-ph",
			DescQN: oneType.DefSpec.TypeQN,
		},
		ClientVSes: []descvar.VarSpec{
			{ChnlPH: "caller-injectee-ph", DescQN: oneType.DefSpec.TypeQN},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	err = s.PoolExecAPI.Take(poolstep.StepSpec{
		ImplRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Hire,
			Hire: &poolexp.HireSpec{
				ProcDescQN: callerProcQN,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	callerProcExec, err := s.PoolExecAPI.Spawn(poolstep.StepSpec{
		ImplRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Spawn,
			Spawn: &poolexp.SpawnSpec{
				ProcDescRef:  callerProcDec.DescRef,
				ProcImplRefs: []implsem.SemRef{injecteeProcExec},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	calleeProcQN := "callee-proc-qn"
	_, err = s.ProcDecAPI.Create(procdec.DecSpec{
		DescQN: calleeProcQN,
		ProviderVS: descvar.VarSpec{
			ChnlPH: "callee-provider-ph",
			DescQN: oneType.DefSpec.TypeQN,
		},
		ClientVSes: []descvar.VarSpec{
			{ChnlPH: "callee-injectee-ph", DescQN: oneType.DefSpec.TypeQN},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	callerCalleePH := "caller-callee-ph"
	// when
	err = s.ProcExecAPI.Take(procstep.StepSpec{
		ExecRef: callerProcExec,
		ProcES: procexp.ExpSpec{
			K: procexp.Call,
			Call: &procexp.CallSpec{
				BindChnlPH: callerCalleePH,
				ProcDescQN: calleeProcQN,
				ValChnlPHs: []string{callerProcDec.ClientVRs[0].ChnlPH},
				ContES: procexp.ExpSpec{
					K: procexp.Wait,
					Wait: &procexp.WaitSpec{
						CommChnlPH: callerCalleePH,
						ContES: procexp.ExpSpec{
							K: procexp.Close,
							Close: &procexp.CloseSpec{
								CommChnlPH: callerProcDec.ProviderVR.ChnlPH,
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// then
	// TODO добавить проверку
}

func (s *suite) fwd(t *testing.T) {
	s.beforeEach(t)
	// given
	myPoolQN := "my-pool-qn"
	myPoolExec, err := s.PoolExecAPI.Create(poolexec.ExecSpec{
		DescQN: myPoolQN,
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	oneType, err := s.TypeDefAPI.Create(typedef.DefSpec{
		TypeQN: "one-type-qn",
		TypeES: typeexp.ExpSpec{K: typeexp.One},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	closerProcQN := "closer-proc-qn"
	closerProcDec, err := s.ProcDecAPI.Create(procdec.DecSpec{
		DescQN: closerProcQN,
		ProviderVS: descvar.VarSpec{
			ChnlPH: "closer-provider-ph",
			DescQN: oneType.DefSpec.TypeQN,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	forwarderProcQN := "forwarder-proc-qn"
	forwarderProcDec, err := s.ProcDecAPI.Create(procdec.DecSpec{
		DescQN: forwarderProcQN,
		ProviderVS: descvar.VarSpec{
			ChnlPH: "forwarder-provider-ph",
			DescQN: oneType.DefSpec.TypeQN,
		},
		ClientVSes: []descvar.VarSpec{
			{ChnlPH: "forwarder-closer-ph", DescQN: oneType.DefSpec.TypeQN},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	waiterProcQN := "waiter-proc-qn"
	waiterProcDec, err := s.ProcDecAPI.Create(procdec.DecSpec{
		DescQN: waiterProcQN,
		ProviderVS: descvar.VarSpec{
			ChnlPH: "waiter-provider-ph",
			DescQN: oneType.DefSpec.TypeQN,
		},
		ClientVSes: []descvar.VarSpec{
			{ChnlPH: "waiter-forwarder-ph", DescQN: oneType.DefSpec.TypeQN},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	err = s.PoolExecAPI.Take(poolstep.StepSpec{
		ImplRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Hire,
			Hire: &poolexp.HireSpec{
				ProcDescQN: closerProcQN,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	closerProcExec, err := s.PoolExecAPI.Spawn(poolstep.StepSpec{
		ImplRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Spawn,
			Spawn: &poolexp.SpawnSpec{
				ProcDescRef: closerProcDec.DescRef,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	err = s.PoolExecAPI.Take(poolstep.StepSpec{
		ImplRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Hire,
			Hire: &poolexp.HireSpec{
				ProcDescQN: forwarderProcQN,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	forwarderProcExec, err := s.PoolExecAPI.Spawn(poolstep.StepSpec{
		ImplRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Spawn,
			Spawn: &poolexp.SpawnSpec{
				ProcDescRef:  forwarderProcDec.DescRef,
				ProcImplRefs: []implsem.SemRef{closerProcExec},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	err = s.PoolExecAPI.Take(poolstep.StepSpec{
		ImplRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Hire,
			Hire: &poolexp.HireSpec{
				ProcDescQN: waiterProcQN,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	waiterProcExec, err := s.PoolExecAPI.Spawn(poolstep.StepSpec{
		ImplRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Spawn,
			Spawn: &poolexp.SpawnSpec{
				ProcDescRef:  waiterProcDec.DescRef,
				ProcImplRefs: []implsem.SemRef{forwarderProcExec},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	err = s.ProcExecAPI.Take(procstep.StepSpec{
		ExecRef: closerProcExec,
		ProcES: procexp.ExpSpec{
			K: procexp.Close,
			Close: &procexp.CloseSpec{
				CommChnlPH: closerProcDec.ProviderVR.ChnlPH,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// when
	err = s.ProcExecAPI.Take(procstep.StepSpec{
		ExecRef: forwarderProcExec,
		ProcES: procexp.ExpSpec{
			K: procexp.Fwd,
			Fwd: &procexp.FwdSpec{
				CommChnlPH: forwarderProcDec.ProviderVR.ChnlPH,
				ContChnlPH: forwarderProcDec.ClientVRs[0].ChnlPH,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	err = s.ProcExecAPI.Take(procstep.StepSpec{
		ExecRef: waiterProcExec,
		ProcES: procexp.ExpSpec{
			K: procexp.Wait,
			Wait: &procexp.WaitSpec{
				CommChnlPH: waiterProcDec.ClientVRs[0].ChnlPH,
				ContES: procexp.ExpSpec{
					K: procexp.Close,
					Close: &procexp.CloseSpec{
						CommChnlPH: waiterProcDec.ProviderVR.ChnlPH,
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// then
	// TODO добавить проверку
}
