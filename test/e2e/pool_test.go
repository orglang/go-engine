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
	"github.com/orglang/go-sdk/adt/implsubst"
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
		"pool_decs", "pool_execs", "pool_liabs",
		"type_defs", "type_exps",
		"proc_decs", "proc_binds", "proc_steps",
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
		ProviderSS: implsubst.SubstSpec{
			ChnlPH: poolProviderPH,
			ImplQN: poolImplQN,
		},
		ClientSSes: []implsubst.SubstSpec{{
			ChnlPH: poolClientPH,
			ImplQN: poolImplQN,
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	poolCloserPH := "pool-closer-ph"
	err = s.PoolExecAPI.Take(poolstep.StepSpec{
		ExecRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Hire, // пул делает предложение поработать в качестве closerProcQN
			Hire: &poolexp.HireSpec{
				CommPH: poolClientPH, // пул выступает в качестве клиента самого себя
				ProcQN: closerProcQN,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	err = s.PoolExecAPI.Take(poolstep.StepSpec{
		ExecRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Apply, // пул принимает предложение поработать в качестве closerProcQN
			Apply: &poolexp.ApplySpec{
				CommPH: poolProviderPH, // пул выступает в качестве провайдера самому себе
				ProcQN: closerProcQN,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	closerProcExec, err := s.PoolExecAPI.Spawn(poolstep.StepSpec{
		ExecRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Spawn,
			Spawn: &poolexp.SpawnSpec{
				BindPH: poolCloserPH,
				ProcQN: closerProcQN,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	poolWaiterPH := "pool-waiter-ph"
	err = s.PoolExecAPI.Take(poolstep.StepSpec{
		ExecRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Hire,
			Hire: &poolexp.HireSpec{
				CommPH: poolClientPH,
				ProcQN: waiterProcQN,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	err = s.PoolExecAPI.Take(poolstep.StepSpec{
		ExecRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Apply,
			Apply: &poolexp.ApplySpec{
				CommPH: poolProviderPH,
				ProcQN: waiterProcQN,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	waiterProcExec, err := s.PoolExecAPI.Spawn(poolstep.StepSpec{
		ExecRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Spawn,
			Spawn: &poolexp.SpawnSpec{
				BindPH: poolWaiterPH,
				ProcQN: waiterProcQN,
				ValPHs: []string{poolCloserPH},
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
	// when
	err = s.ProcExecAPI.Take(procstep.StepSpec{
		ExecRef: closerProcExec,
		ProcES: procexp.ExpSpec{
			K: procexp.Close,
			Close: &procexp.CloseSpec{
				CommPH: closerProcDec.ProviderVR.ChnlPH,
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
				CommPH: waiterProcDec.ClientVRs[0].ChnlPH,
				ContES: procexp.ExpSpec{
					K: procexp.Close,
					Close: &procexp.CloseSpec{
						CommPH: waiterProcDec.ProviderVR.ChnlPH,
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
	_, err = s.ProcDecAPI.Create(procdec.DecSpec{
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
	poolReceiverPH := "pool-receiver-ph"
	err = s.PoolExecAPI.Take(poolstep.StepSpec{
		ExecRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Hire,
			Hire: &poolexp.HireSpec{
				CommPH: poolClientPH,
				ProcQN: receiverProcQN,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	receiverProcExec, err := s.PoolExecAPI.Spawn(poolstep.StepSpec{
		ExecRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Spawn,
			Spawn: &poolexp.SpawnSpec{
				BindPH: poolReceiverPH,
				ProcQN: receiverProcQN,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	poolMessagePH := "pool-message-ph"
	err = s.PoolExecAPI.Take(poolstep.StepSpec{
		ExecRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Hire,
			Hire: &poolexp.HireSpec{
				CommPH: poolClientPH,
				ProcQN: messageProcQN,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.PoolExecAPI.Spawn(poolstep.StepSpec{
		ExecRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Spawn,
			Spawn: &poolexp.SpawnSpec{
				BindPH: poolMessagePH,
				ProcQN: messageProcQN,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	poolSenderPH := "pool-sender-ph"
	err = s.PoolExecAPI.Take(poolstep.StepSpec{
		ExecRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Hire,
			Hire: &poolexp.HireSpec{
				CommPH: poolClientPH,
				ProcQN: senderProcQN,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	senderProcExec, err := s.PoolExecAPI.Spawn(poolstep.StepSpec{
		ExecRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Spawn,
			Spawn: &poolexp.SpawnSpec{
				BindPH: poolSenderPH,
				ProcQN: senderProcQN,
				ValPHs: []string{poolReceiverPH, poolMessagePH},
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
				CommPH: receiverProcDec.ProviderVR.ChnlPH,
				BindPH: "receiver-message-ph",
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
				CommPH: senderProcDec.ClientVRs[0].ChnlPH,
				ValPH:  senderProcDec.ClientVRs[1].ChnlPH,
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
	labelQN := "label-qn"
	// and
	withType, err := s.TypeDefAPI.Create(typedef.DefSpec{
		TypeQN: "with-type-qn",
		TypeES: typeexp.ExpSpec{
			With: &typeexp.SumSpec{
				Choices: []typeexp.ChoiceSpec{
					{LabQN: labelQN, Cont: typeexp.ExpSpec{K: typeexp.One}},
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
	_, err = s.ProcDecAPI.Create(followerProcDec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	deciderProcQN := "decider-proc-qn"
	_, err = s.ProcDecAPI.Create(procdec.DecSpec{
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
		ExecRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Hire,
			Hire: &poolexp.HireSpec{
				ProcQN: followerProcQN,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	followerProcExec, err := s.PoolExecAPI.Spawn(poolstep.StepSpec{
		ExecRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Spawn,
			Spawn: &poolexp.SpawnSpec{
				BindPH: poolFollowerPH,
				ProcQN: followerProcQN,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	poolDeciderPH := "pool-decider-ph"
	err = s.PoolExecAPI.Take(poolstep.StepSpec{
		ExecRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Hire,
			Hire: &poolexp.HireSpec{
				ProcQN: deciderProcQN,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	deciderProcExec, err := s.PoolExecAPI.Spawn(poolstep.StepSpec{
		ExecRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Spawn,
			Spawn: &poolexp.SpawnSpec{
				BindPH: poolDeciderPH,
				ProcQN: deciderProcQN,
				ValPHs: []string{poolFollowerPH},
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
				CommPH: poolFollowerPH,
				ContBSs: []procexp.BranchSpec{
					{LabQN: labelQN, ContES: procexp.ExpSpec{
						K: procexp.Close,
						Close: &procexp.CloseSpec{
							CommPH: poolFollowerPH,
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
				CommPH: poolFollowerPH,
				InfoQN: labelQN,
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
	_, err = s.ProcDecAPI.Create(procdec.DecSpec{
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
		ExecRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Hire,
			Hire: &poolexp.HireSpec{
				ProcQN: injecteeProcQN,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	poolInjecteePH := "pool-injectee-ph"
	_, err = s.PoolExecAPI.Spawn(poolstep.StepSpec{
		ExecRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Spawn,
			Spawn: &poolexp.SpawnSpec{
				BindPH: poolInjecteePH,
				ProcQN: injecteeProcQN,
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
		ExecRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Hire,
			Hire: &poolexp.HireSpec{
				ProcQN: callerProcQN,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	callerProcExec, err := s.PoolExecAPI.Spawn(poolstep.StepSpec{
		ExecRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Spawn,
			Spawn: &poolexp.SpawnSpec{
				BindPH: "pool-caller-ph",
				ProcQN: callerProcQN,
				ValPHs: []string{poolInjecteePH},
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
				BindPH: callerCalleePH,
				ProcQN: calleeProcQN,
				ValPHs: []string{callerProcDec.ClientVRs[0].ChnlPH},
				ContES: procexp.ExpSpec{
					K: procexp.Wait,
					Wait: &procexp.WaitSpec{
						CommPH: callerCalleePH,
						ContES: procexp.ExpSpec{
							K: procexp.Close,
							Close: &procexp.CloseSpec{
								CommPH: callerProcDec.ProviderVR.ChnlPH,
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
		ExecRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Hire,
			Hire: &poolexp.HireSpec{
				ProcQN: closerProcQN,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	poolCloserPH := "pool-closer-ph"
	closerProcExec, err := s.PoolExecAPI.Spawn(poolstep.StepSpec{
		ExecRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Spawn,
			Spawn: &poolexp.SpawnSpec{
				BindPH: poolCloserPH,
				ProcQN: closerProcQN,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	err = s.PoolExecAPI.Take(poolstep.StepSpec{
		ExecRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Hire,
			Hire: &poolexp.HireSpec{
				ProcQN: forwarderProcQN,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	poolForwarderPH := "pool-forwarder-ph"
	forwarderProcExec, err := s.PoolExecAPI.Spawn(poolstep.StepSpec{
		ExecRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Spawn,
			Spawn: &poolexp.SpawnSpec{
				BindPH: poolForwarderPH,
				ProcQN: forwarderProcQN,
				ValPHs: []string{poolCloserPH},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	err = s.PoolExecAPI.Take(poolstep.StepSpec{
		ExecRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Hire,
			Hire: &poolexp.HireSpec{
				ProcQN: waiterProcQN,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	waiterProcExec, err := s.PoolExecAPI.Spawn(poolstep.StepSpec{
		ExecRef: myPoolExec,
		PoolES: poolexp.ExpSpec{
			K: poolexp.Spawn,
			Spawn: &poolexp.SpawnSpec{
				BindPH: "pool-waiter-ph",
				ProcQN: waiterProcQN,
				ValPHs: []string{poolForwarderPH},
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
				CommPH: closerProcDec.ProviderVR.ChnlPH,
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
				CommPH: forwarderProcDec.ProviderVR.ChnlPH,
				ContPH: forwarderProcDec.ClientVRs[0].ChnlPH,
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
				CommPH: waiterProcDec.ClientVRs[0].ChnlPH,
				ContES: procexp.ExpSpec{
					K: procexp.Close,
					Close: &procexp.CloseSpec{
						CommPH: waiterProcDec.ProviderVR.ChnlPH,
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
