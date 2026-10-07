# Future Work

Design-improvement backlog from a code review. Roughly ordered by leverage;
tackle one at a time. None are blocking — the app works today.

> Shipped (removed from this list): the Wails boundary typing and the HTTP-status-code result
> type — bound `App` methods now return concrete typed values + `error`; the alias-method
> cleanup — the bound API now uses one vocabulary (`Registration`, `Call`, `Waitlist`) with no
> aliases. See `CLAUDE.md` for the current contract. The semester-split policy was also
> extracted out of the CSV parser into `types.Course.SemesterForRanking`, with even seat counts
> enforced by the `types.Seats` value type. The call-entries rework made semesters fixed seeded
> data (removing the "explicit semester creation" and `FindOrCreateSemester` items) and moved the
> cycle state into `call_entries`, so every action is undoable.
>
> Decided against: auto-splitting odd seat counts (Semester 1 taking the extra seat). An odd
> seat count keeps failing the import with `ErrOddSeatsCount` on purpose.

No open items right now.
