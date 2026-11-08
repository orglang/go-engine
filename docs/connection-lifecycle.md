# Connection lifecycle

A Computation is the aggregate root. A Connection is a nested entity that links the provider and client Computations; it is not an independent aggregate root. Publication and Subscription are Acts owned by the Connection.

## Identity and creation

- `conn_id` is the stable, globally unique Connection identity.
- When a provider Computation is created, it reserves a `conn_id` in its liability variable. No Connection row is created at this point.
- When a client Computation is created, it obtains the provider's reserved `conn_id` from the referenced liability variable. The Connection row is created in the same transaction as the client Computation and its variables.
- A Connection therefore enters persistence only after both endpoints exist. Both sides use the same `conn_id`.
- A reserved `conn_id` without a Connection row is not a usable communication stream. Operations that read or append Acts require the Connection row to exist.
- `conn_rn` is the next Act revision to allocate.
- An Act is persisted with the current `conn_rn`; the Connection revision is then incremented in the same transaction.
- `offset_nr` is the revision of the last Act consumed by the reader.
- A snapshot returns Acts with `conn_rn > offset_nr`, ordered by `conn_rn`.

## Interaction Channels

Each interaction retains its own physical `chnl_id`. The first Act establishes that channel and its endpoint association; the complementary Publication or Subscription is matched against the same Connection and channel. Acts from different channels are not interchangeable.

## Consuming and appending Acts

Consuming an Act advances `offset_nr` to that Act's revision. The update compares both the snapshot's `conn_rn` and previous `offset_nr`; if either changed, the operation returns a concurrency conflict and must be retried from a fresh snapshot.

Appending an Act increments `conn_rn` with a compare-and-update before inserting the Act. The Act insertion, revision increment, offset update (when present), and related Execution modification run in the same caller-owned transaction. A failure rolls back the complete operation.
