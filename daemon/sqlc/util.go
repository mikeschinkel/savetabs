package sqlc

import (
	"database/sql"
	"fmt"
	"io"
	"os"
)

func panicf(msg string, args ...any) {
	panic(fmt.Sprintf(msg, args...))
}

func Close(c io.Closer, f func(err error)) {
	f(c.Close())
}

func WarnOnError(err error) {
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
	}
}

func NewNullInt64(n int64) sql.NullInt64 {
	return sql.NullInt64{Int64: n, Valid: true}
}

func NewNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: true}
}
