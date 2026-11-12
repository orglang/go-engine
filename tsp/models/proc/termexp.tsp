namespace OrgLang.Engine;

/** Process term expression specification (go-sdk/proc/termexp). */
model ProcTermExprSpec {
  kind: "close" | "wait" | "send" | "recv" | "lab" | "case" | "link" | "fwd" | "call";
  close?: ProcCloseSpec;
  wait?: ProcWaitSpec;
  send?: ProcSendSpec;
  recv?: ProcRecvSpec;
  lab?: ProcLabSpec;
  case?: ProcCaseSpec;
  call?: ProcCallSpec;
  fwd?: ProcFwdSpec;
}

model ProcCloseSpec { comm_ph: string; }
model ProcWaitSpec { comm_ph: string; cont_es: ProcTermExprSpec; }
model ProcSendSpec { comm_ph: string; val_ph: string; }
model ProcRecvSpec { comm_ph: string; bind_ph: string; cont_es: ProcTermExprSpec; }
model ProcLabSpec { comm_ph: string; pattern_qn: string; }
model ProcCaseSpec { comm_ph: string; cont_bses: ProcBranchSpec[]; }
model ProcBranchSpec { pattern_qn: string; cont_es: ProcTermExprSpec; }
model ProcCallSpec {
  bind_ph: string;
  proc_qn: string;
  val_phs: string[];
  cont_es: ProcTermExprSpec;
}
model ProcFwdSpec { comm_ph: string; cont_ph: string; }
