package rules

import "errors"

type (
	ContextKeys int
)

const (
	UserIDFromToken ContextKeys = iota
)

const (
	NameMaxLength        = 32
	NameMinLength        = 1
	DomainMaxLength      = 32
	DomainMinLength      = 4
	AppSettingsMaxLength = 512
	MaxCountOwnedChats   = 128
	NoteMaxLength        = 64
	MessageBodyMaxLength = 4096
	RefreshTokenLength   = 16
	MaxCountRooms        = 128
	MaxUserChats         = 128
)

// todo: коды ошибок {"error":{"message":"...","code":1}} либо уникальное имя ошибки и ее описание
var (
	// Errors
	ErrOutOfRange          = errors.New("parameter value is out of range")
	ErrInvalidValue        = errors.New("invalid parameter value")
	ErrAccessingDatabase   = errors.New("internal server error when accessing the database")
	ErrUserNotFound        = errors.New("user was not found")
	ErrChatNotFound        = errors.New("chat was not found")
	ErrRoomNotFound        = errors.New("room was not found")
	ErrDataRetrieved       = errors.New("server failed to process the data successfully")
	ErrOccupiedDomain      = errors.New("domain is occupied by someone")
	ErrLimitHasBeenReached = errors.New("maximum limit has been reached")
	ErrNoAccess            = errors.New("there is not enough access for this action")
	ErrBadRequestBody      = errors.New("bad request body")
)
