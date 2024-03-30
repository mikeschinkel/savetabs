package sqlc

import (
	_ "embed"
)

//go:generate sqlc generate -f ./sqlc.yaml

//go:embed schema.sql
var ddl string

func DDL() string {
	return ddl
}
