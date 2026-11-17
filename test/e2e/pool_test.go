package e2e

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	"orglang/go-engine/internal/lib/e2e/api"
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
	api *generatedAPI
	db  *sql.DB
}

func (s *suite) beforeAll(t *testing.T) {
	db, err := sql.Open("pgx", "postgres://orglang:orglang@localhost:5432/orglang")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Logf("stop failed: %s", err)
		}
	})
	s.db = db
	s.api, err = newGeneratedAPI("http://localhost:8080")
	if err != nil {
		t.Fatal(err)
	}
}

func (s *suite) beforeEach(t *testing.T) {
	tables := []string{
		"pool_desc_binds", "pool_type_defs", "pool_type_exps", "pool_term_decs",
		"pool_impl_binds", "pool_comp_execs", "pool_comp_vars", "pool_comp_conns", "pool_conn_acts",
		"proc_desc_binds", "proc_type_defs", "proc_type_exps", "proc_term_decs",
		"proc_impl_binds", "proc_comp_execs", "proc_comp_vars", "proc_comp_conns", "proc_conn_acts",
	}
	for _, table := range tables {
		if _, err := s.db.ExecContext(t.Context(), fmt.Sprintf("truncate table %v", table)); err != nil {
			t.Fatal(err)
		}
	}
}

// The generated OpenAPI union types intentionally hide their discriminator
// storage. Decode a contract-shaped value to construct them without SDK types.
func poolTypeExpr(kind string, value any) api.PoolTypeExprSpec {
	return decodeUnion[api.PoolTypeExprSpec](kind, value)
}

func poolTermExpr(kind string, value any) api.PoolTermExprSpec {
	return decodeUnion[api.PoolTermExprSpec](kind, value)
}

func decodeUnion[T any](kind string, value any) T {
	var result T
	payload, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(payload, &fields); err != nil {
		panic(err)
	}
	fields["kind"], _ = json.Marshal(kind)
	payload, err = json.Marshal(fields)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(payload, &result); err != nil {
		panic(err)
	}
	return result
}

func (s *suite) waitClose(t *testing.T) {
	s.beforeEach(t)
	closerProcQN, waiterProcQN := "closer-proc-qn", "waiter-proc-qn"
	poolTypeQN := "pool-type-qn"

	_, err := s.api.createPoolTypeDefinition(api.PoolTypeDefSpec{
		Qn: poolTypeQN,
		Expr: poolTypeExpr("up", api.PoolTypeUpSpec{
			ContExpr: poolTypeExpr("with", api.PoolTypeWithSpec{
				ProcQns: []string{closerProcQN, waiterProcQN},
				ContExpr: poolTypeExpr("down", api.PoolTypeDownSpec{
					ContExpr: poolTypeExpr("link", api.PoolTypeLinkSpec{
						Link: api.AnyTypeLinkSpec{TypeQn: poolTypeQN},
					}),
				}),
			}),
		}),
	})
	if err != nil {
		t.Fatal(err)
	}

	poolTermQN, poolCompQN := "pool-term-qn", "pool-comp-qn"
	poolAssetPH, poolLiabPH := "pool-asset-ph", "pool-liab-ph"
	_, err = s.api.createPoolDeclaration(api.PoolDecSpec{
		DescQn:  poolTermQN,
		LiabVar: api.AnyConnVarSpec{ConnPh: poolLiabPH, TypeQn: poolTypeQN},
		AssetVars: []api.AnyConnVarSpec{
			{ConnPh: poolAssetPH, TypeQn: poolTypeQN},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	poolExecRef, err := s.api.createPoolExecution(api.AnyExecSpec{
		TermQn:  poolTermQN,
		LiabVar: api.AnyExecVarSpec{ConnPh: poolLiabPH, TermQn: poolCompQN},
		AssetVars: []api.AnyExecVarSpec{
			{ConnPh: poolAssetPH, TermQn: poolCompQN},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	err = s.api.stepPoolExecution(api.PoolStepSpec{
		CompRef: poolExecRef,
		PoolExpr: poolTermExpr("acquire", api.PoolTermAcquireSpec{
			ConnPh: poolAssetPH,
			ContExpr: poolTermExpr("hire", api.PoolTermHireSpec{
				ConnPh:   poolAssetPH,
				ProcQn:   closerProcQN,
				ContExpr: poolTermExpr("release", api.PoolTermReleaseSpec{ConnPh: poolAssetPH}),
			}),
		}),
	})
	if err != nil {
		t.Fatal(err)
	}

	err = s.api.stepPoolExecution(api.PoolStepSpec{
		CompRef: poolExecRef,
		PoolExpr: poolTermExpr("accept", api.PoolTermAcceptSpec{
			ConnPh: poolLiabPH,
			ContExpr: poolTermExpr("apply", api.PoolTermApplySpec{
				ConnPh:   poolLiabPH,
				ProcQn:   closerProcQN,
				ContExpr: poolTermExpr("detach", api.PoolTermDetachSpec{ConnPh: poolLiabPH}),
			}),
		}),
	})
	if err != nil {
		t.Fatal(err)
	}
}

// The following scenarios are retained verbatim while they are migrated to the generated API.
// They intentionally remain comments: they still use the removed Resty/go-sdk test client.
// 	// and
// 	time.Sleep(1 * time.Second)
// 	// and
// 	oneTypeQN := "one-type-qn"
// 	_, err = s.TypeDefAPI.Create(typedef.DefSpec{
// 		TypeQN:  oneTypeQN,
// 		TypeExp: typeexp.ExpSpec{K: typeexp.One},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	// and
// 	err = s.PoolExecAPI.Take(compstep.StepSpec{
// 		CompRef: poolExecRef,
// 		PoolExp: pooltermexp.ExpSpec{
// 			K: pooltermexp.Acquire,
// 			Acquire: &pooltermexp.AcquireSpec{
// 				CommChnlPH: poolAssetPH,
// 				ContExp: pooltermexp.ExpSpec{
// 					K: pooltermexp.Hire,
// 					Hire: &pooltermexp.HireSpec{
// 						CommChnlPH: poolAssetPH,
// 						ProcTermQN: waiterProcQN,
// 						ContExp: pooltermexp.ExpSpec{
// 							K: pooltermexp.Release,
// 							Release: &pooltermexp.ReleaseSpec{
// 								CommChnlPH: poolAssetPH,
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	err = s.PoolExecAPI.Take(compstep.StepSpec{
// 		CompRef: poolExecRef,
// 		PoolExp: pooltermexp.ExpSpec{
// 			K: pooltermexp.Accept,
// 			Accept: &pooltermexp.AcceptSpec{
// 				CommChnlPH: poolAssetPH,
// 				ContExp: pooltermexp.ExpSpec{
// 					K: pooltermexp.Apply,
// 					Apply: &pooltermexp.ApplySpec{
// 						CommChnlPH: poolLiabPH,
// 						ProcTermQN: waiterProcQN,
// 						ContExp: pooltermexp.ExpSpec{
// 							K: pooltermexp.Detach,
// 							Detach: &pooltermexp.DetachSpec{
// 								CommChnlPH: poolAssetPH,
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	// and
// 	// 	waiterProcExec, err := s.PoolExecAPI.Spawn(compstep.StepSpec{
// 	// 		CompRef: poolExecRef,
// 	// 		PoolExp: pooltermexp.ExpSpec{
// 	// 			K: pooltermexp.Spawn,
// 	// 			Spawn: &pooltermexp.SpawnSpec{
// 	// 				ProcTermRef:  waiterProcDec.TermRef,
// 	// 				ProcCompRefs: []compsem.SemRef{closerProcExec},
// 	// 			},
// 	// 		},
// 	// 	})
// 	// 	if err != nil {
// 	// 		t.Fatal(err)
// 	// 	}
// 	// 	// when
// 	// 	err = s.ProcExecAPI.Take(proccompstep.StepSpec{
// 	// 		CompRef: closerProcExec,
// 	// 		ProcExp: proctermexp.ExpSpec{
// 	// 			K: proctermexp.Close,
// 	// 			Close: &proctermexp.CloseSpec{
// 	// 				CommChnlPH: closerProcDec.LiabVar.ChnlPH,
// 	// 			},
// 	// 		},
// 	// 	})
// 	// 	if err != nil {
// 	// 		t.Fatal(err)
// 	// 	}
// 	// 	// and
// 	// 	err = s.ProcExecAPI.Take(proccompstep.StepSpec{
// 	// 		CompRef: waiterProcExec,
// 	// 		ProcExp: proctermexp.ExpSpec{
// 	// 			K: proctermexp.Wait,
// 	// 			Wait: &proctermexp.WaitSpec{
// 	// 				CommChnlPH: waiterProcDec.AssetVars[0].ChnlPH,
// 	// 				ContES: proctermexp.ExpSpec{
// 	// 					K: proctermexp.Close,
// 	// 					Close: &proctermexp.CloseSpec{
// 	// 						CommChnlPH: waiterProcDec.LiabVar.ChnlPH,
// 	// 					},
// 	// 				},
// 	// 			},
// 	// 		},
// 	// 	})
// 	// 	if err != nil {
// 	// 		t.Fatal(err)
// 	// 	}
// 	// 	// then
// 	// 	// TODO добавить проверку
// }
//
//

// func (s *suite) recvSend(t *testing.T) {
// 	s.beforeEach(t)
// 	// given
// 	senderProcQN := "sender-proc-qn"
// 	messageProcQN := "message-proc-qn"
// 	receiverProcQN := "receiver-proc-qn"
// 	// and
// 	withXactQN := "with-xact-qn"
// 	_, err := s.XactDefAPI.Create(pooltypedef.DefSpec{
// 		TypeQN: withXactQN,
// 		TypeExp: pooltypeexp.ExpSpec{
// 			K: pooltypeexp.Up,
// 			Up: &pooltypeexp.ShiftSpec{
// 				ContExp: pooltypeexp.ExpSpec{
// 					K: pooltypeexp.With,
// 					With: &pooltypeexp.LaborSpec{
// 						ProcQNs: []string{senderProcQN, receiverProcQN, messageProcQN},
// 						ContExp: pooltypeexp.ExpSpec{
// 							K: pooltypeexp.Down,
// 							Down: &pooltypeexp.ShiftSpec{
// 								ContExp: pooltypeexp.ExpSpec{
// 									K:    pooltypeexp.Link,
// 									Link: &pooltypeexp.LinkSpec{TypeQN: withXactQN},
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	myPoolQN := "my-pool-qn"
// 	poolProviderPH := "pool-provider-ph"
// 	poolClientPH := "pool-client-ph"
// 	_, err = s.PoolDecAPI.Create(pooltermdec.DecSpec{
// 		TermQN: myPoolQN,
// 		LiabVar: termvar.VarSpec{
// 			ChnlPH: poolProviderPH,
// 			TypeQN: withXactQN,
// 		},
// 		AssetVars: []termvar.VarSpec{{
// 			ChnlPH: poolClientPH,
// 			TypeQN: withXactQN,
// 		}},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	myPoolExec, err := s.PoolExecAPI.Create(compexec.ExecSpec{
// 		TermQN: myPoolQN,
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	lolliTypeQN := "lolli-type-qn"
// 	_, err = s.TypeDefAPI.Create(typedef.DefSpec{
// 		TypeQN: lolliTypeQN,
// 		TypeExp: typeexp.ExpSpec{
// 			K: typeexp.Lolli,
// 			Lolli: &typeexp.ProdSpec{
// 				Val:  typeexp.ExpSpec{K: typeexp.One},
// 				Cont: typeexp.ExpSpec{K: typeexp.One},
// 			},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	oneTypeQN := "one-type-qn"
// 	_, err = s.TypeDefAPI.Create(typedef.DefSpec{
// 		TypeQN:  oneTypeQN,
// 		TypeExp: typeexp.ExpSpec{K: typeexp.One},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	receiverProcDec, err := s.ProcDecAPI.Create(proctermdec.DecSpec{
// 		TermQN: receiverProcQN,
// 		LiabVar: termvar.VarSpec{
// 			ChnlPH: "receiver-provider-ph",
// 			TypeQN: lolliTypeQN,
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	messageProcDec, err := s.ProcDecAPI.Create(proctermdec.DecSpec{
// 		TermQN: messageProcQN,
// 		LiabVar: termvar.VarSpec{
// 			ChnlPH: "message-provider-ph",
// 			TypeQN: oneTypeQN,
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	senderProcDec, err := s.ProcDecAPI.Create(proctermdec.DecSpec{
// 		TermQN: senderProcQN,
// 		LiabVar: termvar.VarSpec{
// 			ChnlPH: "sender-provider-ph",
// 			TypeQN: oneTypeQN,
// 		},
// 		AssetVars: []termvar.VarSpec{
// 			{ChnlPH: "sender-receiver-ph", TypeQN: lolliTypeQN},
// 			{ChnlPH: "sender-message-ph", TypeQN: oneTypeQN},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	err = s.PoolExecAPI.Take(compstep.StepSpec{
// 		CompRef: myPoolExec,
// 		PoolExp: pooltermexp.ExpSpec{
// 			K: pooltermexp.Hire,
// 			Hire: &pooltermexp.HireSpec{
// 				CommChnlPH: poolClientPH,
// 				ProcTermQN: receiverProcQN,
// 			},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	receiverProcExec, err := s.PoolExecAPI.Spawn(compstep.StepSpec{
// 		CompRef: myPoolExec,
// 		PoolExp: pooltermexp.ExpSpec{
// 			K: pooltermexp.Spawn,
// 			Spawn: &pooltermexp.SpawnSpec{
// 				ProcTermRef: receiverProcDec.TermRef,
// 			},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	err = s.PoolExecAPI.Take(compstep.StepSpec{
// 		CompRef: myPoolExec,
// 		PoolExp: pooltermexp.ExpSpec{
// 			K: pooltermexp.Hire,
// 			Hire: &pooltermexp.HireSpec{
// 				CommChnlPH: poolClientPH,
// 				ProcTermQN: messageProcQN,
// 			},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	messageProcExec, err := s.PoolExecAPI.Spawn(compstep.StepSpec{
// 		CompRef: myPoolExec,
// 		PoolExp: pooltermexp.ExpSpec{
// 			K: pooltermexp.Spawn,
// 			Spawn: &pooltermexp.SpawnSpec{
// 				ProcTermRef: messageProcDec.TermRef,
// 			},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	err = s.PoolExecAPI.Take(compstep.StepSpec{
// 		CompRef: myPoolExec,
// 		PoolExp: pooltermexp.ExpSpec{
// 			K: pooltermexp.Hire,
// 			Hire: &pooltermexp.HireSpec{
// 				CommChnlPH: poolClientPH,
// 				ProcTermQN: senderProcQN,
// 			},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	senderProcExec, err := s.PoolExecAPI.Spawn(compstep.StepSpec{
// 		CompRef: myPoolExec,
// 		PoolExp: pooltermexp.ExpSpec{
// 			K: pooltermexp.Spawn,
// 			Spawn: &pooltermexp.SpawnSpec{
// 				ProcTermRef:  senderProcDec.TermRef,
// 				ProcCompRefs: []compsem.SemRef{receiverProcExec, messageProcExec},
// 			},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// when
// 	err = s.ProcExecAPI.Take(proccompstep.StepSpec{
// 		CompRef: receiverProcExec,
// 		ProcExp: proctermexp.ExpSpec{
// 			K: proctermexp.Recv,
// 			Recv: &proctermexp.RecvSpec{
// 				CommChnlPH: receiverProcDec.LiabVar.ChnlPH,
// 				BindChnlPH: "receiver-message-ph",
// 			},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	err = s.ProcExecAPI.Take(proccompstep.StepSpec{
// 		CompRef: senderProcExec,
// 		ProcExp: proctermexp.ExpSpec{
// 			K: proctermexp.Send,
// 			Send: &proctermexp.SendSpec{
// 				CommChnlPH: senderProcDec.AssetVars[0].ChnlPH,
// 				ValChnlPH:  senderProcDec.AssetVars[1].ChnlPH,
// 			},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// then
// 	// TODO добавить проверку
// }
//
// func (s *suite) caseLab(t *testing.T) {
// 	s.beforeEach(t)
// 	// given
// 	myPoolQN := "my-pool-qn"
// 	myPoolExec, err := s.PoolExecAPI.Create(compexec.ExecSpec{
// 		TermQN: myPoolQN,
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	patternQN := "pattern-qn"
// 	// and
// 	withType, err := s.TypeDefAPI.Create(typedef.DefSpec{
// 		TypeQN: "with-type-qn",
// 		TypeExp: typeexp.ExpSpec{
// 			With: &typeexp.SumSpec{
// 				Choices: []typeexp.ChoiceSpec{
// 					{LabQN: patternQN, Cont: typeexp.ExpSpec{K: typeexp.One}},
// 				},
// 			},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	oneType, err := s.TypeDefAPI.Create(typedef.DefSpec{
// 		TypeQN:  "one-type-qn",
// 		TypeExp: typeexp.ExpSpec{K: typeexp.One},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	followerProcQN := "follower-proc-qn"
// 	followerProcDec := proctermdec.DecSpec{
// 		TermQN: followerProcQN,
// 		LiabVar: termvar.VarSpec{
// 			ChnlPH: "follower-provider-ph",
// 			TypeQN: withType.DefSpec.TypeQN,
// 		},
// 	}
// 	followerProcDec2, err := s.ProcDecAPI.Create(followerProcDec)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	deciderProcQN := "decider-proc-qn"
// 	deciderProcDec, err := s.ProcDecAPI.Create(proctermdec.DecSpec{
// 		TermQN:    deciderProcQN,
// 		AssetVars: []termvar.VarSpec{followerProcDec.LiabVar},
// 		LiabVar: termvar.VarSpec{
// 			ChnlPH: "decider-provider-ph",
// 			TypeQN: oneType.DefSpec.TypeQN,
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	poolFollowerPH := "pool-follower-ph"
// 	err = s.PoolExecAPI.Take(compstep.StepSpec{
// 		CompRef: myPoolExec,
// 		PoolExp: pooltermexp.ExpSpec{
// 			K: pooltermexp.Hire,
// 			Hire: &pooltermexp.HireSpec{
// 				ProcTermQN: followerProcQN,
// 			},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	followerProcExec, err := s.PoolExecAPI.Spawn(compstep.StepSpec{
// 		CompRef: myPoolExec,
// 		PoolExp: pooltermexp.ExpSpec{
// 			K: pooltermexp.Spawn,
// 			Spawn: &pooltermexp.SpawnSpec{
// 				ProcTermRef: followerProcDec2.TermRef,
// 			},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	err = s.PoolExecAPI.Take(compstep.StepSpec{
// 		CompRef: myPoolExec,
// 		PoolExp: pooltermexp.ExpSpec{
// 			K: pooltermexp.Hire,
// 			Hire: &pooltermexp.HireSpec{
// 				ProcTermQN: deciderProcQN,
// 			},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	deciderProcExec, err := s.PoolExecAPI.Spawn(compstep.StepSpec{
// 		CompRef: myPoolExec,
// 		PoolExp: pooltermexp.ExpSpec{
// 			K: pooltermexp.Spawn,
// 			Spawn: &pooltermexp.SpawnSpec{
// 				ProcTermRef:  deciderProcDec.TermRef,
// 				ProcCompRefs: []compsem.SemRef{followerProcExec},
// 			},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// when
// 	err = s.ProcExecAPI.Take(proccompstep.StepSpec{
// 		CompRef: followerProcExec,
// 		ProcExp: proctermexp.ExpSpec{
// 			K: proctermexp.Case,
// 			Case: &proctermexp.CaseSpec{
// 				CommChnlPH: poolFollowerPH,
// 				ContBSes: []proctermexp.BranchSpec{
// 					{PatternQN: patternQN, ContES: proctermexp.ExpSpec{
// 						K: proctermexp.Close,
// 						Close: &proctermexp.CloseSpec{
// 							CommChnlPH: poolFollowerPH,
// 						},
// 					},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	err = s.ProcExecAPI.Take(proccompstep.StepSpec{
// 		CompRef: deciderProcExec,
// 		ProcExp: proctermexp.ExpSpec{
// 			K: proctermexp.Lab,
// 			Lab: &proctermexp.LabSpec{
// 				CommChnlPH: poolFollowerPH,
// 				PatternQN:  patternQN,
// 			},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// then
// 	// TODO добавить проверку
// }
//
// func (s *suite) call(t *testing.T) {
// 	s.beforeEach(t)
// 	// given
// 	myPoolQN := "my-pool-qn"
// 	myPoolExec, err := s.PoolExecAPI.Create(compexec.ExecSpec{
// 		TermQN: myPoolQN,
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	oneType, err := s.TypeDefAPI.Create(
// 		typedef.DefSpec{
// 			TypeQN:  "one-type-qn",
// 			TypeExp: typeexp.ExpSpec{K: typeexp.One},
// 		},
// 	)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	injecteeProcQN := "injectee-proc-qn"
// 	injecteeProcDec, err := s.ProcDecAPI.Create(proctermdec.DecSpec{
// 		TermQN: injecteeProcQN,
// 		LiabVar: termvar.VarSpec{
// 			ChnlPH: "injectee-provider-ph",
// 			TypeQN: oneType.DefSpec.TypeQN,
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	err = s.PoolExecAPI.Take(compstep.StepSpec{
// 		CompRef: myPoolExec,
// 		PoolExp: pooltermexp.ExpSpec{
// 			K: pooltermexp.Hire,
// 			Hire: &pooltermexp.HireSpec{
// 				ProcTermQN: injecteeProcQN,
// 			},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	injecteeProcExec, err := s.PoolExecAPI.Spawn(compstep.StepSpec{
// 		CompRef: myPoolExec,
// 		PoolExp: pooltermexp.ExpSpec{
// 			K: pooltermexp.Spawn,
// 			Spawn: &pooltermexp.SpawnSpec{
// 				ProcTermRef: injecteeProcDec.TermRef,
// 			},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	callerProcQN := "caller-proc-qn"
// 	callerProcDec, err := s.ProcDecAPI.Create(proctermdec.DecSpec{
// 		TermQN: callerProcQN,
// 		LiabVar: termvar.VarSpec{
// 			ChnlPH: "caller-provider-ph",
// 			TypeQN: oneType.DefSpec.TypeQN,
// 		},
// 		AssetVars: []termvar.VarSpec{
// 			{ChnlPH: "caller-injectee-ph", TypeQN: oneType.DefSpec.TypeQN},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	err = s.PoolExecAPI.Take(compstep.StepSpec{
// 		CompRef: myPoolExec,
// 		PoolExp: pooltermexp.ExpSpec{
// 			K: pooltermexp.Hire,
// 			Hire: &pooltermexp.HireSpec{
// 				ProcTermQN: callerProcQN,
// 			},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	callerProcExec, err := s.PoolExecAPI.Spawn(compstep.StepSpec{
// 		CompRef: myPoolExec,
// 		PoolExp: pooltermexp.ExpSpec{
// 			K: pooltermexp.Spawn,
// 			Spawn: &pooltermexp.SpawnSpec{
// 				ProcTermRef:  callerProcDec.TermRef,
// 				ProcCompRefs: []compsem.SemRef{injecteeProcExec},
// 			},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	calleeProcQN := "callee-proc-qn"
// 	_, err = s.ProcDecAPI.Create(proctermdec.DecSpec{
// 		TermQN: calleeProcQN,
// 		LiabVar: termvar.VarSpec{
// 			ChnlPH: "callee-provider-ph",
// 			TypeQN: oneType.DefSpec.TypeQN,
// 		},
// 		AssetVars: []termvar.VarSpec{
// 			{ChnlPH: "callee-injectee-ph", TypeQN: oneType.DefSpec.TypeQN},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	callerCalleePH := "caller-callee-ph"
// 	// when
// 	err = s.ProcExecAPI.Take(proccompstep.StepSpec{
// 		CompRef: callerProcExec,
// 		ProcExp: proctermexp.ExpSpec{
// 			K: proctermexp.Call,
// 			Call: &proctermexp.CallSpec{
// 				BindChnlPH: callerCalleePH,
// 				ProcTermQN: calleeProcQN,
// 				ValChnlPHs: []string{callerProcDec.AssetVars[0].ChnlPH},
// 				ContES: proctermexp.ExpSpec{
// 					K: proctermexp.Wait,
// 					Wait: &proctermexp.WaitSpec{
// 						CommChnlPH: callerCalleePH,
// 						ContES: proctermexp.ExpSpec{
// 							K: proctermexp.Close,
// 							Close: &proctermexp.CloseSpec{
// 								CommChnlPH: callerProcDec.LiabVar.ChnlPH,
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// then
// 	// TODO добавить проверку
// }
//
// func (s *suite) fwd(t *testing.T) {
// 	s.beforeEach(t)
// 	// given
// 	myPoolQN := "my-pool-qn"
// 	myPoolExec, err := s.PoolExecAPI.Create(compexec.ExecSpec{
// 		TermQN: myPoolQN,
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	oneType, err := s.TypeDefAPI.Create(typedef.DefSpec{
// 		TypeQN:  "one-type-qn",
// 		TypeExp: typeexp.ExpSpec{K: typeexp.One},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	closerProcQN := "closer-proc-qn"
// 	closerProcDec, err := s.ProcDecAPI.Create(proctermdec.DecSpec{
// 		TermQN: closerProcQN,
// 		LiabVar: termvar.VarSpec{
// 			ChnlPH: "closer-provider-ph",
// 			TypeQN: oneType.DefSpec.TypeQN,
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	forwarderProcQN := "forwarder-proc-qn"
// 	forwarderProcDec, err := s.ProcDecAPI.Create(proctermdec.DecSpec{
// 		TermQN: forwarderProcQN,
// 		LiabVar: termvar.VarSpec{
// 			ChnlPH: "forwarder-provider-ph",
// 			TypeQN: oneType.DefSpec.TypeQN,
// 		},
// 		AssetVars: []termvar.VarSpec{
// 			{ChnlPH: "forwarder-closer-ph", TypeQN: oneType.DefSpec.TypeQN},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	waiterProcQN := "waiter-proc-qn"
// 	waiterProcDec, err := s.ProcDecAPI.Create(proctermdec.DecSpec{
// 		TermQN: waiterProcQN,
// 		LiabVar: termvar.VarSpec{
// 			ChnlPH: "waiter-provider-ph",
// 			TypeQN: oneType.DefSpec.TypeQN,
// 		},
// 		AssetVars: []termvar.VarSpec{
// 			{ChnlPH: "waiter-forwarder-ph", TypeQN: oneType.DefSpec.TypeQN},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	err = s.PoolExecAPI.Take(compstep.StepSpec{
// 		CompRef: myPoolExec,
// 		PoolExp: pooltermexp.ExpSpec{
// 			K: pooltermexp.Hire,
// 			Hire: &pooltermexp.HireSpec{
// 				ProcTermQN: closerProcQN,
// 			},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	closerProcExec, err := s.PoolExecAPI.Spawn(compstep.StepSpec{
// 		CompRef: myPoolExec,
// 		PoolExp: pooltermexp.ExpSpec{
// 			K: pooltermexp.Spawn,
// 			Spawn: &pooltermexp.SpawnSpec{
// 				ProcTermRef: closerProcDec.TermRef,
// 			},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	err = s.PoolExecAPI.Take(compstep.StepSpec{
// 		CompRef: myPoolExec,
// 		PoolExp: pooltermexp.ExpSpec{
// 			K: pooltermexp.Hire,
// 			Hire: &pooltermexp.HireSpec{
// 				ProcTermQN: forwarderProcQN,
// 			},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	forwarderProcExec, err := s.PoolExecAPI.Spawn(compstep.StepSpec{
// 		CompRef: myPoolExec,
// 		PoolExp: pooltermexp.ExpSpec{
// 			K: pooltermexp.Spawn,
// 			Spawn: &pooltermexp.SpawnSpec{
// 				ProcTermRef:  forwarderProcDec.TermRef,
// 				ProcCompRefs: []compsem.SemRef{closerProcExec},
// 			},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	err = s.PoolExecAPI.Take(compstep.StepSpec{
// 		CompRef: myPoolExec,
// 		PoolExp: pooltermexp.ExpSpec{
// 			K: pooltermexp.Hire,
// 			Hire: &pooltermexp.HireSpec{
// 				ProcTermQN: waiterProcQN,
// 			},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	waiterProcExec, err := s.PoolExecAPI.Spawn(compstep.StepSpec{
// 		CompRef: myPoolExec,
// 		PoolExp: pooltermexp.ExpSpec{
// 			K: pooltermexp.Spawn,
// 			Spawn: &pooltermexp.SpawnSpec{
// 				ProcTermRef:  waiterProcDec.TermRef,
// 				ProcCompRefs: []compsem.SemRef{forwarderProcExec},
// 			},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	err = s.ProcExecAPI.Take(proccompstep.StepSpec{
// 		CompRef: closerProcExec,
// 		ProcExp: proctermexp.ExpSpec{
// 			K: proctermexp.Close,
// 			Close: &proctermexp.CloseSpec{
// 				CommChnlPH: closerProcDec.LiabVar.ChnlPH,
// 			},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// when
// 	err = s.ProcExecAPI.Take(proccompstep.StepSpec{
// 		CompRef: forwarderProcExec,
// 		ProcExp: proctermexp.ExpSpec{
// 			K: proctermexp.Fwd,
// 			Fwd: &proctermexp.FwdSpec{
// 				CommChnlPH: forwarderProcDec.LiabVar.ChnlPH,
// 				ContChnlPH: forwarderProcDec.AssetVars[0].ChnlPH,
// 			},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// and
// 	err = s.ProcExecAPI.Take(proccompstep.StepSpec{
// 		CompRef: waiterProcExec,
// 		ProcExp: proctermexp.ExpSpec{
// 			K: proctermexp.Wait,
// 			Wait: &proctermexp.WaitSpec{
// 				CommChnlPH: waiterProcDec.AssetVars[0].ChnlPH,
// 				ContES: proctermexp.ExpSpec{
// 					K: proctermexp.Close,
// 					Close: &proctermexp.CloseSpec{
// 						CommChnlPH: waiterProcDec.LiabVar.ChnlPH,
// 					},
// 				},
// 			},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// then
// 	// TODO добавить проверку
// }
//
