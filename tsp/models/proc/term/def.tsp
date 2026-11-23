namespace OrgLang.Proc.Term;

/** Declarative process term Definition persisted independently from execution expressions. */
@discriminated(#{ envelope: "none" })
union DefinitionSpec {
  close: DefCloseSpec,
  wait: DefWaitSpec,
  send: DefSendSpec,
  recv: DefRecvSpec,
  lab: DefLabSpec,
  case: DefCaseSpec,
  link: DefLinkSpec,
  fwd: DefFwdSpec,
  call: DefCallSpec,
  spawn: DefSpawnSpec,
  acquire: DefAcquireSpec,
  accept: DefAcceptSpec,
  detach: DefDetachSpec,
  release: DefReleaseSpec,
}

model DefCloseSpec {
  cont_ph: string;
}
model DefWaitSpec {
  cont_ph: string;
  cont_expr: DefinitionSpec;
}
model DefSendSpec {
  conn_ph: string;
  val_ph: string;
}
model DefRecvSpec {
  conn_ph: string;
  new_ph: string;
  cont_expr: DefinitionSpec;
}
model DefLabSpec {
  conn_ph: string;
  val_lab_qn: string;
  cont_expr: DefinitionSpec;
}
model DefCaseSpec {
  conn_ph: string;
  cont_exps: DefBranchSpec[];
}
model DefBranchSpec {
  label_qn: string;
  cont_expr: DefinitionSpec;
}
model DefLinkSpec {
  proc_qn: string;
  x: string;
  ys: string[];
}
model DefFwdSpec {
  conn_ph: string;
  cont_ph: string;
}
model DefCallSpec {
  conn_ph: string;
  new_ph: string;
  proc_qn: string;
  val_phs: string[];
  cont_expr: DefinitionSpec;
}
model DefSpawnSpec {
  conn_ph: string;
  proc_qn: string;
  new_phs: string[];
  cont_expr: DefinitionSpec;
}
model DefAcquireSpec {
  conn_ph: string;
  cont_expr: DefinitionSpec;
}
model DefAcceptSpec {
  conn_ph: string;
  cont_expr: DefinitionSpec;
}
model DefDetachSpec {
  conn_ph: string;
}
model DefReleaseSpec {
  conn_ph: string;
}
