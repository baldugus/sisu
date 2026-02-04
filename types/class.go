package types

type Class struct {
	ID           int64 `csv:"ClassID"`
	Seats        int
	MinimumScore float64 `db:"minimum_score"`
	Period       `db:"period"`
	Quota        `db:"quota"`
}

type Period struct {
	ID   int64  `csv:"PeriodID"`
	Name string `json:"Period" csv:"PeriodName"`
}
type Quota struct {
	ID   int64  `csv:"QuotaID"`
	Name string `json:"Quota" csv:"QuotaName"`
}

type ClassRepository interface {
	FindClassByID(id int64) (*Class, error)
	FindClassesByPeriodID(periodID int64) ([]*Class, error)
	FindClasses() ([]*Class, error)
	CreateClass(class *Class) error
	DeleteClass(id int64) error
}
