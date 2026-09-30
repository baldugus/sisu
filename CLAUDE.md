# CLAUDE.md - AI Coding Agent Guidelines

This document provides guidelines for AI coding agents working on the SISU codebase.

## Project Overview

SISU is a full-stack desktop application for managing student admissions from Brazil's SiSU (Sistema de Selecao Unificada). It uses:
- **Backend**: Go 1.26+ with Wails v2 framework
- **Frontend**: React 18 + TypeScript + Vite — lives in `frontend/` in this same repository (monorepo; formerly the `sisu-frontend` submodule)
- **Database**: SQLite with go-jet for queries and golang-migrate for migrations
- **UI**: Tailwind CSS v4 + shadcn/ui (migration off Material Tailwind is complete)

## Build/Development Commands

### Backend (Go)

#### Linting

**Important**: Always run linters before building or committing code.

```bash
# Format Go code (run first)
gofmt -w .

# Go linting (golangci-lint)
golangci-lint run

# Same lint as CI (golangci-lint v2; config lives in .github/ since the root .golangci.yml is gitignored)
golangci-lint run --config .github/golangci.yml

# Known nolint directives used: funlen, tagalign, mnd, godox, varnamelen, exhaustruct, lll, wrapcheck
```

#### Generate, Test and Build

```bash
# Run development server (frontend + backend with hot reload)
wails dev

# Build production binary
wails build

# Run all Go tests
# (main.go embeds frontend/dist — run `cd frontend && npm run build` once first,
# otherwise the root package fails with "pattern all:frontend/dist: no matching files found")
go test ./...

# Run tests in a specific package
go test ./csvparser
go test ./database
go test ./commands

# Run a single test by name
go test -run TestParseCSVFile ./csvparser
go test -v -run TestParse ./csvparser

# Run tests with verbose output
go test -v ./...

# Generate code (enums using go-enum)
go generate ./...
```

### Frontend (Node.js)

```bash
# Install dependencies (run from frontend/ directory)
cd frontend && npm install

# Run frontend dev server (usually run via wails dev)
cd frontend && npm run dev

# Build frontend for production
cd frontend && npm run build

# Type check
cd frontend && npx tsc --noEmit
```

## Project Structure

```
sisu/
├── main.go              # Application entry point
├── app.go               # Wails app struct and API methods exposed to frontend
├── commands/            # Command pattern implementations
├── csvparser/           # CSV parsing for SiSU data files
├── database/            # Database layer (go-jet queries, migrations)
│   ├── migrations/      # SQL migration files
│   └── .gen/            # Generated go-jet code (do not edit manually)
├── pdfbuilder/          # PDF report generation
├── types/               # Domain types and enums
├── testutil/            # Test fixtures, helpers, and assertions for integration tests
└── frontend/            # React/TypeScript frontend (see frontend/CLAUDE.md)
    ├── src/
    │   ├── components/  # Reusable UI components (Tailwind v4 + shadcn/ui)
    │   ├── pages/       # Page components (Portuguese names, e.g. Painel, Chamadas, Dados)
    │   ├── hooks/       # Custom React hooks
    │   ├── lib/         # Frontend utilities (e.g. wailsCall wrapper)
    │   └── mocks/       # Mock backend for frontend-only dev
    └── wailsjs/         # Auto-generated Wails JS bindings (do not edit)
```

## Code Style Guidelines

### Go

#### Imports
- Group imports: stdlib, external packages, internal packages (separated by blank lines)
- Use the module path `github.com/baldugus/sisu` for internal imports

```go
import (
    "context"
    "errors"
    "fmt"

    "github.com/go-jet/jet/v2/sqlite"
    "go.uber.org/zap"

    "github.com/baldugus/sisu/database"
    "github.com/baldugus/sisu/types"
)
```

#### Naming Conventions
- Types: PascalCase (`Selection`, `ParsedCsv`, `LoadSelectionCommand`)
- Functions/Methods: PascalCase for exported, camelCase for unexported
- Variables: camelCase (`parsedCsv`, `selectionRepo`)
- Constants: PascalCase for exported errors, camelCase otherwise
- Acronyms: Keep capitalized when at start (`ID`, `CPF`, `CSV`), lowercase otherwise (`selectionID`)

#### Error Handling
- Use custom error types as structs implementing the `error` interface
- Wrap errors with context using `fmt.Errorf("context: %w", err)`
- Use `errors.As()` for type checking, `errors.Is()` for value checking
- Return errors early, avoid deep nesting

```go
// Custom error type pattern
type ErrFileNotFound struct {
    Path string
    Err  error
}

func (e *ErrFileNotFound) Error() string {
    return fmt.Sprintf("file not found: %s", e.Path)
}

// Error wrapping pattern
if err != nil {
    return fmt.Errorf("fetch selection kind: %w", err)
}
```

#### Comments
- When commenting out more than 5 sequential lines, use block comments (`/* */`) instead of line comments (`//`)

#### Structs
- Use JSON tags for API responses: `json:"fieldName"`
- Use CSV tags for parsing: `csv:"COLUMN_NAME"`
- Align struct tags when practical

#### Testing
- Use table-driven tests with descriptive names
- Test files: `*_test.go` in the same package
- Use `t.Helper()` in test helper functions
- Test data goes in `testdata/` subdirectories

```go
func TestParse(t *testing.T) {
    tests := []struct {
        name        string
        input       string
        wantLen     int
        expectedErr func(error) bool
    }{
        {
            name:    "valid CSV with data",
            input:   "header1;header2\nvalue1;value2",
            wantLen: 1,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test logic
        })
    }
}
```

### TypeScript/React

#### Imports
- React imports first, then external libraries, then local imports
- Use named exports from barrel files (index.tsx)

```typescript
import { Route, Routes } from "react-router-dom";
import { Toaster } from "react-hot-toast";

import { NavBar, SideBar } from "./components";
import { HomePage, ApprovedPage } from "./pages";
```

#### Naming Conventions
- Components: PascalCase (`DataImportBox`, `NavBar`)
- Props interfaces: `I{ComponentName}` or descriptive type (`IDataImportBox`)
- Functions: camelCase
- Files: PascalCase for components (`DataImportBox.tsx`), lowercase for utilities

#### Component Structure
- Functional components with arrow functions
- Props destructured in function signature
- Tailwind CSS for styling (utility-first)

```typescript
type IDataImportBox = {
  dataType: "inscritos" | "aprovados";
  text: string;
  mockHasData: boolean;
  actionfunction: () => void;
};

const DataImportBox = ({ dataType, text, mockHasData, actionfunction }: IDataImportBox) => {
  return (
    <CardWrapper>
      {/* component content */}
    </CardWrapper>
  );
};

export default DataImportBox;
```

#### TypeScript Configuration
- Strict mode enabled
- Target: ESNext
- JSX: react-jsx

## Generated Code

**Do not manually edit:**
- `database/.gen/` - Generated by go-jet from SQL schema
- `frontend/wailsjs/` - Generated by Wails for Go/JS bindings
- `types/*_enum.go` - Generated by go-enum from comments

To regenerate:
```bash
# Regenerate enums
go generate ./types/...

# Regenerate go-jet models (requires running migrations first)
jet -dsn="path/to/db.sqlite" -schema=main -path=./database/.gen

# Regenerate frontend Wails JS/TS bindings (after adding/changing bound App methods in app.go)
wails generate module
```

## Architecture

### Domain Types (`types/`)

Domain types are intentionally kept flat without nested relationships:

- `Selection` — Yearly batch import metadata (name, kind, year, institution, degree). **No semester field** — selections are annual; semester is a separate entity. (`types/selection.go`)
- `Semester` — One of the two fixed intakes (`Number` 1|2). Both rows are **seeded by the migration** and always exist; only their closure changes (`ClosedAfterCall`, derived `Status` open|closed). `FetchSemesters` also returns `Seats` and `Occupied`. (`types/semester.go`)
- `Registration` — Candidate's application: **imported data only** (scores, ranking, candidate). `Status` and `Semester` are **derived** from its call entries, never stored. (`types/registration.go`)
- `Course` — Academic program (period, seats, quota, minimum score). `Seats` is a `types.Seats` value (`types/seats.go`) that stores the per-semester count, so an odd total is unrepresentable; build it with `types.NewSeats(total)` (returns `ErrOddSeatsCount` on odd). The database backs this with a `CHECK (seats % 2 = 0)` on `courses`. Each semester gets `Seats.PerSemester()`.
- `Call` — Enrollment call (status, number). A call covers **both** semesters; each entry says which. (`types/call.go`)
- `CallEntry` — One registration taking part in one call: `Kind` (initial|waitlist|promotion), `Semester`, `Outcome` (pending|enrolled|absent|declined), `WantsPromotion`. (`types/call_entry.go`)
- `Candidate` — Personal data (name, CPF, address, contact)

**Design principle**: Domain types don't embed related entities (e.g., `Registration` doesn't contain `Course` or `Call`). Relationships are managed at the database level via foreign keys.

### Admission cycle state: call entries (`call_entries` table)

Everything that happens during the cycle is recorded as **call entries owned by the call that
caused it**; registrations are never updated. This is what makes every action undoable in
reverse order — deleting a call deletes its entries (`ON DELETE CASCADE`) and all of its
effects disappear with them, with nothing to restore.

- **Placement** (view `registration_placements`): a registration's current state is its entry
  from the latest call, ignoring promotion offers that were not accepted. No entry = waitlisted;
  `pending` → `approved` (called); `enrolled`; `absent`.
- **Occupied seats** per (course, semester) = placements that are `pending` or `enrolled`.
- **Undo stack**: only the last call can be reopened/deleted; call 1 is only removed by deleting
  the approved selection; a semester closure (recorded as `closed_after_call`) must be undone
  before the calls it follows. Undoing everything returns to the state right after the import.

### CSV Parser (`csvparser/`)

The CSV parser uses intermediate types to group related data during parsing:

```go
// ParsedRegistration groups a registration with its related course and call
type ParsedRegistration struct {
    Registration *types.Registration
    Course       *types.Course
    Call         *types.Call
}

// ParsedSelection is returned by ToSelectionDomain()
type ParsedSelection struct {
    Selection     *types.Selection
    Registrations []*ParsedRegistration
}
```

### Database Layer (`database/`)

Database operations use `qrm.DB` interface (from go-jet) to work with both transactions and direct connections:

```go
// Functions accept qrm.DB, allowing use inside or outside transactions
func CreateSelection(db qrm.DB, selection *types.Selection) (int32, error)
func CreateRegistration(db qrm.DB, args *CreateRegistrationArgs) error
```

Transaction handling is done via `Database.RunInTx()`:

```go
db.RunInTx(func(tx qrm.DB) error {
    selectionID, _ := database.CreateSelection(tx, selection)
    // ... more operations in same transaction
    return nil
})
```

#### CASCADE DELETE and Referential Integrity

**IMPORTANT**: Registrations are NEVER deleted directly. The database enforces CASCADE DELETE:

```sql
-- registrations.candidate_id foreign key has ON DELETE CASCADE
candidate_id INTEGER NOT NULL REFERENCES candidates ON DELETE CASCADE
```

**Deletion Pattern:**
- Delete candidates → Database automatically cascades to delete their registrations
- Never call `DeleteRegistrations*()` functions when deleting selections
- Only `Candidates → Registrations` uses CASCADE (tightly coupled entities)
- All other relationships use RESTRICT (explicit control required)

**Rationale:**
- Candidates and registrations are essentially the same data split for organization
- Prevents orphaned registrations at the database level
- Simplifies deletion logic and prevents bugs
- Maintains explicit control over other relationships (selections, courses, calls)

**Example:**
```go
// CORRECT - Delete candidates, CASCADE handles registrations
candidateIDs, _ := database.FetchCandidateIDsBySelectionID(tx, selectionID)
database.DeleteCandidatesByIDs(tx, candidateIDs)  // Registrations auto-deleted

// WRONG - Never delete registrations directly in selection deletion
database.DeleteRegistrationsBySelectionID(tx, selectionID)  // Function doesn't exist
database.DeleteCandidatesByIDs(tx, candidateIDs)
```

### Command Layer (`commands/`)

Commands orchestrate business logic and transactions. Example flow for `LoadSelectionCommand`:

1. Validate business rules (check existing selections)
2. Parse CSV into `ParsedSelection`
3. Run transaction: create selection, call, courses, candidates, registrations

**Approved import specifics (`commands/load_selection.go`):**
- Creates call 1 with an `initial` entry for every approved student.
- Splits approved candidates 50/50 by ranking within each course (time slot × quota): top half → Semester 1, bottom half → Semester 2. The rule lives in the domain method `Course.SemesterForRanking` (`types/course.go`), called per registration by the command when building the call-1 entries.
- Every course must have an **even** total seat count (approved and waitlist files alike): the parser builds seats with `types.NewSeats`, so an odd count fails the import with `types.ErrOddSeatsCount` ("O número total de vagas deve ser par para divisão entre semestres."). This is intentional — odd counts are not auto-split.

**Call creation specifics (`commands/create_call.go`, rule in `allocation/`):**
- `PreviewCallCommand` and `CreateCallCommand` run the same plan, so the preview is exactly what gets created.
- Per course: semester-1 vacancies go first to semester-2 students who are enrolled, marked `WantsPromotion`, and were never offered a promotion (`promotion` entries, best ranked first); the remaining semester-1 vacancies, then semester-2 vacancies, go to never-called waitlist registrations (ranking ascending; missing ranking last). A closed semester receives no one. A seat freed by an accepted promotion is only refilled in the **next** call.
- Promotion offers are answered with `enrolled` (accepted → moves to semester 1) or `declined` (stays in semester 2, never offered again).
- Returns `ErrOpenCallExists`, `ErrAllSemestersClosed`, `ErrAllCoursesFull` or `ErrNoCandidatesToCall` when nothing can be created.

### Wails Boundary (`app.go`)

Every bound `App` method returns a concrete typed value plus `error` — there is no
`Response{ Status; Msg; Data any }` wrapper. On success, Wails resolves the JS promise with
the typed value (`*types.Selection`, `[]*types.Registration`, etc., or nothing for action
methods); on failure it rejects with an `Error` whose `.message` is the Portuguese
user-facing string produced by `translateError()` (`app.go`).

- One vocabulary, no aliases: bound names use the domain terms `Registration` (never
  "Application"), `Call` (never "RollCall") and `Waitlist` (never "Interested"), e.g.
  `SetCallEntryOutcome`, `FetchCalls`, `LoadWaitlistSelection`. UI-only identifiers (e.g. the
  `useRollCallRows` hook) may use screen language, but don't add bound aliases.
- Enums (`RegistrationStatus`, `SelectionKind`, `CallStatus`, `CoursePeriod`, `SemesterStatus`,
  `CallEntryKind`, `CallEntryOutcome`) serialize as
  string literals (e.g. `"approved"`), not numeric codes.
- `*Score` fields (`types/score.go`) serialize as a formatted string (e.g. `"655,16"`), not
  a numeric struct. `Seats` (`types/seats.go`) serializes as the total seat count (a number).
- `wails generate module` regenerates `frontend/wailsjs/go/main/App.d.ts` and `models.ts`
  with these concrete TS types, so a shape change now fails at `tsc` time on the frontend
  instead of silently drifting.

## Domain Notes

- **Selection**: A yearly batch import of candidates (approved or waitlist). The approved import opens call 1 with candidates split 50/50 by ranking between the semesters. Only one waitlist file per cycle.
- **Semester**: One of the two fixed intakes (1 or 2), seeded by the migration. The operator can close one when it should receive no one else ("Fechar semestre"); later calls then only fill the other.
- **Registration**: A candidate's application to a course. Its status and semester are derived from its call entries.
- **Promotion**: during a call, a semester-2 student can ask to move to semester 1 (`WantsPromotion` on their entry). The next call offers them free semester-1 seats before the waitlist; the operator records accepted (`enrolled`) or `declined`.
- **Call/Rollcall**: An enrollment call covering both semesters where called candidates enroll or are marked absent (and promotion offers are accepted or declined).
- **Course**: Academic program with period (morning/evening) and quota info
- Messages and UI are in Portuguese (pt-BR)
- See **`docs/future-work.md`** for the design-improvement backlog and **`docs/testing.md`** for the integration-testing guide.
