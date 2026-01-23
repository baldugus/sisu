package main

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/baldugus/sisu/database"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"github.com/alecthomas/kong"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	_ "modernc.org/sqlite"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed database/migrations
var migrations embed.FS

// This needs to be done somewhere else.
func main() {
	os.Exit(run())
}

// FIXME: too long fsr.
func run() int { //nolint: funlen
	var CLI struct {
		LogFile  string `help:"Path to log file."        short:"o"`
		LogLevel string `help:"Log level."               short:"l" default:"info"` //nolint:tagalign
		DBFile   string `help:"Path to sqlite database." short:"d"`
		File     string `help:"file"                     short:"f"`
	}

	kong.Parse(&CLI, kong.ConfigureHelp(kong.HelpOptions{Compact: true, Summary: false}))

	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	configDir := path.Join(userConfigDir, "sisu")

	if err := os.MkdirAll(configDir, 0o755); err != nil { //nolint: mnd
		panic(err)
	}

	if CLI.LogFile == "" {
		CLI.LogFile = path.Join(configDir, "sisu.log")
	}

	logger := configLogger(CLI.LogFile, CLI.LogLevel)
	zap.ReplaceGlobals(logger)
	defer func() { _ = logger.Sync() }()

	logger.Info("SISU application starting",
		zap.String("log_file", CLI.LogFile),
		zap.String("log_level", CLI.LogLevel),
		zap.String("config_dir", configDir),
	)

	if CLI.DBFile == "" {
		CLI.DBFile = path.Join(configDir, "sisu.db")
	}

	logger.Info("initializing database", zap.String("db_file", CLI.DBFile))
	sqliteDB, err := configDB(CLI.DBFile)
	if err != nil {
		logger.Sugar().Errorw("config db", "error", err)

		return 1
	}
	logger.Info("database initialized successfully")

	qrm.GlobalConfig.StrictScan = true

	newDB := database.NewDatabase(sqliteDB, CLI.DBFile)
	defer func() { _ = newDB.Close() }()
	sisu := SISU{
		database: newDB,
	}

	// --- CMD LoadSelectionCommand
	/*
		contents, err := os.ReadFile(CLI.File)
		if err != nil {
			panic(err)
		}

		cmd := commands.LoadSelectionCommand{
			Year:     2025,
			Semester: 1,
			Path:     CLI.File,
			Content:  contents,
			Kind:     types.SelectionKindApproved,
		}

		if err := cmd.Execute(sisu.database); err != nil {
			panic(err)
		}
	*/
	// ---

	// --- CMD Generic Fetch
	//cmd := commands.FetchSelectionCommand{
	//	Kind: types.SelectionKindApproved,
	//}

	//r, err := cmd.Execute(sisu.database)
	//if err != nil {
	//	panic(err)
	//}

	//j, err := json.MarshalIndent(r, "", "\t")
	//if err != nil {
	//	panic(err)
	//}

	//fmt.Println(string(j))

	// ---

	// --- CMD CloseCall
	//cmd := commands.CloseCallCommand{
	//	ID: 1,
	//}
	//err = cmd.Execute(sisu.database)
	//if err != nil {
	//	panic(err)
	//}

	var app App
	app.sisu = sisu

	var assetServerOptions assetserver.Options
	assetServerOptions.Assets = assets

	// Create application with wailsOptions
	var wailsOptions options.App
	wailsOptions.Title = "sisu"
	wailsOptions.Width = 1024
	wailsOptions.Height = 768
	wailsOptions.AssetServer = &assetServerOptions
	wailsOptions.BackgroundColour = &options.RGBA{R: 27, G: 38, B: 54, A: 1}
	wailsOptions.OnStartup = app.startup
	wailsOptions.Bind = []any{&app}

	logger.Info("starting Wails application")
	if err := wails.Run(&wailsOptions); err != nil {
		logger.Sugar().Errorw("app start error", "error", err)
		return 1
	}

	logger.Info("application shutdown complete")
	return 0
}

func configLogger(logFile string, logLevel string) *zap.Logger {
	level, err := zap.ParseAtomicLevel(logLevel)
	if err != nil {
		panic(err)
	}

	cfg := zap.NewProductionConfig()

	// Use human-readable timestamp format instead of epoch
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// On Windows, GUI apps launched via .exe don't have stdout/stderr attached
	// Writing to them causes errors, so only write to the log file
	if runtime.GOOS == "windows" {
		cfg.OutputPaths = []string{logFile}
		cfg.ErrorOutputPaths = []string{logFile}
	} else {
		cfg.OutputPaths = []string{"stdout", logFile}
		cfg.ErrorOutputPaths = []string{"stderr", logFile}
	}

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

	// Configure connection pool for SQLite
	// SQLite works best with a single connection to avoid database locking
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)

	if err := initDB(db); err != nil {
		return nil, fmt.Errorf("init repository: %w", err)
	}

	return db, nil
}

func initDB(db *sqlx.DB) error {
	zap.L().Info("running database migrations")
	filesystem, err := iofs.New(migrations, "database/migrations")
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

	zap.L().Info("database migrations completed")

	return nil
}
