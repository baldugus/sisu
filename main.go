package main

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"path"

	"changeme/repository"

	"github.com/alecthomas/kong"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	"github.com/wailsapp/wails/v2"
	wailsoptions "github.com/wailsapp/wails/v2/pkg/options"
	wailsassetserver "github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"go.uber.org/zap"
	_ "modernc.org/sqlite"
)

//go:embed all:frontend/dist
var assets embed.FS

//o:embed repository/migrations
//var migrations embed.FS

// This needs to be done somewhere else.
//func main() {
//	os.Exit(run())
//}

// FIXME: too long fsr.
func run() int { //nolint: funlen
	var CLI struct {
		logFile  string `help:"Path to log file."        short:"o"`
		logLevel string `help:"Log level."               short:"l" default:"error"` //nolint:tagalign
		dbFile   string `help:"Path to sqlite database." short:"d"`
	}

	kong.Parse(&CLI)

	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	configDir := path.Join(userConfigDir, "sisu")

	if err := os.MkdirAll(configDir, 0o755); err != nil { //nolint: gomnd
		panic(err)
	}

	if CLI.logFile == "" {
		CLI.logFile = path.Join(configDir, "sisu.log")
	}

	logger := configLogger(CLI.logFile, CLI.logLevel)
	defer func() { _ = logger.Sync() }()

	logger.Info("starting")

	if CLI.dbFile == "" {
		CLI.dbFile = path.Join(configDir, "sisu.db")
	}

	sqliteDB, err := configDB(CLI.dbFile)
	repositoryService := NewRepositoryService(CLI.dbFile, sqliteDB)
	defer func() { _ = repositoryService.Close() }()

	if err != nil {
		logger.Sugar().Errorw("config db", "error", err)

		return 1
	}

	DB := repository.NewDB(sqliteDB)
	applicantRepo := repository.NewApplicantRepository(DB)
	applicationRepo := repository.NewApplicationRepository(DB)
	classRepo := repository.NewClassRepository(DB)
	rollcallRepo := repository.NewRollcallRepository(DB)
	selectionRepo := repository.NewSelectionRepository(DB)

	// defer applicationRepository.Close(sqliteDB)
	sisu := SISU{
		applicantRepo:   *applicantRepo,
		applicationRepo: *applicationRepo,
		classRepo:       *classRepo,
		rollcallRepo:    *rollcallRepo,
		selectionRepo:   *selectionRepo,
		service:         repositoryService,
		l:               logger.Sugar(),
	}

	var app App
	app.sisu = sisu

	var assetServerOptions wailsassetserver.Options
	assetServerOptions.Assets = assets

	// Create application with options
	var options wailsoptions.App
	options.Title = "sisu"
	options.Width = 1024
	options.Height = 768
	options.AssetServer = &assetServerOptions
	options.BackgroundColour = &wailsoptions.RGBA{R: 27, G: 38, B: 54, A: 1}
	options.OnStartup = app.startup
	options.Bind = []any{&app}

	if err := wails.Run(&options); err != nil {
		logger.Sugar().Errorw("app start error", "error", err)
	}

	return 0
}

func configLogger(logFile string, logLevel string) *zap.Logger {
	level, err := zap.ParseAtomicLevel(logLevel)
	if err != nil {
		panic(err)
	}

	cfg := zap.NewProductionConfig()

	cfg.OutputPaths = []string{"stdout", logFile}
	cfg.ErrorOutputPaths = []string{"stderr", logFile}
	cfg.Level = level

	return zap.Must(cfg.Build())
}

func configDB(dbFile string) (*sqlx.DB, error) {
	db, err := sqlx.Open(
		"sqlite",
		fmt.Sprintf("%s?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)", dbFile),
	)
	if err != nil {
		return nil, fmt.Errorf("sqlx open: %w", err)
	}

	if err := initDB(db); err != nil {
		return nil, fmt.Errorf("init repository: %w", err)
	}

	return db, nil
}

func initDB(db *sqlx.DB) error {
	filesystem, err := iofs.New(migrations, "repository/migrations")
	if err != nil {
		return fmt.Errorf("new iofs: %w", err)
	}

	var sqliteConfig sqlite.Config

	s, err := sqlite.WithInstance(db.DB, &sqliteConfig)
	if err != nil {
		return fmt.Errorf("sqlite with instance: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", filesystem, "sqlite", s)
	if err != nil {
		return fmt.Errorf("migrate new with instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate up: %w", err)
	}

	return nil
}
