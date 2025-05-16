package pool_test

import (
	"database/sql"
	"fmt"
	"os"
	"slices"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	"smecalculus/rolevod/lib/sym"

	pooldec "smecalculus/rolevod/app/pool/dec"
	poolexec "smecalculus/rolevod/app/pool/exec"
	procdec "smecalculus/rolevod/app/proc/dec"
	procdef "smecalculus/rolevod/app/proc/def"
	procexec "smecalculus/rolevod/app/proc/exec"
	typedef "smecalculus/rolevod/app/type/def"
)

var (
	poolDecAPI  = pooldec.NewAPI()
	poolExecAPI = poolexec.NewAPI()
	procDecAPI  = procdec.NewAPI()
	procExecAPI = procexec.NewAPI()
	typeDefAPI  = typedef.NewAPI()
	tc          *testCase
)

func TestMain(m *testing.M) {
	ts := testSuite{}
	tc = ts.Setup()
	code := m.Run()
	ts.Teardown()
	os.Exit(code)
}

type testSuite struct {
	db *sql.DB
}

func (ts *testSuite) Setup() *testCase {
	db, err := sql.Open("pgx", "postgres://rolevod:rolevod@localhost:5432/rolevod")
	if err != nil {
		panic(err)
	}
	ts.db = db
	return &testCase{db}
}

func (ts *testSuite) Teardown() {
	err := ts.db.Close()
	if err != nil {
		panic(err)
	}
}

type testCase struct {
	db *sql.DB
}

func (tc *testCase) Setup(t *testing.T) {
	tables := []string{
		"aliases",
		"pool_roots", "pool_liabs", "proc_bnds", "proc_steps",
		"sig_roots", "sig_pes", "sig_ces",
		"role_roots", "role_states",
		"states"}
	for _, table := range tables {
		_, err := tc.db.Exec(fmt.Sprintf("truncate table %v", table))
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestCreation(t *testing.T) {

	t.Run("CreateRetreive", func(t *testing.T) {
		// given
		poolSpec1 := poolexec.PoolSpec{PoolQN: "ts1"}
		poolRef1, err := poolExecAPI.Create(poolSpec1)
		if err != nil {
			t.Fatal(err)
		}
		// and
		poolSpec2 := poolexec.PoolSpec{PoolQN: "ts2", SupID: poolRef1.PoolID}
		poolRef2, err := poolExecAPI.Create(poolSpec2)
		if err != nil {
			t.Fatal(err)
		}
		// when
		poolSnap1, err := poolExecAPI.Retrieve(poolRef1.PoolID)
		if err != nil {
			t.Fatal(err)
		}
		// then
		if !slices.Contains(poolSnap1.Subs, poolRef2) {
			t.Errorf("unexpected subs in %q; want: %+v, got: %+v",
				poolSpec1.PoolQN, poolRef2, poolSnap1.Subs)
		}
	})
}

func TestTaking(t *testing.T) {

	t.Run("WaitClose", func(t *testing.T) {
		tc.Setup(t)
		// given
		oneTypeSN := sym.New("one-type-sn")
		oneTypeTS := typedef.OneSpec{}
		_, err := typeDefAPI.Create(
			typedef.TypeSpec{
				TypeSN: oneTypeSN,
				TypeTS: oneTypeTS,
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		fundTypeSN := sym.New("fund-type-sn")
		closerProcSN := sym.New("closer-proc-sn")
		waiterProcSN := sym.New("waiter-proc-sn")
		_, err = typeDefAPI.Create(
			typedef.TypeSpec{
				TypeSN: fundTypeSN,
				TypeTS: typedef.UpSpec{
					Z: typedef.WithSpec{
						Zs: map[sym.ADT]typedef.TermSpec{
							closerProcSN: typedef.TensorSpec{
								Y: oneTypeTS,
								Z: typedef.DownSpec{
									Z: typedef.LinkSpec{TypeQN: fundTypeSN},
								},
							},
							waiterProcSN: typedef.TensorSpec{
								Y: oneTypeTS,
								Z: typedef.DownSpec{
									Z: typedef.LinkSpec{TypeQN: fundTypeSN},
								},
							},
						},
					},
				},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		mainPoolSN := sym.New("main-pool-sn")
		mainPoolPH := sym.New("main-pool-ph")
		_, err = poolDecAPI.Create(
			pooldec.PoolSpec{
				X: pooldec.ChnlSpec{
					CommPH: mainPoolPH,
					TypeQN: fundTypeSN,
				},
				PoolSN: mainPoolSN,
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		mainPoolRef, err := poolExecAPI.Create(
			poolexec.PoolSpec{
				PoolQN: mainPoolSN,
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		closerSigSpec := procdec.ProcSpec{
			X: procdec.ChnlSpec{
				CommPH: sym.New("x"),
				TypeQN: oneTypeSN,
			},
			ProcSN: closerProcSN,
		}
		_, err = procDecAPI.Create(closerSigSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		waiterSigSpec := procdec.ProcSpec{
			X:      procdec.ChnlSpec{CommPH: sym.New("x"), TypeQN: oneTypeSN},
			ProcSN: waiterProcSN,
			Ys: []procdec.ChnlSpec{
				{CommPH: sym.New("y"), TypeQN: oneTypeSN},
			},
		}
		_, err = procDecAPI.Create(waiterSigSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		closerProcPH := sym.New("closer-proc-ph")
		closerProcRef, err := procExecAPI.Create(
			procexec.ProgSpec{
				PoolID: mainPoolRef.PoolID,
				ExecID: mainPoolRef.ProcID,
				ProcTS: procdef.AcqureSpec{
					CommPH: mainPoolPH,
					ContTS: procdef.CallSpec2{
						CommPH: mainPoolPH,
						ProcSN: closerProcSN,
						ContTS: procdef.RecvSpec{
							CommPH: mainPoolPH,
							BindPH: closerProcPH,
							ContTS: procdef.ReleaseSpec{X: mainPoolPH},
						},
					},
				},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		waiterProcPH := sym.New("waiter-proc-ph")
		waiterProcRef, err := procExecAPI.Create(
			procexec.ProgSpec{
				PoolID: mainPoolRef.PoolID,
				ExecID: mainPoolRef.ProcID,
				ProcTS: procdef.AcqureSpec{
					CommPH: mainPoolPH,
					ContTS: procdef.CallSpec2{
						CommPH: mainPoolPH,
						ProcSN: waiterProcSN,
						ValPHs: []sym.ADT{closerProcPH},
						ContTS: procdef.RecvSpec{
							CommPH: mainPoolPH,
							BindPH: waiterProcPH,
							ContTS: procdef.ReleaseSpec{X: mainPoolPH},
						},
					},
				},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// when
		err = poolExecAPI.Take(
			poolexec.StepSpec{
				PoolID: mainPoolRef.PoolID,
				ProcID: closerProcRef.ExecID,
				Term: procdef.CloseSpec{
					X: closerSigSpec.X.CommPH,
				},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		waitStepSpec := poolexec.StepSpec{
			PoolID: mainPoolRef.PoolID,
			ProcID: waiterProcRef.ExecID,
			Term: procdef.WaitSpec{
				X: waiterSigSpec.Ys[0].CommPH,
				Cont: procdef.CloseSpec{
					X: sym.Blank,
				},
			},
		}
		// and
		err = poolExecAPI.Take(waitStepSpec)
		if err != nil {
			t.Fatal(err)
		}
		// then
		// TODO добавить проверку
	})

	t.Run("RecvSend", func(t *testing.T) {
		tc.Setup(t)
		// given
		lolliRoleSpec := typedef.TypeSpec{
			TypeSN: "lolli-role",
			TypeTS: typedef.LolliSpec{
				Y: typedef.OneSpec{},
				Z: typedef.OneSpec{},
			},
		}
		lolliRole, err := typeDefAPI.Create(lolliRoleSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		oneRoleSpec := typedef.TypeSpec{
			TypeSN: "one-role",
			TypeTS: typedef.OneSpec{},
		}
		oneRole, err := typeDefAPI.Create(oneRoleSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		lolliSigSpec := procdec.ProcSpec{
			ProcSN: "sig-1",
			X: procdec.ChnlSpec{
				CommPH: "chnl-1",
				TypeQN: lolliRole.TypeQN,
			},
		}
		_, err = procDecAPI.Create(lolliSigSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		oneSigSpec1 := procdec.ProcSpec{
			ProcSN: "sig-2",
			X: procdec.ChnlSpec{
				CommPH: "chnl-2",
				TypeQN: oneRole.TypeQN,
			},
		}
		oneSig1, err := procDecAPI.Create(oneSigSpec1)
		if err != nil {
			t.Fatal(err)
		}
		// and
		oneSigSpec2 := procdec.ProcSpec{
			ProcSN: "sig-3",
			Ys:     []procdec.ChnlSpec{lolliSigSpec.X, oneSig1.X},
			X: procdec.ChnlSpec{
				CommPH: "chnl-3",
				TypeQN: oneRole.TypeQN,
			},
		}
		_, err = procDecAPI.Create(oneSigSpec2)
		if err != nil {
			t.Fatal(err)
		}
		// and
		poolImpl, err := poolExecAPI.Create(
			poolexec.PoolSpec{
				PoolQN: "pool-1",
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		receiverChnlPH := sym.New("receiver")
		receiver, err := procExecAPI.Create(
			procexec.ProgSpec{
				PoolID: poolImpl.PoolID,
				ExecID: poolImpl.ProcID,
				ProcTS: procdef.CallSpec2{
					CommPH: receiverChnlPH,
					ProcSN: "tbd",
				},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		messageChnlPH := sym.New("message")
		_, err = procExecAPI.Create(
			procexec.ProgSpec{
				PoolID: poolImpl.PoolID,
				ExecID: poolImpl.ProcID,
				ProcTS: procdef.CallSpec2{
					CommPH: messageChnlPH,
					ProcSN: "tbd",
				},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		senderChnlPH := sym.New("sender")
		sender, err := procExecAPI.Create(
			procexec.ProgSpec{
				PoolID: poolImpl.PoolID,
				ExecID: poolImpl.ProcID,
				ProcTS: procdef.CallSpec2{
					CommPH: senderChnlPH,
					ProcSN: "tbd",
					ValPHs: []sym.ADT{receiverChnlPH, senderChnlPH},
				},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		recvSpec := poolexec.StepSpec{
			PoolID: poolImpl.PoolID,
			ProcID: receiver.ExecID,
			Term: procdef.RecvSpec{
				CommPH: receiverChnlPH,
				BindPH: messageChnlPH,
				ContTS: procdef.WaitSpec{
					X: messageChnlPH,
					Cont: procdef.CloseSpec{
						X: receiverChnlPH,
					},
				},
			},
		}
		// when
		err = poolExecAPI.Take(recvSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		sendSpec := poolexec.StepSpec{
			PoolID: poolImpl.PoolID,
			ProcID: sender.ExecID,
			Term: procdef.SendSpec{
				CommPH: receiverChnlPH,
				ValPH:  messageChnlPH,
			},
		}
		// and
		err = poolExecAPI.Take(sendSpec)
		if err != nil {
			t.Fatal(err)
		}
		// then
		// TODO добавить проверку
	})

	t.Run("CaseLab", func(t *testing.T) {
		tc.Setup(t)
		// given
		label := sym.ADT("label-1")
		// and
		withRoleSpec := typedef.TypeSpec{
			TypeSN: "with-role",
			TypeTS: typedef.WithSpec{
				Zs: map[sym.ADT]typedef.TermSpec{
					label: typedef.OneSpec{},
				},
			},
		}
		withRole, err := typeDefAPI.Create(withRoleSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		oneRoleSpec := typedef.TypeSpec{
			TypeSN: "one-role",
			TypeTS: typedef.OneSpec{},
		}
		oneRole, err := typeDefAPI.Create(oneRoleSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		withSigSpec := procdec.ProcSpec{
			ProcSN: "sig-1",
			X: procdec.ChnlSpec{
				CommPH: "chnl-1",
				TypeQN: withRole.TypeQN,
			},
		}
		withSig, err := procDecAPI.Create(withSigSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		oneSigSpec := procdec.ProcSpec{
			ProcSN: "sig-2",
			Ys:     []procdec.ChnlSpec{withSig.X},
			X: procdec.ChnlSpec{
				CommPH: "chnl-2",
				TypeQN: oneRole.TypeQN,
			},
		}
		_, err = procDecAPI.Create(oneSigSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		poolSpec := poolexec.PoolSpec{
			PoolQN: "pool-1",
		}
		poolImpl, err := poolExecAPI.Create(poolSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		followerPH := sym.New("follower")
		follower, err := procExecAPI.Create(
			procexec.ProgSpec{
				PoolID: poolImpl.PoolID,
				ExecID: poolImpl.ProcID,
				ProcTS: procdef.CallSpec2{
					CommPH: followerPH,
					ProcSN: "tbd",
				},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		deciderPH := sym.New("decider")
		decider, err := procExecAPI.Create(
			procexec.ProgSpec{
				PoolID: poolImpl.PoolID,
				ExecID: poolImpl.ProcID,
				ProcTS: procdef.CallSpec2{
					CommPH: deciderPH,
					ProcSN: "tbd",
					ValPHs: []sym.ADT{followerPH},
				},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		caseSpec := poolexec.StepSpec{
			PoolID: poolImpl.PoolID,
			ProcID: follower.ExecID,
			Term: procdef.CaseSpec{
				X: followerPH,
				Conts: map[sym.ADT]procdef.TermSpec{
					label: procdef.CloseSpec{
						X: followerPH,
					},
				},
			},
		}
		// when
		err = poolExecAPI.Take(caseSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		labSpec := poolexec.StepSpec{
			PoolID: poolImpl.PoolID,
			ProcID: decider.ExecID,
			Term: procdef.LabSpec{
				X:     followerPH,
				Label: label,
			},
		}
		// and
		err = poolExecAPI.Take(labSpec)
		if err != nil {
			t.Fatal(err)
		}
		// then
		// TODO добавить проверку
	})

	t.Run("Spawn", func(t *testing.T) {
		tc.Setup(t)
		// given
		oneRole, err := typeDefAPI.Create(
			typedef.TypeSpec{
				TypeSN: "one-role",
				TypeTS: typedef.OneSpec{},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		oneSig1, err := procDecAPI.Create(
			procdec.ProcSpec{
				ProcSN: "sig-1",
				X: procdec.ChnlSpec{
					CommPH: "chnl-1",
					TypeQN: oneRole.TypeQN,
				},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		_, err = procDecAPI.Create(
			procdec.ProcSpec{
				ProcSN: "sig-2",
				Ys:     []procdec.ChnlSpec{oneSig1.X},
				X: procdec.ChnlSpec{
					CommPH: "chnl-2",
					TypeQN: oneRole.TypeQN,
				},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		oneSig3, err := procDecAPI.Create(
			procdec.ProcSpec{
				ProcSN: "sig-3",
				Ys:     []procdec.ChnlSpec{oneSig1.X},
				X: procdec.ChnlSpec{
					CommPH: "chnl-3",
					TypeQN: oneRole.TypeQN,
				},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		poolImpl, err := poolExecAPI.Create(
			poolexec.PoolSpec{
				PoolQN: "pool-1",
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		injecteePH := sym.New("injectee")
		_, err = procExecAPI.Create(
			procexec.ProgSpec{
				PoolID: poolImpl.PoolID,
				ExecID: poolImpl.ProcID,
				ProcTS: procdef.CallSpec2{
					CommPH: injecteePH,
					ProcSN: "tbd",
				},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		spawnerPH := sym.New("spawner")
		spawner, err := procExecAPI.Create(
			procexec.ProgSpec{
				PoolID: poolImpl.PoolID,
				ExecID: poolImpl.ProcID,
				ProcTS: procdef.CallSpec2{
					CommPH: spawnerPH,
					ProcSN: "tbd",
					ValPHs: []sym.ADT{injecteePH},
				},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		x := sym.New("x")
		spawnSpec := poolexec.StepSpec{
			PoolID: poolImpl.PoolID,
			ProcID: spawner.ExecID,
			Term: procdef.SpawnSpec{
				SigID: oneSig3.DecID,
				Ys:    []sym.ADT{injecteePH},
				X:     x,
				Cont: procdef.WaitSpec{
					X: x,
					Cont: procdef.CloseSpec{
						X: spawnerPH,
					},
				},
			},
		}
		// when
		err = poolExecAPI.Take(spawnSpec)
		if err != nil {
			t.Fatal(err)
		}
		// then
		// TODO добавить проверку
	})

	t.Run("Fwd", func(t *testing.T) {
		tc.Setup(t)
		// given
		oneRole, err := typeDefAPI.Create(
			typedef.TypeSpec{
				TypeSN: "one-role",
				TypeTS: typedef.OneSpec{},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		oneSig1, err := procDecAPI.Create(
			procdec.ProcSpec{
				ProcSN: "sig-1",
				X: procdec.ChnlSpec{
					CommPH: "chnl-1",
					TypeQN: oneRole.TypeQN,
				},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		_, err = procDecAPI.Create(
			procdec.ProcSpec{
				ProcSN: "sig-2",
				Ys:     []procdec.ChnlSpec{oneSig1.X},
				X: procdec.ChnlSpec{
					CommPH: "chnl-2",
					TypeQN: oneRole.TypeQN,
				},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		_, err = procDecAPI.Create(
			procdec.ProcSpec{
				ProcSN: "sig-3",
				Ys:     []procdec.ChnlSpec{oneSig1.X},
				X: procdec.ChnlSpec{
					CommPH: "chnl-3",
					TypeQN: oneRole.TypeQN,
				},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		poolImpl, err := poolExecAPI.Create(
			poolexec.PoolSpec{
				PoolQN: "pool-1",
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		closerChnlPH := sym.New("closer")
		closerSpec := procexec.ProgSpec{
			PoolID: poolImpl.PoolID,
			ExecID: poolImpl.ProcID,
			ProcTS: procdef.CallSpec2{
				CommPH: closerChnlPH,
				ProcSN: "tbd",
			},
		}
		closer, err := procExecAPI.Create(closerSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		forwarderChnlPH := sym.New("forwarder")
		forwarderSpec := procexec.ProgSpec{
			PoolID: poolImpl.PoolID,
			ExecID: poolImpl.ProcID,
			ProcTS: procdef.CallSpec2{
				CommPH: forwarderChnlPH,
				ProcSN: "tbd",
				ValPHs: []sym.ADT{closerChnlPH},
			},
		}
		forwarder, err := procExecAPI.Create(forwarderSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		waiterChnlPH := sym.New("waiter")
		waiterSpec := procexec.ProgSpec{
			PoolID: poolImpl.PoolID,
			ExecID: poolImpl.ProcID,
			ProcTS: procdef.CallSpec2{
				CommPH: waiterChnlPH,
				ProcSN: "tbd",
				ValPHs: []sym.ADT{forwarderChnlPH},
			},
		}
		waiter, err := procExecAPI.Create(waiterSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		closeSpec := poolexec.StepSpec{
			PoolID: poolImpl.PoolID,
			ProcID: closer.ExecID,
			Term: procdef.CloseSpec{
				X: closerChnlPH,
			},
		}
		err = poolExecAPI.Take(closeSpec)
		if err != nil {
			t.Fatal(err)
		}
		// when
		fwdSpec := poolexec.StepSpec{
			PoolID: poolImpl.PoolID,
			ProcID: forwarder.ExecID,
			Term: procdef.FwdSpec{
				X: forwarderChnlPH,
				Y: closerChnlPH,
			},
		}
		err = poolExecAPI.Take(fwdSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		waitSpec := poolexec.StepSpec{
			PoolID: poolImpl.PoolID,
			ProcID: waiter.ExecID,
			Term: procdef.WaitSpec{
				X: forwarderChnlPH,
				Cont: procdef.CloseSpec{
					X: waiterChnlPH,
				},
			},
		}
		err = poolExecAPI.Take(waitSpec)
		if err != nil {
			t.Fatal(err)
		}
		// then
		// TODO добавить проверку
	})
}
