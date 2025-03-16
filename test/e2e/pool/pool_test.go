package pool_test

import (
	"database/sql"
	"fmt"
	"os"
	"slices"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/ph"

	"smecalculus/rolevod/internal/chnl"
	"smecalculus/rolevod/internal/state"
	"smecalculus/rolevod/internal/step"

	"smecalculus/rolevod/app/pool"
	"smecalculus/rolevod/app/role"
	"smecalculus/rolevod/app/sig"
)

var (
	roleAPI = role.NewAPI()
	sigAPI  = sig.NewAPI()
	poolAPI = pool.NewAPI()
	tc      *testCase
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
		"pool_roots",
		"sig_roots", "sig_pes", "sig_ces",
		"role_roots", "role_states",
		"states", "channels", "steps", "clientships"}
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
		poolSpec1 := pool.Spec{Title: "ts1"}
		poolRoot1, err := poolAPI.Create(poolSpec1)
		if err != nil {
			t.Fatal(err)
		}
		// and
		poolSpec2 := pool.Spec{Title: "ts2", SupID: poolRoot1.PoolID}
		poolRoot2, err := poolAPI.Create(poolSpec2)
		if err != nil {
			t.Fatal(err)
		}
		// when
		poolSnap1, err := poolAPI.Retrieve(poolRoot1.PoolID)
		if err != nil {
			t.Fatal(err)
		}
		// then
		extectedSub := pool.ConvertRootToRef(poolRoot2)
		if !slices.Contains(poolSnap1.Subs, extectedSub) {
			t.Errorf("unexpected subs in %q; want: %+v, got: %+v",
				poolSpec1.Title, extectedSub, poolSnap1.Subs)
		}
	})
}

func TestTaking(t *testing.T) {

	t.Run("WaitClose", func(t *testing.T) {
		tc.Setup(t)
		// given
		oneRoleSpec := role.Spec{
			QN:    "one-role",
			State: state.OneSpec{},
		}
		oneRole, err := roleAPI.Create(oneRoleSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		closerSigSpec := sig.Spec{
			QN: "closer",
			X: chnl.Spec{
				ChnlPH: "closing-1",
				RoleQN: oneRole.QN,
			},
		}
		closerSig, err := sigAPI.Create(closerSigSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		waiterSigSpec := sig.Spec{
			QN: "waiter",
			Ys: []chnl.Spec{closerSig.X},
			X: chnl.Spec{
				ChnlPH: "closing-2",
				RoleQN: oneRole.QN,
			},
		}
		waiterSig, err := sigAPI.Create(waiterSigSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		poolSpec := pool.Spec{
			Title: "pool-1",
		}
		poolImpl, err := poolAPI.Create(poolSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		closerPH := ph.New("closer")
		closerSpec := pool.TranSpec{
			PoolID: poolImpl.PoolID,
			ProcID: poolImpl.ProcID,
			Term: step.SpawnSpec{
				SigID: closerSig.SigID,
				X:     closerPH,
			},
		}
		closerID, err := poolAPI.Spawn(closerSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		waiterPH := ph.New("waiter")
		waiterSpec := pool.TranSpec{
			PoolID: poolImpl.PoolID,
			ProcID: poolImpl.ProcID,
			Term: step.SpawnSpec{
				SigID: waiterSig.SigID,
				Ys:    []ph.ADT{closerPH},
				X:     waiterPH,
			},
		}
		waiterID, err := poolAPI.Spawn(waiterSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		closeSpec := pool.TranSpec{
			PoolID: poolImpl.PoolID,
			ProcID: closerID,
			Term: step.CloseSpec{
				X: closerPH,
			},
		}
		// when
		err = poolAPI.Take(closeSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		waitSpec := pool.TranSpec{
			PoolID: poolImpl.PoolID,
			ProcID: waiterID,
			Term: step.WaitSpec{
				X: closerPH,
				Cont: step.CloseSpec{
					X: waiterPH,
				},
			},
		}
		// and
		err = poolAPI.Take(waitSpec)
		if err != nil {
			t.Fatal(err)
		}
		// then
		// TODO добавить проверку
	})

	t.Run("RecvSend", func(t *testing.T) {
		tc.Setup(t)
		// given
		lolliRoleSpec := role.Spec{
			QN: "lolli-role",
			State: state.LolliSpec{
				Y: state.OneSpec{},
				Z: state.OneSpec{},
			},
		}
		lolliRole, err := roleAPI.Create(lolliRoleSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		oneRoleSpec := role.Spec{
			QN:    "one-role",
			State: state.OneSpec{},
		}
		oneRole, err := roleAPI.Create(oneRoleSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		lolliSigSpec := sig.Spec{
			QN: "sig-1",
			X: chnl.Spec{
				ChnlPH: "chnl-1",
				RoleQN: lolliRole.QN,
			},
		}
		lolliSig, err := sigAPI.Create(lolliSigSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		oneSigSpec1 := sig.Spec{
			QN: "sig-2",
			X: chnl.Spec{
				ChnlPH: "chnl-2",
				RoleQN: oneRole.QN,
			},
		}
		oneSig1, err := sigAPI.Create(oneSigSpec1)
		if err != nil {
			t.Fatal(err)
		}
		// and
		oneSigSpec2 := sig.Spec{
			QN: "sig-3",
			Ys: []chnl.Spec{lolliSigSpec.X, oneSig1.X},
			X: chnl.Spec{
				ChnlPH: "chnl-3",
				RoleQN: oneRole.QN,
			},
		}
		oneSig2, err := sigAPI.Create(oneSigSpec2)
		if err != nil {
			t.Fatal(err)
		}
		// and
		poolImpl, err := poolAPI.Create(
			pool.Spec{
				Title: "pool-1",
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		receiverPH := ph.New("receiver")
		receiverID, err := poolAPI.Spawn(
			pool.TranSpec{
				PoolID: poolImpl.PoolID,
				ProcID: poolImpl.ProcID,
				Term: step.SpawnSpec{
					SigID: lolliSig.SigID,
					X:     receiverPH,
				},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		messagePH := ph.New("message")
		_, err = poolAPI.Spawn(
			pool.TranSpec{
				PoolID: poolImpl.PoolID,
				ProcID: poolImpl.ProcID,
				Term: step.SpawnSpec{
					SigID: oneSig1.SigID,
					X:     messagePH,
				},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		senderPH := ph.New("sender")
		senderID, err := poolAPI.Spawn(
			pool.TranSpec{
				PoolID: poolImpl.PoolID,
				ProcID: poolImpl.ProcID,
				Term: step.SpawnSpec{
					SigID: oneSig2.SigID,
					Ys:    []ph.ADT{receiverPH, senderPH},
					X:     senderPH,
				},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		recvSpec := pool.TranSpec{
			PoolID: poolImpl.PoolID,
			ProcID: receiverID,
			Term: step.RecvSpec{
				X: receiverPH,
				Y: messagePH,
				Cont: step.WaitSpec{
					X: messagePH,
					Cont: step.CloseSpec{
						X: receiverPH,
					},
				},
			},
		}
		// when
		err = poolAPI.Take(recvSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		sendSpec := pool.TranSpec{
			PoolID: poolImpl.PoolID,
			ProcID: senderID,
			Term: step.SendSpec{
				X: receiverPH,
				Y: messagePH,
			},
		}
		// and
		err = poolAPI.Take(sendSpec)
		if err != nil {
			t.Fatal(err)
		}
		// then
		// TODO добавить проверку
	})

	t.Run("CaseLab", func(t *testing.T) {
		tc.Setup(t)
		// given
		label := core.Label("label-1")
		// and
		withRoleSpec := role.Spec{
			QN: "with-role",
			State: state.WithSpec{
				Choices: map[core.Label]state.Spec{
					label: state.OneSpec{},
				},
			},
		}
		withRole, err := roleAPI.Create(withRoleSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		oneRoleSpec := role.Spec{
			QN:    "one-role",
			State: state.OneSpec{},
		}
		oneRole, err := roleAPI.Create(oneRoleSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		withSigSpec := sig.Spec{
			QN: "sig-1",
			X: chnl.Spec{
				ChnlPH: "chnl-1",
				RoleQN: withRole.QN,
			},
		}
		withSig, err := sigAPI.Create(withSigSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		oneSigSpec := sig.Spec{
			QN: "sig-2",
			Ys: []chnl.Spec{withSig.X},
			X: chnl.Spec{
				ChnlPH: "chnl-2",
				RoleQN: oneRole.QN,
			},
		}
		oneSig, err := sigAPI.Create(oneSigSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		poolSpec := pool.Spec{
			Title: "pool-1",
		}
		poolImpl, err := poolAPI.Create(poolSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		followerPH := ph.New("follower")
		followerID, err := poolAPI.Spawn(
			pool.TranSpec{
				PoolID: poolImpl.PoolID,
				ProcID: poolImpl.ProcID,
				Term: step.SpawnSpec{
					SigID: withSig.SigID,
					X:     followerPH,
				},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		deciderPH := ph.New("decider")
		deciderID, err := poolAPI.Spawn(
			pool.TranSpec{
				PoolID: poolImpl.PoolID,
				ProcID: poolImpl.ProcID,
				Term: step.SpawnSpec{
					SigID: oneSig.SigID,
					Ys:    []ph.ADT{followerPH},
					X:     deciderPH,
				},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		caseSpec := pool.TranSpec{
			PoolID: poolImpl.PoolID,
			ProcID: followerID,
			Term: step.CaseSpec{
				X: followerPH,
				Conts: map[core.Label]step.Term{
					label: step.CloseSpec{
						X: followerPH,
					},
				},
			},
		}
		// when
		err = poolAPI.Take(caseSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		labSpec := pool.TranSpec{
			PoolID: poolImpl.PoolID,
			ProcID: deciderID,
			Term: step.LabSpec{
				X: followerPH,
				L: label,
			},
		}
		// and
		err = poolAPI.Take(labSpec)
		if err != nil {
			t.Fatal(err)
		}
		// then
		// TODO добавить проверку
	})

	t.Run("Spawn", func(t *testing.T) {
		tc.Setup(t)
		// given
		oneRole, err := roleAPI.Create(
			role.Spec{
				QN:    "one-role",
				State: state.OneSpec{},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		oneSig1, err := sigAPI.Create(
			sig.Spec{
				QN: "sig-1",
				X: chnl.Spec{
					ChnlPH: "chnl-1",
					RoleQN: oneRole.QN,
				},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		oneSig2, err := sigAPI.Create(
			sig.Spec{
				QN: "sig-2",
				Ys: []chnl.Spec{oneSig1.X},
				X: chnl.Spec{
					ChnlPH: "chnl-2",
					RoleQN: oneRole.QN,
				},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		oneSig3, err := sigAPI.Create(
			sig.Spec{
				QN: "sig-3",
				Ys: []chnl.Spec{oneSig1.X},
				X: chnl.Spec{
					ChnlPH: "chnl-3",
					RoleQN: oneRole.QN,
				},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		poolImpl, err := poolAPI.Create(
			pool.Spec{
				Title: "pool-1",
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		injecteePH := ph.New("injectee")
		_, err = poolAPI.Spawn(
			pool.TranSpec{
				PoolID: poolImpl.PoolID,
				ProcID: poolImpl.ProcID,
				Term: step.SpawnSpec{
					SigID: oneSig1.SigID,
					X:     injecteePH,
				},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		spawnerPH := ph.New("spawner")
		spawnerID, err := poolAPI.Spawn(
			pool.TranSpec{
				PoolID: poolImpl.PoolID,
				ProcID: poolImpl.ProcID,
				Term: step.SpawnSpec{
					SigID: oneSig2.SigID,
					Ys:    []ph.ADT{injecteePH},
					X:     spawnerPH,
				},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		x := ph.New("x")
		spawnSpec := pool.TranSpec{
			PoolID: poolImpl.PoolID,
			ProcID: spawnerID,
			Term: step.SpawnSpec{
				SigID: oneSig3.SigID,
				Ys:    []ph.ADT{injecteePH},
				X:     x,
				Cont: step.WaitSpec{
					X: x,
					Cont: step.CloseSpec{
						X: spawnerPH,
					},
				},
			},
		}
		// when
		err = poolAPI.Take(spawnSpec)
		if err != nil {
			t.Fatal(err)
		}
		// then
		// TODO добавить проверку
	})

	t.Run("Fwd", func(t *testing.T) {
		tc.Setup(t)
		// given
		oneRole, err := roleAPI.Create(
			role.Spec{
				QN:    "one-role",
				State: state.OneSpec{},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		oneSig1, err := sigAPI.Create(
			sig.Spec{
				QN: "sig-1",
				X: chnl.Spec{
					ChnlPH: "chnl-1",
					RoleQN: oneRole.QN,
				},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		oneSig2, err := sigAPI.Create(
			sig.Spec{
				QN: "sig-2",
				Ys: []chnl.Spec{oneSig1.X},
				X: chnl.Spec{
					ChnlPH: "chnl-2",
					RoleQN: oneRole.QN,
				},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		oneSig3, err := sigAPI.Create(
			sig.Spec{
				QN: "sig-3",
				Ys: []chnl.Spec{oneSig1.X},
				X: chnl.Spec{
					ChnlPH: "chnl-3",
					RoleQN: oneRole.QN,
				},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		poolImpl, err := poolAPI.Create(
			pool.Spec{
				Title: "pool-1",
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		closerPH := ph.New("closer")
		closerSpec := pool.TranSpec{
			PoolID: poolImpl.PoolID,
			ProcID: poolImpl.ProcID,
			Term: step.SpawnSpec{
				SigID: oneSig1.SigID,
				X:     closerPH,
			},
		}
		closerID, err := poolAPI.Spawn(closerSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		forwarderPH := ph.New("forwarder")
		forwarderSpec := pool.TranSpec{
			PoolID: poolImpl.PoolID,
			ProcID: poolImpl.ProcID,
			Term: step.SpawnSpec{
				SigID: oneSig2.SigID,
				Ys:    []ph.ADT{closerPH},
				X:     forwarderPH,
			},
		}
		forwarderID, err := poolAPI.Spawn(forwarderSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		waiterPH := ph.New("waiter")
		waiterSpec := pool.TranSpec{
			PoolID: poolImpl.PoolID,
			ProcID: poolImpl.ProcID,
			Term: step.SpawnSpec{
				SigID: oneSig3.SigID,
				Ys:    []ph.ADT{forwarderPH},
				X:     waiterPH,
			},
		}
		waiterID, err := poolAPI.Spawn(waiterSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		closeSpec := pool.TranSpec{
			PoolID: poolImpl.PoolID,
			ProcID: closerID,
			Term: step.CloseSpec{
				X: closerPH,
			},
		}
		err = poolAPI.Take(closeSpec)
		if err != nil {
			t.Fatal(err)
		}
		// when
		fwdSpec := pool.TranSpec{
			PoolID: poolImpl.PoolID,
			ProcID: forwarderID,
			Term: step.FwdSpec{
				X: forwarderPH,
				Y: closerPH,
			},
		}
		err = poolAPI.Take(fwdSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		waitSpec := pool.TranSpec{
			PoolID: poolImpl.PoolID,
			ProcID: waiterID,
			Term: step.WaitSpec{
				X: forwarderPH,
				Cont: step.CloseSpec{
					X: waiterPH,
				},
			},
		}
		err = poolAPI.Take(waitSpec)
		if err != nil {
			t.Fatal(err)
		}
		// then
		// TODO добавить проверку
	})
}
