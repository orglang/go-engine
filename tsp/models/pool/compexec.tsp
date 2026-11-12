import "../comp.tsp";
namespace OrgLang.Engine.Pool;

model ExecSpec {
  term_qn: string;
  liab_var: Comp.VarSpec;
  asset_vars: Comp.VarSpec[];
}

model Snap {
  comp_ref: Comp.Ref;
}
