package pool

import (
	"database/sql"
	"fmt"
	"slices"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/fx"

	"orglang/go-runtime/lib/e2e"

	"github.com/orglang/go-sdk/lib/rc"

	"github.com/orglang/go-sdk/adt/pooldec"
	"github.com/orglang/go-sdk/adt/poolexec"
	"github.com/orglang/go-sdk/adt/procbind"
	"github.com/orglang/go-sdk/adt/procdec"
	"github.com/orglang/go-sdk/adt/procexec"
	"github.com/orglang/go-sdk/adt/procexp"
	"github.com/orglang/go-sdk/adt/typedef"
	"github.com/orglang/go-sdk/adt/typeexp"
)

func TestPool(t *testing.T) {
	s := suite{}
	s.beforeAll(t)
	t.Run("CreateRetreive", s.createRetreive)
	t.Run("WaitClose", s.waitClose)
	t.Run("RecvSend", s.recvSend)
	t.Run("CaseLab", s.caseLab)
	t.Run("SpawnCall", s.spawnCall)
	t.Run("Fwd", s.fwd)

}

type suite struct {
	poolDecAPI  e2e.PoolDecAPI
	poolExecAPI e2e.PoolExecAPI
	procDecAPI  e2e.ProcDecAPI
	procExecAPI e2e.ProcExecAPI
	typeDefAPI  e2e.TypeDefAPI
	db          *sql.DB
}

func (s *suite) beforeAll(t *testing.T) {
	db, err := sql.Open("pgx", "postgres://orglang:orglang@localhost:5432/orglang")
	if err != nil {
		panic(err)
	}
	t.Cleanup(func() { db.Close() })
	s.db = db
	app := fx.New(rc.Module, e2e.Module,
		fx.Populate(
			s.poolDecAPI,
			s.poolExecAPI,
			s.procDecAPI,
			s.procExecAPI,
			s.typeDefAPI,
		))
	t.Cleanup(func() { app.Stop(t.Context()) })
}

func (s *suite) beforeEach(t *testing.T) {
	tables := []string{
		"aliases",
		"pool_roots", "pool_liabs", "proc_bnds", "proc_steps",
		"sig_roots", "sig_pes", "sig_ces",
		"type_def_roots", "type_term_states",
		"type_term_states"}
	for _, table := range tables {
		_, err := s.db.Exec(fmt.Sprintf("truncate table %v", table))
		if err != nil {
			t.Fatal(err)
		}
	}
}

func (s *suite) createRetreive(t *testing.T) {
	// given
	poolSpec1 := poolexec.ExecSpec{PoolQN: "pool-1"}
	poolRef1, err := s.poolExecAPI.Create(poolSpec1)
	if err != nil {
		t.Fatal(err)
	}
	// and
	poolSpec2 := poolexec.ExecSpec{PoolQN: "pool-2", SupID: poolRef1.ID}
	poolRef2, err := s.poolExecAPI.Create(poolSpec2)
	if err != nil {
		t.Fatal(err)
	}
	// when
	poolSnap1, err := s.poolExecAPI.Retrieve(poolRef1.ID)
	if err != nil {
		t.Fatal(err)
	}
	// then
	if !slices.Contains(poolSnap1.Subs, poolRef2) {
		t.Errorf("unexpected subs in %q; want: %+v, got: %+v",
			poolSpec1.PoolQN, poolRef2, poolSnap1.Subs)
	}
}

func (s *suite) waitClose(t *testing.T) {
	s.beforeEach(t)
	// given
	mainTypeQN := "main-type-qn"
	closerProcQN := "closer-proc-qn"
	waiterProcQN := "waiter-proc-qn"
	_, err := s.typeDefAPI.Create(typedef.DefSpec{
		TypeQN: mainTypeQN,
		TypeES: typeexp.ExpSpec{
			K: typeexp.UpExp,
			Up: &typeexp.FooSpec{
				ContES: typeexp.ExpSpec{
					K: typeexp.XactExp,
					Xact: &typeexp.XactSpec{
						ContESs: map[string]typeexp.ExpSpec{
							closerProcQN: {
								K: typeexp.DownExp,
								Down: &typeexp.FooSpec{
									ContES: typeexp.ExpSpec{
										K:    typeexp.LinkExp,
										Link: &typeexp.LinkSpec{TypeQN: mainTypeQN},
									},
								},
							},
							waiterProcQN: {
								K: typeexp.DownExp,
								Down: &typeexp.FooSpec{
									ContES: typeexp.ExpSpec{
										K:    typeexp.LinkExp,
										Link: &typeexp.LinkSpec{TypeQN: mainTypeQN},
									},
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
	mainPoolQN := "main-pool-qn"
	mainProvisionPH := "main-provision-ph"
	mainReceptionPH := "main-reception-ph"
	_, err = s.poolDecAPI.Create(pooldec.DecSpec{
		PoolSN: mainPoolQN,
		InsiderProvisionBC: procbind.BindSpec{
			ChnlPH: mainProvisionPH,
			TypeQN: mainTypeQN,
		},
		InsiderReceptionBC: procbind.BindSpec{
			ChnlPH: mainReceptionPH,
			TypeQN: mainTypeQN,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	mainExecRef, err := s.poolExecAPI.Create(poolexec.ExecSpec{
		PoolQN: mainPoolQN,
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	oneTypeQN := "one-type-qn"
	_, err = s.typeDefAPI.Create(typedef.DefSpec{
		TypeQN: oneTypeQN,
		TypeES: typeexp.ExpSpec{K: typeexp.OneExp},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	closerDecSpec := procdec.DecSpec{
		ProcQN: closerProcQN,
		X: procbind.BindSpec{
			ChnlPH: "closer-provision-ph",
			TypeQN: oneTypeQN,
		},
	}
	_, err = s.procDecAPI.Create(closerDecSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	closerProcPH := "closer-proc-ph"
	err = s.procExecAPI.Run(procexec.ExecSpec{
		PoolID: mainExecRef.ID,
		ExecID: mainExecRef.ID,
		ProcES: procexp.ExpSpec{
			Acquire: &procexp.AcqureSpec{
				CommPH: mainReceptionPH,
				ContES: procexp.ExpSpec{
					Call: &procexp.CallSpec{
						CommPH: mainReceptionPH,
						BindPH: closerProcPH,
						ProcQN: closerProcQN,
						ContES: procexp.ExpSpec{
							Release: &procexp.ReleaseSpec{
								CommPH: mainReceptionPH,
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
	err = s.procExecAPI.Run(procexec.ExecSpec{
		PoolID: mainExecRef.ID,
		ExecID: mainExecRef.ID,
		ProcES: procexp.ExpSpec{
			Accept: &procexp.AcceptSpec{
				CommPH: mainProvisionPH,
				ContES: procexp.ExpSpec{
					Spawn: &procexp.SpawnSpec{
						CommPH: mainProvisionPH,
						ProcQN: closerProcQN,
						ContES: procexp.ExpSpec{
							Detach: &procexp.DetachSpec{
								CommPH: mainProvisionPH,
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
	closerExecRef, err := s.poolExecAPI.Poll(poolexec.PollSpec{
		ExecID: mainExecRef.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	waiterDecSpec := procdec.DecSpec{
		ProcQN: waiterProcQN,
		X:      procbind.BindSpec{ChnlPH: "waiter-provision-ph", TypeQN: oneTypeQN},
		Ys: []procbind.BindSpec{
			{ChnlPH: "closer-reception-ph", TypeQN: oneTypeQN},
		},
	}
	_, err = s.procDecAPI.Create(waiterDecSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	err = s.procExecAPI.Run(procexec.ExecSpec{
		PoolID: mainExecRef.ID,
		ExecID: mainExecRef.ID,
		ProcES: procexp.ExpSpec{
			Acquire: &procexp.AcqureSpec{
				CommPH: mainProvisionPH,
				ContES: procexp.ExpSpec{
					Call: &procexp.CallSpec{
						BindPH: mainProvisionPH,
						CommPH: "waiter-proc-ph",
						ProcQN: waiterProcQN,
						ValPHs: []string{closerProcPH},
						ContES: procexp.ExpSpec{
							Release: &procexp.ReleaseSpec{
								CommPH: mainProvisionPH,
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
	err = s.procExecAPI.Run(procexec.ExecSpec{
		PoolID: mainExecRef.ID,
		ExecID: mainExecRef.ID,
		ProcES: procexp.ExpSpec{
			Accept: &procexp.AcceptSpec{
				CommPH: mainProvisionPH,
				ContES: procexp.ExpSpec{
					Spawn: &procexp.SpawnSpec{
						CommPH: mainProvisionPH,
						ProcQN: waiterProcQN,
						ContES: procexp.ExpSpec{
							Detach: &procexp.DetachSpec{
								CommPH: mainProvisionPH,
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
	waiterExecRef, err := s.poolExecAPI.Poll(poolexec.PollSpec{
		ExecID: mainExecRef.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	// when
	err = s.procExecAPI.Run(procexec.ExecSpec{
		PoolID: mainExecRef.ID,
		ExecID: closerExecRef.ID,
		ProcES: procexp.ExpSpec{
			Close: &procexp.CloseSpec{
				CommPH: closerDecSpec.X.ChnlPH,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	err = s.procExecAPI.Run(procexec.ExecSpec{
		PoolID: mainExecRef.ID,
		ExecID: waiterExecRef.ID,
		ProcES: procexp.ExpSpec{
			Wait: &procexp.WaitSpec{
				CommPH: waiterDecSpec.Ys[0].ChnlPH,
				ContES: procexp.ExpSpec{
					Close: &procexp.CloseSpec{
						CommPH: waiterDecSpec.X.ChnlPH,
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
	mainTypeQN := "main-type-qn"
	senderProcQN := "sender-proc-qn"
	receiverProcQN := "receiver-proc-qn"
	messageProcQN := "message-proc-qn"
	_, err := s.typeDefAPI.Create(typedef.DefSpec{
		TypeQN: mainTypeQN,
		TypeES: typeexp.ExpSpec{
			K: typeexp.UpExp,
			Up: &typeexp.FooSpec{
				ContES: typeexp.ExpSpec{
					K: typeexp.XactExp,
					Xact: &typeexp.XactSpec{
						ContESs: map[string]typeexp.ExpSpec{
							senderProcQN: typeexp.ExpSpec{
								K: typeexp.DownExp,
								Down: &typeexp.FooSpec{
									ContES: typeexp.ExpSpec{
										K:    typeexp.LinkExp,
										Link: &typeexp.LinkSpec{TypeQN: mainTypeQN},
									},
								},
							},
							receiverProcQN: typeexp.ExpSpec{
								K: typeexp.DownExp,
								Down: &typeexp.FooSpec{
									ContES: typeexp.ExpSpec{
										K:    typeexp.LinkExp,
										Link: &typeexp.LinkSpec{TypeQN: mainTypeQN},
									},
								},
							},
							messageProcQN: typeexp.ExpSpec{
								K: typeexp.DownExp,
								Down: &typeexp.FooSpec{
									ContES: typeexp.ExpSpec{
										K:    typeexp.LinkExp,
										Link: &typeexp.LinkSpec{TypeQN: mainTypeQN},
									},
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
	mainPoolQN := "main-pool-qn"
	mainProvisionPH := "main-provision-ph"
	mainReceptionPH := "main-reception-ph"
	_, err = s.poolDecAPI.Create(pooldec.DecSpec{
		PoolSN: mainPoolQN,
		InsiderProvisionBC: procbind.BindSpec{
			ChnlPH: mainProvisionPH,
			TypeQN: mainPoolQN,
		},
		InsiderReceptionBC: procbind.BindSpec{
			ChnlPH: mainReceptionPH,
			TypeQN: mainPoolQN,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	mainExecRef, err := s.poolExecAPI.Create(poolexec.ExecSpec{
		PoolQN: mainPoolQN,
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	lolliTypeQN := "lolli-type-qn"
	_, err = s.typeDefAPI.Create(typedef.DefSpec{
		TypeQN: lolliTypeQN,
		TypeES: typeexp.ExpSpec{
			K: typeexp.LolliExp,
			Lolli: &typeexp.ProdSpec{
				ValES:  typeexp.ExpSpec{K: typeexp.OneExp},
				ContES: typeexp.ExpSpec{K: typeexp.OneExp},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	oneTypeQN := "one-type-qn"
	_, err = s.typeDefAPI.Create(typedef.DefSpec{
		TypeQN: oneTypeQN,
		TypeES: typeexp.ExpSpec{K: typeexp.OneExp},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	receiverDecSpec := procdec.DecSpec{
		ProcQN: receiverProcQN,
		X: procbind.BindSpec{
			ChnlPH: "receiver-provision-ph",
			TypeQN: lolliTypeQN,
		},
	}
	_, err = s.procDecAPI.Create(receiverDecSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	messageDecSpec := procdec.DecSpec{
		ProcQN: messageProcQN,
		X: procbind.BindSpec{
			ChnlPH: "message-provision-ph",
			TypeQN: oneTypeQN,
		},
	}
	_, err = s.procDecAPI.Create(messageDecSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	senderDecSpec := procdec.DecSpec{
		ProcQN: senderProcQN,
		X: procbind.BindSpec{
			ChnlPH: "sender-provision-ph",
			TypeQN: oneTypeQN,
		},
		Ys: []procbind.BindSpec{
			{ChnlPH: "receiver-reception-ph", TypeQN: lolliTypeQN},
			{ChnlPH: "message-reception-ph", TypeQN: oneTypeQN},
		},
	}
	_, err = s.procDecAPI.Create(senderDecSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	receiverProcPH := "receiver-proc-ph"
	err = s.procExecAPI.Run(procexec.ExecSpec{
		PoolID: mainExecRef.ID,
		ExecID: mainExecRef.ID,
		ProcES: procexp.ExpSpec{
			Call: &procexp.CallSpec{
				BindPH: mainReceptionPH,
				CommPH: receiverProcPH,
				ProcQN: receiverProcQN,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	messageProcPH := "message-proc-ph"
	err = s.procExecAPI.Run(procexec.ExecSpec{
		PoolID: mainExecRef.ID,
		ExecID: mainExecRef.ID,
		ProcES: procexp.ExpSpec{
			Call: &procexp.CallSpec{
				BindPH: mainReceptionPH,
				CommPH: messageProcPH,
				ProcQN: messageProcQN,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	senderProcPH := "sender-proc-ph"
	err = s.procExecAPI.Run(procexec.ExecSpec{
		PoolID: mainExecRef.ID,
		ExecID: mainExecRef.ID,
		ProcES: procexp.ExpSpec{
			Call: &procexp.CallSpec{
				BindPH: mainReceptionPH,
				CommPH: senderProcPH,
				ProcQN: senderProcQN,
				ValPHs: []string{receiverProcPH, messageProcPH},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	receiverExecRef, err := s.poolExecAPI.Poll(poolexec.PollSpec{
		ExecID: mainExecRef.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	// when
	err = s.procExecAPI.Run(procexec.ExecSpec{
		PoolID: mainExecRef.ID,
		ExecID: receiverExecRef.ID,
		ProcES: procexp.ExpSpec{
			Recv: &procexp.RecvSpec{
				BindPH: receiverDecSpec.X.ChnlPH,
				CommPH: "message-reception-ph",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	senderExecRef, err := s.poolExecAPI.Poll(poolexec.PollSpec{
		ExecID: mainExecRef.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	err = s.procExecAPI.Run(procexec.ExecSpec{
		PoolID: mainExecRef.ID,
		ExecID: senderExecRef.ID,
		ProcES: procexp.ExpSpec{
			Send: &procexp.SendSpec{
				CommPH: senderDecSpec.Ys[0].ChnlPH,
				ValPH:  senderDecSpec.Ys[1].ChnlPH,
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
	mainPoolQN := "main-pool-qn"
	mainExecRef, err := s.poolExecAPI.Create(poolexec.ExecSpec{
		PoolQN: mainPoolQN,
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	label := "label-1"
	// and
	withTypeSpec := typedef.DefSpec{
		TypeQN: "with-role",
		TypeES: typeexp.ExpSpec{
			With: &typeexp.SumSpec{
				Choices: []typeexp.ChoiceSpec{
					{LabQN: label, ContES: typeexp.ExpSpec{K: typeexp.OneExp}},
				},
			},
		},
	}
	withTypeDef, err := s.typeDefAPI.Create(withTypeSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	oneTypeSpec := typedef.DefSpec{
		TypeQN: "one-role",
		TypeES: typeexp.ExpSpec{K: typeexp.OneExp},
	}
	oneTypeSnap, err := s.typeDefAPI.Create(oneTypeSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	withDecSpec := procdec.DecSpec{
		ProcQN: "sig-1",
		X: procbind.BindSpec{
			ChnlPH: "chnl-1",
			TypeQN: withTypeDef.Spec.TypeQN,
		},
	}
	withSig, err := s.procDecAPI.Create(withDecSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	oneSigSpec := procdec.DecSpec{
		ProcQN: "sig-2",
		Ys:     []procbind.BindSpec{withSig.Spec.X},
		X: procbind.BindSpec{
			ChnlPH: "chnl-2",
			TypeQN: oneTypeSnap.Spec.TypeQN,
		},
	}
	_, err = s.procDecAPI.Create(oneSigSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	followerPH := "follower"
	err = s.procExecAPI.Run(procexec.ExecSpec{
		PoolID: mainExecRef.ID,
		ExecID: mainExecRef.ID,
		ProcES: procexp.ExpSpec{
			Call: &procexp.CallSpec{
				BindPH: followerPH,
				ProcQN: "tbd",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	deciderPH := "decider"
	err = s.procExecAPI.Run(procexec.ExecSpec{
		PoolID: mainExecRef.ID,
		ExecID: mainExecRef.ID,
		ProcES: procexp.ExpSpec{
			Call: &procexp.CallSpec{
				BindPH: deciderPH,
				ProcQN: "tbd",
				ValPHs: []string{followerPH},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	followerExecRef, err := s.poolExecAPI.Poll(poolexec.PollSpec{
		ExecID: mainExecRef.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	// when
	err = s.procExecAPI.Run(procexec.ExecSpec{
		PoolID: mainExecRef.ID,
		ExecID: followerExecRef.ID,
		ProcES: procexp.ExpSpec{
			Case: &procexp.CaseSpec{
				CommPH: followerPH,
				ContBSs: []procexp.BranchSpec{
					procexp.BranchSpec{
						LabQN: label,
						ContES: procexp.ExpSpec{
							Close: &procexp.CloseSpec{
								CommPH: followerPH,
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
	deciderExecRef, err := s.poolExecAPI.Poll(poolexec.PollSpec{
		ExecID: mainExecRef.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	err = s.procExecAPI.Run(procexec.ExecSpec{
		PoolID: mainExecRef.ID,
		ExecID: deciderExecRef.ID,
		ProcES: procexp.ExpSpec{
			Lab: &procexp.LabSpec{
				CommPH: followerPH,
				LabQN:  label,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// then
	// TODO добавить проверку
}

func (s *suite) spawnCall(t *testing.T) {
	s.beforeEach(t)
	// given
	mainPoolQN := "main-pool-qn"
	mainExecRef, err := s.poolExecAPI.Create(poolexec.ExecSpec{
		PoolQN: mainPoolQN,
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	oneDef, err := s.typeDefAPI.Create(
		typedef.DefSpec{
			TypeQN: "one-type",
			TypeES: typeexp.ExpSpec{K: typeexp.OneExp},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	// and
	oneDec1, err := s.procDecAPI.Create(procdec.DecSpec{
		ProcQN: "one-proc-1",
		X: procbind.BindSpec{
			ChnlPH: "chnl-1",
			TypeQN: oneDef.Spec.TypeQN,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	_, err = s.procDecAPI.Create(procdec.DecSpec{
		ProcQN: "one-proc-2",
		Ys:     []procbind.BindSpec{oneDec1.Spec.X},
		X: procbind.BindSpec{
			ChnlPH: "chnl-2",
			TypeQN: oneDef.Spec.TypeQN,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	oneDec3, err := s.procDecAPI.Create(procdec.DecSpec{
		ProcQN: "one-proc-3",
		Ys:     []procbind.BindSpec{oneDec1.Spec.X},
		X: procbind.BindSpec{
			ChnlPH: "chnl-3",
			TypeQN: oneDef.Spec.TypeQN,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	poolExecRef, err := s.poolExecAPI.Create(poolexec.ExecSpec{
		PoolQN: "pool-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	injecteePH := "injectee"
	err = s.procExecAPI.Run(procexec.ExecSpec{
		PoolID: poolExecRef.ID,
		ExecID: poolExecRef.ID,
		ProcES: procexp.ExpSpec{
			Call: &procexp.CallSpec{
				BindPH: injecteePH,
				ProcQN: "tbd",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	spawnerPH := "spawner"
	err = s.procExecAPI.Run(procexec.ExecSpec{
		PoolID: poolExecRef.ID,
		ExecID: poolExecRef.ID,
		ProcES: procexp.ExpSpec{
			Call: &procexp.CallSpec{
				BindPH: spawnerPH,
				ProcQN: "tbd",
				ValPHs: []string{injecteePH},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	x := "x"
	// and
	spawnerExecRef, err := s.poolExecAPI.Poll(poolexec.PollSpec{
		ExecID: mainExecRef.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	// when
	err = s.procExecAPI.Run(procexec.ExecSpec{
		PoolID: poolExecRef.ID,
		ExecID: spawnerExecRef.ID,
		ProcES: procexp.ExpSpec{
			Spawn: &procexp.SpawnSpec{
				DecID:   oneDec3.Ref.ID,
				BindPHs: []string{injecteePH},
				CommPH:  x,
				ContES: procexp.ExpSpec{
					Wait: &procexp.WaitSpec{
						CommPH: x,
						ContES: procexp.ExpSpec{
							Close: &procexp.CloseSpec{
								CommPH: spawnerPH,
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
	mainPoolQN := "main-pool-qn"
	mainExecRef, err := s.poolExecAPI.Create(poolexec.ExecSpec{
		PoolQN: mainPoolQN,
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	oneDefSnap, err := s.typeDefAPI.Create(typedef.DefSpec{
		TypeQN: "one-role",
		TypeES: typeexp.ExpSpec{K: typeexp.OneExp},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	oneDecSnap, err := s.procDecAPI.Create(procdec.DecSpec{
		ProcQN: "proc-dec-1",
		X: procbind.BindSpec{
			ChnlPH: "chnl-1",
			TypeQN: oneDefSnap.Spec.TypeQN,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	_, err = s.procDecAPI.Create(procdec.DecSpec{
		ProcQN: "proc-dec-2",
		Ys:     []procbind.BindSpec{oneDecSnap.Spec.X},
		X: procbind.BindSpec{
			ChnlPH: "chnl-2",
			TypeQN: oneDefSnap.Spec.TypeQN,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	_, err = s.procDecAPI.Create(procdec.DecSpec{
		ProcQN: "proc-dec-3",
		Ys:     []procbind.BindSpec{oneDecSnap.Spec.X},
		X: procbind.BindSpec{
			ChnlPH: "chnl-3",
			TypeQN: oneDefSnap.Spec.TypeQN,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	closerChnlPH := "closer"
	err = s.procExecAPI.Run(procexec.ExecSpec{
		PoolID: mainExecRef.ID,
		ExecID: mainExecRef.ID,
		ProcES: procexp.ExpSpec{
			Call: &procexp.CallSpec{
				BindPH: closerChnlPH,
				ProcQN: "tbd",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	forwarderChnlPH := "forwarder"
	err = s.procExecAPI.Run(procexec.ExecSpec{
		PoolID: mainExecRef.ID,
		ExecID: mainExecRef.ID,
		ProcES: procexp.ExpSpec{
			Call: &procexp.CallSpec{
				BindPH: forwarderChnlPH,
				ProcQN: "tbd",
				ValPHs: []string{closerChnlPH},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	waiterChnlPH := "waiter"
	err = s.procExecAPI.Run(procexec.ExecSpec{
		PoolID: mainExecRef.ID,
		ExecID: mainExecRef.ID,
		ProcES: procexp.ExpSpec{
			Call: &procexp.CallSpec{
				BindPH: waiterChnlPH,
				ProcQN: "tbd",
				ValPHs: []string{forwarderChnlPH},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	closerExecRef, err := s.poolExecAPI.Poll(poolexec.PollSpec{
		ExecID: mainExecRef.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	err = s.procExecAPI.Run(procexec.ExecSpec{
		PoolID: mainExecRef.ID,
		ExecID: closerExecRef.ID,
		ProcES: procexp.ExpSpec{
			Close: &procexp.CloseSpec{
				CommPH: closerChnlPH,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	forwarderExecRef, err := s.poolExecAPI.Poll(poolexec.PollSpec{
		ExecID: mainExecRef.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	// when
	err = s.procExecAPI.Run(procexec.ExecSpec{
		PoolID: mainExecRef.ID,
		ExecID: forwarderExecRef.ID,
		ProcES: procexp.ExpSpec{
			Fwd: &procexp.FwdSpec{
				CommPH: forwarderChnlPH,
				ContPH: closerChnlPH,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	waiterExecRef, err := s.poolExecAPI.Poll(poolexec.PollSpec{
		ExecID: mainExecRef.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	// and
	err = s.procExecAPI.Run(procexec.ExecSpec{
		PoolID: mainExecRef.ID,
		ExecID: waiterExecRef.ID,
		ProcES: procexp.ExpSpec{
			Wait: &procexp.WaitSpec{
				CommPH: forwarderChnlPH,
				ContES: procexp.ExpSpec{
					Close: &procexp.CloseSpec{
						CommPH: waiterChnlPH,
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
