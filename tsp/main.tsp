import "@typespec/http";
import "@typespec/openapi3";

using TypeSpec.Http;

@service(#{ title: "OrgLang Go Engine API" })
@server("http://localhost:8080", "Local development server")
namespace OrgLang.Engine;

/** A semantic reference to a term declaration. */
model TermRef {
  term_id: string;
  term_rn: int64;
}

/** A semantic reference to a type definition. */
model TypeRef {
  type_id: string;
  type_rn: int64;
}

/** A semantic reference to a process or pool component. */
model ComponentRef {
  comp_id: string;
  comp_rn: int64;
}

model TermVarSpec {
  chnl_ph: string;
  type_qn: string;
}

model CompVarSpec {
  chnl_ph: string;
  term_qn: string;
}

model TermVarRec {
  type_ref: TypeRef;
  chnl_ph: string;
  exp_vk: int64;
}

/** Process type expressions (proc/typeexp). */
model ProcTypeExpr {
  kind: "one" | "link" | "tensor" | "lolli" | "plus" | "with" | "up" | "down";
  link?: TypeLinkSpec;
  tensor?: ProcTypeProductSpec;
  lolli?: ProcTypeProductSpec;
  plus?: ProcTypeSumSpec;
  with?: ProcTypeSumSpec;
  up?: ProcTypeShiftSpec;
  down?: ProcTypeShiftSpec;
}

model TypeLinkSpec {
  type_qn: string;
}

model ProcTypeProductSpec {
  val_es: ProcTypeExpr;
  cont_es: ProcTypeExpr;
}

model ProcTypeSumSpec {
  choices: ProcTypeChoiceSpec[];
}

model ProcTypeChoiceSpec {
  lab_qn: string;
  cont_es: ProcTypeExpr;
}

model ProcTypeShiftSpec {
  cont_es: ProcTypeExpr;
}

/** Process term expressions (proc/termexp). */
model ProcTermExpr {
  kind: "close" | "wait" | "send" | "recv" | "lab" | "case" | "link" | "fwd" | "call";
  close?: ProcCloseSpec;
  wait?: ProcWaitSpec;
  send?: ProcSendSpec;
  recv?: ProcRecvSpec;
  lab?: ProcLabSpec;
  case?: ProcCaseSpec;
  call?: ProcCallSpec;
  fwd?: ProcFwdSpec;
}

model ProcCloseSpec { comm_ph: string; }
model ProcWaitSpec { comm_ph: string; cont_es: ProcTermExpr; }
model ProcSendSpec { comm_ph: string; val_ph: string; }
model ProcRecvSpec { comm_ph: string; bind_ph: string; cont_es: ProcTermExpr; }
model ProcLabSpec { comm_ph: string; pattern_qn: string; }
model ProcCaseSpec { comm_ph: string; cont_bses: ProcBranchSpec[]; }
model ProcBranchSpec { pattern_qn: string; cont_es: ProcTermExpr; }
model ProcCallSpec {
  bind_ph: string;
  proc_qn: string;
  val_phs: string[];
  cont_es: ProcTermExpr;
}
model ProcFwdSpec { comm_ph: string; cont_ph: string; }

/** Pool type expressions (pool/typeexp). */
model PoolTypeExpr {
  kind: "one" | "link" | "plus" | "with" | "tensor" | "lolli" | "up" | "down";
  link?: TypeLinkSpec;
  plus?: PoolTypeLaborSpec;
  with?: PoolTypeLaborSpec;
  tensor?: PoolTypeResourceSpec;
  lolli?: PoolTypeResourceSpec;
  up?: PoolTypeShiftSpec;
  down?: PoolTypeShiftSpec;
}

model PoolTypeLaborSpec {
  proc_qns: string[];
  cont_exp: PoolTypeExpr;
}

model PoolTypeResourceSpec {
  val_exp: PoolTypeExpr;
  cont_exp: PoolTypeExpr;
}

model PoolTypeShiftSpec {
  cont_exp: PoolTypeExpr;
}

/** Pool term expressions (pool/termexp). */
model PoolTermExpr {
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

model PoolHireSpec { comm_ph: string; proc_qn: string; cont_exp: PoolTermExpr; }
model PoolFireSpec { comm_ph: string; }
model PoolApplySpec { comm_ph: string; proc_qn: string; cont_exp: PoolTermExpr; }
model PoolQuitSpec { comm_ph: string; proc_qn: string; }
model PoolAcquireSpec { comm_ph: string; cont_exp: PoolTermExpr; }
model PoolReleaseSpec { comm_ph: string; }
model PoolAcceptSpec { comm_ph: string; cont_exp: PoolTermExpr; }
model PoolDetachSpec { comm_ph: string; }
model PoolSpawnSpec { proc_dr: TermRef; proc_irs: ComponentRef[]; }

model ProcTypeSpec {
  type_qn: string;
  type_exp: ProcTypeExpr;
}

model ProcTypeSnap {
  ref: TypeRef;
  spec: ProcTypeSpec;
}

model PoolTypeSpec {
  qn: string;
  exp: PoolTypeExpr;
}

model PoolTypeSnap {
  ref: TypeRef;
  exp: Expr;
}

model ProcDeclSpec {
  term_qn: string;
  liab_var: TermVarSpec;
  asset_vars: TermVarSpec[];
}

model ProcDeclSnap {
  term_ref: TermRef;
  liab_var: TermVarRec;
  asset_vars: TermVarRec[];
}

model PoolDeclSpec {
  desc_qn: string;
  liab_var: TermVarSpec;
  asset_vars: TermVarSpec[];
}

model ProcStepSpec {
  comp_ref: ComponentRef;
  proc_exp: ProcTermExpr;
}

model PoolExecSpec {
  term_qn: string;
  liab_var: CompVarSpec;
  asset_vars: CompVarSpec[];
}

model PoolStepSpec {
  comp_ref: ComponentRef;
  pool_exp: PoolTermExpr;
}

model CompSnap {
  comp_ref: ComponentRef;
}

model Created<T> {
  @statusCode statusCode: 201;
  @body body: T;
}

model CreatedNoBody {
  @statusCode statusCode: 201;
}

model Ok<T> {
  @statusCode statusCode: 200;
  @body body: T;
}

model NoContent {
  @statusCode statusCode: 204;
}

model OkNoBody {
  @statusCode statusCode: 200;
}

@route("/api/v1")
interface EngineApi {
  /** Submit program source text. */
  @post
  @route("/progs")
  createProgram(@body body: string): CreatedNoBody;

  @post
  @route("/xacts/defs")
  createPoolTypeDef(@body body: PoolTypeSpec): Created<TypeRef>;

  @post
  @route("/procs/decs")
  createProcDecl(@body body: ProcDeclSpec): Created<ProcDeclSnap>;

  /** Omit revision to request the latest available declaration revision. */
  @get
  @route("/procs/decs/{id}")
  getProcDecl(
    @path id: string,
    @query revision?: int64
  ): Ok<ProcDeclSnap>;

  /** Omit revision to request the latest available process revision. */
  @get
  @route("/procs/{id}")
  getProc(
    @path id: string,
    @query revision?: int64
  ): Ok<CompSnap>;

  @post
  @route("/procs/{id}/steps")
  takeProcStep(@path id: string, @body body: ProcStepSpec): OkNoBody;

  @post
  @route("/types")
  createProcType(@body body: ProcTypeSpec): Created<ProcTypeSnap>;

  /** Omit revision to request the latest available type revision. */
  @get
  @route("/types/{id}")
  getProcType(
    @path id: string,
    @query revision?: int64
  ): Ok<ProcTypeSnap>;

  @patch
  @route("/types/{id}")
  modifyProcType(@path id: string, @body body: ProcTypeSnap): Ok<ProcTypeSnap>;

  @post
  @route("/pools/decs")
  createPoolDecl(@body body: PoolDeclSpec): Created<TermRef>;

  @post
  @route("/pools/execs")
  createPoolExec(@body body: PoolExecSpec): Created<ComponentRef>;

  @post
  @route("/pools/execs/steps")
  takePoolStep(@body body: PoolStepSpec): NoContent;

  @post
  @route("/pools/execs/spawns")
  spawnPoolExec(@body body: PoolStepSpec): Created<ComponentRef>;
}
