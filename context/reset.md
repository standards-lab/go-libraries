# reset · data-cqrs-roadmap

- **Status:** closeout
- **Session:** plan
- **Branch:** data-cqrs-roadmap

## Disposition

- **Captured:** `concepts/database.md` — the settled direction and open questions for the
  database capability, the postgres and migrate sub-modules, and the web additions, with the four
  library build slices. This repository leaves its resting state; the data composition and CQRS
  layer planned with the reference service brings the next work here.
- **Retained:** `concepts/module-set.md` unchanged; its reserved database section and open
  questions (the members of `database.DB`, the query vocabulary) now carry candidate answers in
  `concepts/database.md` and are settled per build session.

## Next-focus

Rung 1 of the build ladder, a coordinated session with the reference service: the base `database`
package skeleton (wrapper, dialect seam, config block, lifecycle and readiness) and
`database/postgres` open and ping. Plan the API in plan mode from `concepts/database.md`.
