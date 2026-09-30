# Future Work

Design-improvement backlog from a code review. Roughly ordered by leverage;
tackle one at a time. None are blocking — the app works today.

> Shipped (removed from this list): the Wails boundary typing and the HTTP-status-code result
> type — bound `App` methods now return concrete typed values + `error`; the alias-method
> cleanup — the bound API now uses one vocabulary (`Registration`, `Call`, `Waitlist`) with no
> aliases; and the call-entries rework — semesters are fixed seeded data (the "explicit semester
> creation" and `FindOrCreateSemester` items are gone with it), and the cycle state lives in
> `call_entries` so every action is undoable. See `CLAUDE.md` for the current contract.

## 1. Handle odd seat counts instead of erroring

- **Problem:** Approved import rejects any course whose seat count is odd (`ErrOddSeatsCount`). Seat counts come from the SiSU CSV — data the operator doesn't control — so an odd count leaves them stuck with no recourse.
- **Where:** `csvparser/mapper.go:134` (raises it), `csvparser/errors.go:121`.
- **Fix:** Replace the hard error with a deterministic rule, e.g. on an odd count Semester 1 gets the extra seat (ceil for sem 1, floor for sem 2). Allocation already uses `types.SemesterSeats` (ceil for sem 1); the parser split is the one place left. Becomes a one-line change once item 2 is done. (Confirm the desired tie-breaking with the domain.)
- **Effort:** Small.

## 2. Extract the semester-split policy out of the CSV parser

- **Problem:** The institutional admissions rule — one annual SiSU selection feeds two semester intakes, top-ranked half → Sem 1, rest → Sem 2 — lives *inside the CSV parser*. The parser does two unrelated jobs: deserialize CSV rows **and** apply admissions policy. They change for different reasons (CSV format vs intake rule), the policy is hard to find, and it can't be unit-tested in isolation.
- **Where:** `csvparser/mapper.go` — `csvCandidate.Parse` enforces even seats and assigns `Semester` per-row by ranking against the course's seat count.
- **Fix:** Make the parser policy-free (faithful ranked registrations, no `Semester`). Add one named pure function that owns the rule, e.g.:
  ```go
  // SplitApprovedBySemester divides a course's ranked approved candidates across
  // the year's two semester intakes (SiSU runs one annual selection; universities
  // keep two intakes). `ranked` must be sorted by ranking ascending (best first).
  func SplitApprovedBySemester(ranked []*types.Registration) (sem1, sem2 []*types.Registration)
  ```
  `commands/load_selection.go` calls it when building the call-1 entries. Pure function → trivial table-driven tests (even / odd / single seat / empty). Gives item 1 exactly one home. Full Strategy-pattern (interface + impls) is **not** needed unless the rule ever becomes configurable per institution/year — YAGNI for now.
- **Effort:** Medium.
