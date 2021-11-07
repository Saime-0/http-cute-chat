package rules

const Year int64 = 31536000

type BindKeys string

const (
	BindUserDomain BindKeys = "user-domain"
	BindChatDomain BindKeys = "chat-domain"
	BindUserID     BindKeys = "user-id"
	BindChatID     BindKeys = "chat-id"
	BindRoomID     BindKeys = "room-id"
	BindMessageID  BindKeys = "message-id"
	BindInviteCode BindKeys = "invite-code"
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
	MinPasswordLength        = 6
	MaxPasswordLength        = 32
	MaxInviteLinks           = 3
	MaxRolesInChat           = 128
	Max
)

// Errors ...
type AdvancedError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	// http status const
}

func (p *AdvancedError) Error() string {
	return p.Message
}

var (
	// Errors
	ErrOutOfRange = AdvancedError{
		Code:    "ErrOutOfRange",
		Message: "parameter value is out of range",
	}
	ErrInvalidValue = AdvancedError{
		Code:    "ErrInvalidValue",
		Message: "invalid parameter value",
	}
	ErrAccessingDatabase = AdvancedError{
		Code:    "ErrAccessingDatabase",
		Message: "internal server error when accessing the database",
	}
	ErrUserNotFound = AdvancedError{
		Code:    "ErrUserNotFound",
		Message: "user was not found",
	}
	ErrChatNotFound = AdvancedError{
		Code:    "ErrChatNotFound",
		Message: "chat was not found",
	}
	ErrRoomNotFound = AdvancedError{
		Code:    "ErrRoomNotFound",
		Message: "room was not found",
	}
	ErrDialogNotFound = AdvancedError{
		Code:    "ErrDialogNotFound",
		Message: "dialog was not found",
	}
	ErrMessageNotFound = AdvancedError{
		Code:    "ErrMessageNotFound",
		Message: "message was not found",
	}
	ErrDataRetrieved = AdvancedError{
		Code:    "ErrDataRetrieved",
		Message: "server failed to process the data successfully",
	}
	ErrOccupiedDomain = AdvancedError{
		Code:    "ErrOccupiedDomain",
		Message: "domain is occupied by someone",
	}
	ErrLimitHasBeenReached = AdvancedError{
		Code:    "ErrLimitHasBeenReached",
		Message: "maximum limit has been reached",
	}
	ErrMembersLimitHasBeenReached = AdvancedError{
		Code:    "ErrMembersLimitHasBeenReached",
		Message: "seats are occupied in this chat",
	}
	ErrNoAccess = AdvancedError{
		Code:    "ErrNoAccess",
		Message: "there is not enough access for this action",
	}
	ErrBadRequestBody = AdvancedError{
		Code:    "ErrBadRequestBody",
		Message: "bad request body",
	}
	ErrDialogWithYourself = AdvancedError{
		Code:    "ErrDialogWithYourself",
		Message: "there is a special section for this",
	}
	ErrInvalidName = AdvancedError{
		Code:    "ErrInvalidUserName",
		Message: "invalid name",
	}
	ErrInvalidDomain = AdvancedError{
		Code:    "ErrInvalidDomain",
		Message: "invalid domain",
	}
	ErrInvalidEmail = AdvancedError{
		Code:    "ErrInvalidEmail",
		Message: "invalid email",
	}
	ErrInvalidPassword = AdvancedError{
		Code:    "ErrInvalidPassword",
		Message: "invalid password",
	}
	ErrInviteLinkNotFound = AdvancedError{
		Code:    "ErrInviteLinkNotFound",
		Message: "invite link was not found",
	}
	ErrUserBannedInChat = AdvancedError{
		Code:    "ErrUserBannedInChat",
		Message: "user was banned in this chat",
	}
	ErrUserIsNotChatMember = AdvancedError{
		Code:    "ErrUserIsNotChatMember",
		Message: "user is not a member of the chat",
	}
	ErrRoleHidden = AdvancedError{
		Code:    "ErrRoleHidden",
		Message: "this role is hidden",
	}
	ErrPrivateRoom = AdvancedError{
		Code:    "ErrRoomPrivate",
		Message: "no right to send messages to the room",
	}
)
