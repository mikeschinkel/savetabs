package storage

import (
	"errors"
	"fmt"
)

var (
	ErrFailedToUnmarshal        = errors.New("failed to unmarshal JSON")
	ErrFailedUpsertLinks        = errors.New("failed to upsert links")
	ErrFailedUpsertLink         = errors.New("failed to upsert link")
	ErrFailedInsertLinkContent  = errors.New("failed to insert link content")
	ErrFailedLoadLinkByUrl      = errors.New("failed to load link by URL")
	ErrFailedUpsertLinkGroups   = errors.New("failed to upsert link-groups")
	ErrFailedUpsertGroups       = errors.New("failed to upsert groups")
	ErrFailedToArchiveLinks     = errors.New("failed to archive links")
	ErrFailedToMarkLinksDeleted = errors.New("failed to delete links")
	ErrInvalidGroupId           = errors.New("invalid group Id")
	ErrInvalidLinkIds           = errors.New("invalid Link Ids")
	ErrFailedConvertToAbsPath   = errors.New("failed to convert to absolute path")
	ErrFailedToInitDataStore    = errors.New("failed to initialize data store")
)

func Error(namedErr, actualErr error, args ...string) error {
	var arg string
	if len(args)%2 == 1 {
		panicf("Cannot call Error with mismatched keys and values for args: %v.", args)
	}
	if len(args) == 0 {
		goto end
	}
	for i := 0; i < len(args); i += 2 {
		arg += args[i] + "=" + args[i+1] + ","
	}
	arg = arg[:len(arg)-1]
	actualErr = fmt.Errorf("%w [%s]", actualErr, arg)
end:
	return errors.Join(namedErr, actualErr)
}
