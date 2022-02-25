package rules

const (
	NameMaxLength                = 32
	NameMinLength                = 1
	DomainMaxLength              = 32
	DomainMinLength              = 4
	MaxCountOwnedChats           = 128
	NoteMaxLength                = 64
	MaxMessagesCount             = 100
	MessageBodyMaxLength         = 4096
	RefreshTokenLength           = 28
	RefreshTokenBytesLength      = 16
	MaxCountRooms                = 128
	MaxUserChats                 = 128
	MaxMembersOnChat             = 2_097_152
	LimitOnShowUnitsInSearch     = 20
	LimitOnShowMessages          = 20
	LimitOnShowChats             = 20
	LimitOnShowDialogs           = 20
	LimitOnShowMembers           = 20
	MaxLimit                     = 20
	MinPasswordLength            = 6
	MaxPasswordLength            = 32
	MaxInviteLinks               = 3
	MaxRolesInChat               = 128
	MaxFormFields                = 16
	MaxFielditems                = 16
	RefreshTokenLiftime          = int64(60 * 60 * 24 * 60) // 60 days
	MaxRefreshSession            = 5
	LifetimeOfMarkedClient       = int64(60)      // s.
	LiftimeOfRegistrationSession = int64(60 * 60) // 1 hour
	DurationOfScheduleInterval   = int64(60)      // 1 hour

	AllowedConnectionShutdownDuration = 120
)

type AdvancedError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (p *AdvancedError) Error() string {
	return p.Message
}

var (
	ErrDataRetrieved = &AdvancedError{
		Code:    "ErrDataRetrieved",
		Message: "server failed to process the data successfully",
	}
	ErrInvalidEmail = &AdvancedError{
		Code:    "ErrInvalidEmail",
		Message: "invalid email",
	}
	ErrInvalidLink = &AdvancedError{
		Code:    "ErrInvalidLink",
		Message: "invalid link",
	}
	ErrInvalidChoiceDate = &AdvancedError{
		Code:    "ErrInvalidChoiceDate",
		Message: "invalid date",
	}
	ErrChoiceValueLength = &AdvancedError{
		Code:    "ErrChoiceValueLength",
		Message: "exceeding the value length limit",
	}
	ErrInvalidChoiceValue = &AdvancedError{
		Code:    "ErrInvalidChoiceValue",
		Message: "the key value does not match the template key type",
	}
	ErrMissingChoicePair = &AdvancedError{
		Code:    "ErrMissingChoicePair",
		Message: "the mandatory key-value pair is missing in the submitted form",
	}
)
