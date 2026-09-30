# Future Work

Design-improvement backlog from a code review. Roughly ordered by leverage;
tackle one at a time. None are blocking — the app works today.

> Shipped (removed from this list): the Wails boundary typing and the HTTP-status-code result
> type — bound `App` methods now return concrete typed values + `error`; and the alias-method
> cleanup — the bound API now uses one vocabulary (`Registration`, `Call`, `Waitlist`) with no
> aliases. See `CLAUDE.md` for the current contract. The semester-split policy was also
> extracted out of the CSV parser into `commands.SplitApprovedBySemester`.
>
> Decided against: auto-splitting odd seat counts (Semester 1 taking the extra seat). An odd
> seat count keeps failing the approved import with `ErrOddSeatsCount` on purpose.

## 1. (Open question) Make semester creation explicit, not a side effect

- **Problem:** Both semesters are auto-created whenever an approved selection is imported. This couples "I have an approved CSV" with "this year has exactly two intakes." Works for the current institution; the model would lie if a year ever had a single intake.
- **Where:** `commands/load_selection.go` (find-or-create sem 1 & 2, lines 55-83).
- **Fix:** Mostly a design question to revisit if single-intake years ever happen. Option: create semesters from an explicit action rather than as an import side effect. No action needed unless the requirement appears.
- **Effort:** N/A (decision, not a task).

## 2. Small cleanup: `FindOrCreateSemester` helper

- **Problem:** `commands/load_selection.go` lines 55-83 copy-pastes the same ~15-line find-or-create block for Semester 1 and Semester 2.
- **Where:** `commands/load_selection.go:55-83`.
- **Fix:** Extract `FindOrCreateSemester(tx qrm.DB, year, number int32) (int32, error)` and call it twice. Halves the block.
- **Effort:** Trivial.
