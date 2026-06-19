# Frontend migration: typed Wails boundary

This doc accompanies the backend change that replaced `Response{ Status; Msg; Data any }`
with concrete, typed return values on every bound `App` method. Apply these steps to
the `sisu-frontend` repo.

---

## 1. Regenerate Wails bindings

From the parent (`sisu`) repo root:

```bash
wails generate module
```

This overwrites `frontend/wailsjs/go/main/App.js`, `App.d.ts`, and `models.ts`.

**What changes in the generated files:**

- `models.ts` now exports real classes: `Selection`, `Registration`, `Candidate`,
  `Course`, `Call`, `RegistrationDetail` — plus string enums for `CallStatus`,
  `CoursePeriod`, `RegistrationStatus`, `SelectionKind`.
- Every method in `App.d.ts` now returns a concrete `Promise<T>` instead of
  `Promise<main.Response>`. Action methods (load, delete, enroll, PDF, etc.) return
  `Promise<void>`.
- The `Response` class is gone.

Do **not** hand-edit any file under `wailsjs/`.

---

## 2. Rewrite `src/lib/wailsCall.ts`

The old wrapper unpacked `{ status, msg, data }`. The new contract is:
- Success → method resolves with the typed value (or `null` for not-found fetch).
- Failure → method rejects with an `Error` whose `.message` is the Portuguese
  user-facing string (e.g. `"Inscrição não encontrada."`).

Replace the file with a thin pass-through:

```typescript
// Calls a bound Go method and returns its resolved value.
// On error, rejects with an Error whose .message is the user-facing message.
export async function wailsCall<T>(fn: () => Promise<T>): Promise<T> {
  return fn();
}
```

Call sites that previously accessed `res.data` become:

```typescript
// Before
const res = await wailsCall(FetchRollCalls);
const calls = res.data as Call[];

// After
const calls = await wailsCall(FetchRollCalls);
```

Error handling that previously checked `res.status !== 200`:

```typescript
// Before
const res = await wailsCall(EnrollRegistration, id);
if (res.status !== 200) toast.error(res.msg);

// After
try {
  await wailsCall(() => EnrollRegistration(id));
} catch (e) {
  toast.error((e as Error).message);
}
```

---

## 3. Call-site changes by page

All imports stay the same (`from ".../wailsjs/go/main/App"`).

### Pattern: fetch with data

```typescript
// Before
const res = await wailsCall(FetchApprovedSelection);
const selection = res.data as Selection | null;

// After  (null means not found — same semantics as before)
const selection = await wailsCall(FetchApprovedSelection);
```

### Pattern: fetch list

```typescript
// Before
const res = await wailsCall(FetchRollCalls);
const calls = res.data as Call[];

// After
const calls = await wailsCall(FetchRollCalls) ?? [];
```

### Pattern: action

```typescript
// Before
const res = await wailsCall(DeleteApprovedSelection);
if (res.status !== 200) throw new Error(res.msg);

// After
await wailsCall(DeleteApprovedSelection);  // throws on error automatically
```

### Affected files (search for `wailsCall` or direct App imports)

- `src/pages/SubscribePage.tsx`
- `src/pages/ApprovedPage.tsx`
- `src/pages/CallsPage.tsx`
- `src/pages/CallPage.tsx`
- `src/pages/ReportsPage.tsx`
- `src/pages/DataManagementPage.tsx`
- `src/pages/components/InterestedImportModal.tsx`
- `src/components/CallDataTable/DialogDataShow/index.tsx`

---

## 4. Type changes to expect at each call site

| Field / method | Was | Now |
|---|---|---|
| `Registration.Status` | `number` (0/1/2/3) | `"approved"` \| `"waitlisted"` \| `"absent"` \| `"enrolled"` |
| `Selection.Kind` | `number` | `"approved"` \| `"waitlist"` |
| `Call.Status` | `number` | `"calling"` \| `"done"` |
| `Course.Period` | `number` | `"morning"` \| `"evening"` |
| `Registration.*Score` | `{ Value: number } \| null` | `string \| null` (e.g. `"655,16"`) |
| `Course.MinimumScore` | `{ Value: number } \| null` | `string \| null` |
| `FetchApproved/InterestedSelection` | `Promise<Response>` | `Promise<Selection \| null>` |
| `FetchRollCalls` | `Promise<Response>` | `Promise<Call[]>` |
| `FetchRegistration*` | `Promise<Response>` | `Promise<Registration[]>` |
| `FetchRegistration(id)` | `Promise<Response>` | `Promise<RegistrationDetail>` |
| Action methods | `Promise<Response>` | `Promise<void>` |

**Score display:** `*Score` fields now arrive as the formatted string `"NNN,DD"` (e.g.
`"655,16"`). Remove any client-side cents-to-decimal conversion; render the string
directly.

**Enum comparisons:** replace numeric comparisons with string literals:

```typescript
// Before
if (reg.Status === 2) { /* absent */ }

// After
if (reg.Status === "absent") { }
```

---

## 5. TypeScript check

```bash
npx tsc --noEmit
```

This should pass with zero errors once the above changes are applied. Any remaining
`res.data` or `res.status` access will fail at compile time — that is the intended
outcome of this change.

---

## 6. `docs/api-changelog.md` is retired

The manual changelog was a workaround for `Data any` erasing types. With concrete
return types, breakage surfaces at `tsc` time. The file can be deleted.
