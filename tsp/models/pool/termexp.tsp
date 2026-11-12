import "../common.tsp";
namespace OrgLang.Engine;

/** Pool term expression specification (go-sdk/pool/termexp). */
model PoolTermExprSpec {
  kind: "hire" | "fire" | "apply" | "quit" | "acquire" | "release" | "accept" | "detach" | "spawn";
  hire?: PoolHireSpec;
  fire?: PoolFireSpec;
  apply?: PoolApplySpec;
  quit?: PoolQuitSpec;
  acquire?: PoolAcquireSpec;
  release?: PoolReleaseSpec;
  accept?: PoolAcceptSpec;
  detach?: PoolDetachSpec;
  spawn?: PoolSpawnSpec;
}

model PoolHireSpec { comm_ph: string; proc_qn: string; cont_exp: PoolTermExprSpec; }
model PoolFireSpec { CommChnlPH: string; }
model PoolApplySpec { comm_ph: string; proc_qn: string; cont_exp: PoolTermExprSpec; }
model PoolQuitSpec { comm_ph: string; proc_qn: string; }
model PoolAcquireSpec { comm_ph: string; cont_exp: PoolTermExprSpec; }
model PoolReleaseSpec { comm_ph: string; }
model PoolAcceptSpec { comm_ph: string; cont_exp: PoolTermExprSpec; }
model PoolDetachSpec { comm_ph: string; }
model PoolSpawnSpec { proc_dr: TermRef; proc_irs: ComputationRef[]; }
