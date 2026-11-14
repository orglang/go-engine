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
@route("/api/v1")
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


  // Programs

  /** Submit program source text for parsing and registration. */
  @post
  @route("/programs")
op createProgram(@body source: string): CreatedNoBody;

  // Pool type definitions

  /** Create a pool-level type definition and return its identity. */
  @post
  @route("/pools/types")
op createPoolTypeDefinition(@body spec: Pool.Type.DefSpec): Created<Any.TypeRef>;

  /** Get a pool type definition snapshot; omit rn to select the latest revision. */
  @get
  @route("/pools/types/{id}")
op getPoolTypeDefinition(@path id: string, @query rn?: int64): Ok<Pool.Type.DefSnap>;

  // Process declarations and executions

  /** Create a process declaration and return its identity. */
  @post
  @route("/processes/declarations")
op createProcessDeclaration(@body spec: Proc.DecSpec): Created<Any.TermRef>;

  /** Get a process declaration by identity and revision. */
  @get
  @route("/processes/declarations/{id}")
op getProcessDeclaration(@path id: string, @query rn?: int64): Ok<Proc.DecSnap>;

  /** Get the declaration associated with a process. */
  @get
  @route("/processes/{id}/declaration")
op getProcessAssociatedDeclaration(@path id: string, @query rn?: int64): Ok<Proc.DecSnap>;

  /** List execution identities associated with a process. */
  @get
  @route("/processes/{id}/executions")
op listProcessExecutions(@path id: string): Ok<Any.CompRef[]>;

  /** Get a process execution; omit rn to select the latest available revision. */
  @get
  @route("/processes/executions/{id}")
op getProcessExecution(@path id: string, @query rn?: int64): Ok<Proc.CompSnap>;

  /** Apply one step to a process execution. */
  @post
  @route("/processes/executions/{id}/steps")
op stepProcessExecution(@path id: string, @body spec: Proc.StepSpec): OkNoBody;

  // Process type definitions

  /** Create a process type definition and return its identity. */
  @post
  @route("/processes/types")
op createProcessTypeDefinition(@body spec: Proc.Type.DefSpec): Created<Any.TypeRef>;

  /** Get a process type definition snapshot; omit rn to select the latest revision. */
  @get
  @route("/processes/types/{id}")
op getProcessTypeDefinition(@path id: string, @query rn?: int64): Ok<Proc.Type.DefSnap>;

  /** Get the term definition of a specific process; omit rn to select its latest revision. */
  @get
  @route("/processes/{id}/definition")
op getProcessDefinition(@path id: string, @query rn?: int64): Ok<Proc.Term.ExprSpec>;

  /** Extend a process term definition without replacing or removing existing expression parts. */
  @patch
  @route("/processes/{id}/definition")
op extendProcessDefinition(@path id: string, @body extension: Proc.Term.ExprSpec): Ok<Any.TermRef>;

  // Pool declarations and executions

  /** Create a pool declaration and return its identity. */
  @post
  @route("/pools/declarations")
op createPoolDeclaration(@body spec: Pool.DecSpec): Created<Any.TermRef>;

  /** Get the declaration associated with a pool. */
  @get
  @route("/pools/{id}/declaration")
op getPoolAssociatedDeclaration(@path id: string, @query rn?: int64): Ok<Pool.DecSnap>;

  /** Get the term definition of a specific pool; omit rn to select its latest revision. */
  @get
  @route("/pools/{id}/definition")
op getPoolDefinition(@path id: string, @query rn?: int64): Ok<Pool.Term.ExprSpec>;

  /** Extend a pool term definition without replacing or removing existing expression parts. */
  @patch
  @route("/pools/{id}/definition")
op extendPoolDefinition(@path id: string, @body extension: Pool.Term.ExprSpec): Ok<Any.TermRef>;

  /** List execution identities associated with a pool. */
  @get
  @route("/pools/{id}/executions")
op listPoolExecutions(@path id: string): Ok<Any.CompRef[]>;

  /** Create a pool execution and return its identity. */
  @post
  @route("/pools/executions")
op createPoolExecution(@body spec: Any.ExecSpec): Created<Any.CompRef>;

  /** Apply a batch of steps to pool executions. */
  @post
  @route("/pools/executions/steps")
op stepPoolExecution(@body spec: Pool.StepSpec): NoContent;

  /** Spawn a child execution from a pool execution and return its identity. */
  @post
  @route("/pools/executions/spawns")
op spawnPoolExecution(@body spec: Pool.StepSpec): Created<Any.CompRef>;
