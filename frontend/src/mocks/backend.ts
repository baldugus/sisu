/**
 * Mock implementations of every Wails-bound method, for frontend-only dev.
 * State lives in memory (reset on reload). Registration status and semester are
 * derived from call entries with the same rule as the Go backend.
 */
import { types } from '../../wailsjs/go/models';
import * as f from './fixtures';

// ── Derived state ────────────────────────────────────────────────────────────

const callNumber = (callID: number) => f.mockCalls.find((c) => c.ID === callID)?.Number ?? 0;

function placement(regID: number) {
  // Latest entry, ignoring promotion offers that were not accepted.
  return f.mockEntries
    .filter((e) => e.RegistrationID === regID && (e.Kind !== 'promotion' || e.Outcome === 'enrolled'))
    .sort((a, b) => callNumber(b.CallID) - callNumber(a.CallID))[0];
}

function registration(regID: number): types.Registration {
  const r = f.mockRegistrations.find((x) => x.ID === regID)!;
  const p = placement(regID);
  const status = !p ? 'waitlisted' : p.Outcome === 'pending' ? 'approved' : p.Outcome;
  const semester = p && (p.Outcome === 'pending' || p.Outcome === 'enrolled') ? p.Semester : undefined;
  return types.Registration.createFrom({ ...r, Status: status, Semester: semester });
}

const course = (regID: number) =>
  types.Course.createFrom(f.mockCourses.find((c) => c.ID === f.mockRegistrations.find((r) => r.ID === regID)!.CourseID));

const entry = (e: (typeof f.mockEntries)[number]) =>
  types.CallEntry.createFrom({ ...e, CallNumber: callNumber(e.CallID) });

function semesters(): types.Semester[] {
  return f.mockSemesters.map((s) => {
    const seats = f.mockCourses.reduce(
      (sum, c) => sum + (s.Number === 1 ? Math.ceil(c.Seats / 2) : Math.floor(c.Seats / 2)),
      0
    );
    const occupied = f.mockRegistrations.filter((r) => registration(r.ID).Semester === s.Number).length;
    return types.Semester.createFrom({
      Number: s.Number,
      Status: s.ClosedAfterCall ? 'closed' : 'open',
      ClosedAfterCall: s.ClosedAfterCall,
      Seats: seats,
      Occupied: occupied,
    });
  });
}

// ── Mutations ────────────────────────────────────────────────────────────────

export const Backup = async (_: string): Promise<void> => {};
export const CloseCall = async (id: number): Promise<void> => {
  if (f.mockEntries.some((e) => e.CallID === id && e.Outcome === 'pending')) {
    throw new Error('Não é possível fechar a chamada com inscrições pendentes.');
  }
  f.mockCalls.find((c) => c.ID === id)!.Status = 'done';
};
export const OpenCall = async (id: number): Promise<void> => {
  f.mockCalls.find((c) => c.ID === id)!.Status = 'calling';
};
export const CloseSemester = async (n: number): Promise<void> => {
  f.mockSemesters[n - 1].ClosedAfterCall = f.mockCalls.length;
};
export const ReopenSemester = async (n: number): Promise<void> => {
  f.mockSemesters[n - 1].ClosedAfterCall = undefined;
};
// Fixed plan for the next call (the mock does not run the allocation rule).
const NEXT_CALL = [
  { RegistrationID: 4, CourseID: 1, Kind: 'promotion', Semester: 1 },
];

function requireNoOpenCall() {
  if (f.mockCalls.some((c) => c.Status === 'calling')) {
    throw new Error('Não é possível criar nova chamada enquanto outra está aberta.');
  }
}

export const CreateCall = async (): Promise<void> => {
  requireNoOpenCall();
  const id = Math.max(0, ...f.mockCalls.map((c) => c.ID)) + 1;
  f.mockCalls.push({ ID: id, Number: f.mockCalls.length + 1, Status: 'calling' });
  for (const p of NEXT_CALL) {
    f.mockEntries.push({ CallID: id, RegistrationID: p.RegistrationID, Kind: p.Kind, Semester: p.Semester, Outcome: 'pending', WantsPromotion: false });
  }
};
export const DeleteCall = async (id: number): Promise<void> => {
  const i = f.mockCalls.findIndex((c) => c.ID === id);
  f.mockCalls.splice(i, 1);
  for (let j = f.mockEntries.length - 1; j >= 0; j--) {
    if (f.mockEntries[j].CallID === id) f.mockEntries.splice(j, 1);
  }
};
export const SetCallEntryOutcome = async (callID: number, regID: number, outcome: string): Promise<void> => {
  const e = f.mockEntries.find((x) => x.CallID === callID && x.RegistrationID === regID)!;
  e.Outcome = outcome;
  if (outcome === 'absent') e.WantsPromotion = false;
};
export const SetWantsPromotion = async (callID: number, regID: number, wants: boolean): Promise<void> => {
  f.mockEntries.find((x) => x.CallID === callID && x.RegistrationID === regID)!.WantsPromotion = wants;
};
export const DeleteApprovedSelection = async (): Promise<void> => {};
export const DeleteWaitlistSelection = async (): Promise<void> => {};
export const Destroy = async (): Promise<void> => {};
export const EmailPDF = async (_1: number, _2: string, _3: number, _4: string): Promise<void> => {};
export const EnrollmentPDF = async (_1: number, _2: string, _3: number, _4: string): Promise<void> => {};
export const ExportCSV = async (_: string): Promise<void> => {};
export const LoadApprovedSelection = async (_1: number, _2: string): Promise<void> => {};
export const LoadWaitlistSelection = async (_1: number, _2: string): Promise<void> => {};
export const Restore = async (_: string): Promise<void> => {};
export const TeacherPDF = async (_1: string, _2: number, _3: string): Promise<void> => {};
export const WebsitePDF = async (_1: number, _2: string, _3: number, _4: string): Promise<void> => {};

// ── File dialogs (stub paths) ────────────────────────────────────────────────

export const OpenFileDialog = async (_1: string, _2: string, _3: string): Promise<string> =>
  '/mock/arquivo.csv';

export const SaveFileDialog = async (_1: string, _2: string, _3: string, _4: string): Promise<string> =>
  '/mock/arquivo.csv';

// ── Read methods ─────────────────────────────────────────────────────────────

export const FetchApprovedSelection = async () => types.Selection.createFrom(f.approvedSelection);

export const FetchWaitlistSelection = async () => types.Selection.createFrom(f.waitlistSelection);

export const FetchSemesters = async () => semesters();

export const FetchCalls = async () =>
  f.mockCalls.map((c) => {
    const entries = f.mockEntries.filter((e) => e.CallID === c.ID);
    const count = (s: number, k: string) => entries.filter((e) => e.Semester === s && e.Kind === k).length;
    return types.CallSummary.createFrom({
      ...c,
      Pending: entries.filter((e) => e.Outcome === 'pending').length,
      Semesters: [1, 2].map((s) => ({
        Semester: s,
        Initial: count(s, 'initial'),
        Waitlist: count(s, 'waitlist'),
        Promotion: count(s, 'promotion'),
      })),
    });
  });

export const FetchCallEntries = async (callID: number) =>
  f.mockEntries
    .filter((e) => e.CallID === callID)
    .map((e) =>
      types.CallEntryDetail.createFrom({
        Entry: entry(e),
        Registration: registration(e.RegistrationID),
        Course: course(e.RegistrationID),
      })
    );

const selectionRegs = (kind: string) =>
  f.mockRegistrations.filter((r) => r.SelectionKind === kind).map((r) => registration(r.ID));

export const FetchRegistrations = async () => f.mockRegistrations.map((r) => registration(r.ID));

export const FetchRegistrationsBySelectionID = async (id: number) =>
  selectionRegs(id === f.approvedSelection.ID ? 'approved' : 'waitlist');

export const FetchRegistrationsByCourseID = async (id: number) =>
  f.mockRegistrations.filter((r) => r.CourseID === id).map((r) => registration(r.ID));

export const FetchRegistration = async (id: number) =>
  types.RegistrationDetail.createFrom({
    Registration: registration(id),
    Course: course(id),
    History: f.mockEntries
      .filter((e) => e.RegistrationID === id)
      .sort((a, b) => callNumber(a.CallID) - callNumber(b.CallID))
      .map(entry),
  });

export const PreviewCall = async () => {
  requireNoOpenCall();
  const pick = (courseID: number, kind: string, semester: number) =>
    NEXT_CALL.filter((p) => p.CourseID === courseID && p.Kind === kind && p.Semester === semester)
      .map((p) => registration(p.RegistrationID));
  return types.CallPlan.createFrom({
    Number: f.mockCalls.length + 1,
    Semesters: semesters(),
    Courses: f.mockCourses.map((c) => ({
      Course: c,
      Vacancies1: c.ID === 1 ? 1 : 0,
      Vacancies2: 0,
      Promoted: pick(c.ID, 'promotion', 1),
      Waitlist1: pick(c.ID, 'waitlist', 1),
      Waitlist2: pick(c.ID, 'waitlist', 2),
    })),
    Promoted: 1,
    Waitlist1: 0,
    Waitlist2: 0,
  });
};
