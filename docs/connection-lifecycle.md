# Connection lifecycle

A Connection is the aggregate root for one logical communication stream. Publication and Subscription are Acts owned by that Connection; they do not have independent repositories or persistence lifecycles. Execution remains a separate aggregate.

## Identity and sequence

- `conn_id` is the stable Connection identity.
- `conn_rn` is the next Act revision to allocate.
- An Act is persisted with the current `conn_rn`; the Connection revision is then incremented in the same transaction.
- `offset_nr` is the revision of the last Act consumed by the reader.
- A snapshot returns Acts with `conn_rn > offset_nr`, ordered by `conn_rn`.

## Consuming and appending Acts

Consuming an Act advances `offset_nr` to that Act's revision. The update compares both the snapshot's `conn_rn` and previous `offset_nr`; if either changed, the operation returns a concurrency conflict and must be retried from a fresh snapshot.

Appending an Act increments `conn_rn` with a compare-and-update before inserting the Act. The Act insertion, revision increment, offset update (when present), and related Execution modification run in the same caller-owned transaction. A failure rolls back the complete operation.

The process-side Connection row is created on first Act persistence. Pool-side Connections are created with their owning execution lifecycle.
