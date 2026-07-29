# Architecture Refactor Plan

This plan evolves the existing application rather than rewriting it. It is divided into bounded steps that can be given individually to a fresh agent.

Do not begin this plan until every item in `TODO.md` is complete. `doc/admission-workflow-spec.md` and the approved external-contract specifications are authoritative. If the code, this plan, and an approved specification disagree, follow the specification and report the discrepancy.

## Working rules for every step

Give a fresh agent only one step at a time. The agent must:

1. Read `AGENTS.md`, `TODO.md`, the approved specifications, and the requested step completely.
2. Inspect the current code and migrations before proposing changes.
3. Preserve unrelated user changes and avoid broad cleanup outside the step.
4. Add or update tests for behavior changed in the step.
5. Run the relevant backend and frontend verification recorded in `AGENTS.md`.
6. Update generated Wails bindings when a bound Go API changes.
7. Report schema changes, compatibility consequences, tests run, and remaining known failures.
8. Stop after the step's completion criteria are met; do not begin the next step.

Each step should leave the application buildable and the database in a valid state.

## Step 1 — Establish the verified baseline

### Prompt for a fresh agent

> Establish a safe baseline for the SISU architecture refactor. Do not change product behavior or redesign the schema. Verify the commands in `AGENTS.md`, inventory current database migrations and public Wails methods, and add missing characterization tests for the approved import, initial call, call outcomes, subsequent-call allocation, backup/restore, CSV export, and all PDF populations described by the approved specs. Characterization tests must clearly distinguish behavior that is intentionally preserved from known defects that later steps will replace. Update test documentation only where it records project-specific facts.

### Deliverables

- A recorded baseline of passing backend and frontend commands.
- Characterization tests around every workflow that later steps will change.
- Contract tests for objective SiSU-file parsing behavior.
- Population tests for each report in addition to existing visual golden tests.
- A concise inventory of migrations and bound frontend APIs, kept in the step handoff or an appropriate maintained document.

### Completion criteria

- Existing behavior has not intentionally changed.
- All added tests pass.
- Known incorrect behavior is documented without being asserted as desired behavior.

## Step 2 — Introduce domain vocabulary and pure policies

### Prompt for a fresh agent

> Introduce domain-level types and pure policies matching `doc/admission-workflow-spec.md`, without replacing persistence or UI workflows yet. Represent the admission cycle, course offering, competition group, semester capacity, selection result, call outcome, enrollment state, and semester assignment as distinct concepts. Extract annual-seat splitting and vacancy calculation into pure functions. Do not put CSV parsing or database IDs inside these policies. Keep adapters to the existing model where necessary.

### Deliverables

- Domain types with names taken from the approved terminology.
- A pure annual-seat allocation policy, including the approved odd-seat rule.
- A pure semester-vacancy calculation.
- Explicit invariants for competition groups and semester capacity.
- Table-driven unit tests covering boundaries and invalid states.

### Completion criteria

- The new policy package has no dependency on CSV, Wails, SQLite, or generated database models.
- Existing application behavior still builds and passes its baseline tests.

## Step 3 — Add the admission-cycle and capacity schema

### Prompt for a fresh agent

> Add persistence for the approved admission-cycle model and semester-specific capacities. Follow the data-compatibility decision in `TODO.md`. Introduce normalized records for admission cycles, course offerings, competition groups, semesters, and their capacities. Do not yet replace calls or call outcomes. Add repositories that map database rows to the domain types from Step 2. If existing databases must survive, provide and test a forward migration using representative legacy data.

### Deliverables

- Forward database migration or documented clean-database transition, as previously decided.
- Persistence mappings for admission cycles, offerings, groups, semesters, and capacities.
- Database constraints reflecting domain uniqueness and capacity invariants.
- Migration and repository integration tests.

### Completion criteria

- No annual seat count is used ambiguously as a semester capacity.
- Migration behavior matches the approved compatibility decision.
- Backup/restore still works with the new schema.

## Step 4 — Replace the import pipeline

### Prompt for a fresh agent

> Refactor approved and waitlist import into three explicit stages: external SiSU records, validated import data, and domain persistence. Follow `doc/sisu-file-spec.md` exactly. CSV parsing must not assign semesters or create database-shaped objects. Apply the pure seat-allocation policy after validation, then persist the admission cycle, competition groups, capacities, candidates, registrations, selection results, and semester assignments transactionally. Preserve import provenance if required by the approved decisions.

### Deliverables

- External DTOs named after SiSU fields.
- File-format normalization isolated from semantic validation.
- Group-level validation for institution, course, shift, modality, seats, cutoff score, and ranking.
- Transactional approved and waitlist imports using the new model.
- Tests for malformed files, inconsistent group data, duplicate identifiers, rollback, and the real anonymized waitlist fixture.

### Completion criteria

- Semester assignment no longer occurs inside the CSV parser.
- The imported model can explain the source and capacity of every registration.
- Old import persistence is no longer the authoritative write path.

## Step 5 — Implement calls and call entries

### Prompt for a fresh agent

> Implement the call model selected in `doc/admission-workflow-spec.md`. Separate a registration's selection result and semester assignment from its participation and outcome in a call. Introduce call entries with explicit pending/enrolled/absent outcomes and implement the specified numbering, concurrency, close, reopen, and undo rules. Migrate the initial-call workflow and its frontend screens to this model. Do not implement semester movement or waitlist allocation in this step.

### Deliverables

- Schema and domain model for calls and call entries.
- Explicit commands such as recording enrollment, recording absence, undoing an outcome, closing a call, and reopening a call where permitted.
- Updated Wails API and frontend call-management UI.
- Migration or compatibility handling for legacy call data.
- State-transition and complete-call integration tests.

### Completion criteria

- Selection result is never overwritten to record a call outcome.
- Historical outcomes remain attached to their call entries.
- The initial call represents both semesters exactly as specified.

## Step 6 — Implement semester lifecycle and operator-controlled movement

### Prompt for a fresh agent

> Implement semester lifecycle and operator-controlled movement according to `doc/admission-workflow-spec.md`. The application must validate capacity, eligibility, competition group, semester state, and movement timing, but must not infer candidate interest or automatically choose a candidate unless the approved specification requires it. Persist the required audit record. Add a frontend action that lets the operator select and confirm an eligible movement.

### Deliverables

- Operational open/closed semester state and commands.
- Eligibility and capacity queries for movement.
- A transactional semester-movement command.
- Movement history/audit persistence if specified.
- UI for selecting a candidate, reviewing the target vacancy, and confirming the movement.
- Tests for successful movement, full capacity, wrong group, wrong semester, ineligible state, closed semester, and rollback.

### Completion criteria

- A movement changes the actual semester assignment.
- No movement silently occurs during call creation.
- Vacancy counts are correct in both the source and target semesters.

## Step 7 — Implement deterministic waitlist allocation

### Prompt for a fresh agent

> Replace the existing subsequent-call promotion query with a deterministic waitlist allocation service. It must run only after the operator-controlled movement phase defined by the workflow specification. Calculate vacancies per semester and competition group, select eligible waitlist registrations using the confirmed ranking direction, assign semesters, and create the next call and its entries in one transaction. Update the UI to expose the required movement-complete or allocation action.

### Deliverables

- Pure waitlist allocation policy.
- Queries that provide vacancies and eligible candidates without embedding selection policy in SQL ordering alone.
- Transactional next-call creation.
- UI flow matching the approved operational sequence.
- Tests for multiple groups, both semesters, exhausted groups, ties or missing rankings as specified, partial fills, and rollback.

### Completion criteria

- Waitlist candidates are selected in the specified order within their group.
- Annual capacity cannot be confused with semester vacancy.
- Re-running or failing allocation cannot duplicate assignments.

## Step 8 — Support admission-cycle lifecycle and history

### Prompt for a fresh agent

> Complete admission-cycle lifecycle behavior according to the decisions in `TODO.md`. Support either retained historical cycles or the explicitly approved single-cycle model. Scope all queries, uniqueness constraints, calls, semesters, reports, and imports to the active cycle. Implement safe cycle selection, deletion, and archival behavior where required.

### Deliverables

- Admission-cycle selection and lifecycle commands.
- Cycle-scoped uniqueness and queries.
- Frontend cycle context if history is retained.
- Safe deletion/archive behavior with explicit confirmation and database constraints.
- Tests proving that one cycle cannot leak candidates, capacities, calls, or reports into another.

### Completion criteria

- The `year` field is not decorative; it belongs to a coherent admission cycle.
- All behavior matches the chosen historical-data policy.

## Step 9 — Build report projections and conform to external specs

### Prompt for a fresh agent

> Refactor reporting around explicit read-only projections and make every generated output conform to the approved specifications. Construct call-report and teacher-report data before rendering. Require a semester for the teacher report, show it in the header and filename, and include only enrolled candidates from that semester and shift. Apply the approved legal/social-name rule. Add or update the CSV export specification and implementation if it was declared an external contract.

### Deliverables

- Dedicated report projection types and queries.
- PDF renderers that consume projections rather than perform domain selection.
- Semester-aware teacher report API and frontend controls.
- Updated Wails bindings and default filenames.
- Projection tests and updated PDF golden tests.
- CSV export contract tests where applicable.

### Completion criteria

- Report population, headers, ordering, grouping, and filenames match the approved specs.
- Rendering code contains no vacancy, eligibility, or call-allocation policy.

## Step 10 — Remove the legacy model and compatibility surface

### Prompt for a fresh agent

> Remove domain and persistence paths superseded by Steps 2–9. Delete the automatic semester-2 promotion logic, overloaded registration-status transitions, obsolete schema columns/tables after safe migration, duplicate Wails aliases, dead frontend compatibility code, and obsolete generated bindings. Use the canonical vocabulary from `AGENTS.md`. Do not combine this with new product behavior.

### Deliverables

- Removal migration where needed.
- Deleted legacy commands, queries, aliases, and unused frontend paths.
- Regenerated bindings.
- Updated architecture, command, and rule sections in `AGENTS.md`.
- No references to the superseded interactive-promotion design.

### Completion criteria

- There is one authoritative model and write path for each domain concept.
- Searches for old aliases and obsolete statuses return only intentional migration history.
- All contract, integration, migration, frontend, and golden tests pass.

## Step 11 — Final acceptance and documentation audit

### Prompt for a fresh agent

> Perform final acceptance against `doc/admission-workflow-spec.md`, `doc/sisu-file-spec.md`, `doc/pdf-file-spec.md`, and any CSV export specification. Do not redesign functionality. Exercise the complete workflow with representative approved and waitlist data: import, initial call, outcomes, semester movement, waitlist allocation, subsequent calls, semester closure, reports, backup, restore, and admission-cycle isolation. Correct only discrepancies required by the approved specifications and document the verification evidence.

### Deliverables

- End-to-end automated coverage of the approved workflow.
- Manual acceptance checklist for desktop-only interactions.
- Clean documentation with no stale or contradictory design files.
- Final migration and recovery verification.

### Completion criteria

- Every approved specification has corresponding automated or explicitly manual verification.
- Backend and frontend builds/tests pass from a fresh clone.
- No item remains in `TODO.md`.
