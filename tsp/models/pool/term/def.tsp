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
  conn_ph: string;
  cont_expr: DefinitionSpec;
}
model DefAcceptSpec {
  conn_ph: string;
  cont_expr: DefinitionSpec;
}
model DefHireSpec {
  conn_ph: string;
  proc_qn: string;
  cont_expr: DefinitionSpec;
}
model DefApplySpec {
  conn_ph: string;
  proc_qn: string;
  cont_expr: DefinitionSpec;
}
model DefReleaseSpec {
  conn_ph: string;
}
model DefDetachSpec {
  conn_ph: string;
}
model DefFireSpec {
  conn_ph: string;
  proc_qn: string;
}
model DefQuitSpec {
  conn_ph: string;
  proc_qn: string;
}
model DefSpawnSpec {
  proc_qn: string;
  proc_comp_qns: string[];
}
