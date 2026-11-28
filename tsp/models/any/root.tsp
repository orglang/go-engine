namespace OrgLang.Any;

/** Shared root identity used by process and pool execution aggregates. */
model RootRef {
  root_id: string;
  root_rn: int64;
}
