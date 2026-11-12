import "@typespec/http";
import "@typespec/openapi3";

import "./models/type.tsp";
import "./models/any/type.tsp";
import "./models/any/term.tsp";
import "./models/any/comp.tsp";
import "./models/term.tsp";
import "./models/proc/typeexp.tsp";
import "./models/proc/termexp.tsp";
import "./models/proc/typedef.tsp";
import "./models/proc/termdec.tsp";
import "./models/proc/compexec.tsp";
import "./models/proc/compstep.tsp";
import "./models/pool/typeexp.tsp";
import "./models/pool/termexp.tsp";
import "./models/pool/typedef.tsp";
import "./models/pool/termdec.tsp";
import "./models/pool/compexec.tsp";
import "./models/pool/compstep.tsp";

using TypeSpec.Http;

@service(#{ title: "OrgLang Go Engine API" })
@server("http://localhost:8080", "Local development server")
namespace OrgLang.Engine;

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
  createPoolTypeDef(@body body: Pool.TypeSpec): Created<Any.Type.Ref>;

  @post
  @route("/procs/decs")
  createProcDecl(@body body: Proc.DeclSpec): Created<Proc.DeclSnap>;

  /** Omit revision to request the latest available declaration revision. */
  @get
  @route("/procs/decs/{id}")
  getProcDecl(@path id: string, @query revision?: int64): Ok<Proc.DeclSnap>;

  /** Omit revision to request the latest available process computation revision. */
  @get
  @route("/procs/{id}")
  getProc(@path id: string, @query revision?: int64): Ok<Proc.CompSnap>;

  @post
  @route("/procs/{id}/steps")
  takeProcStep(@path id: string, @body body: Proc.StepSpec): OkNoBody;

  @post
  @route("/types")
  createProcType(@body body: Proc.TypeSpec): Created<Proc.TypeSnap>;

  /** Omit revision to request the latest available type revision. */
  @get
  @route("/types/{id}")
  getProcType(@path id: string, @query revision?: int64): Ok<Proc.TypeSnap>;

  @patch
  @route("/types/{id}")
  modifyProcType(@path id: string, @body body: Proc.TypeSnap): Ok<Proc.TypeSnap>;

  @post
  @route("/pools/decs")
  createPoolDecl(@body body: Pool.DeclSpec): Created<Any.Term.Ref>;

  @post
  @route("/pools/execs")
  createPoolExec(@body body: Pool.ExecSpec): Created<Any.Comp.Ref>;

  @post
  @route("/pools/execs/steps")
  takePoolStep(@body body: Pool.StepSpec): NoContent;

  @post
  @route("/pools/execs/spawns")
  spawnPoolExec(@body body: Pool.StepSpec): Created<Any.Comp.Ref>;
}
