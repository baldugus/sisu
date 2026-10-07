package types

// CallPlan is what opening the next call would do: who is called, per course.
// The same plan is shown as a preview and then applied, so the operator sees
// exactly what will happen.
type CallPlan struct {
	Number    int32
	Semesters []*Semester
	Courses   []*CoursePlan
	// Totals across all courses.
	Promoted  int32
	Waitlist1 int32
	Waitlist2 int32
}

// Total is the number of registrations the plan calls.
func (p *CallPlan) Total() int32 {
	return p.Promoted + p.Waitlist1 + p.Waitlist2
}

// Vacancies is the number of free seats across open semesters.
func (p *CallPlan) Vacancies() int32 {
	var total int32
	for _, c := range p.Courses {
		total += c.Vacancies1 + c.Vacancies2
	}

	return total
}

type CoursePlan struct {
	Course     *Course
	Vacancies1 int32
	Vacancies2 int32
	// Promoted are semester-2 students offered a seat in semester 1.
	Promoted  []*Registration
	Waitlist1 []*Registration
	Waitlist2 []*Registration
}
