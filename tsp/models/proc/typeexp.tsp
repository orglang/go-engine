import "../type.tsp";
import "../any/type.tsp";
namespace OrgLang.Engine.Proc;

/** Process type expression specification (go-sdk/proc/typeexp). */
model TypeExprSpec {
  kind: "one" | "link" | "tensor" | "lolli" | "plus" | "with" | "up" | "down";
  link?: Type.LinkSpec;
  tensor?: TypeProductSpec;
  lolli?: TypeProductSpec;
  plus?: TypeSumSpec;
  with?: TypeSumSpec;
  up?: TypeShiftSpec;
  down?: TypeShiftSpec;
}

model TypeProductSpec { val_es: TypeExprSpec; cont_es: TypeExprSpec; }
model TypeSumSpec { choices: TypeChoiceSpec[]; }
model TypeChoiceSpec { lab_qn: string; cont_es: TypeExprSpec; }
model TypeShiftSpec { cont_es: TypeExprSpec; }
