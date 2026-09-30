# Aprovados Page — Design

> Status: design agreed, not yet implemented. Implement against this doc.
> See `../CLAUDE.md` for repo-wide conventions (shadcn setup, aliases, adding a page, Wails data).

## Purpose

The **Aprovados** page is the **master roster** of every candidate approved in the current
yearly selection, across both semesters. It is **read-only** and exists for the
"review → contact" half of the workflow:

1. User imports the approved file (Data page) → backend splits candidates 50/50 by ranking into
   **Semestre 1** and **Semestre 2**.
2. **Aprovados page:** review who got in, see their standing, grab their **emails** to contact them.
3. After summoning, the user marks present/absent/enrolled — that happens on **Chamadas**, not here.

### Division of labor

| Page          | Role                                              | Editable? |
|---------------|---------------------------------------------------|-----------|
| **Aprovados** | Master roster + overview. Review, search, emails. | Read-only |
| **Chamadas**  | Per-call enrollment. Mark present/absent/enrolled.| Editable  |

Chamadas only shows people once summoned into a specific call; Aprovados shows **everyone**,
regardless of call. That is its distinct value.

## Layout

Follow the established shadcn page shell (see `DashboardPage.tsx`): page shell → KPI cards →
card-wrapped table.

```
<div className="flex-1 space-y-4 p-8 pt-6">
  ── Header: "Aprovados" + subtitle (selection context: instituição · curso · ano · N aprovados)
  ── KPI cards: grid gap-4 md:grid-cols-2 lg:grid-cols-4
       • Total aprovados
       • Matriculados (with <Progress/> = % enrolled)
       • Convocados aguardando (approved, not yet enrolled)
       • Faltosos (or Em espera)
  ── Card wrapping the roster:
       • Semester tabs: Todos / Semestre 1 / Semestre 2
       • Filter row: Turno, Cota, Status + search (nome / CPF)  [reuse current filter values]
       • shadcn <Table>:
           # Classificação · Nome · Turno · Cota · Status (badge) · Email · ⋯(detail)
       • Row click → detail <Dialog> (full ENEM scores, address, all contacts)
</div>
```

### Notes

- **Email column** is the key workflow win — make it copyable (click to copy + toast). Today
  email is buried in the detail modal.
- **Status** renders as a colored badge. Status map (label → meaning): `APPROVED` = Convocado(a),
  `ENROLLED` = Matriculado(a), `ABSENT` = Faltoso(a), `WAITLISTED` = Em espera.
  (`declined_promotion` may appear for Semester-2 candidates who declined promotion — surface it
  explicitly, e.g. "Declinou promoção".)
- **Read-only:** no checkboxes, no bulk edit, no status `<select>` in the dialog (dialog shows a
  "Voltar" button only).

## Filters (reuse existing values)

- **Turno:** Todos / Matutino (`morning`) / Noturno (`evening`)
- **Cota:** Todos / AC / C1 / C2 / C3 — the option *values* are the full quota description strings
  (see the old `CallDataTable` `TABS_3` for exact text).
- **Status:** Todos / Convocado / Matriculado / Faltoso / Em espera
- **Search:** matches Nome or CPF (case-insensitive)
- **Semester:** Todos / 1 / 2 (new — via `Registration.SemesterID`)

## Data

All via Wails bindings (`wailsjs/go/main/App`), each returns `Promise<main.Response>`
(`{ status, msg, data: any }`). Real shapes are Go structs (see `CLAUDE.md` → "Data shapes").

Fetch flow (same as today's `ApprovedPage.handleGetData`, plus two extra fields):

1. `FetchApprovedSelection()` → `data: Selection | null` (null when nothing imported → render empty state).
2. `FetchRegistrationsBySelectionID(selection.ID)` → `data: Registration[]`.
3. Per registration: `FetchRegistration(reg.ID)` (parallel `Promise.all`) → `data: RegistrationDetail`.
4. Map each detail to a row:

```ts
{
  ID, Name, CPF, Period,        // "Matutino" | "Noturno"
  Quota, Status,                // uppercased; default "APPROVED"
  Ranking,
  Email,                        // NEW — detail.Registration.Candidate.Email
  SemesterID,                   // NEW — detail.Registration.SemesterID (number | null)
  EnrollmentID,
}
```

> The current `handleGetData` already maps everything except `Email` and `SemesterID` — add those two.

## Implementation approach

Build a **self-contained** read-only page; do **not** reuse the shared `CallDataTable`
(it is still Material Tailwind and shared with Chamadas — refactor that later, separately).

- New page: `src/pages/ApprovedPage.tsx` (replace current thin wrapper).
- Page-local subcomponents under `src/pages/components/` if helpful (e.g. `ApprovedDetailDialog`).
- shadcn primitives to add (not yet installed): **`tabs`, `dialog`, `input`, `badge`**.
  Already present: `button`, `card`, `progress`, `table`. Add via `npx shadcn@latest add <name>`.
- Icons: `lucide-react`.
- Route + nav are already wired (`/approved-page` in `App.tsx`, link in `TopNav.tsx`) — no nav work needed.

## Out of scope / future

- Email/report export (will be refactored later; user will evaluate whether the in-table email makes it redundant).
- Refactoring the shared `CallDataTable` / the Chamadas page.
- Any status mutation from this page.
