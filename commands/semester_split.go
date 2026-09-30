package commands

import (
	"github.com/baldugus/sisu/csvparser"
	"github.com/baldugus/sisu/types"
)

// SplitApprovedBySemester divides a course's approved registrations across the
// year's two semester intakes (SiSU runs one annual selection; the institution
// keeps two intakes): candidates ranked within the first half of the course's
// seats go to Semester 1, the rest to Semester 2. The seat count must be even,
// otherwise ErrOddSeatsCount is returned. Input order is preserved in each half.
func SplitApprovedBySemester(seats int32, registrations []*types.Registration) (sem1, sem2 []*types.Registration, err error) {
	if seats%2 != 0 {
		return nil, nil, ErrOddSeatsCount{Count: seats}
	}

	for _, reg := range registrations {
		if reg.Ranking <= seats/2 {
			sem1 = append(sem1, reg)
		} else {
			sem2 = append(sem2, reg)
		}
	}

	return sem1, sem2, nil
}

type courseKey struct {
	period types.CoursePeriod
	quota  string
}

// approvedSemesterNumbers groups parsed approved registrations by course and
// applies SplitApprovedBySemester to each group, returning the semester number
// (1 or 2) assigned to every registration.
func approvedSemesterNumbers(parsed []*csvparser.ParsedRegistration) (map[*types.Registration]int32, error) {
	var keys []courseKey
	seatsByCourse := make(map[courseKey]int32)
	regsByCourse := make(map[courseKey][]*types.Registration)

	for _, p := range parsed {
		key := courseKey{period: p.Course.Period, quota: p.Course.Quota}
		if _, ok := regsByCourse[key]; !ok {
			keys = append(keys, key)
			seatsByCourse[key] = p.Course.Seats
		}
		regsByCourse[key] = append(regsByCourse[key], p.Registration)
	}

	numbers := make(map[*types.Registration]int32, len(parsed))
	for _, key := range keys {
		sem1, sem2, err := SplitApprovedBySemester(seatsByCourse[key], regsByCourse[key])
		if err != nil {
			return nil, err
		}

		for _, reg := range sem1 {
			numbers[reg] = 1
		}
		for _, reg := range sem2 {
			numbers[reg] = 2
		}
	}

	return numbers, nil
}
