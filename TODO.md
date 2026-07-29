# Prerequisites for the Architecture Refactor

Complete every item in this file before starting `doc/architecture-refactor-plan.md`. The implementation plan assumes these decisions are settled and recorded; agents must not invent answers while implementing it.

## Repository preparation

- [ ] Merge the frontend submodule into this repository as a normal directory and remove the submodule relationship. There is no practical reason for a single-developer project to keep the frontend in a separate repository, and the split makes cross-stack changes and agent work unnecessarily difficult.
- [ ] Preserve the frontend repository's useful history or record the commit from which it was imported.
- [ ] Confirm that a fresh clone contains the frontend source without requiring submodule initialization.
- [ ] Confirm that the root project can build the embedded frontend after the merge.
- [ ] Move the useful project-specific test commands and conventions into the `Commands` and `Rules` sections of `AGENTS.md`.
- [ ] Delete or replace `doc/testing.md` after migrating its small amount of project-specific guidance.
- [ ] Delete `doc/future-work.md` after migrating the unresolved odd-seat decision into the workflow specification.
- [ ] Delete `doc/single-annual-two-semester-gaps.md`; its interactive-promotion design is superseded.
- [ ] Delete or replace `frontend/docs/single-annual-two-semester-gaps.md` when merging the frontend; it describes the same superseded design.
- [ ] Review `frontend/docs/aprovados-page-design.md` and either keep it as an active UI contract, migrate its still-valid requirements, or delete it.

## Authoritative workflow specification

- [ ] Create `doc/admission-workflow-spec.md`. It must describe desired institutional behavior independently of the current code and answer every decision below.
- [ ] Decide whether a call is one institution-wide administrative round that can contain candidates assigned to both semesters, or a semester-specific event.
- [ ] Decide how calls are numbered if they can be associated with different semesters.
- [ ] Decide whether more than one call may be open at once and, if so, under what scope.
- [ ] Define exactly when a call may be closed and whether a closed call may be reopened.
- [ ] Define the distinct meanings of selection result, call outcome, enrollment state, and semester assignment. Do not use one ambiguous status to answer all four.

## Seats and semesters

- [ ] Confirm that `QT_VAGAS_CONCORRENCIA` is the annual seat allocation for one course, shift, and competition modality.
- [ ] Decide how an odd annual seat allocation is divided between semesters.
- [ ] Define what counts as an occupied seat in each semester.
- [ ] Define when an absent, withdrawn, cancelled, or moved candidate releases a seat.
- [ ] Decide whether semester capacities are fixed after import or can be adjusted by an operator.
- [ ] Define what opening and closing a semester means.
- [ ] Define which actions become forbidden after a semester is closed.

## Semester movement

- [ ] Confirm that candidate interest and consent remain outside the application and that the application records only an operator-authorized movement.
- [ ] Define which semester-2 candidates are eligible to move: pending approved candidates, enrolled candidates, or both.
- [ ] Confirm that a movement must remain within the same course, shift, and competition modality.
- [ ] Define whether an operator may choose any eligible candidate or must respect a ranking rule.
- [ ] Define the point after which semester movements are no longer allowed.
- [ ] Define how the operator indicates that desired movements are finished and remaining vacancies may be offered to the waitlist.
- [ ] Decide whether movements require a persisted audit record and which facts it must contain.

## Waitlist

- [ ] Obtain and preserve a representative real waitlist file, with personal information anonymized if necessary.
- [ ] Update `doc/sisu-file-spec.md` with an objective waitlist-file specification instead of assuming it matches the approved file.
- [ ] Confirm that lower `NU_CLASSIFICACAO` values are better and that ranking is scoped to course, shift, and competition modality.
- [ ] Define how a waitlist candidate is assigned to a semester.
- [ ] Define how vacancies are calculated before creating a subsequent call.
- [ ] Decide whether a subsequent call may fill vacancies in both semesters.
- [ ] Define what happens when a competition group has vacancies but no remaining eligible waitlist candidates.

## Admission-cycle scope and terminology

- [ ] Decide whether the application must retain historical admission cycles or intentionally operates on only one cycle at a time.
- [ ] Decide whether one import will remain restricted to one institution and one degree/course.
- [ ] Decide whether multiple campuses or course offerings must be representable even if the current institution does not need them yet.
- [ ] Choose canonical terms and use them in the workflow specification: for example, `Call` rather than both `Call` and `RollCall`, and `Registration` rather than both `Registration` and `Application`.
- [ ] Decide whether the current entity called `Course` should be named `CompetitionGroup`, `SeatGroup`, or another domain-accurate term.

## Data compatibility

- [ ] Decide whether existing production databases must be migrated or may be discarded when the new model ships.
- [ ] If migration is required, obtain an anonymized representative database and define the minimum data that must survive.
- [ ] Decide what backup or recovery step is required before applying structural migrations.
- [ ] Decide whether imports need provenance such as source filename, checksum, import timestamp, and file kind.
- [ ] Decide whether operators may delete an admission cycle and what must happen to its candidates, calls, movements, and reports.

## External contracts

- [ ] Review and approve `doc/sisu-file-spec.md`.
- [ ] Review and approve `doc/pdf-file-spec.md`.
- [ ] Decide how the semester is selected for the teacher attendance report and include the semester in its default filename if files for both semesters must coexist.
- [ ] Decide whether the CSV export is an external contract; if it is, create `doc/csv-export-spec.md`.
- [ ] Decide whether public reports should use legal name, social name, or a defined fallback rule.
- [ ] Decide whether reports require page numbers, generation timestamps, signatures, or other audit metadata.

## Readiness gate

- [ ] Ensure `AGENTS.md` links to every approved specification and no longer points to superseded design documents.
- [ ] Run and record the current backend and frontend build/test commands after the frontend merge.
- [ ] Confirm that every checkbox above is complete.
- [ ] Only then begin Step 1 of `doc/architecture-refactor-plan.md`.
