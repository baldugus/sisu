# API Changelog — Semester-Based SiSU Imports

This document describes all backend API changes introduced by the semester splitting feature.
Feed this to the frontend agent so it can update the Wails JS bindings and component calls accordingly.

---

## Summary

Selections are now **yearly** (no longer per-semester). A new `Semester` entity was introduced.
Approved candidates are automatically split 50/50 across Semester 1 and Semester 2 by ranking.
Calls and registrations are now associated with a specific semester.

---

## Domain Model Changes

### `Selection` — field removed

| Field      | Change  | Notes                                           |
|------------|---------|-------------------------------------------------|
| `Semester` | REMOVED | Selections are now yearly. Semester lives on registrations and calls instead. |

The JSON shape returned by `FetchApprovedSelection()` and `FetchInterestedSelection()` **no longer contains a `semester` field**.

### `Registration` — field added

| Field        | Change | Type     | Notes                                                     |
|--------------|--------|----------|------------------------------------------------------------|
| `SemesterID` | ADDED  | `*int32` | Nullable. Set for approved registrations (1 or 2). `null` for waitlisted. |

### `Call` — field added

| Field        | Change | Type    | Notes                                    |
|--------------|--------|---------|------------------------------------------|
| `SemesterID` | ADDED  | `int32` | Required. The semester this call belongs to. |

### `RegistrationStatus` — new value

| Value                | Change | Notes                                                         |
|----------------------|--------|----------------------------------------------------------------|
| `declined_promotion` | ADDED  | When a Sem 2 candidate declines promotion to Sem 1. Can transition back to `approved`. |

### `Semester` — new entity

```json
{
  "ID": 1,
  "Year": 2025,
  "Number": 1,
  "Status": "open"
}
```

`Status` is one of: `"open"`, `"closed"`.
`Number` is one of: `1`, `2`.

---

## API Method Changes

### `LoadApprovedSelection` — parameter removed

```diff
-LoadApprovedSelection(year: number, semester: number, filePath: string): Response
+LoadApprovedSelection(year: number, filePath: string): Response
```

The `semester` parameter was **removed**. The backend automatically creates both Semester 1 and Semester 2 records and splits approved candidates evenly by ranking (top half → Sem 1, bottom half → Sem 2).

> **Important:** The CSV file must have an **even** number of seats per course/quota. If odd, the backend returns:
> `"O número total de vagas deve ser par para divisão entre semestres."`

### `LoadWaitlistSelection` — parameter removed

```diff
-LoadWaitlistSelection(year: number, semester: number, filePath: string): Response
+LoadWaitlistSelection(year: number, filePath: string): Response
```

The `semester` parameter was **removed**. Waitlist registrations do not get assigned a semester.

### `LoadInterestedSelection` — parameter removed

```diff
-LoadInterestedSelection(year: number, semester: number, filePath: string): Response
+LoadInterestedSelection(year: number, filePath: string): Response
```

Same as `LoadWaitlistSelection` (this is an alias).

### `CreateRollCall` — parameter added

```diff
-CreateRollCall(): Response
+CreateRollCall(semesterID: number): Response
```

Now requires a `semesterID` parameter indicating which semester the new call is for. The frontend needs a way for the user to select which semester. This could be a dropdown/selector showing available semesters.

### `FetchApprovedSelection()` / `FetchInterestedSelection()` — response shape changed

The `data` field in the response no longer contains `semester`. Example:
```diff
 {
   "id": 1,
   "kind": "approved",
   "name": "FAETERJ-Rio",
   "year": 2025,
-  "semester": 1,
   "institution": "FAETERJ-Rio",
   "degree": "Tecnológico"
 }
```

### `FetchRegistrations*()` methods — response shape changed

Each registration object in the response now includes `SemesterID`:
```diff
 {
   "ID": 1,
   "EnrollmentID": "241077018570",
   ...
   "Status": "approved",
+  "SemesterID": 1
 }
```

`SemesterID` is `null` for waitlisted registrations.

### `FetchRollCalls()` — response shape changed

Each call object now includes `SemesterID`:
```diff
 {
   "ID": 1,
   "Number": 1,
   "Status": "calling",
+  "SemesterID": 1
 }
```

`SemesterID` is always set (non-nullable).

---

## New Error Messages

| Error                    | Portuguese Message                                                 | When                                                  |
|--------------------------|--------------------------------------------------------------------|-------------------------------------------------------|
| `ErrOddSeatsCount`       | "O número total de vagas deve ser par para divisão entre semestres." | Importing approved CSV with odd seat count per course |

---

## Suggested New APIs (not yet implemented)

The frontend will likely need these to support the semester selector UI:

| Method                      | Signature                        | Returns              | Purpose                                    |
|-----------------------------|----------------------------------|----------------------|--------------------------------------------|
| `FetchSemesters`            | `() Response`                    | `[]Semester`         | List all semesters for the dropdown         |
| `FetchSemestersByYear`      | `(year int32) Response`          | `[]Semester`         | List semesters for a specific year          |

---

## Migration Notes

- The database schema is a single migration (`000001_first_migration.up.sqlite`).
- The `semesters` table is new. The `selections` table no longer has a `semester` column.
- `registrations` has a new nullable `semester_id` column.
- `calls` has a new **non-nullable** `semester_id` column (foreign key to `semesters`).
