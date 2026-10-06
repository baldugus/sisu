package database

import (
	"slices"

	"github.com/baldugus/sisu/database/.gen/model"
	. "github.com/baldugus/sisu/database/.gen/table"
	"github.com/baldugus/sisu/database/.gen/view"
	"github.com/baldugus/sisu/types"
	"github.com/go-jet/jet/v2/qrm"
	. "github.com/go-jet/jet/v2/sqlite"
)

// callEntriesInsertBatch keeps each INSERT well under SQLite's bound-variable limit.
const callEntriesInsertBatch = 500

func CreateCallEntries(db qrm.DB, entries []*types.CallEntry) error {
	for batch := range slices.Chunk(entries, callEntriesInsertBatch) {
		models := make([]*model.CallEntries, len(batch))
		for i, e := range batch {
			models[i] = toCallEntryModel(e)
		}

		stmt := CallEntries.INSERT(CallEntries.AllColumns).
			MODELS(models)

		if _, err := stmt.Exec(db); err != nil {
			return err
		}
	}

	return nil
}

type callEntryResult struct {
	model.CallEntries

	Call model.Calls
}

// FetchCallEntry returns one registration's entry in one call.
func FetchCallEntry(db qrm.DB, callID, registrationID int32) (*types.CallEntry, error) {
	stmt := SELECT(
		CallEntries.AllColumns,
		Calls.AllColumns,
	).FROM(
		CallEntries.
			INNER_JOIN(Calls, Calls.ID.EQ(CallEntries.CallID)),
	).WHERE(
		CallEntries.CallID.EQ(Int32(callID)).
			AND(CallEntries.RegistrationID.EQ(Int32(registrationID))),
	)

	var result callEntryResult

	err := stmt.Query(db, &result)
	if err != nil {
		return nil, err
	}

	return toCallEntryDomain(&result.CallEntries, result.Call.Number), nil
}

// FetchRegistrationHistory returns every call entry of a registration, oldest call first.
func FetchRegistrationHistory(db qrm.DB, registrationID int32) ([]*types.CallEntry, error) {
	stmt := SELECT(
		CallEntries.AllColumns,
		Calls.AllColumns,
	).FROM(
		CallEntries.
			INNER_JOIN(Calls, Calls.ID.EQ(CallEntries.CallID)),
	).WHERE(
		CallEntries.RegistrationID.EQ(Int32(registrationID)),
	).ORDER_BY(
		Calls.Number.ASC(),
	)

	var result []callEntryResult

	err := stmt.Query(db, &result)
	if err != nil {
		return nil, err
	}

	history := make([]*types.CallEntry, len(result))
	for i := range result {
		history[i] = toCallEntryDomain(&result[i].CallEntries, result[i].Call.Number)
	}

	return history, nil
}

type callEntryDetailResult struct {
	model.CallEntries

	Call         model.Calls
	Registration registrationResult
	Course       model.Courses
	Quota        model.Quotas
}

// FetchCallEntryDetails returns the entries of a call with their registrations
// and courses. Semester and course filters are optional.
func FetchCallEntryDetails(
	db qrm.DB,
	callID int32,
	courseID *int32,
	semester *int32,
) ([]*types.CallEntryDetail, error) {
	where := CallEntries.CallID.EQ(Int32(callID))
	if courseID != nil {
		where = where.AND(Registrations.CourseID.EQ(Int32(*courseID)))
	}

	if semester != nil {
		where = where.AND(CallEntries.Semester.EQ(Int32(*semester)))
	}

	stmt := SELECT(
		CallEntries.AllColumns,
		Calls.AllColumns,
		Registrations.AllColumns,
		Candidates.AllColumns,
		view.RegistrationPlacements.AllColumns,
		Courses.AllColumns,
		Quotas.AllColumns,
	).FROM(
		CallEntries.
			INNER_JOIN(Calls, Calls.ID.EQ(CallEntries.CallID)).
			INNER_JOIN(Registrations, Registrations.ID.EQ(CallEntries.RegistrationID)).
			INNER_JOIN(Candidates, Candidates.ID.EQ(Registrations.CandidateID)).
			LEFT_JOIN(view.RegistrationPlacements, view.RegistrationPlacements.RegistrationID.EQ(Registrations.ID)).
			INNER_JOIN(Courses, Courses.ID.EQ(Registrations.CourseID)).
			INNER_JOIN(Quotas, Quotas.ID.EQ(Courses.QuotaID)),
	).WHERE(
		where,
	).ORDER_BY(
		Registrations.Ranking.ASC(),
		Registrations.ID.ASC(),
	)

	var result []callEntryDetailResult

	err := stmt.Query(db, &result)
	if err != nil {
		return nil, err
	}

	details := make([]*types.CallEntryDetail, len(result))
	for i := range result {
		r := &result[i]

		course, err := toCourseDomain(&r.Course, r.Quota.Name)
		if err != nil {
			return nil, err
		}

		details[i] = &types.CallEntryDetail{
			Entry:        toCallEntryDomain(&r.CallEntries, r.Call.Number),
			Registration: r.Registration.toRegistrationDomain(),
			Course:       course,
		}
	}

	return details, nil
}

func UpdateCallEntryOutcome(db qrm.DB, callID, registrationID int32, outcome types.CallEntryOutcome) error {
	stmt := CallEntries.UPDATE().
		SET(
			CallEntries.Outcome.SET(String(outcome.String())),
		).
		WHERE(
			CallEntries.CallID.EQ(Int32(callID)).
				AND(CallEntries.RegistrationID.EQ(Int32(registrationID))),
		)

	_, err := stmt.Exec(db)
	return err
}

func (d *Database) CallHasPendingEntries(callID int32) (bool, error) {
	stmt := SELECT(
		CallEntries.RegistrationID,
	).FROM(
		CallEntries,
	).WHERE(
		CallEntries.CallID.EQ(Int32(callID)).
			AND(CallEntries.Outcome.EQ(String(types.CallEntryOutcomePending.String()))),
	).LIMIT(1)

	var result []int32

	err := stmt.Query(d.db, &result)
	if err != nil {
		return false, err
	}

	return len(result) > 0, nil
}

// SelectionHasModifiedEntries reports whether any registration of the selection
// has an entry the operator already acted on (a recorded outcome).
func (d *Database) SelectionHasModifiedEntries(kind types.SelectionKind) (bool, error) {
	stmt := SELECT(
		CallEntries.RegistrationID,
	).FROM(
		CallEntries.
			INNER_JOIN(Registrations, Registrations.ID.EQ(CallEntries.RegistrationID)).
			INNER_JOIN(Selections, Selections.ID.EQ(Registrations.SelectionID)),
	).WHERE(
		Selections.Kind.EQ(String(kind.String())).
			AND(CallEntries.Outcome.NOT_EQ(String(types.CallEntryOutcomePending.String()))),
	).LIMIT(1)

	var result []int32

	err := stmt.Query(d.db, &result)
	if err != nil {
		return false, err
	}

	return len(result) > 0, nil
}

// CourseSemesterCount is the number of registrations placed in one course and semester.
type CourseSemesterCount struct {
	CourseID int32
	Semester int32
	Count    int32
}

// FetchOccupiedSeats counts, per course and semester, the registrations that
// currently hold a seat (placement pending or enrolled).
func FetchOccupiedSeats(db qrm.DB) ([]CourseSemesterCount, error) {
	stmt := SELECT(
		Registrations.CourseID.AS("course_semester_count.course_id"),
		view.RegistrationPlacements.Semester.AS("course_semester_count.semester"),
		COUNT(Registrations.ID).AS("course_semester_count.count"),
	).FROM(
		Registrations.
			INNER_JOIN(view.RegistrationPlacements, view.RegistrationPlacements.RegistrationID.EQ(Registrations.ID)),
	).WHERE(
		view.RegistrationPlacements.Outcome.IN(
			String(types.CallEntryOutcomePending.String()),
			String(types.CallEntryOutcomeEnrolled.String()),
		),
	).GROUP_BY(
		Registrations.CourseID,
		view.RegistrationPlacements.Semester,
	)

	var result []CourseSemesterCount

	err := stmt.Query(db, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// FetchUncalledWaitlist returns waitlist registrations that were never called.
func FetchUncalledWaitlist(db qrm.DB) ([]model.Registrations, error) {
	stmt := SELECT(
		Registrations.AllColumns,
	).FROM(
		Registrations.
			INNER_JOIN(Selections, Selections.ID.EQ(Registrations.SelectionID)),
	).WHERE(
		Selections.Kind.EQ(String(types.SelectionKindWaitlist.String())).
			AND(NOT(EXISTS(
				SELECT(CallEntries.RegistrationID).
					FROM(CallEntries).
					WHERE(CallEntries.RegistrationID.EQ(Registrations.ID)),
			))),
	)

	var result []model.Registrations

	err := stmt.Query(db, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type callEntryCountRow struct {
	CallID   int32
	Semester int32
	Kind     string
	Outcome  string
	Count    int32
}

// FetchCallSummaries returns every call (by number) with per-semester entry counts.
func (d *Database) FetchCallSummaries() ([]*types.CallSummary, error) {
	callsStmt := SELECT(
		Calls.AllColumns,
	).FROM(
		Calls,
	).ORDER_BY(
		Calls.Number.ASC(),
	)

	var calls []model.Calls

	err := callsStmt.Query(d.db, &calls)
	if err != nil {
		return nil, err
	}

	countsStmt := SELECT(
		CallEntries.CallID.AS("call_entry_count_row.call_id"),
		CallEntries.Semester.AS("call_entry_count_row.semester"),
		CallEntries.Kind.AS("call_entry_count_row.kind"),
		CallEntries.Outcome.AS("call_entry_count_row.outcome"),
		COUNT(CallEntries.RegistrationID).AS("call_entry_count_row.count"),
	).FROM(
		CallEntries,
	).GROUP_BY(
		CallEntries.CallID,
		CallEntries.Semester,
		CallEntries.Kind,
		CallEntries.Outcome,
	)

	var counts []callEntryCountRow

	err = countsStmt.Query(d.db, &counts)
	if err != nil {
		return nil, err
	}

	summaries := make([]*types.CallSummary, len(calls))
	byID := make(map[int32]*types.CallSummary, len(calls))

	for i := range calls {
		call := toCallDomain(&calls[i])
		summaries[i] = &types.CallSummary{
			ID:     call.ID,
			Status: call.Status,
			Number: call.Number,
			Semesters: []*types.CallSemesterSummary{
				{Semester: 1},
				{Semester: 2},
			},
		}
		byID[call.ID] = summaries[i]
	}

	for _, c := range counts {
		summary, ok := byID[c.CallID]
		if !ok || c.Semester < 1 || c.Semester > 2 {
			continue
		}

		if c.Outcome == types.CallEntryOutcomePending.String() {
			summary.Pending += c.Count
		}

		sem := summary.Semesters[c.Semester-1]

		switch types.CallEntryKind(c.Kind) {
		case types.CallEntryKindInitial:
			sem.Initial += c.Count
		case types.CallEntryKindWaitlist:
			sem.Waitlist += c.Count
		}
	}

	return summaries, nil
}
