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

/** Recursive expression payload; its discriminated variants follow the SDK wire model. */
model Expr {
  kind: string;
  ...Record<unknown>;
}

model ProcTypeSpec {
  type_qn: string;
  type_exp: Expr;
}

model ProcTypeSnap {
  ref: TypeRef;
  spec: ProcTypeSpec;
}

model PoolTypeSpec {
  qn: string;
  exp: Expr;
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
  proc_exp: Expr;
}

model PoolExecSpec {
  term_qn: string;
  liab_var: CompVarSpec;
  asset_vars: CompVarSpec[];
}

model PoolStepSpec {
  comp_ref: ComponentRef;
  pool_exp: Expr;
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
