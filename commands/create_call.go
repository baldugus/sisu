package commands

import (
	"fmt"

	"github.com/go-jet/jet/v2/qrm"

	"github.com/baldugus/sisu/allocation"
	"github.com/baldugus/sisu/database"
	"github.com/baldugus/sisu/database/.gen/model"
	"github.com/baldugus/sisu/types"
)

// CreateCallCommand opens the next call: promotions to semester 1 first, then
// the waitlist (see package allocation).
type CreateCallCommand struct{}

func (cmd *CreateCallCommand) Execute(db *database.Database) error {
	return db.RunInTx(func(tx qrm.DB) error {
		next, err := planNextCall(db, tx)
		if err != nil {
			return err
		}

		if next.vacancies == 0 {
			return ErrAllCoursesFull{}
		}

		if len(next.entries) == 0 {
			return ErrNoCandidatesToCall{}
		}

		callID, err := database.CreateCall(tx, &types.Call{
			Number: next.number,
			Status: types.CallStatusCalling,
		})
		if err != nil {
			return fmt.Errorf("create call: %w", err)
		}

		for _, e := range next.entries {
			e.CallID = callID
		}

		return database.CreateCallEntries(tx, next.entries)
	})
}

// nextCall is what opening the next call would do.
type nextCall struct {
	number int32
	// vacancies is the number of free seats across both semesters.
	vacancies int32
	// entries are the entries to insert (without CallID).
	entries []*types.CallEntry
}

// planNextCall gathers the allocation input and runs the rule.
func planNextCall(db *database.Database, tx qrm.DB) (*nextCall, error) {
	hasOpenCall, err := db.HasOpenCall()
	if err != nil {
		return nil, err
	}

	if hasOpenCall {
		return nil, ErrOpenCallExists{}
	}

	lastCallNumber, err := db.GetLastCallNumber()
	if err != nil {
		return nil, err
	}

	if lastCallNumber == 0 {
		return nil, ErrNoCalls{}
	}

	in, err := allocationInput(db, tx)
	if err != nil {
		return nil, err
	}

	next := &nextCall{number: lastCallNumber + 1}

	for _, r := range allocation.Plan(in) {
		next.vacancies += r.Vacancies[0] + r.Vacancies[1]

		next.entries = append(next.entries, newEntries(r.Promoted, types.CallEntryKindPromotion, 1)...)
		next.entries = append(next.entries, newEntries(r.Waitlist[0], types.CallEntryKindWaitlist, 1)...)
		next.entries = append(next.entries, newEntries(r.Waitlist[1], types.CallEntryKindWaitlist, 2)...) //nolint: mnd
	}

	return next, nil
}

func allocationInput(db *database.Database, tx qrm.DB) (allocation.Input, error) {
	courses, err := db.FetchCourses()
	if err != nil {
		return allocation.Input{}, fmt.Errorf("fetch courses: %w", err)
	}

	occupied, err := database.FetchOccupiedSeats(tx)
	if err != nil {
		return allocation.Input{}, fmt.Errorf("fetch occupied seats: %w", err)
	}

	promotion, err := database.FetchPromotionCandidates(tx)
	if err != nil {
		return allocation.Input{}, fmt.Errorf("fetch promotion candidates: %w", err)
	}

	waitlist, err := database.FetchUncalledWaitlist(tx)
	if err != nil {
		return allocation.Input{}, fmt.Errorf("fetch waitlist: %w", err)
	}

	byCourse := make(map[int32]*allocation.Course, len(courses))
	in := allocation.Input{Courses: make([]allocation.Course, len(courses))}

	for i, c := range courses {
		in.Courses[i] = allocation.Course{ID: c.ID, SeatsPerSemester: c.Seats.PerSemester()}
		byCourse[c.ID] = &in.Courses[i]
	}

	for _, o := range occupied {
		if c, ok := byCourse[o.CourseID]; ok && o.Semester >= 1 && o.Semester <= 2 {
			c.Occupied[o.Semester-1] = o.Count
		}
	}

	for _, r := range promotion {
		if c, ok := byCourse[r.CourseID]; ok {
			c.Promotion = append(c.Promotion, toCandidate(r))
		}
	}

	for _, r := range waitlist {
		if c, ok := byCourse[r.CourseID]; ok {
			c.Waitlist = append(c.Waitlist, toCandidate(r))
		}
	}

	return in, nil
}

func toCandidate(r model.Registrations) allocation.Candidate {
	return allocation.Candidate{
		RegistrationID: r.ID,
		Ranking:        r.Ranking,
		CompositeScore: r.CompositeScore,
	}
}

func newEntries(ids []int32, kind types.CallEntryKind, semester int32) []*types.CallEntry {
	entries := make([]*types.CallEntry, len(ids))
	for i, id := range ids {
		entries[i] = &types.CallEntry{
			RegistrationID: id,
			Kind:           kind,
			Semester:       semester,
			Outcome:        types.CallEntryOutcomePending,
		}
	}

	return entries
}
