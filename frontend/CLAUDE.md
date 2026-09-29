# CLAUDE.md — SISU Frontend

Guidelines for AI agents working in the SISU **frontend** (this repo). The frontend is a
**git submodule** of the main SISU repo, versioned independently.

## Overview

Desktop UI for SISU (Brazilian SiSU admissions management). Renders inside a **Wails v2**
(Go) desktop app — this is not a standalone web app; the browser dev server talks to Go bindings.

- **Stack:** React 18 + TypeScript + Vite
- **Styling:** Tailwind CSS (CSS-variable theme tokens, `darkMode: "class"`)
- **UI:** mid-migration **Material Tailwind → shadcn/ui**. New work uses shadcn.
- **Language:** all user-facing copy is **Portuguese (pt-BR)**.

## Commands

Run from this directory:

```bash
npm install          # install deps
npm run dev          # Vite dev server (usually launched via `wails dev` from the parent repo)
npm run build        # tsc + vite build
npx tsc --noEmit     # type-check only
```

There is no test or lint script configured in `package.json`.

## Structure

```
src/
  pages/             # one component per route (PascalCasePage.tsx), barrel-exported via index.tsx
    components/      # page-local subcomponents (e.g. import modals)
  components/        # shared components
    ui/              # shadcn primitives (lowercase: button.tsx, card.tsx, …)
    TopNav.tsx       # the active top navigation (see "Routing")
  lib/utils.ts       # cn() = clsx + tailwind-merge
  style.css          # Tailwind entry + CSS theme variables
wailsjs/             # GENERATED Go↔JS bindings — DO NOT EDIT
  go/main/App.{d.ts,js}   # bound backend methods
  go/models.ts            # generated TS models (main.Response, etc.)
components.json      # shadcn config
docs/                # design docs for individual pages/features
```

## shadcn/ui

- Config in `components.json`: style `base-nova`, baseColor `neutral`, `cssVariables: true`,
  icon library **lucide**.
- **Path alias `@/*` → `./src/*`** (configured in both `vite.config.ts` and `tsconfig.json`).
  Import primitives from `@/components/ui/...`, the helper from `@/lib/utils`.
- Add a primitive: `npx shadcn@latest add <name>` (e.g. `dialog`, `tabs`, `input`, `badge`).
  Installed so far: `button`, `card`, `progress`, `table`.
- Use the `cn()` helper (`@/lib/utils`) to merge classes. Theme via tokens
  (`bg-background`, `text-foreground`, `text-muted-foreground`, `bg-card`, `text-primary`,
  `border-destructive`, …) — defined in `src/style.css`.
- **Icons:** new code uses `lucide-react`; legacy components use `react-icons`. Prefer lucide.

### Refactored-page template (the pattern to copy)

`DashboardPage.tsx` (cards/progress/table) and `DataManagementPage.tsx` (cards/buttons/toasts)
are the canonical shadcn pages. Page shell + header:

```tsx
<div className="flex-1 space-y-4 p-8 pt-6">
  <div className="flex items-center justify-between space-y-2">
    <h2 className="text-3xl font-bold tracking-tight text-foreground">Título</h2>
  </div>
  {/* KPI cards: grid gap-4 md:grid-cols-2 lg:grid-cols-4 */}
  {/* content cards / tables */}
</div>
```

## Routing & navigation

- Router: `react-router-dom` `<Routes>` in `src/App.tsx`.
- **Navigation is a top bar: `src/components/TopNav.tsx`** (sticky, shadcn-styled, active link via
  a `getLinkClass(path)` helper).
- `SideBar.tsx` and `NavBar.tsx` are **legacy/dead** (Material Tailwind, not mounted in `App.tsx`). Ignore them.

### Adding a page

1. Create `src/pages/XxxPage.tsx` (default export).
2. Re-export it from `src/pages/index.tsx`.
3. Add a `<Route>` in `src/App.tsx`.
4. Add a `<Link>` to `src/components/TopNav.tsx` using the `getLinkClass` pattern + a lucide icon.

## Backend data (Wails bindings)

- Bound methods live in `wailsjs/go/main/App` and each returns `Promise<main.Response>`:
  ```ts
  class Response { status: number; msg: string; data: any }
  ```
- **`data` is untyped (`any`).** The real shapes are the Go structs in the parent repo's `types/`
  package. When unsure of a field, read the Go source (`../types/*.go`) — do not guess.
- After adding/changing a bound Go method, bindings are regenerated with
  `wails generate module` (run from the parent repo). Never hand-edit `wailsjs/`.

### Data shapes (from `../types/`)

- **Selection** — yearly batch: `ID, Name, Kind ("approved"|"waitlist"), Year, Institution, Degree`.
- **Registration** — an application: `ID, EnrollmentID, Option, *Score fields (.Value),
  Ranking, Status, Candidate, SemesterID (*int32)`.
  `Status` values: `approved | waitlisted | absent | enrolled | declined_promotion`.
- **RegistrationDetail** — `{ Registration, Course, Call }`.
- **Candidate** — `ID, CPF, Name, SocialName, BirthDate, Sex, MotherName, AddressLine,
  AddressLine2, HouseNumber, Neighborhood, Municipality, State, CEP, Email, Phone1, Phone2`.
- **Course** — `ID, Seats, MinimumScore, Period ("morning"|"evening"), Quota`.
- **Call** — enrollment call: `Status, Number, SemesterID`.

## Conventions

- Components: `PascalCase`; props types `I{Name}` or a descriptive type; functions `camelCase`.
- Pages: `PascalCasePage.tsx` in `src/pages/`; shadcn primitives lowercase in `src/components/ui/`.
- Functional components, arrow functions, props destructured in the function signature.
- TypeScript strict mode is on; target ESNext, `jsx: react-jsx`.
- All copy in pt-BR.

## Migration status

| Area                                  | State                                                      |
|---------------------------------------|------------------------------------------------------------|
| DashboardPage, DataManagementPage     | ✅ shadcn                                                   |
| Aprovados (`ApprovedPage`)            | 🔜 redesign — see `docs/aprovados-page-design.md`          |
| Chamadas (`CallDataTable` & co.)      | ❌ still Material Tailwind                                  |
| Subscribe / Reports / import modals   | ❌ still Material Tailwind / plain Tailwind                 |
