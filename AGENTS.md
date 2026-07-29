# What this is

SISU is a desktop application for managing the admission of students to a Brazilian public university through Brazil’s SiSU (Sistema de Seleção Unificada) process, which uses ENEM (Exame Nacional do Ensino Médio) scores.

Before the exam, universities register the courses and shifts they offer and the total number of seats for each course-and-shift combination in MEC’s platform. SiSU distributes those seats among the applicable competition modalities, including broad competition and quota categories. The files later provided by MEC contain the resulting seat count for each shift-and-quota group.

Because universities have two academic semesters, the annual selection supplies two incoming cohorts. SiSU reports the annual seat allocation, and this application divides it between the two semesters.

Eligible candidates may compete through quota categories reserved for particular social, educational, racial, disability, or income groups.

At the end of the exams, when scores are available, candidates can choose universities to compete for a seat within a time frame.

After the selection period, MEC provides the university with an approved-candidates file. This application uses that file as its initial import.

During a later enrollment window, candidates that didn't get approved can enroll in a waitlist, in an effort to fill seats if someone approved gives up. When that period is over, MEC provides another file for the university, the waitlist. This app uses this file to complete remaining seats.

This application shares its name with the government-run SiSU process it helps the university administer.

# Intended Workflow

The application should have the following set of features:

- Import one annual approved-candidates CSV and, optionally afterward, one annual waitlist CSV.
- Create the two semesters for the selected year.
- Group registrations by shift and quota within a single degree/course imported for the admission cycle.
- Use MEC’s supplied ranking and declared seat count to assign approved candidates to semesters: ranks up to half the annual seats go to semester 1 and the remaining ranks go to semester 2.
- Create an initial call and let operators mark called candidates as enrolled or absent.
- Require all registrations in a call to be resolved before the call can be closed.
- Before filling semester-1 vacancies from the waitlist, allow the operator to move eligible candidates from semester 2 to semester 1 within the same course, shift, and quota group. Candidate interest and confirmation are handled administratively outside the application.
- After the operator finishes any semester movements, create subsequent calls to fill remaining vacancies with waitlisted candidates according to their ranking and competition group.
- Generate administrative PDFs, export enrolled candidates as CSV, and support database backup and restore.

# Current Implementation

The approved-candidates import creates both semesters and assigns each registration to semester 1 or semester 2 using its MEC ranking and the quota-specific seat allocation from the CSV. It also creates a single initial call associated with semester 1 and assigns every approved registration to that call, including registrations assigned to semester 2.

Only one call can be open globally at a time, and call numbers are global rather than scoped to a semester. While a call is open, operators can mark its registrations as enrolled or absent and can revert those decisions. A call can only be closed after every registration in it has been resolved.

When creating a subsequent semester-1 call, the application automatically selects approved semester-2 registrations from the same shift-and-quota group before drawing from the waitlist. This changes their call and status but does not change their assigned semester. For other calls, remaining vacancies are drawn directly from the waitlist.

The application also generates administrative PDFs, exports enrolled registrations as CSV, and supports database backup and restore.

# Known Gaps

The current implementation does not yet support the intended independent management of both semesters. The initial import does not create a separate call for semester 2, only one call can be open at a time, call numbering is global, and semester open/closed status is not operationally enforced.

Semester movement also does not match the intended workflow. Instead of allowing the operator to choose an eligible semester-2 candidate after confirming their interest outside the application, subsequent semester-1 call creation automatically selects candidates by ranking. That operation does not update their semester assignment, and it only considers registrations whose status is still approved, which may exclude candidates who were already enrolled into semester 2.

Waitlist selection currently orders ranking in the opposite direction from semester-2 selection. The correct ordering must be confirmed against the MEC file semantics before relying on it to select the best-ranked candidates.

The existing backend and frontend gap-analysis documents describe a superseded design in which marking a semester-1 candidate absent immediately prompts the operator to accept or decline an automatically selected semester-2 candidate. That is not the intended workflow: candidate interest and confirmation remain an informal administrative process, and the application should provide a flexible operator-controlled movement before remaining vacancies are filled from the waitlist.

# Commands
# Rules
# References
