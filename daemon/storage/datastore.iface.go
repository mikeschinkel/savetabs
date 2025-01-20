package storage

import (
	"github.com/mikeschinkel/savetabs/daemon/sqlc"
)

type DataStore interface {
	Open() error
	Query(ctx Context, sql string) error
	DB() sqlc.DBTX
	Initialize(ctx Context) error
	Queries(tx sqlc.DBTX) *sqlc.Queries
}
