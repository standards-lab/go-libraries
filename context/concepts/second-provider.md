# Second provider: proving the provider contracts

Captured 2026-08-13 during the connectivity session, where the question was raised and deliberately
deferred; revised 2026-08-17 in the capability-tiers planning session. A provider contract proven
against one implementation can bake in that implementation's assumptions; a second provider is the test
with teeth. That is what a second provider is for — it proves that the dialect interface, the
constructor contract, and the error mapping hold for more than one engine, and it serves a real
consumer — and it is never built to demonstrate optionality. It does not widen the standard tier: a
feature the second engine happens to share with the first stays native unless the SQL standard defines
it.

At the connectivity slice the contract carried three members, two of them trivial and one the identity,
so a second provider built then would have proven little beyond DSN composition and would have been
reworked at every rung as the surface grew. The stress points arrive with the later slices: placeholder
rendering threaded through the query builder (rung 3); error classification into the base sentinels
plus the optimistic-concurrency contract (rung 4); and paging, which the base emits in the standard
SQL:2008 form and which only an engine that lacks that form can stress. SQL Server accepts the standard
form, so it proves nothing there; SQLite and MySQL do not.

Two candidates, each its own session after the reference service's `v0.1.0`:

- **SQLite** (`database/sqlite`) — the cheapest second provider: in-process, no container, a
  near-stdlib driver. It is the one engine on the list that lacks standard-form paging, so it is the
  engine that makes the dialect interface grow, and it gives hermetic tests real SQL, which a consumer
  whose CI runs no database container would benefit from.
- **SQL Server** (`database/mssql`) — the enterprise engine, with `@p1` placeholders and `OUTPUT` in
  place of `RETURNING`; the provider to build when a real consumer needs it.

Until then, agnosticism is design-guarded: no base port default, provider-owned DSN composition, the
`Options` map as the native tier at configuration level, sentinels declared in the base, standard-form
paging emitted by the base, and a driverless base module.
