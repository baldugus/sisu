/**
 * Backend shim — forwards calls to the real Wails bindings when running inside
 * the desktop host (window.go present), and to mock fixtures otherwise.
 *
 * Import from "@/lib/backend" instead of "wailsjs/go/main/App" everywhere in src/.
 */
import * as real from "../../wailsjs/go/main/App";
import * as mock from "../mocks/backend";

/** True when the Wails host has injected window.go (desktop mode). */
const hasWails = (): boolean => !!(window as any)?.go?.main?.App;

export const AbsentRegistration: typeof real.AbsentRegistration = (...a) =>
  hasWails() ? real.AbsentRegistration(...a) : mock.AbsentRegistration(...a);

export const Backup: typeof real.Backup = (...a) =>
  hasWails() ? real.Backup(...a) : mock.Backup(...a);

export const ClearRegistrationStatus: typeof real.ClearRegistrationStatus = (...a) =>
  hasWails() ? real.ClearRegistrationStatus(...a) : mock.ClearRegistrationStatus(...a);

export const CloseCall: typeof real.CloseCall = (...a) =>
  hasWails() ? real.CloseCall(...a) : mock.CloseCall(...a);

export const CreateCall: typeof real.CreateCall = (...a) =>
  hasWails() ? real.CreateCall(...a) : mock.CreateCall(...a);

export const DeleteApprovedSelection: typeof real.DeleteApprovedSelection = (...a) =>
  hasWails() ? real.DeleteApprovedSelection(...a) : mock.DeleteApprovedSelection(...a);

export const DeleteCall: typeof real.DeleteCall = (...a) =>
  hasWails() ? real.DeleteCall(...a) : mock.DeleteCall(...a);

export const DeleteWaitlistSelection: typeof real.DeleteWaitlistSelection = (...a) =>
  hasWails() ? real.DeleteWaitlistSelection(...a) : mock.DeleteWaitlistSelection(...a);

export const Destroy: typeof real.Destroy = (...a) =>
  hasWails() ? real.Destroy(...a) : mock.Destroy(...a);

export const EmailPDF: typeof real.EmailPDF = (...a) =>
  hasWails() ? real.EmailPDF(...a) : mock.EmailPDF(...a);

export const EnrollRegistration: typeof real.EnrollRegistration = (...a) =>
  hasWails() ? real.EnrollRegistration(...a) : mock.EnrollRegistration(...a);

export const EnrollmentPDF: typeof real.EnrollmentPDF = (...a) =>
  hasWails() ? real.EnrollmentPDF(...a) : mock.EnrollmentPDF(...a);

export const ExportCSV: typeof real.ExportCSV = (...a) =>
  hasWails() ? real.ExportCSV(...a) : mock.ExportCSV(...a);

export const FetchApprovedSelection: typeof real.FetchApprovedSelection = (...a) =>
  hasWails() ? real.FetchApprovedSelection(...a) : mock.FetchApprovedSelection(...a);

export const FetchCalls: typeof real.FetchCalls = (...a) =>
  hasWails() ? real.FetchCalls(...a) : mock.FetchCalls(...a);

export const FetchRegistration: typeof real.FetchRegistration = (...a) =>
  hasWails() ? real.FetchRegistration(...a) : mock.FetchRegistration(...a);

export const FetchRegistrations: typeof real.FetchRegistrations = (...a) =>
  hasWails() ? real.FetchRegistrations(...a) : mock.FetchRegistrations(...a);

export const FetchRegistrationsByCallID: typeof real.FetchRegistrationsByCallID = (...a) =>
  hasWails() ? real.FetchRegistrationsByCallID(...a) : mock.FetchRegistrationsByCallID(...a);

export const FetchRegistrationsByCourseID: typeof real.FetchRegistrationsByCourseID = (...a) =>
  hasWails() ? real.FetchRegistrationsByCourseID(...a) : mock.FetchRegistrationsByCourseID(...a);

export const FetchRegistrationsBySelectionID: typeof real.FetchRegistrationsBySelectionID = (...a) =>
  hasWails() ? real.FetchRegistrationsBySelectionID(...a) : mock.FetchRegistrationsBySelectionID(...a);

export const FetchSemesters: typeof real.FetchSemesters = (...a) =>
  hasWails() ? real.FetchSemesters(...a) : mock.FetchSemesters(...a);

export const FetchWaitlistSelection: typeof real.FetchWaitlistSelection = (...a) =>
  hasWails() ? real.FetchWaitlistSelection(...a) : mock.FetchWaitlistSelection(...a);

export const LoadApprovedSelection: typeof real.LoadApprovedSelection = (...a) =>
  hasWails() ? real.LoadApprovedSelection(...a) : mock.LoadApprovedSelection(...a);

export const LoadWaitlistSelection: typeof real.LoadWaitlistSelection = (...a) =>
  hasWails() ? real.LoadWaitlistSelection(...a) : mock.LoadWaitlistSelection(...a);

export const OpenCall: typeof real.OpenCall = (...a) =>
  hasWails() ? real.OpenCall(...a) : mock.OpenCall(...a);

export const OpenFileDialog: typeof real.OpenFileDialog = (...a) =>
  hasWails() ? real.OpenFileDialog(...a) : mock.OpenFileDialog(...a);

export const Restore: typeof real.Restore = (...a) =>
  hasWails() ? real.Restore(...a) : mock.Restore(...a);

export const SaveFileDialog: typeof real.SaveFileDialog = (...a) =>
  hasWails() ? real.SaveFileDialog(...a) : mock.SaveFileDialog(...a);

export const TeacherPDF: typeof real.TeacherPDF = (...a) =>
  hasWails() ? real.TeacherPDF(...a) : mock.TeacherPDF(...a);

export const WebsitePDF: typeof real.WebsitePDF = (...a) =>
  hasWails() ? real.WebsitePDF(...a) : mock.WebsitePDF(...a);
