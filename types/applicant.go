package types

type Applicant struct {
	ID           int64 `csv:"ApplicantID"`
	CPF          string
	Name         string
	SocialName   string `db:"social_name"`
	BirthDate    string
	Sex          string
	MotherName   string `db:"mother_name"`
	AddressLine  string `db:"address_line"`
	AddressLine2 string `db:"address_line2"`
	HouseNumber  string `db:"house_number"`
	Neighborhood string
	Municipality string
	State        string
	CEP          string
	Email        string
	Phone1       string
	Phone2       string
}

type ApplicantRepository interface {
	CreateApplicant(applicant *Applicant) error
	FindApplicantByID(id int64) (*Applicant, error)
	FindApplicantByApplicationID(id int64) (*Applicant, error)
	FindApplicants() ([]*Applicant, error)
	UpdateApplicant(id int64, update ApplicantUpdate) (*Applicant, error)
	DeleteApplicant(id int64) error
}

type ApplicantUpdate struct {
	Name       *string
	SocialName *string `json:"social_name"`
	CPF        *string
	Email      *string
	Phone1     *string
	Phone2     *string
	MotherName *string `json:"mother_name"`
	Birthdate  *string
}
