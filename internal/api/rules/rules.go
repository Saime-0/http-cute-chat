package rules

type BindKeys string

const (
	BindUserDomain BindKeys = "user-domain"
	BindChatDomain BindKeys = "chat-domain"
	BindUserID     BindKeys = "user-id"
	BindChatID     BindKeys = "chat-id"
	BindRoomID     BindKeys = "room-id"
	BindMessageID  BindKeys = "message-id"
)

type ContextKeys int

const (
	UserIDFromToken ContextKeys = iota
	PipeLineUserID
	PipeLineUserDomain
	PipeLineChatID
	PipeLineChatDomain
	PipeLineFragmentName
	PipeLineOffset
	PipeLineUserUpdateDataModel
	PipeLineUserUpdateSettingsModel
)

const (
	NameMaxLength            = 32
	NameMinLength            = 1
	DomainMaxLength          = 32
	DomainMinLength          = 4
	AppSettingsMaxLength     = 512
	MaxCountOwnedChats       = 128
	NoteMaxLength            = 64
	MessageBodyMaxLength     = 4096
	RefreshTokenLength       = 16
	MaxCountRooms            = 128
	MaxUserChats             = 128
	MaxMembersOnChat         = 2_097_152
	LimitOnShowUnitsInSearch = 20
	LimitOnShowMessages      = 20
	LimitOnShowChats         = 20
	LimitOnShowDialogs       = 20
	LimitOnShowMembers       = 20
)

// Errors ...
type PureErrorModels struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (p *PureErrorModels) Error() string {
	return p.Message
}

var (
	// Errors
	ErrOutOfRange = PureErrorModels{
		Code:    "ErrOutOfRange",
		Message: "parameter value is out of range",
	}
	ErrInvalidValue = PureErrorModels{
		Code:    "ErrInvalidValue",
		Message: "invalid parameter value",
	}
	ErrAccessingDatabase = PureErrorModels{
		Code:    "ErrAccessingDatabase",
		Message: "internal server error when accessing the database",
	}
	ErrUserNotFound = PureErrorModels{
		Code:    "ErrUserNotFound",
		Message: "user was not found",
	}
	ErrChatNotFound = PureErrorModels{
		Code:    "ErrChatNotFound",
		Message: "chat was not found",
	}
	ErrRoomNotFound = PureErrorModels{
		Code:    "ErrRoomNotFound",
		Message: "room was not found",
	}
	ErrDialogNotFound = PureErrorModels{
		Code:    "ErrDialogNotFound",
		Message: "dialog was not found",
	}
	ErrMessageNotFound = PureErrorModels{
		Code:    "ErrMessageNotFound",
		Message: "message was not found",
	}
	ErrDataRetrieved = PureErrorModels{
		Code:    "ErrDataRetrieved",
		Message: "server failed to process the data successfully",
	}
	ErrOccupiedDomain = PureErrorModels{
		Code:    "ErrOccupiedDomain",
		Message: "domain is occupied by someone",
	}
	ErrLimitHasBeenReached = PureErrorModels{
		Code:    "ErrLimitHasBeenReached",
		Message: "maximum limit has been reached",
	}
	ErrMembersLimitHasBeenReached = PureErrorModels{
		Code:    "ErrMembersLimitHasBeenReached",
		Message: "seats are occupied in this chat",
	}
	ErrNoAccess = PureErrorModels{
		Code:    "ErrNoAccess",
		Message: "there is not enough access for this action",
	}
	ErrBadRequestBody = PureErrorModels{
		Code:    "ErrBadRequestBody",
		Message: "bad request body",
	}
	ErrDialogWithYourself = PureErrorModels{
		Code:    "ErrDialogWithYourself",
		Message: "there is a special section for this",
	}
)
