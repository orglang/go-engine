package e2e

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

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
		"comm_sems",
		"xact_defs", "xact_exps",
		"pool_decs", "pool_execs", "pool_conns", "pool_vars", "pool_steps",
		"type_defs", "type_exps",
		"proc_decs", "proc_linear_vars", "proc_steps",
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
		XactExp: xactexp.ExpSpec{
			K: xactexp.Up,
			Up: &xactexp.ShiftSpec{
				ContExp: xactexp.ExpSpec{
					K: xactexp.With,
					With: &xactexp.LaborSpec{
						// пул заявляет компетенции
						ProcQNs: []string{closerProcQN, waiterProcQN},
						ContExp: xactexp.ExpSpec{
							K: xactexp.Down,
							Down: &xactexp.ShiftSpec{
								ContExp: xactexp.ExpSpec{
									K:    xactexp.Link,
									Link: &xactexp.LinkSpec{XactQN: withXactQN},
								},
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
	poolDescQN := "pool-desc-qn"
	poolImplQN := "pool-impl-qn"
	poolAssetPH := "pool-asset-ph"
	poolLiabPH := "pool-liab-ph"
	_, err = s.PoolDecAPI.Create(pooldec.DecSpec{
		DescQN: poolDescQN,
		LiabVar: descvar.VarSpec{
			ChnlPH: poolLiabPH,
			DescQN: withXactQN,
		},
		AssetVars: []descvar.VarSpec{{
			ChnlPH: poolAssetPH,
			DescQN: withXactQN,
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	myPoolExec, err := s.PoolExecAPI.Create(poolexec.ExecSpec{
		DescQN: poolDescQN,
		LiabVar: implvar.VarSpec{
			ChnlPH: poolLiabPH,
			ImplQN: poolImplQN,
		},
		AssetVars: []implvar.VarSpec{{
			ChnlPH: poolAssetPH,
			ImplQN: poolImplQN,
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	err = s.PoolExecAPI.Take(poolstep.StepSpec{
		ImplRef: myPoolExec,
		PoolExp: poolexp.ExpSpec{
			K: poolexp.Acquire, // пул запрашивает доступ к самому себе
			Acquire: &poolexp.AcquireSpec{
				CommChnlPH: poolAssetPH,
				ContExp: poolexp.ExpSpec{
					K: poolexp.Hire, // пул запрашивает компетенцию closerProcQN
					Hire: &poolexp.HireSpec{
						CommChnlPH: poolAssetPH, // пул выступает нанимателем самого себя
						ProcDescQN: closerProcQN,
						ContExp: poolexp.ExpSpec{
							K: poolexp.Release,
							Release: &poolexp.ReleaseSpec{
								CommChnlPH: poolAssetPH,
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
	err = s.PoolExecAPI.Take(poolstep.StepSpec{
		ImplRef: myPoolExec,
		PoolExp: poolexp.ExpSpec{
			K: poolexp.Accept, // пул одобряет доступ к самому себе
			Accept: &poolexp.AcceptSpec{
				CommChnlPH: poolLiabPH,
				ContExp: poolexp.ExpSpec{
					K: poolexp.Apply, // пул предлагает компетенцию closerProcQN
					Apply: &poolexp.ApplySpec{
						CommChnlPH: poolLiabPH, // пул выступает соискателем
						ProcDescQN: closerProcQN,
						ContExp: poolexp.ExpSpec{
							K: poolexp.Detach,
							Detach: &poolexp.DetachSpec{
								CommChnlPH: poolLiabPH,
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
	time.Sleep(1 * time.Second)
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
		LiabVar: descvar.VarSpec{
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
		PoolExp: poolexp.ExpSpec{
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
		PoolExp: poolexp.ExpSpec{
			K: poolexp.Acquire,
			Acquire: &poolexp.AcquireSpec{
				CommChnlPH: poolAssetPH,
				ContExp: poolexp.ExpSpec{
					K: poolexp.Hire,
					Hire: &poolexp.HireSpec{
						CommChnlPH: poolAssetPH,
						ProcDescQN: waiterProcQN,
						ContExp: poolexp.ExpSpec{
							K: poolexp.Release,
							Release: &poolexp.ReleaseSpec{
								CommChnlPH: poolAssetPH,
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
	err = s.PoolExecAPI.Take(poolstep.StepSpec{
		ImplRef: myPoolExec,
		PoolExp: poolexp.ExpSpec{
			K: poolexp.Accept,
			Accept: &poolexp.AcceptSpec{
				CommChnlPH: poolAssetPH,
				ContExp: poolexp.ExpSpec{
					K: poolexp.Apply,
					Apply: &poolexp.ApplySpec{
						CommChnlPH: poolLiabPH,
						ProcDescQN: waiterProcQN,
						ContExp: poolexp.ExpSpec{
							K: poolexp.Detach,
							Detach: &poolexp.DetachSpec{
								CommChnlPH: poolAssetPH,
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
	waiterProcDec, err := s.ProcDecAPI.Create(procdec.DecSpec{
		DescQN:  waiterProcQN,
		LiabVar: descvar.VarSpec{ChnlPH: "waiter-provider-ph", DescQN: oneTypeQN},
		AssetVars: []descvar.VarSpec{
			{ChnlPH: "waiter-closer-ph", DescQN: oneTypeQN},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	waiterProcExec, err := s.PoolExecAPI.Spawn(poolstep.StepSpec{
		ImplRef: myPoolExec,
		PoolExp: poolexp.ExpSpec{
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
		ImplRef: closerProcExec,
		ProcExp: procexp.ExpSpec{
			K: procexp.Close,
			Close: &procexp.CloseSpec{
				CommChnlPH: closerProcDec.LiabVar.ChnlPH,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	err = s.ProcExecAPI.Take(procstep.StepSpec{
		ImplRef: waiterProcExec,
		ProcExp: procexp.ExpSpec{
			K: procexp.Wait,
			Wait: &procexp.WaitSpec{
				CommChnlPH: waiterProcDec.AssetVars[0].ChnlPH,
				ContES: procexp.ExpSpec{
					K: procexp.Close,
					Close: &procexp.CloseSpec{
						CommChnlPH: waiterProcDec.LiabVar.ChnlPH,
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
		XactExp: xactexp.ExpSpec{
			K: xactexp.Up,
			Up: &xactexp.ShiftSpec{
				ContExp: xactexp.ExpSpec{
					K: xactexp.With,
					With: &xactexp.LaborSpec{
						ProcQNs: []string{senderProcQN, receiverProcQN, messageProcQN},
						ContExp: xactexp.ExpSpec{
							K: xactexp.Down,
							Down: &xactexp.ShiftSpec{
								ContExp: xactexp.ExpSpec{
									K:    xactexp.Link,
									Link: &xactexp.LinkSpec{XactQN: withXactQN},
								},
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
	myPoolQN := "my-pool-qn"
	poolProviderPH := "pool-provider-ph"
	poolClientPH := "pool-client-ph"
	_, err = s.PoolDecAPI.Create(pooldec.DecSpec{
		DescQN: myPoolQN,
		LiabVar: descvar.VarSpec{
			ChnlPH: poolProviderPH,
			DescQN: withXactQN,
		},
		AssetVars: []descvar.VarSpec{{
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
		LiabVar: descvar.VarSpec{
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
		LiabVar: descvar.VarSpec{
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
		LiabVar: descvar.VarSpec{
			ChnlPH: "sender-provider-ph",
			DescQN: oneTypeQN,
		},
		AssetVars: []descvar.VarSpec{
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
		PoolExp: poolexp.ExpSpec{
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
		PoolExp: poolexp.ExpSpec{
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
		PoolExp: poolexp.ExpSpec{
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
		PoolExp: poolexp.ExpSpec{
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
		PoolExp: poolexp.ExpSpec{
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
		PoolExp: poolexp.ExpSpec{
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
		ImplRef: receiverProcExec,
		ProcExp: procexp.ExpSpec{
			K: procexp.Recv,
			Recv: &procexp.RecvSpec{
				CommChnlPH: receiverProcDec.LiabVar.ChnlPH,
				BindChnlPH: "receiver-message-ph",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	err = s.ProcExecAPI.Take(procstep.StepSpec{
		ImplRef: senderProcExec,
		ProcExp: procexp.ExpSpec{
			K: procexp.Send,
			Send: &procexp.SendSpec{
				CommChnlPH: senderProcDec.AssetVars[0].ChnlPH,
				ValChnlPH:  senderProcDec.AssetVars[1].ChnlPH,
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
		LiabVar: descvar.VarSpec{
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
		DescQN:    deciderProcQN,
		AssetVars: []descvar.VarSpec{followerProcDec.LiabVar},
		LiabVar: descvar.VarSpec{
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
		PoolExp: poolexp.ExpSpec{
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
		PoolExp: poolexp.ExpSpec{
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
		PoolExp: poolexp.ExpSpec{
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
		PoolExp: poolexp.ExpSpec{
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
		ImplRef: followerProcExec,
		ProcExp: procexp.ExpSpec{
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
		ImplRef: deciderProcExec,
		ProcExp: procexp.ExpSpec{
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
		LiabVar: descvar.VarSpec{
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
		PoolExp: poolexp.ExpSpec{
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
		PoolExp: poolexp.ExpSpec{
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
		LiabVar: descvar.VarSpec{
			ChnlPH: "caller-provider-ph",
			DescQN: oneType.DefSpec.TypeQN,
		},
		AssetVars: []descvar.VarSpec{
			{ChnlPH: "caller-injectee-ph", DescQN: oneType.DefSpec.TypeQN},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	err = s.PoolExecAPI.Take(poolstep.StepSpec{
		ImplRef: myPoolExec,
		PoolExp: poolexp.ExpSpec{
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
		PoolExp: poolexp.ExpSpec{
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
		LiabVar: descvar.VarSpec{
			ChnlPH: "callee-provider-ph",
			DescQN: oneType.DefSpec.TypeQN,
		},
		AssetVars: []descvar.VarSpec{
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
		ImplRef: callerProcExec,
		ProcExp: procexp.ExpSpec{
			K: procexp.Call,
			Call: &procexp.CallSpec{
				BindChnlPH: callerCalleePH,
				ProcDescQN: calleeProcQN,
				ValChnlPHs: []string{callerProcDec.AssetVars[0].ChnlPH},
				ContES: procexp.ExpSpec{
					K: procexp.Wait,
					Wait: &procexp.WaitSpec{
						CommChnlPH: callerCalleePH,
						ContES: procexp.ExpSpec{
							K: procexp.Close,
							Close: &procexp.CloseSpec{
								CommChnlPH: callerProcDec.LiabVar.ChnlPH,
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
		LiabVar: descvar.VarSpec{
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
		LiabVar: descvar.VarSpec{
			ChnlPH: "forwarder-provider-ph",
			DescQN: oneType.DefSpec.TypeQN,
		},
		AssetVars: []descvar.VarSpec{
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
		LiabVar: descvar.VarSpec{
			ChnlPH: "waiter-provider-ph",
			DescQN: oneType.DefSpec.TypeQN,
		},
		AssetVars: []descvar.VarSpec{
			{ChnlPH: "waiter-forwarder-ph", DescQN: oneType.DefSpec.TypeQN},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	err = s.PoolExecAPI.Take(poolstep.StepSpec{
		ImplRef: myPoolExec,
		PoolExp: poolexp.ExpSpec{
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
		PoolExp: poolexp.ExpSpec{
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
		PoolExp: poolexp.ExpSpec{
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
		PoolExp: poolexp.ExpSpec{
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
		PoolExp: poolexp.ExpSpec{
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
		PoolExp: poolexp.ExpSpec{
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
		ImplRef: closerProcExec,
		ProcExp: procexp.ExpSpec{
			K: procexp.Close,
			Close: &procexp.CloseSpec{
				CommChnlPH: closerProcDec.LiabVar.ChnlPH,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// when
	err = s.ProcExecAPI.Take(procstep.StepSpec{
		ImplRef: forwarderProcExec,
		ProcExp: procexp.ExpSpec{
			K: procexp.Fwd,
			Fwd: &procexp.FwdSpec{
				CommChnlPH: forwarderProcDec.LiabVar.ChnlPH,
				ContChnlPH: forwarderProcDec.AssetVars[0].ChnlPH,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	err = s.ProcExecAPI.Take(procstep.StepSpec{
		ImplRef: waiterProcExec,
		ProcExp: procexp.ExpSpec{
			K: procexp.Wait,
			Wait: &procexp.WaitSpec{
				CommChnlPH: waiterProcDec.AssetVars[0].ChnlPH,
				ContES: procexp.ExpSpec{
					K: procexp.Close,
					Close: &procexp.CloseSpec{
						CommChnlPH: waiterProcDec.LiabVar.ChnlPH,
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
