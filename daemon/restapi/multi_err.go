package restapi

import (
	"errors"
)

type multiErr struct {
	errs []error
}

func newMultiErr() *multiErr {
	return &multiErr{
		errs: make([]error, 0),
	}
}
func (r multiErr) IsError() bool {
	return len(r.errs) > 0
}
func (r multiErr) Add(errs ...error) {
	r.errs = append(r.errs, errs...)
}
func (r multiErr) Err() (errs error) {
	if len(r.errs) == 0 {
		return nil
	}
	return errors.Join(r.errs...)
}
