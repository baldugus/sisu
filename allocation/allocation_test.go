package allocation_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/baldugus/sisu/allocation"
)

func cands(ids ...int32) []allocation.Candidate {
	out := make([]allocation.Candidate, len(ids))
	for i, id := range ids {
		// Ranking follows the ID so the expected order is readable.
		out[i] = allocation.Candidate{RegistrationID: id, Ranking: id}
	}

	return out
}

func TestPlan(t *testing.T) {
	tests := []struct {
		name          string
		course        allocation.Course
		wantVacancies [2]int32
		wantPromoted  []int32
		wantWaitlist1 []int32
		wantWaitlist2 []int32
	}{
		{
			name:          "no vacancies calls no one",
			course:        allocation.Course{SeatsPerSemester: 2, Occupied: [2]int32{2, 2}, Promotion: cands(10), Waitlist: cands(20)},
			wantVacancies: [2]int32{0, 0},
		},
		{
			name:          "promotion fills semester 1 before the waitlist",
			course:        allocation.Course{SeatsPerSemester: 5, Occupied: [2]int32{3, 5}, Promotion: cands(11, 10), Waitlist: cands(20, 21)},
			wantVacancies: [2]int32{2, 0},
			wantPromoted:  []int32{10, 11},
		},
		{
			name:          "waitlist fills what promotion leaves, semester 1 first",
			course:        allocation.Course{SeatsPerSemester: 5, Occupied: [2]int32{2, 4}, Promotion: cands(10), Waitlist: cands(22, 20, 21, 23)},
			wantVacancies: [2]int32{3, 1},
			wantPromoted:  []int32{10},
			wantWaitlist1: []int32{20, 21},
			wantWaitlist2: []int32{22},
		},
		{
			name:          "short waitlist leaves semester 2 empty",
			course:        allocation.Course{SeatsPerSemester: 5, Occupied: [2]int32{3, 3}, Waitlist: cands(20, 21, 22)},
			wantVacancies: [2]int32{2, 2},
			wantWaitlist1: []int32{20, 21},
			wantWaitlist2: []int32{22},
		},
		{
			name:          "overbooked semester reports zero vacancies",
			course:        allocation.Course{SeatsPerSemester: 2, Occupied: [2]int32{3, 1}, Waitlist: cands(20)},
			wantVacancies: [2]int32{0, 1},
			wantWaitlist2: []int32{20},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := allocation.Plan(allocation.Input{Courses: []allocation.Course{tt.course}})

			assert.Len(t, got, 1)
			assert.Equal(t, tt.wantVacancies, got[0].Vacancies)
			assert.Equal(t, tt.wantPromoted, got[0].Promoted)
			assert.Equal(t, tt.wantWaitlist1, got[0].Waitlist[0])
			assert.Equal(t, tt.wantWaitlist2, got[0].Waitlist[1])
		})
	}
}

func TestSortByRanking(t *testing.T) {
	in := []allocation.Candidate{
		{RegistrationID: 1, Ranking: 0, CompositeScore: 90000},
		{RegistrationID: 2, Ranking: 5, CompositeScore: 70000},
		{RegistrationID: 3, Ranking: 0, CompositeScore: 80000},
		{RegistrationID: 4, Ranking: 4, CompositeScore: 60000},
		{RegistrationID: 5, Ranking: 5, CompositeScore: 75000},
	}

	got := allocation.SortByRanking(in)

	ids := make([]int32, len(got))
	for i, c := range got {
		ids[i] = c.RegistrationID
	}

	assert.Equal(t, []int32{4, 5, 2, 1, 3}, ids)
	assert.Equal(t, int32(1), in[0].RegistrationID, "input must not be reordered")
}
