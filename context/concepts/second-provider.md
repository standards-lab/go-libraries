# Second provider: proving the dialect seam

Captured 2026-08-13 during the connectivity session, where the question was raised and deliberately
deferred. A dialect seam proven against one implementation can bake in that implementation's
assumptions; a second provider is the test with teeth. The timing is the decision recorded here.

At the connectivity slice the seam carries three members, two of them trivial and one the
identity, so a second provider built now would prove little beyond DSN composition and would be
reworked at every rung as the surface grows. The seam's real stress points arrive with the later
slices: placeholder rendering threaded through the query builder (rung 3), and error
classification into the base sentinels plus the optimistic-concurrency contract (rung 4).

The plan: an mssql provider (`database/mssql`, per the sub-module naming rule) as its own session
after rung 4, once the reference service has proven the surface and the first versions are
tagged. Until then, agnosticism is design-guarded: no base port default, provider-owned DSN
composition, the `Options` map for dialect-specific keys, sentinels declared in the base, and a
driverless base module.
