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

model TermVariableSpec {
  chnl_ph: string;
  type_qn: string;
}

model ComponentVariableSpec {
  chnl_ph: string;
  term_qn: string;
}

model TermVariableRecord {
  type_ref: TypeRef;
  chnl_ph: string;
  exp_vk: int64;
}

/** Recursive expression payload; its discriminated variants follow the SDK wire model. */
model Expression {
  kind: string;
  ...Record<unknown>;
}

model ProcessTypeSpec {
  type_qn: string;
  type_exp: Expression;
}

model ProcessTypeSnapshot {
  ref: TypeRef;
  spec: ProcessTypeSpec;
}

model PoolTypeSpec {
  qn: string;
  exp: Expression;
}

model PoolTypeSnapshot {
  ref: TypeRef;
  exp: Expression;
}

model ProcessDeclarationSpec {
  term_qn: string;
  liab_var: TermVariableSpec;
  asset_vars: TermVariableSpec[];
}

model ProcessDeclarationSnapshot {
  term_ref: TermRef;
  liab_var: TermVariableRecord;
  asset_vars: TermVariableRecord[];
}

model PoolDeclarationSpec {
  desc_qn: string;
  liab_var: TermVariableSpec;
  asset_vars: TermVariableSpec[];
}

model ProcessStepSpec {
  comp_ref: ComponentRef;
  proc_exp: Expression;
}

model PoolExecutionSpec {
  term_qn: string;
  liab_var: ComponentVariableSpec;
  asset_vars: ComponentVariableSpec[];
}

model PoolStepSpec {
  comp_ref: ComponentRef;
  pool_exp: Expression;
}

model ComponentSnapshot {
  comp_ref: ComponentRef;
}

model Created<T> {
  @statusCode statusCode: 201;
  @body body: T;
}

model Ok<T> {
  @statusCode statusCode: 200;
  @body body: T;
}

model NoContent {
  @statusCode statusCode: 204;
}

@route("/api/v1")
interface EngineApi {
  /** Submit program source text. */
  @post
  @route("/progs")
  @statusCode statusCode: 201;
  createProgram(@body body: string): NoContent;

  @post
  @route("/xacts/defs")
  createPoolTypeDefinition(@body body: PoolTypeSpec): Created<TypeRef>;

  @post
  @route("/procs/decs")
  createProcessDeclaration(@body body: ProcessDeclarationSpec): Created<ProcessDeclarationSnapshot>;

  /** Omit revision to request the latest available declaration revision. */
  @get
  @route("/procs/decs/{id}")
  getProcessDeclaration(
    @path id: string,
    @query revision?: int64
  ): Ok<ProcessDeclarationSnapshot>;

  /** Omit revision to request the latest available process revision. */
  @get
  @route("/procs/{id}")
  getProcess(
    @path id: string,
    @query revision?: int64
  ): Ok<ComponentSnapshot>;

  @post
  @route("/procs/{id}/steps")
  takeProcessStep(@path id: string, @body body: ProcessStepSpec): NoContent;

  @post
  @route("/types")
  createProcessType(@body body: ProcessTypeSpec): Created<ProcessTypeSnapshot>;

  /** Omit revision to request the latest available type revision. */
  @get
  @route("/types/{id}")
  getProcessType(
    @path id: string,
    @query revision?: int64
  ): Ok<ProcessTypeSnapshot>;

  @patch
  @route("/types/{id}")
  modifyProcessType(@path id: string, @body body: ProcessTypeSnapshot): Ok<ProcessTypeSnapshot>;

  @post
  @route("/pools/decs")
  createPoolDeclaration(@body body: PoolDeclarationSpec): Created<TermRef>;

  @post
  @route("/pools/execs")
  createPoolExecution(@body body: PoolExecutionSpec): Created<ComponentRef>;

  @post
  @route("/pools/execs/steps")
  takePoolStep(@body body: PoolStepSpec): NoContent;

  @post
  @route("/pools/execs/spawns")
  spawnPoolExecution(@body body: PoolStepSpec): Created<ComponentRef>;
}
