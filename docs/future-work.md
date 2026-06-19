# Future Work

Design-improvement backlog from a code review. Roughly ordered by leverage;
tackle one at a time. None are blocking — the app works today.

> The Wails boundary typing and the HTTP-status-code result type (formerly items 1 and 5
> here) have both shipped — bound `App` methods now return concrete typed values + `error`.
> See `CLAUDE.md` for the current contract.

## 1. Collapse the alias methods and settle on one vocabulary

- **Problem:** Multiple bound names for one operation: `DeleteCall` / `DeleteRollCall` / `DeleteRollcall`; `EnrollRegistration` / `EnrollApplication`; `AbsentRegistration` / `AbsentApplication`; `FetchRegistrationsByCallID` / `FetchApplicationsByRollCall`. This is backend/frontend vocabulary drift ("Registration" vs "Application", "Call" vs "RollCall") papered over with aliases. Bloats the API surface; each alias can drift.
- **Where:** `app.go` (alias methods around lines 351-411).
- **Fix:** Pick one term for each concept, rename on both sides, delete the aliases.
- **Effort:** Medium.

## 2. Handle odd seat counts instead of erroring

- **Problem:** Approved import rejects any course whose seat count is odd (`ErrOddSeatsCount`). Seat counts come from the SiSU CSV — data the operator doesn't control — so an odd count leaves them stuck with no recourse.
- **Where:** `csvparser/mapper.go:134` (raises it), `csvparser/errors.go:121`.
- **Fix:** Replace the hard error with a deterministic rule, e.g. on an odd count Semester 1 gets the extra seat (ceil for sem 1, floor for sem 2). Becomes a one-line change once item 3 is done. (Confirm the desired tie-breaking with the domain.)
- **Effort:** Small.

## 3. Extract the semester-split policy out of the CSV parser

- **Problem:** The institutional admissions rule — one annual SiSU selection feeds two semester intakes, top-ranked half → Sem 1, rest → Sem 2 — lives *inside the CSV parser*. The parser does two unrelated jobs: deserialize CSV rows **and** apply admissions policy. They change for different reasons (CSV format vs intake rule), the policy is hard to find, and it can't be unit-tested in isolation.
- **Where:** `csvparser/mapper.go` — `csvCandidate.Parse` (line 106) enforces even seats (line 134) and assigns `SemesterID` per-row by ranking against the course's seat count (lines 178-185).
- **Fix:** Make the parser policy-free (faithful ranked registrations, no `SemesterID`). Add one named pure function that owns the rule, e.g.:
  ```go
  // SplitApprovedBySemester divides a course's ranked approved candidates across
  // the year's two semester intakes (SiSU runs one annual selection; universities
  // keep two intakes). `ranked` must be sorted by ranking ascending (best first).
  func SplitApprovedBySemester(ranked []*types.Registration) (sem1, sem2 []*types.Registration)
  ```
  `commands/load_selection.go` calls it, then maps results to real semester IDs (the loop it already has). Pure function → trivial table-driven tests (even / odd / single seat / empty). Gives item 2 exactly one home. Full Strategy-pattern (interface + impls) is **not** needed unless the rule ever becomes configurable per institution/year — YAGNI for now.
- **Effort:** Medium.

## 4. (Open question) Make semester creation explicit, not a side effect

- **Problem:** Both semesters are auto-created whenever an approved selection is imported. This couples "I have an approved CSV" with "this year has exactly two intakes." Works for the current institution; the model would lie if a year ever had a single intake.
- **Where:** `commands/load_selection.go` (find-or-create sem 1 & 2, lines 47-75).
- **Fix:** Mostly a design question to revisit if single-intake years ever happen. Option: create semesters from an explicit action rather than as an import side effect. No action needed unless the requirement appears.
- **Effort:** N/A (decision, not a task).

## 5. Small cleanup: `FindOrCreateSemester` helper

- **Problem:** `commands/load_selection.go` lines 47-75 copy-pastes the same ~15-line find-or-create block for Semester 1 and Semester 2.
- **Where:** `commands/load_selection.go:47-75`.
- **Fix:** Extract `FindOrCreateSemester(tx qrm.DB, year, number int32) (int32, error)` and call it twice. Halves the block.
- **Effort:** Trivial.
