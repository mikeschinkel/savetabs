package sqlc

import (
	"context"
	"database/sql"
	_ "embed"
	"log/slog"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var _ DataStore = (*SqliteDataStore)(nil)

type SqliteDataStore struct {
	queries  *Queries
	Filepath string
	db       DBTX
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
	err = ds.Query(ctx, DDL())
	if err != nil {
		goto end
	}
end:
	return err
}

func (ds *SqliteDataStore) Open() (err error) {
	ds.db, err = sql.Open("sqlite3", ds.Filepath)
	if err != nil {
		goto end
	}
	ds.db = NewNestedDBTX(ds.db)
	ds.queries = New(ds.db)
end:
	return err
}

func (ds *SqliteDataStore) Query(ctx context.Context, sql string) (err error) {
	_, err = ds.db.ExecContext(ctx, sql)
	return err
}

func (ds *SqliteDataStore) Queries() *Queries {
	return ds.queries
}

func (ds *SqliteDataStore) DB() DBTX {
	return ds.db
}

func (ds *SqliteDataStore) SetQueries(q *Queries) {
	ds.queries = q
}