package storage

import (
	"context"
	"database/sql"
	_ "embed"
	"log/slog"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"savetabs"
	"savetabs/sqlc"
)

var _ DataStore = (*SqliteDataStore)(nil)

type SqliteDataStore struct {
	Filepath string
	db       *sql.DB
}

func NewSqliteDataStore(dbFile string) DataStore {
	return &SqliteDataStore{
		Filepath: dbFile,
	}
}

func (ds *SqliteDataStore) Initialize(ctx context.Context) (err error) {
	slog.Info("Initializing data store", "data_file", ds.Filepath)

	absFP, err := filepath.Abs(ds.Filepath)
	if err != nil {
		err = Error(ErrFailedConvertToAbsPath, err, "filepath", ds.Filepath)
		goto end
	}
	ds.Filepath = absFP

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
	ds.db, err = sql.Open("sqlite3", ds.Filepath)
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
