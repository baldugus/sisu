package commands

import (
	"fmt"

	"github.com/go-jet/jet/v2/qrm"

	"github.com/baldugus/sisu/allocation"
	"github.com/baldugus/sisu/database"
	"github.com/baldugus/sisu/database/.gen/model"
	"github.com/baldugus/sisu/types"
)

// PreviewCallCommand computes the next call without creating it.
type PreviewCallCommand struct{}

func (cmd *PreviewCallCommand) Execute(db *database.Database) (*types.CallPlan, error) {
	plan, _, err := planNextCall(db, db.DB())
	return plan, err
}

// CreateCallCommand opens the next call, calling exactly who PreviewCallCommand
// shows: promotions to semester 1 first, then the waitlist (see package allocation).
type CreateCallCommand struct{}

func (cmd *CreateCallCommand) Execute(db *database.Database) error {
	return db.RunInTx(func(tx qrm.DB) error {
		plan, entries, err := planNextCall(db, tx)
		if err != nil {
			return err
		}

		if plan.Vacancies() == 0 {
			return ErrAllCoursesFull{}
		}

		if len(entries) == 0 {
			return ErrNoCandidatesToCall{}
		}

		callID, err := database.CreateCall(tx, &types.Call{
			Number: plan.Number,
			Status: types.CallStatusCalling,
		})
		if err != nil {
			return fmt.Errorf("create call: %w", err)
		}

		for _, e := range entries {
			e.CallID = callID
		}

		return database.CreateCallEntries(tx, entries)
	})
}

// planNextCall gathers the allocation input, runs the rule and returns both the
// readable plan and the entries that would be inserted (without CallID).
func planNextCall(db *database.Database, tx qrm.DB) (*types.CallPlan, []*types.CallEntry, error) {
	hasOpenCall, err := db.HasOpenCall()
	if err != nil {
		return nil, nil, err
	}

	if hasOpenCall {
		return nil, nil, ErrOpenCallExists{}
	}

	lastCallNumber, err := db.GetLastCallNumber()
	if err != nil {
		return nil, nil, err
	}

	if lastCallNumber == 0 {
		return nil, nil, ErrNoCalls{}
	}

	semesters, err := database.FetchSemesters(tx)
	if err != nil {
		return nil, nil, fmt.Errorf("fetch semesters: %w", err)
	}

	var open [2]bool
	for _, s := range semesters {
		open[s.Number-1] = s.Status == types.SemesterStatusOpen
	}

	if !open[0] && !open[1] {
		return nil, nil, ErrAllSemestersClosed{}
	}

	in, courses, registrations, err := allocationInput(db, tx, open)
	if err != nil {
		return nil, nil, err
	}

	results := allocation.Plan(in)

	plan := &types.CallPlan{
		Number:    lastCallNumber + 1,
		Semesters: semesters,
	}

	var entries []*types.CallEntry

	for i, r := range results {
		coursePlan := &types.CoursePlan{
			Course:     courses[i],
			Vacancies1: r.Vacancies[0],
			Vacancies2: r.Vacancies[1],
			Promoted:   pick(registrations, r.Promoted),
			Waitlist1:  pick(registrations, r.Waitlist[0]),
			Waitlist2:  pick(registrations, r.Waitlist[1]),
		}
		plan.Courses = append(plan.Courses, coursePlan)
		plan.Promoted += int32(len(r.Promoted))     //nolint: gosec
		plan.Waitlist1 += int32(len(r.Waitlist[0])) //nolint: gosec
		plan.Waitlist2 += int32(len(r.Waitlist[1])) //nolint: gosec

		entries = append(entries, newEntries(r.Promoted, types.CallEntryKindPromotion, 1)...)
		entries = append(entries, newEntries(r.Waitlist[0], types.CallEntryKindWaitlist, 1)...)
		entries = append(entries, newEntries(r.Waitlist[1], types.CallEntryKindWaitlist, 2)...) //nolint: mnd
	}

	return plan, entries, nil
}

func allocationInput(
	db *database.Database,
	tx qrm.DB,
	open [2]bool,
) (allocation.Input, []*types.Course, map[int32]*types.Registration, error) {
	courses, err := db.FetchCourses()
	if err != nil {
		return allocation.Input{}, nil, nil, fmt.Errorf("fetch courses: %w", err)
	}

	occupied, err := database.FetchOccupiedSeats(tx)
	if err != nil {
		return allocation.Input{}, nil, nil, fmt.Errorf("fetch occupied seats: %w", err)
	}

	promotion, err := database.FetchPromotionCandidates(tx)
	if err != nil {
		return allocation.Input{}, nil, nil, fmt.Errorf("fetch promotion candidates: %w", err)
	}

	waitlist, err := database.FetchUncalledWaitlist(tx)
	if err != nil {
		return allocation.Input{}, nil, nil, fmt.Errorf("fetch waitlist: %w", err)
	}

	byCourse := make(map[int32]*allocation.Course, len(courses))
	in := allocation.Input{Open: open, Courses: make([]allocation.Course, len(courses))}

	for i, c := range courses {
		in.Courses[i] = allocation.Course{ID: c.ID, Seats: c.Seats}
		byCourse[c.ID] = &in.Courses[i]
	}

	for _, o := range occupied {
		if c, ok := byCourse[o.CourseID]; ok && o.Semester >= 1 && o.Semester <= 2 {
			c.Occupied[o.Semester-1] = o.Count
		}
	}

	var ids []int32

	for _, r := range promotion {
		if c, ok := byCourse[r.CourseID]; ok {
			c.Promotion = append(c.Promotion, toCandidate(r))
			ids = append(ids, r.ID)
		}
	}

	for _, r := range waitlist {
		if c, ok := byCourse[r.CourseID]; ok {
			c.Waitlist = append(c.Waitlist, toCandidate(r))
			ids = append(ids, r.ID)
		}
	}

	regs, err := database.FetchRegistrationsByIDs(tx, ids)
	if err != nil {
		return allocation.Input{}, nil, nil, fmt.Errorf("fetch registrations: %w", err)
	}

	registrations := make(map[int32]*types.Registration, len(regs))
	for _, r := range regs {
		registrations[r.ID] = r
	}

	return in, courses, registrations, nil
}

func toCandidate(r model.Registrations) allocation.Candidate {
	return allocation.Candidate{
		RegistrationID: r.ID,
		Ranking:        r.Ranking,
		CompositeScore: r.CompositeScore,
	}
}

func pick(registrations map[int32]*types.Registration, ids []int32) []*types.Registration {
	out := make([]*types.Registration, 0, len(ids))
	for _, id := range ids {
		out = append(out, registrations[id])
	}

	return out
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
