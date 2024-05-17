package storage

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"time"
)

func throttle() {
	time.Sleep(250 * time.Millisecond)
}

func NewNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: true}
}

func NewNullInt64(n int64) sql.NullInt64 {
	return sql.NullInt64{Int64: n, Valid: true}
}

func Close(c io.Closer, f func(err error)) {
	f(c.Close())
}

func WarnOnError(err error) {
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
	}
}

func panicf(msg string, args ...any) {
	panic(fmt.Sprintf(msg, args...))
}

func newContext() context.Context {
	return context.Background()
}
