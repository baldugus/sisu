package main

import (
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"os"

	repository2 "changeme/internal/repository"
	"changeme/repository"
	"changeme/types"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	_ "modernc.org/sqlite"
)

//go:embed repository/migrations
var migrations embed.FS

func main() {
	db, err := sqlx.Connect("sqlite", "test.db?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)")
	if err != nil {
		slog.Error("connect to db", slog.Any("err", err))
		os.Exit(1)
	}

	filesystem, err := iofs.New(migrations, "repository/migrations")
	if err != nil {
		slog.Error("new iofs", slog.Any("err", err))
		os.Exit(1)
	}

	var sqliteConfig sqlite.Config

	s, err := sqlite.WithInstance(db.DB, &sqliteConfig)
	if err != nil {
		slog.Error("new migrate sqlite instance", slog.Any("err", err))
		os.Exit(1)
	}

	m, err := migrate.NewWithInstance("iofs", filesystem, "sqlite", s)
	if err != nil {
		slog.Error("new migrate instance", slog.Any("err", err))
		os.Exit(1)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error("migrate up", slog.Any("err", err))
		os.Exit(1)
	}

	DB := repository.NewDB(db)
	applicantRepo := repository.NewApplicantRepository(DB)
	applicationRepo := repository.NewApplicationRepository(DB)
	classRepo := repository.NewClassRepository(DB)
	rollcallRepo := repository.NewRollcallRepository(DB)
	selectionRepo := repository.NewSelectionRepository(DB)

	// new repo
	repo := repository2.Repository{DB: db}

	cfg := zap.NewProductionConfig()
	logger := zap.Must(cfg.Build())
	// defer applicationRepository.Close(sqliteDB)
	sisu := SISU{
		repo:            &repo,
		applicantRepo:   *applicantRepo,
		applicationRepo: *applicationRepo,
		classRepo:       *classRepo,
		rollcallRepo:    *rollcallRepo,
		selectionRepo:   *selectionRepo,
		service:         nil,
		l:               logger.Sugar(),
	}

	err = sisu.LoadSelection("/home/gustavo/Downloads/listagem-alunos-aprovados-ies-5016-8937.csv", types.ApprovedSelection)
	fmt.Println(err)

	os.Exit(0)
}
