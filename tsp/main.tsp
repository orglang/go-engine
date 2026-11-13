import "@typespec/http";
import "@typespec/openapi3";

import "./models/any/type.tsp";
import "./models/any/type/link.tsp";
import "./models/any/term.tsp";
import "./models/any/comp.tsp";
import "./models/any/conn/var.tsp";
import "./models/any/exec.tsp";
import "./models/proc/type/expr.tsp";
import "./models/proc/type/def.tsp";
import "./models/proc/term/expr.tsp";
import "./models/proc/term/dec.tsp";
import "./models/proc/comp.tsp";
import "./models/pool/type/expr.tsp";
import "./models/pool/type/def.tsp";
import "./models/pool/term/expr.tsp";
import "./models/pool/term/dec.tsp";
import "./models/pool/comp.tsp";

using TypeSpec.Http;

@service(#{ title: "OrgLang Go Engine API" })
@server("http://localhost:8080", "Local development server")
namespace OrgLang;

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

model OkNoBody {
  @statusCode statusCode: 200;
}

model NoContent {
  @statusCode statusCode: 204;
}

@route("/api/v1")
interface EngineApi {
  // Programs

  /** Submit program source text for parsing and registration. */
  @post
  @route("/programs")
  createProgram(@body source: string): CreatedNoBody;

  // Pool type definitions

  /** Create a pool-level type definition and return its identity. */
  @post
  @route("/pools/types")
  createPoolTypeDefinition(@body spec: Pool.Type.DefSpec): Created<Any.TypeRef>;

  /** Get a pool type definition snapshot; omit rn to select the latest revision. */
  @get
  @route("/pools/types/{id}")
  getPoolTypeDefinition(@path id: string, @query rn?: int64): Ok<Pool.Type.DefSnap>;

  /** Get the expression of a pool type definition. */
  @get
  @route("/pools/types/{id}/definition")
  getPoolTypeDefinitionExpression(@path id: string, @query rn?: int64): Ok<Pool.Type.ExprSpec>;

  /** Extend a pool type definition without replacing or removing existing expression parts. */
  @patch
  @route("/pools/types/{id}/definition")
  extendPoolTypeDefinition(@path id: string, @body patch: Pool.Type.DefPatch): Ok<Any.TypeRef>;

  // Process declarations and executions

  /** Create a process declaration and return its identity. */
  @post
  @route("/processes/declarations")
  createProcessDeclaration(@body spec: Proc.DecSpec): Created<Any.TermRef>;

  /** Get a process declaration by identity and revision. */
  @get
  @route("/processes/declarations/{id}")
  getProcessDeclaration(@path id: string, @query rn?: int64): Ok<Proc.DecSnap>;

  /** Get the declaration associated with a process. */
  @get
  @route("/processes/{id}/declaration")
  getProcessAssociatedDeclaration(@path id: string, @query rn?: int64): Ok<Proc.DecSnap>;

  /** List execution identities associated with a process. */
  @get
  @route("/processes/{id}/executions")
  listProcessExecutions(@path id: string): Ok<Any.CompRef[]>;

  /** Get a process execution; omit rn to select the latest available revision. */
  @get
  @route("/processes/executions/{id}")
  getProcessExecution(@path id: string, @query rn?: int64): Ok<Proc.CompSnap>;

  /** Apply one step to a process execution. */
  @post
  @route("/processes/executions/{id}/steps")
  stepProcessExecution(@path id: string, @body spec: Proc.StepSpec): OkNoBody;

  // Process type definitions

  /** Create a process type definition and return its identity. */
  @post
  @route("/processes/types")
  createProcessTypeDefinition(@body spec: Proc.Type.DefSpec): Created<Any.TypeRef>;

  /** Get a process type definition snapshot; omit rn to select the latest revision. */
  @get
  @route("/processes/types/{id}")
  getProcessTypeDefinition(@path id: string, @query rn?: int64): Ok<Proc.Type.DefSnap>;

  /** Get the expression of a process type definition. */
  @get
  @route("/processes/types/{id}/definition")
  getProcessTypeDefinitionExpression(@path id: string, @query rn?: int64): Ok<Proc.Type.ExprSpec>;

  /** Extend a process type definition without replacing or removing existing expression parts. */
  @patch
  @route("/processes/types/{id}/definition")
  extendProcessTypeDefinition(@path id: string, @body patch: Proc.Type.DefPatch): Ok<Any.TypeRef>;

  // Pool declarations and executions

  /** Create a pool declaration and return its identity. */
  @post
  @route("/pools/declarations")
  createPoolDeclaration(@body spec: Pool.DecSpec): Created<Any.TermRef>;

  /** Get the declaration associated with a pool. */
  @get
  @route("/pools/{id}/declaration")
  getPoolAssociatedDeclaration(@path id: string, @query rn?: int64): Ok<Pool.DecSnap>;

  /** List execution identities associated with a pool. */
  @get
  @route("/pools/{id}/executions")
  listPoolExecutions(@path id: string): Ok<Any.CompRef[]>;

  /** Create a pool execution and return its identity. */
  @post
  @route("/pools/executions")
  createPoolExecution(@body spec: Any.ExecSpec): Created<Any.CompRef>;

  /** Apply a batch of steps to pool executions. */
  @post
  @route("/pools/executions/steps")
  stepPoolExecution(@body spec: Pool.StepSpec): NoContent;

  /** Spawn a child execution from a pool execution and return its identity. */
  @post
  @route("/pools/executions/spawns")
  spawnPoolExecution(@body spec: Pool.StepSpec): Created<Any.CompRef>;
}
