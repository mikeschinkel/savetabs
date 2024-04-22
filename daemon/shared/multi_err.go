package shared

import (
	"errors"
)

type MultiErr struct {
	errs []error
}

func NewMultiErr() *MultiErr {
	return &MultiErr{
		errs: make([]error, 0),
	}
}
func (r MultiErr) IsError() bool {
	return len(r.errs) > 0
}
func (r MultiErr) Add(errs ...error) {
	r.errs = append(r.errs, errs...)
}
func (r MultiErr) Err() (errs error) {
	if len(r.errs) == 0 {
		return nil
	}
	return errors.Join(r.errs...)
}
