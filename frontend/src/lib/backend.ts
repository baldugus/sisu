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

export const AbsentApplication: typeof real.AbsentApplication = (...a) =>
  hasWails() ? real.AbsentApplication(...a) : mock.AbsentApplication(...a);

export const AbsentRegistration: typeof real.AbsentRegistration = (...a) =>
  hasWails() ? real.AbsentRegistration(...a) : mock.AbsentRegistration(...a);

export const Backup: typeof real.Backup = (...a) =>
  hasWails() ? real.Backup(...a) : mock.Backup(...a);

export const ClearApplicationStatus: typeof real.ClearApplicationStatus = (...a) =>
  hasWails() ? real.ClearApplicationStatus(...a) : mock.ClearApplicationStatus(...a);

export const ClearRegistrationStatus: typeof real.ClearRegistrationStatus = (...a) =>
  hasWails() ? real.ClearRegistrationStatus(...a) : mock.ClearRegistrationStatus(...a);

export const CloseRollCall: typeof real.CloseRollCall = (...a) =>
  hasWails() ? real.CloseRollCall(...a) : mock.CloseRollCall(...a);

export const CreateRollCall: typeof real.CreateRollCall = (...a) =>
  hasWails() ? real.CreateRollCall(...a) : mock.CreateRollCall(...a);

export const DeleteApprovedSelection: typeof real.DeleteApprovedSelection = (...a) =>
  hasWails() ? real.DeleteApprovedSelection(...a) : mock.DeleteApprovedSelection(...a);

export const DeleteCall: typeof real.DeleteCall = (...a) =>
  hasWails() ? real.DeleteCall(...a) : mock.DeleteCall(...a);

export const DeleteInterestedSelection: typeof real.DeleteInterestedSelection = (...a) =>
  hasWails() ? real.DeleteInterestedSelection(...a) : mock.DeleteInterestedSelection(...a);

export const DeleteRollCall: typeof real.DeleteRollCall = (...a) =>
  hasWails() ? real.DeleteRollCall(...a) : mock.DeleteRollCall(...a);

export const DeleteRollcall: typeof real.DeleteRollcall = (...a) =>
  hasWails() ? real.DeleteRollcall(...a) : mock.DeleteRollcall(...a);

export const Destroy: typeof real.Destroy = (...a) =>
  hasWails() ? real.Destroy(...a) : mock.Destroy(...a);

export const EmailPDF: typeof real.EmailPDF = (...a) =>
  hasWails() ? real.EmailPDF(...a) : mock.EmailPDF(...a);

export const EnrollApplication: typeof real.EnrollApplication = (...a) =>
  hasWails() ? real.EnrollApplication(...a) : mock.EnrollApplication(...a);

export const EnrollRegistration: typeof real.EnrollRegistration = (...a) =>
  hasWails() ? real.EnrollRegistration(...a) : mock.EnrollRegistration(...a);

export const EnrollmentPDF: typeof real.EnrollmentPDF = (...a) =>
  hasWails() ? real.EnrollmentPDF(...a) : mock.EnrollmentPDF(...a);

export const ExportCSV: typeof real.ExportCSV = (...a) =>
  hasWails() ? real.ExportCSV(...a) : mock.ExportCSV(...a);

export const FetchApplicationsByRollCall: typeof real.FetchApplicationsByRollCall = (...a) =>
  hasWails() ? real.FetchApplicationsByRollCall(...a) : mock.FetchApplicationsByRollCall(...a);

export const FetchApprovedSelection: typeof real.FetchApprovedSelection = (...a) =>
  hasWails() ? real.FetchApprovedSelection(...a) : mock.FetchApprovedSelection(...a);

export const FetchInterestedSelection: typeof real.FetchInterestedSelection = (...a) =>
  hasWails() ? real.FetchInterestedSelection(...a) : mock.FetchInterestedSelection(...a);

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

export const FetchRollCalls: typeof real.FetchRollCalls = (...a) =>
  hasWails() ? real.FetchRollCalls(...a) : mock.FetchRollCalls(...a);

export const FetchSemesters: typeof real.FetchSemesters = (...a) =>
  hasWails() ? real.FetchSemesters(...a) : mock.FetchSemesters(...a);

export const LoadApprovedSelection: typeof real.LoadApprovedSelection = (...a) =>
  hasWails() ? real.LoadApprovedSelection(...a) : mock.LoadApprovedSelection(...a);

export const LoadInterestedSelection: typeof real.LoadInterestedSelection = (...a) =>
  hasWails() ? real.LoadInterestedSelection(...a) : mock.LoadInterestedSelection(...a);

export const LoadWaitlistSelection: typeof real.LoadWaitlistSelection = (...a) =>
  hasWails() ? real.LoadWaitlistSelection(...a) : mock.LoadWaitlistSelection(...a);

export const OpenCall: typeof real.OpenCall = (...a) =>
  hasWails() ? real.OpenCall(...a) : mock.OpenCall(...a);

export const OpenFileDialog: typeof real.OpenFileDialog = (...a) =>
  hasWails() ? real.OpenFileDialog(...a) : mock.OpenFileDialog(...a);

export const OpenRollCall: typeof real.OpenRollCall = (...a) =>
  hasWails() ? real.OpenRollCall(...a) : mock.OpenRollCall(...a);

export const Restore: typeof real.Restore = (...a) =>
  hasWails() ? real.Restore(...a) : mock.Restore(...a);

export const SaveFileDialog: typeof real.SaveFileDialog = (...a) =>
  hasWails() ? real.SaveFileDialog(...a) : mock.SaveFileDialog(...a);

export const TeacherPDF: typeof real.TeacherPDF = (...a) =>
  hasWails() ? real.TeacherPDF(...a) : mock.TeacherPDF(...a);

export const WebsitePDF: typeof real.WebsitePDF = (...a) =>
  hasWails() ? real.WebsitePDF(...a) : mock.WebsitePDF(...a);
