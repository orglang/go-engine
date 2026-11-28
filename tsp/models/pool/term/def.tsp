namespace OrgLang.Pool.Term;

/** Declarative pool term Definition persisted independently from execution expressions. */
@discriminated(#{ envelope: "none" })
union DefinitionSpec {
  acquire: DefAcquireSpec,
  accept: DefAcceptSpec,
  hire: DefHireSpec,
  apply: DefApplySpec,
  release: DefReleaseSpec,
  detach: DefDetachSpec,
  fire: DefFireSpec,
  quit: DefQuitSpec,
  spawn: DefSpawnSpec,
}

model DefAcquireSpec {
  comm_ph: string;
  cont_expr: DefinitionSpec;
}
model DefAcceptSpec {
  comm_ph: string;
  cont_expr: DefinitionSpec;
}
model DefHireSpec {
  comm_ph: string;
  proc_qn: string;
  cont_expr: DefinitionSpec;
}
model DefApplySpec {
  comm_ph: string;
  proc_qn: string;
  cont_expr: DefinitionSpec;
}
model DefReleaseSpec {
  comm_ph: string;
}
model DefDetachSpec {
  comm_ph: string;
}
model DefFireSpec {
  comm_ph: string;
  proc_qn: string;
}
model DefQuitSpec {
  comm_ph: string;
  proc_qn: string;
}
model DefSpawnSpec {
  proc_qn: string;
  proc_comp_qns: string[];
}
