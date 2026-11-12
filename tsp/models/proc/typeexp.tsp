import "../common.tsp";

namespace OrgLang.Engine;

/** Process type expression specification (go-sdk/proc/typeexp). */
model ProcTypeExprSpec {
  kind: "one" | "link" | "tensor" | "lolli" | "plus" | "with" | "up" | "down";
  link?: TypeLinkSpec;
  tensor?: ProcTypeProductSpec;
  lolli?: ProcTypeProductSpec;
  plus?: ProcTypeSumSpec;
  with?: ProcTypeSumSpec;
  up?: ProcTypeShiftSpec;
  down?: ProcTypeShiftSpec;
}

model ProcTypeProductSpec { val_es: ProcTypeExprSpec; cont_es: ProcTypeExprSpec; }
model ProcTypeSumSpec { choices: ProcTypeChoiceSpec[]; }
model ProcTypeChoiceSpec { lab_qn: string; cont_es: ProcTypeExprSpec; }
model ProcTypeShiftSpec { cont_es: ProcTypeExprSpec; }
