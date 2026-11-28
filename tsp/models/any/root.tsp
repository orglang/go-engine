namespace OrgLang.Any;

/** Shared root identity used by process and pool execution aggregates. */
model RootRef {
  exec_id: string;
  exec_rn: int64;
}
