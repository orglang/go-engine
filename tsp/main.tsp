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

  /** Create a pool-level type definition. */
  @post
  @route("/pools/types")
  createPoolTypeDefinition(@body spec: Pool.Type.DefSpec): Created<Any.TypeRef>;

  /** Get a pool type definition; omit rn to select the latest available revision. */
  @get
  @route("/pools/types/{id}")
  getPoolTypeDefinition(@path id: string, @query rn?: int64): Ok<Pool.Type.DefSnap>;

  /** Replace a pool type definition with a complete snapshot. */
  @put
  @route("/pools/types/{id}")
  updatePoolTypeDefinition(@path id: string, @body snap: Pool.Type.DefSnap): Ok<Pool.Type.DefSnap>;

  // Process declarations and executions

  /** Create a process declaration. */
  @post
  @route("/processes/declarations")
  createProcessDeclaration(@body spec: Proc.DecSpec): Created<Proc.DecSnap>;

  /** Get a process declaration; omit rn to select the latest available revision. */
  @get
  @route("/processes/declarations/{id}")
  getProcessDeclaration(@path id: string, @query rn?: int64): Ok<Proc.DecSnap>;

  /** List declarations associated with a process. */
  @get
  @route("/processes/{id}/declarations")
  listProcessDeclarations(@path id: string): Ok<Proc.DecSnap[]>;

  /** Get a process execution; omit rn to select the latest available revision. */
  @get
  @route("/processes/executions/{id}")
  getProcessExecution(@path id: string, @query rn?: int64): Ok<Proc.CompSnap>;

  /** Apply one step to a process execution. */
  @post
  @route("/processes/executions/{id}/steps")
  stepProcessExecution(@path id: string, @body spec: Proc.StepSpec): OkNoBody;

  // Process type definitions

  /** Create a process type definition. */
  @post
  @route("/processes/types")
  createProcessTypeDefinition(@body spec: Proc.Type.DefSpec): Created<Proc.Type.DefSnap>;

  /** Get a process type definition; omit rn to select the latest available revision. */
  @get
  @route("/processes/types/{id}")
  getProcessTypeDefinition(@path id: string, @query rn?: int64): Ok<Proc.Type.DefSnap>;

  /** Replace a process type definition with a complete snapshot. */
  @put
  @route("/processes/types/{id}")
  updateProcessTypeDefinition(@path id: string, @body snap: Proc.Type.DefSnap): Ok<Proc.Type.DefSnap>;

  // Pool declarations and executions

  /** Create a pool declaration. */
  @post
  @route("/pools/declarations")
  createPoolDeclaration(@body spec: Pool.DecSpec): Created<Any.TermRef>;

  /** Create a pool execution. */
  @post
  @route("/pools/executions")
  createPoolExecution(@body spec: Any.ExecSpec): Created<Any.CompRef>;

  /** Apply a batch of steps to pool executions. */
  @post
  @route("/pools/executions/steps")
  stepPoolExecution(@body spec: Pool.StepSpec): NoContent;

  /** Spawn a child execution from a pool execution. */
  @post
  @route("/pools/executions/spawns")
  spawnPoolExecution(@body spec: Pool.StepSpec): Created<Any.CompRef>;
}
