package storage

import (
	"context"
	"database/sql"
	_ "embed"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mikeschinkel/savetabs/daemon"
	"github.com/mikeschinkel/savetabs/daemon/sqlc"
)

var _ DataStore = (*SqliteDataStore)(nil)

const macOSConfigSubdir = "/Library/Application Support"
const desiredConfigSubdir = ".config"

type SqliteDataStore struct {
	path    string
	AppName string
	db      *sql.DB
}

func (ds *SqliteDataStore) Filepath() string {
	return filepath.Join(ds.path, ds.AppName+".db")
}

func NewSqliteDataStore(args Args) DataStore {
	return &SqliteDataStore{
		AppName: args.AppName,
	}
}

func (ds *SqliteDataStore) Initialize(ctx context.Context) (err error) {
	var configDir string

	configDir, err = os.UserConfigDir()
	if err != nil {
		err = ErrFailedToGetConfigPath
		goto end
	}
	// Move macOS config dir to be ~/.config vs. ~/Library/Application Support
	if strings.HasSuffix(configDir, macOSConfigSubdir) {
		configDir = filepath.Join(
			configDir[:len(configDir)-len(macOSConfigSubdir)],
			desiredConfigSubdir,
			ds.AppName,
		)
	}
	ds.path = configDir

	slog.Info("Initializing data store",
		"data_store", relativeToHomeDir(ds.Filepath()),
	)

	err = ds.Open()
	if err != nil {
		goto end
	}
	err = ds.Query(ctx, daemon.DDL())
	if err != nil {
		goto end
	}
end:
	return err
}

func (ds *SqliteDataStore) Open() (err error) {
	ds.db, err = sql.Open("sqlite3", ds.Filepath())
	return err
}

func (ds *SqliteDataStore) Query(ctx context.Context, sql string) (err error) {
	_, err = ds.db.ExecContext(ctx, sql)
	return err
}

func (ds *SqliteDataStore) Queries(dbtx sqlc.DBTX) *sqlc.Queries {
	return sqlc.New(dbtx)
}

func (ds *SqliteDataStore) DB() sqlc.DBTX {
	return ds.db
}

func relativeToHomeDir(fp string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic("Unable to get user home directory")
	}
	rel, err := filepath.Rel(home, fp)
	if err != nil {
		panicf("Unable to get relative path to %s", fp)
	}
	return "~/" + rel
}
