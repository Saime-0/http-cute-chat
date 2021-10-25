package v1

import "errors"

const (
	// Users
	NameMaxLength = 32
	NameMinLength = 1

	DomainMaxLength = 32
	DomainMinLength = 4

	AppSettingsMaxLength = 512
)

var (
	// Errors
	ErrOutOfRange        = errors.New("parameter value is out of range")
	ErrInvalidValue      = errors.New("invalid parameter value")
	ErrAccessingDatabase = errors.New("internal server error when accessing the database")
	ErrUserNotFound      = errors.New("user was not found")
	ErrChatNotFound      = errors.New("chat was not found")
	ErrDataRetrieved     = errors.New("server failed to process the data successfully")
)
