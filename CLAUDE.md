# CLAUDE.md - AI Coding Agent Guidelines

This document provides guidelines for AI coding agents working on the SISU codebase.

## Project Overview

SISU is a full-stack desktop application for managing student admissions from Brazil's SiSU (Sistema de Selecao Unificada). It uses:
- **Backend**: Go 1.24+ with Wails v2 framework
- **Frontend**: React 18 + TypeScript + Vite — lives in a **git submodule** (`frontend/`), committed/tracked independently from the main repo
- **Database**: SQLite with go-jet for queries and golang-migrate for migrations
- **UI**: Tailwind CSS + Material Tailwind React (frontend submodule is mid-migration to shadcn/ui)

## Build/Development Commands

### Backend (Go)

#### Linting

**Important**: Always run linters before building or committing code.

```bash
# Format Go code (run first)
gofmt -w .

# Go linting (golangci-lint)
golangci-lint run

# Known nolint directives used: funlen, tagalign, mnd, godox, varnamelen, exhaustruct, lll, wrapcheck
```

#### Generate, Test and Build

```bash
# Run development server (frontend + backend with hot reload)
wails dev

# Build production binary
wails build

# Run all Go tests
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
└── frontend/            # React/TypeScript frontend (git submodule — separate repo)
    ├── src/
    │   ├── components/  # Reusable UI components
    │   └── pages/       # Page components
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
- `Semester` — Academic term within a selection year (ID, year, number 1|2, status open|closed). Auto-created when an approved selection is imported. (`types/semester.go`)
- `Registration` — Candidate's application (scores, ranking, status, candidate, nullable `SemesterID`). Approved registrations carry a `SemesterID` (1 or 2); waitlisted registrations have `SemesterID = nil`. (`types/registration.go`)
- `Course` — Academic program (period, seats, quota, minimum score)
- `Call` — Enrollment call (status, number, `SemesterID`). Every call belongs to exactly one semester. (`types/call.go`)
- `Candidate` — Personal data (name, CPF, address, contact)

**Design principle**: Domain types don't embed related entities (e.g., `Registration` doesn't contain `Course` or `Call`). Relationships are managed at the database level via foreign keys.

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
- Automatically finds-or-creates **Semester 1** and **Semester 2** records for the given year.
- Splits approved candidates 50/50 by ranking: top half → Semester 1, bottom half → Semester 2.
- The CSV must have an **even** total seat count per course; an odd count returns `ErrOddSeatsCount` ("O número total de vagas deve ser par para divisão entre semestres.").

**Call creation specifics (`commands/create_call.go`):**
- Requires a `SemesterID`. For a Semester-1 call, priority Semester-2 registrations (declined promotions) are promoted first, then remaining seats are filled from the waitlist.
- Returns `ErrOpenCallExists` if another call is already open.

## Domain Notes

- **Selection**: A yearly batch import of candidates (approved or waitlist). Approved selections auto-create Semesters 1 and 2 and split candidates 50/50 by ranking.
- **Semester**: An academic term (Number 1 or 2) within a selection year. Calls and approved registrations are tied to a semester. Created automatically on approved import — not created by the user directly.
- **Registration**: A candidate's application to a course. Approved registrations have a `SemesterID`; waitlisted registrations do not.
- **`declined_promotion`** status: a Semester-2 candidate who declined promotion to Semester 1. Can transition back to `approved`. Tracked via `RegistrationStatus` enum in `types/registration.go`.
- **Call/Rollcall**: An enrollment call for a specific semester where approved candidates can enroll or be marked absent.
- **Course**: Academic program with period (morning/evening) and quota info
- Messages and UI are in Portuguese (pt-BR)
- See **`docs/api-changelog.md`** for the detailed API diff of the semester/yearly-import rework.
