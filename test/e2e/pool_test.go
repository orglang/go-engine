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
		"pool_prog_binds", "pool_type_defs", "pool_type_exps", "pool_prog_terms", "pool_term_exps",
		"pool_impl_binds", "pool_exec_evals", "pool_exec_vars", "pool_exec_comms", "pool_comm_acts",
		"proc_prog_binds", "proc_type_defs", "proc_type_exps", "proc_prog_terms", "proc_term_exps",
		"proc_impl_binds", "proc_exec_evals", "proc_exec_vars", "proc_exec_comms", "proc_comm_acts",
	}
	for _, table := range tables {
		if _, err := s.db.ExecContext(t.Context(), fmt.Sprintf("truncate table %v", table)); err != nil {
			t.Fatal(err)
		}
	}
}

// The generated OpenAPI union types intentionally hide their discriminator
// storage. Decode a contract-shaped value to construct them from JSON.
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

	poolTermQN, poolExecQN := "pool-term-qn", "pool-exec-qn"
	poolAssetPH, poolLiabPH := "pool-asset-ph", "pool-liab-ph"
	_, err = s.api.createPoolDeclaration(api.PoolDecSpec{
		ProgQn:  poolTermQN,
		LiabVar: api.AnyCommVarSpec{CommPh: poolLiabPH, TypeQn: poolTypeQN},
		AssetVars: []api.AnyCommVarSpec{
			{CommPh: poolAssetPH, TypeQn: poolTypeQN},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	poolExecRef, err := s.api.createPoolExecution(api.AnyEvalSpec{
		TermQn:  poolTermQN,
		LiabVar: api.AnyEvalVarSpec{CommPh: poolLiabPH, TermQn: poolExecQN},
		AssetVars: []api.AnyEvalVarSpec{
			{CommPh: poolAssetPH, TermQn: poolExecQN},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	err = s.api.stepPoolExecution(api.PoolStepSpec{
		RootRef: poolExecRef,
		PoolExpr: poolTermExpr("acquire", api.PoolTermAcquireSpec{
			CommPh: poolAssetPH,
			ContExpr: poolTermExpr("hire", api.PoolTermHireSpec{
				CommPh:   poolAssetPH,
				ProcQn:   closerProcQN,
				ContExpr: poolTermExpr("release", api.PoolTermReleaseSpec{CommPh: poolAssetPH}),
			}),
		}),
	})
	if err != nil {
		t.Fatal(err)
	}

	err = s.api.stepPoolExecution(api.PoolStepSpec{
		RootRef: poolExecRef,
		PoolExpr: poolTermExpr("accept", api.PoolTermAcceptSpec{
			CommPh: poolLiabPH,
			ContExpr: poolTermExpr("apply", api.PoolTermApplySpec{
				CommPh:   poolLiabPH,
				ProcQn:   closerProcQN,
				ContExpr: poolTermExpr("detach", api.PoolTermDetachSpec{CommPh: poolLiabPH}),
			}),
		}),
	})
	if err != nil {
		t.Fatal(err)
	}
}
