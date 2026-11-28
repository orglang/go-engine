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
  comm_ph: string;
  val_ph: string;
}
model DefRecvSpec {
  comm_ph: string;
  new_ph: string;
  cont_expr: DefinitionSpec;
}
model DefLabSpec {
  comm_ph: string;
  val_lab_qn: string;
  cont_expr: DefinitionSpec;
}
model DefCaseSpec {
  comm_ph: string;
  cont_exps: Record<DefinitionSpec>;
}
model DefLinkSpec {
  proc_qn: string;
  x: string;
  ys: string[];
}
model DefFwdSpec {
  comm_ph: string;
  cont_ph: string;
}
model DefCallSpec {
  comm_ph: string;
  new_ph: string;
  proc_qn: string;
  val_phs: string[];
  cont_expr: DefinitionSpec;
}
model DefSpawnSpec {
  comm_ph: string;
  proc_qn: string;
  new_phs: string[];
  cont_expr: DefinitionSpec;
}
model DefAcquireSpec {
  comm_ph: string;
  cont_expr: DefinitionSpec;
}
model DefAcceptSpec {
  comm_ph: string;
  cont_expr: DefinitionSpec;
}
model DefDetachSpec {
  comm_ph: string;
}
model DefReleaseSpec {
  comm_ph: string;
}
