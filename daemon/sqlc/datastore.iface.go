package sqlc

import (
	"context"
)

type DataStore interface {
	Open() error
	Query(ctx context.Context, sql string) error
	DB() DBTX
	Initialize(ctx context.Context) error
	Queries() *Queries
	SetQueries(*Queries)
}
