package sqlc

import (
	"errors"
	"fmt"
)

var (
	ErrFailedConvertToAbsPath  = errors.New("failed to convert to absolute path")
	ErrFailedToInitDataStore   = errors.New("failed to initialize data store")
	ErrTooManyRequests         = errors.New("too many requests")
	ErrCannotUpdateStatusField = errors.New("cannot update status field")
	ErrDBNotANestedDBTX        = errors.New("db is not a *NestedDBTX")
)

func Error(namedErr, actualErr error, args ...string) (err error) {
	var arg string
	if len(args)%2 == 1 {
		panicf("Cannot call Error with mismatched keys and values for args: %v.", args)
	}
	if len(args) == 0 {
		err = actualErr
		goto end
	}
	for i := 0; i < len(args); i += 2 {
		arg += args[i] + "=" + args[i+1] + ","
	}
	arg = arg[:len(arg)-1]
	err = fmt.Errorf("%w [%s]", err, arg)
end:
	return errors.Join(namedErr, actualErr)
}
