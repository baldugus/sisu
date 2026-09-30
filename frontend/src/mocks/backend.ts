/**
 * Mock implementations of every Wails-bound method.
 * Read methods return fixture data; mutations are no-ops.
 * State does NOT change between calls (read-only fixtures).
 */
import { types } from "../../wailsjs/go/models";
import * as f from "./fixtures";

// ── Mutations (no-op) ────────────────────────────────────────────────────────

export const AbsentRegistration = async (_: number): Promise<void> => {};
export const Backup = async (_: string): Promise<void> => {};
export const ClearRegistrationStatus = async (_: number): Promise<void> => {};
export const CloseCall = async (_: number): Promise<void> => {};
export const CreateCall = async (_: number): Promise<void> => {};
export const DeleteApprovedSelection = async (): Promise<void> => {};
export const DeleteCall = async (_: number): Promise<void> => {};
export const DeleteWaitlistSelection = async (): Promise<void> => {};
export const Destroy = async (): Promise<void> => {};
export const EmailPDF = async (_1: number, _2: string, _3: string): Promise<void> => {};
export const EnrollRegistration = async (_: number): Promise<void> => {};
export const EnrollmentPDF = async (_1: number, _2: string, _3: string): Promise<void> => {};
export const ExportCSV = async (_: string): Promise<void> => {};
export const LoadApprovedSelection = async (_1: number, _2: string): Promise<void> => {};
export const LoadWaitlistSelection = async (_1: number, _2: string): Promise<void> => {};
export const OpenCall = async (_: number): Promise<void> => {};
export const Restore = async (_: string): Promise<void> => {};
export const TeacherPDF = async (_1: string, _2: string): Promise<void> => {};
export const WebsitePDF = async (_1: number, _2: string, _3: string): Promise<void> => {};

// ── File dialogs (stub paths) ────────────────────────────────────────────────

export const OpenFileDialog = async (_1: string, _2: string, _3: string): Promise<string> =>
  "/mock/arquivo.csv";

export const SaveFileDialog = async (_1: string, _2: string, _3: string, _4: string): Promise<string> =>
  "/mock/arquivo.csv";

// ── Read methods (fixture data) ──────────────────────────────────────────────

const registrations = () => f.mockRegistrations.map((r) => types.Registration.createFrom(r));

export const FetchApprovedSelection = async () => types.Selection.createFrom(f.approvedSelection);

export const FetchWaitlistSelection = async () => types.Selection.createFrom(f.waitlistSelection);

export const FetchCalls = async () => f.mockCalls.map((c) => types.Call.createFrom(c));

export const FetchSemesters = async () => f.mockSemesters.map((s) => types.Semester.createFrom(s));

export const FetchRegistrations = async () => registrations();

export const FetchRegistrationsBySelectionID = async (_: number) => registrations();

export const FetchRegistrationsByCallID = async (_: number) => registrations();

export const FetchRegistrationsByCourseID = async (_: number) => registrations();

export const FetchRegistration = async (id: number) => {
  const detail =
    f.mockRegistrationDetails.find((d) => d.Registration.ID === id) ??
    f.mockRegistrationDetails[0];
  return types.RegistrationDetail.createFrom(detail);
};
