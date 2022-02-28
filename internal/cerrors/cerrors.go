package cerrors

import (
	"github.com/pkg/errors"
)

func Wrap(err error, message string) error            { return errors.Wrap(err, message) }
func New(message string) error                        { return errors.New(message) }
func Errorf(format string, args ...interface{}) error { return errors.Errorf(format, args) }
func As(err error, target interface{}) bool           { return errors.As(err, target) }
func Cause(err error) error                           { return errors.Cause(err) }
func Is(err error, target error) bool                 { return errors.Is(err, target) }
func Unwrap(err error) error                          { return errors.Unwrap(err) }
func WithMessage(err error, message string) error     { return errors.WithMessage(err, message) }
func WithMessagef(err error, format string, args ...interface{}) error {
	return errors.WithMessagef(err, format, args)
}
func WithStack(err error) error { return errors.WithStack(err) }
func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args)
}
