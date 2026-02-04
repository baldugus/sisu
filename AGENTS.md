# AGENTS.md - AI Coding Agent Guidelines

This document provides guidelines for AI coding agents working on the SISU codebase.

## Project Overview

SISU is a full-stack desktop application for managing student admissions from Brazil's SiSU (Sistema de Selecao Unificada). It uses:
- **Backend**: Go 1.24+ with Wails v2 framework
- **Frontend**: React 18 + TypeScript + Vite
- **Database**: SQLite with go-jet for queries and golang-migrate for migrations
- **UI**: Tailwind CSS + Material Tailwind React

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
└── frontend/            # React/TypeScript frontend
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
```

## Architecture

### Domain Types (`types/`)

Domain types are intentionally kept flat without nested relationships:

- `Selection` - Batch import metadata (name, kind, year, semester)
- `Registration` - Candidate's application (scores, ranking, status, candidate)
- `Course` - Academic program (period, seats, quota, minimum score)
- `Call` - Enrollment call (status, number)
- `Candidate` - Personal data (name, CPF, address, contact)

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

## Domain Notes

- **Selection**: A batch import of candidates (approved or waitlist)
- **Registration**: A candidate's application to a course
- **Call/Rollcall**: An enrollment call where approved candidates can enroll
- **Course**: Academic program with period (morning/evening) and quota info
- Messages and UI are in Portuguese (pt-BR)
