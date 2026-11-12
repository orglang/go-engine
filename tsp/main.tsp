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
  createPoolTypeDef(@body body: Pool.Type.DefSpec): Created<Any.TypeRef>;

  @post
  @route("/procs/decs")
  createProcDecl(@body body: Proc.DecSpec): Created<Proc.DecSnap>;

  /** Omit revision to request the latest available declaration revision. */
  @get
  @route("/procs/decs/{id}")
  getProcDecl(@path id: string, @query revision?: int64): Ok<Proc.DecSnap>;

  /** Omit revision to request the latest available process computation revision. */
  @get
  @route("/procs/{id}")
  getProc(@path id: string, @query revision?: int64): Ok<Proc.CompSnap>;

  @post
  @route("/procs/{id}/steps")
  takeProcStep(@path id: string, @body body: Proc.StepSpec): OkNoBody;

  @post
  @route("/types")
  createProcType(@body body: Proc.Type.DefSpec): Created<Proc.Type.DefSnap>;

  /** Omit revision to request the latest available type revision. */
  @get
  @route("/types/{id}")
  getProcType(@path id: string, @query revision?: int64): Ok<Proc.Type.DefSnap>;

  @patch
  @route("/types/{id}")
  modifyProcType(@path id: string, @body body: Proc.Type.DefSnap): Ok<Proc.Type.DefSnap>;

  @post
  @route("/pools/decs")
  createPoolDecl(@body body: Pool.DecSpec): Created<Any.TermRef>;

  @post
  @route("/pools/execs")
  createPoolExec(@body body: Any.ExecSpec): Created<Any.CompRef>;

  @post
  @route("/pools/execs/steps")
  takePoolStep(@body body: Pool.StepSpec): NoContent;

  @post
  @route("/pools/execs/spawns")
  spawnPoolExec(@body body: Pool.StepSpec): Created<Any.CompRef>;
}
