package repository_test

/*
import (
	"database/sql"
	"embed"
	"errors"
	"testing"

	"changeme/repository"
	"changeme/types"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

//go:embed migrations
var migrations embed.FS

// TODO: add tests for all cases and restrictions.
func TestRepository(t *testing.T) {
	t.Parallel()

	t.Run("simple insert", SimpleInsert)
	t.Run("insert 2 before 1", InsertTwoBeforeOne)
	t.Run("insert 2 times", InsertTwoTimes)
	t.Run("select approved", SelectApproved)
	t.Run("select not approved", SelectNotApproved)
	t.Run("delete application", DeleteApplication)
	t.Run("update rollcall", RollCallUpdate)
}

// TODO: add tests for when there is and isnt enrolled.
func DeleteApplication(t *testing.T) {
	t.Parallel()

	db, err := sqlx.Open(
		"sqlite",
		":memory:?_pragma=foreign_keys(1)",
	)
	if err != nil {
		t.Fatalf("sqlx open: %v", err)
	}

	initializeDatabase(t, db)

	applications := applicationTestData()

	testSelection := &types.Selection{
		Name:         "foo",
		Kind:         1,
		Date:         "2023-06-01",
		Applications: applications[1:2],
	}

	repo := repository.NewRepository(db)
	if err := repo.SaveSelection(testSelection); err != nil {
		t.Fatalf("save selection: %v", err)
	}

	err = repo.DeleteSelection(1)
	if err != nil && !errors.Is(err, repository.ErrDoneRollCallExists) {
		t.Fatalf("delete selection: %v", err)
	}

	_, err = repo.FindSelection(1)
	if err == nil || !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("find selection returned values")
	}
}

func SelectNotApproved(t *testing.T) {
	t.Parallel()

	db, err := sqlx.Open(
		"sqlite",
		":memory:?_pragma=foreign_keys(1)",
	)
	if err != nil {
		t.Fatalf("sqlx open: %v", err)
	}

	initializeDatabase(t, db)

	applications := applicationTestData()

	testSelection := &types.Selection{
		Name:         "foo",
		Kind:         1,
		Date:         "2023-06-01",
		Applications: applications,
	}

	repo := repository.NewRepository(db)
	if err := repo.SaveSelection(testSelection); err != nil {
		t.Fatalf("save selection: %v", err)
	}

	result, err := repo.SearchApplications(nil, []string{"ABSENT"}, nil)
	if err != nil {
		t.Fatalf("find selection: %v", err)
	}

	if !cmp.Equal(testSelection.Applications[2:3], result) {
		t.Error(cmp.Diff(testSelection.Applications[2:3], result))
	}
}

// TODO: status test should get all statuses.
func SelectApproved(t *testing.T) {
	t.Parallel()

	db, err := sqlx.Open(
		"sqlite",
		":memory:?_pragma=foreign_keys(1)",
	)
	if err != nil {
		t.Fatalf("sqlx open: %v", err)
	}

	initializeDatabase(t, db)

	applications := applicationTestData()

	testSelection := &types.Selection{
		Name:         "foo",
		Kind:         1,
		Date:         "2023-06-01",
		Applications: applications,
	}

	repo := repository.NewRepository(db)
	if err := repo.SaveSelection(testSelection); err != nil {
		t.Fatalf("save selection: %v", err)
	}

	result, err := repo.SearchApplications(nil, []string{"APPROVED"}, nil)
	if err != nil {
		t.Fatalf("find selection: %v", err)
	}

	if !cmp.Equal(testSelection.Applications[1:2], result) {
		t.Error(cmp.Diff(testSelection.Applications[1:2], result))
	}
}

func InsertTwoBeforeOne(t *testing.T) {
	t.Parallel()

	db, err := sqlx.Open(
		"sqlite",
		":memory:?_pragma=foreign_keys(1)",
	)
	if err != nil {
		t.Fatalf("sqlx open: %v", err)
	}

	initializeDatabase(t, db)

	testSelection := &types.Selection{
		Name:         "foo",
		Kind:         2,
		Date:         "2023",
		Applications: []*types.Application{},
	}

	repo := repository.NewRepository(db)

	err = repo.SaveSelection(testSelection)
	if err == nil {
		t.Fatalf("expected save selection to err but it didn't")
	}

	if err != nil && !errors.Is(err, repository.ErrImportMissing) {
		t.Fatalf("wrong err received from save selection: %v", err)
	}
}

func InsertTwoTimes(t *testing.T) {
	t.Parallel()

	db, err := sqlx.Open(
		"sqlite",
		":memory:?_pragma=foreign_keys(1)",
	)
	if err != nil {
		t.Fatalf("sqlx open: %v", err)
	}

	initializeDatabase(t, db)

	testSelection := &types.Selection{
		Name:         "foo",
		Kind:         1,
		Date:         "2023",
		Applications: []*types.Application{},
	}

	repo := repository.NewRepository(db)

	err = repo.SaveSelection(testSelection)
	if err != nil {
		t.Fatalf("save selection: %v", err)
	}

	err = repo.SaveSelection(testSelection)
	if err == nil {
		t.Fatalf("expected save selection to err but it didn't")
	}

	if err != nil && !errors.Is(err, repository.ErrImportExists) {
		t.Fatalf("wrong err received from save selection: %v", err)
	}

	testSelection.Kind = 2

	err = repo.SaveSelection(testSelection)
	if err != nil {
		t.Fatalf("save selection: %v", err)
	}

	err = repo.SaveSelection(testSelection)
	if err == nil {
		t.Fatalf("expected save selection to err but it didn't")
	}

	if err != nil && !errors.Is(err, repository.ErrImportExists) {
		t.Fatalf("wrong err received from save selection: %v", err)
	}
}

// FIXME: too long and its wrong.
func SimpleInsert(t *testing.T) { //nolint: funlen, cyclop
	t.Parallel()

	db, err := sqlx.Open(
		"sqlite",
		":memory:?_pragma=foreign_keys(1)",
	)
	if err != nil {
		t.Fatalf("sqlx open: %v", err)
	}

	initializeDatabase(t, db)

	applications := applicationTestData()

	testSelection := &types.Selection{
		Name:         "foo",
		Kind:         1,
		Date:         "2023-06-01",
		Applications: applications,
	}

	repo := repository.NewRepository(db)
	if err := repo.SaveSelection(testSelection); err != nil {
		t.Fatalf("save selection: %v", err)
	}

	result, err := repo.FindSelection(1)
	if err != nil {
		t.Fatalf("find selection: %v", err)
	}

	if !cmp.Equal(testSelection, result) {
		t.Error(cmp.Diff(testSelection, result))
	}

	rollcallApplications, err := repo.SearchApplications(nil, nil, []int64{1})
	if err != nil {
		t.Fatalf("find applications by rollcall number: %v", err)
	}

	if !cmp.Equal(applications, rollcallApplications) {
		t.Error(cmp.Diff(applications, rollcallApplications))
	}

	rollcalls, err := repo.SearchRollCalls([]int64{1}, nil, nil)
	if err != nil {
		t.Fatalf("search rollcalls: %v", err)
	}

	if len(rollcalls) != 1 {
		t.Errorf("wrong rollcall count")
	}

	if rollcalls[0].Status != "CALLING" {
		t.Error("rollcall is in wrong status")
	}

	kinds := []types.SelectionKind{types.ApprovedSelection}
	statuses := []string{"ENROLLED", "APPROVED"}

	moreApplications, err := repo.SearchApplications(kinds, statuses, nil)
	if err != nil {
		t.Fatalf("search applications: %v", err)
	}

	if !cmp.Equal(applications[0:2], moreApplications) {
		t.Error(cmp.Diff(applications[0:2], moreApplications))
	}

	selection := repository.SelectionArgument{types.ApprovedSelection}
	status := repository.StatusArgument{"CALLING"}
	number := repository.NumberArgument{}

	rollcalls, err = repo.SearchRollCalls(number, selection, status)
	if err != nil {
		t.Fatalf("find rollcalls by selection: %v", err)
	}

	if len(rollcalls) != 1 {
		t.Errorf("failed")
	}
}

func initializeDatabase(tb testing.TB, db *sqlx.DB) {
	tb.Helper()

	filesystem, err := iofs.New(migrations, "migrations")
	if err != nil {
		tb.Fatalf("iofs new: %v", err)
	}

	s, err := sqlite.WithInstance(db.DB, &sqlite.Config{}) //nolint: exhaustruct
	if err != nil {
		tb.Fatalf("sqlite with instance: %v", err)
	}

	m, err := migrate.NewWithInstance("iofs", filesystem, "sqlite", s)
	if err != nil {
		tb.Fatalf("migrate new with instance: %v", err)
	}

	if err := m.Up(); err != nil {
		tb.Fatalf("migrate up: %v", err)
	}
}

func applicationTestData() []*types.Application { //nolint: funlen
	return []*types.Application{
		{
			Status:       "ENROLLED",
			EnrollmentID: "221035551624",
			Option:       1,
			Applicant: types.Applicant{
				CPF:          "11122233345",
				Name:         "Gustavo Balduino",
				BirthDate:    "2001-10-02",
				Sex:          "M",
				MotherName:   "Senhora Balduino",
				AddressLine:  "Rua Rio",
				AddressLine2: "Rio",
				HouseNumber:  "22",
				Neighborhood: "Rio",
				Municipality: "Rio",
				State:        "Rio",
				CEP:          "22222222",
				Email:        "gustavo@balduino.com.br",
				Phone1:       "222222222",
				Phone2:       "222222222",
			},
			Scores: types.Scores{
				Languages:       599.0,
				Humanities:      599.1,
				NaturalSciences: 599.2,
				Mathematics:     599.3,
				Essay:           599.4,
				Composite:       599.5,
				Ranking:         2,
			},
			Class: types.Class{
				Period:       "Manhã",
				Quota:        "Ampla Concorrência",
				Seats:        26,
				MinimumScore: 500.0,
			},
		},
		{
			Status:       "APPROVED",
			EnrollmentID: "221035551623",
			Option:       1,
			Applicant: types.Applicant{
				CPF:          "11122233344",
				Name:         "Andre Lemos",
				BirthDate:    "2000-24-04",
				Sex:          "M",
				MotherName:   "Senhora Lemos",
				AddressLine:  "Rua Maricá",
				AddressLine2: "Maricá",
				HouseNumber:  "13",
				Neighborhood: "Maricá",
				Municipality: "Maricá",
				State:        "Maricá",
				CEP:          "13131313",
				Email:        "andre@lemos.com.br",
				Phone1:       "313131313",
				Phone2:       "131313131",
			},
			Scores: types.Scores{
				Languages:       600.0,
				Humanities:      600.1,
				NaturalSciences: 600.2,
				Mathematics:     600.3,
				Essay:           600.4,
				Composite:       600.5,
				Ranking:         1,
			},
			Class: types.Class{
				Period:       "Manhã",
				Quota:        "Ampla Concorrência",
				Seats:        26,
				MinimumScore: 500.0,
			},
		},
		{
			Status:       "ABSENT",
			EnrollmentID: "221035551625",
			Option:       1,
			Applicant: types.Applicant{
				CPF:          "11122233346",
				Name:         "Maria Claudia",
				BirthDate:    "2005-05-05",
				Sex:          "F",
				MotherName:   "Senhora Claudia",
				AddressLine:  "Rua Clarimundo de Melo",
				AddressLine2: "Quintino",
				HouseNumber:  "01",
				Neighborhood: "Quintino",
				Municipality: "Quintino",
				State:        "Quintino",
				CEP:          "11111111",
				Email:        "maria@claudia.com.br",
				Phone1:       "111111111",
				Phone2:       "222222222",
			},
			Scores: types.Scores{
				Languages:       602.0,
				Humanities:      602.1,
				NaturalSciences: 602.2,
				Mathematics:     602.3,
				Essay:           602.4,
				Composite:       602.5,
				Ranking:         4,
			},
			Class: types.Class{
				Period:       "Manhã",
				Quota:        "Ampla Concorrência",
				Seats:        26,
				MinimumScore: 500.0,
			},
		},
	}
}

// FIXME: why is it this big?
func RollCallUpdate(t *testing.T) { //nolint: funlen
	gofakeit.Seed(11)

	applicants := make([]types.Applicant, 3)

	for i := range applicants {
		if err := gofakeit.Struct(&applicants[i]); err != nil {
			t.Fatalf("gofakeit struct: %v", err)
		}
	}

	type Test struct {
		selection *types.Selection
		Expect    []*types.RollCall
		status    string
		err       error
	}

	tests := map[string]Test{
		"success": {
			err:    nil,
			status: "DONE",
			selection: &types.Selection{
				Name: "approved",
				Kind: types.ApprovedSelection,
				Date: "", // FIXME: unused for now
				Applications: []*types.Application{
					{
						Status:       "ABSENT",
						EnrollmentID: "1",
						Option:       1,
						Applicant:    applicants[0],
					},
					{
						Status:       "ENROLLED",
						EnrollmentID: "2",
						Option:       2,
						Applicant:    applicants[1],
					},
					{
						Status:       "ENROLLED",
						EnrollmentID: "3",
						Option:       1,
						Applicant:    applicants[2],
					},
				},
			},
			Expect: []*types.RollCall{
				{
					Status: "DONE",
					Number: 1,
				},
			},
		},
		"fail": {
			err:    repository.ErrRollCallClose,
			status: "DONE",
			selection: &types.Selection{
				Name: "approved",
				Kind: types.ApprovedSelection,
				Date: "", // FIXME: unused for now
				Applications: []*types.Application{
					{
						Status:       "APPROVED",
						EnrollmentID: "1",
						Option:       1,
						Applicant:    applicants[0],
					},
					{
						Status:       "ENROLLED",
						EnrollmentID: "2",
						Option:       2,
						Applicant:    applicants[1],
					},
					{
						Status:       "ENROLLED",
						EnrollmentID: "3",
						Option:       1,
						Applicant:    applicants[2],
					},
				},
			},
			Expect: []*types.RollCall{
				{
					Status: "CALLING",
					Number: 1,
				},
			},
		},
	}

	for name, test := range tests {
		func(name string, test Test) {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				db, err := sqlx.Open(
					"sqlite",
					":memory:?_pragma=foreign_keys(1)",
				)
				if err != nil {
					t.Fatalf("sqlx open: %v", err)
				}

				initializeDatabase(t, db)

				repo := repository.NewRepository(db)
				if err := repo.SaveSelection(test.selection); err != nil {
					t.Fatalf("save selection: %v", err)
				}

				if err := repo.UpdateRollCallStatus(1, test.status); !errors.Is(err, test.err) {
					t.Errorf("expected error (%v) got (%v)", test.err, err)
				}

				rollcall, err := repo.SearchRollCalls(repository.NumberArgument{1}, nil, nil)
				if err != nil {
					t.Fatalf("search rollcalls: %v", err)
				}

				if !cmp.Equal(rollcall, test.Expect) {
					t.Errorf(cmp.Diff(rollcall, test.Expect))
				}
			})
		}(name, test)
	}
}
*/
