import "@typespec/http";
import "@typespec/openapi3";

import "./models/any/type.tsp";
import "./models/any/type/link.tsp";
import "./models/any/term.tsp";
import "./models/any/comp.tsp";
import "./models/any/conn.tsp";
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
  createProgram(@body body: string): CreatedNoBody;

  // Pool type definitions

  /** Create a pool-level type definition. */
  @post
  @route("/pools/types")
  createPoolTypeDefinition(@body body: Pool.Type.DefSpec): Created<Any.TypeRef>;

  // Process declarations and computations

  /** Create a process declaration. */
  @post
  @route("/processes/declarations")
  createProcessDeclaration(@body body: Proc.DecSpec): Created<Proc.DecSnap>;

  /** Get a process declaration; omit revision to select the latest available revision. */
  @get
  @route("/processes/declarations/{id}")
  getProcessDeclaration(@path id: string, @query revision?: int64): Ok<Proc.DecSnap>;

  /** Get a process computation; omit revision to select the latest available revision. */
  @get
  @route("/processes/computations/{id}")
  getProcessComputation(@path id: string, @query revision?: int64): Ok<Proc.CompSnap>;

  /** Apply one step to a process computation. */
  @post
  @route("/processes/computations/{id}/steps")
  stepProcessComputation(@path id: string, @body body: Proc.StepSpec): OkNoBody;

  // Process type definitions

  /** Create a process type definition. */
  @post
  @route("/processes/types")
  createProcessTypeDefinition(@body body: Proc.Type.DefSpec): Created<Proc.Type.DefSnap>;

  /** Get a process type definition; omit revision to select the latest available revision. */
  @get
  @route("/processes/types/{id}")
  getProcessTypeDefinition(@path id: string, @query revision?: int64): Ok<Proc.Type.DefSnap>;

  /** Update a process type definition. */
  @put
  @route("/processes/types/{id}")
  updateProcessTypeDefinition(@path id: string, @body body: Proc.Type.DefSnap): Ok<Proc.Type.DefSnap>;

  // Pool declarations and executions

  /** Create a pool declaration. */
  @post
  @route("/pools/declarations")
  createPoolDeclaration(@body body: Pool.DecSpec): Created<Any.TermRef>;

  /** Create a pool execution. */
  @post
  @route("/pools/executions")
  createPoolExecution(@body body: Any.ExecSpec): Created<Any.CompRef>;

  /** Apply one step to a pool execution. */
  @post
  @route("/pools/executions/steps")
  stepPoolExecution(@body body: Pool.StepSpec): NoContent;

  /** Spawn a child execution from a pool execution. */
  @post
  @route("/pools/executions/spawns")
  spawnPoolExecution(@body body: Pool.StepSpec): Created<Any.CompRef>;
}
