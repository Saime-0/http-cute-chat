// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type AllowsResult interface {
	IsAllowsResult()
}

type ChatResult interface {
	IsChatResult()
}

type ChatRolesResult interface {
	IsChatRolesResult()
}

type ChatsResult interface {
	IsChatsResult()
}

type CountMembersResult interface {
	IsCountMembersResult()
}

type InviteInfoResult interface {
	IsInviteInfoResult()
}

type InvitesResult interface {
	IsInvitesResult()
}

type JoinByInviteResult interface {
	IsJoinByInviteResult()
}

type JoinToChatResult interface {
	IsJoinToChatResult()
}

type LoginResult interface {
	IsLoginResult()
}

type MeResult interface {
	IsMeResult()
}

type MemberResult interface {
	IsMemberResult()
}

type MembersResult interface {
	IsMembersResult()
}

type MessageInfoResult interface {
	IsMessageInfoResult()
}

type MutationResult interface {
	IsMutationResult()
}

type RefreshTokensResult interface {
	IsRefreshTokensResult()
}

type RegisterResult interface {
	IsRegisterResult()
}

type RoleResult interface {
	IsRoleResult()
}

type RolesResult interface {
	IsRolesResult()
}

type RoomFormResult interface {
	IsRoomFormResult()
}

type RoomMessagesResult interface {
	IsRoomMessagesResult()
}

type RoomResult interface {
	IsRoomResult()
}

type RoomWhiteListResult interface {
	IsRoomWhiteListResult()
}

type RoomsResult interface {
	IsRoomsResult()
}

type SendMessageToRoomResult interface {
	IsSendMessageToRoomResult()
}

type UnitResult interface {
	IsUnitResult()
}

type UnitsResult interface {
	IsUnitsResult()
}

type UpdateChatResult interface {
	IsUpdateChatResult()
}

type UpdateMeDataResult interface {
	IsUpdateMeDataResult()
}

type UpdateRoleResult interface {
	IsUpdateRoleResult()
}

type UpdateRoomResult interface {
	IsUpdateRoomResult()
}

type UserResult interface {
	IsUserResult()
}

type UserRoleResult interface {
	IsUserRoleResult()
}

type UsersResult interface {
	IsUsersResult()
}

type AdvancedError struct {
	Code  string `json:"code"`
	Error string `json:"error"`
}

func (AdvancedError) IsJoinByInviteResult()      {}
func (AdvancedError) IsJoinToChatResult()        {}
func (AdvancedError) IsLoginResult()             {}
func (AdvancedError) IsRefreshTokensResult()     {}
func (AdvancedError) IsRegisterResult()          {}
func (AdvancedError) IsSendMessageToRoomResult() {}
func (AdvancedError) IsUpdateChatResult()        {}
func (AdvancedError) IsUpdateMeDataResult()      {}
func (AdvancedError) IsUpdateRoleResult()        {}
func (AdvancedError) IsUpdateRoomResult()        {}
func (AdvancedError) IsChatRolesResult()         {}
func (AdvancedError) IsChatsResult()             {}
func (AdvancedError) IsInviteInfoResult()        {}
func (AdvancedError) IsMeResult()                {}
func (AdvancedError) IsMessageInfoResult()       {}
func (AdvancedError) IsRoomFormResult()          {}
func (AdvancedError) IsRoomMessagesResult()      {}
func (AdvancedError) IsRoomResult()              {}
func (AdvancedError) IsRoomWhiteListResult()     {}
func (AdvancedError) IsUnitResult()              {}
func (AdvancedError) IsUnitsResult()             {}
func (AdvancedError) IsUserRoleResult()          {}
func (AdvancedError) IsMutationResult()          {}
func (AdvancedError) IsUserResult()              {}
func (AdvancedError) IsRoomsResult()             {}
func (AdvancedError) IsCountMembersResult()      {}
func (AdvancedError) IsMembersResult()           {}
func (AdvancedError) IsRolesResult()             {}
func (AdvancedError) IsInvitesResult()           {}
func (AdvancedError) IsUsersResult()             {}
func (AdvancedError) IsChatResult()              {}
func (AdvancedError) IsRoleResult()              {}
func (AdvancedError) IsMemberResult()            {}
func (AdvancedError) IsAllowsResult()            {}

type Allows struct {
	Room       *Room              `json:"room"`
	AllowRead  *PermissionHolders `json:"allow_read"`
	AllowWrite *PermissionHolders `json:"allow_write"`
}

func (Allows) IsAllowsResult() {}

type Chars struct {
	Chars []Char `json:"chars"`
}

type Chat struct {
	Unit         *Unit              `json:"unit"`
	Owner        UserResult         `json:"owner"`
	Rooms        RoomsResult        `json:"rooms"`
	Private      bool               `json:"private"`
	CountMembers CountMembersResult `json:"count_members"`
	Members      MembersResult      `json:"members"`
	Roles        RolesResult        `json:"roles"`
	Invites      InvitesResult      `json:"invites"`
	Banlist      UsersResult        `json:"banlist"`
	Me           MemberResult       `json:"me"`
}

func (Chat) IsJoinByInviteResult() {}
func (Chat) IsJoinToChatResult()   {}
func (Chat) IsUpdateChatResult()   {}
func (Chat) IsChatResult()         {}

type ChatArray struct {
	Chats []*Chat `json:"chats"`
}

func (ChatArray) IsChatsResult() {}

type CreateChatInput struct {
	Domain  string `json:"domain"`
	Name    string `json:"name"`
	Private bool   `json:"private"`
}

type CreateInviteInput struct {
	Code     string `json:"code"`
	Aliens   *int   `json:"aliens"`
	Duration *int64 `json:"duration"`
}

type CreateMessageInput struct {
	ReplyTo int    `json:"reply_to"`
	Body    string `json:"body"`
}

type CreateRoleInput struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type CreateRoomInput struct {
	Name      string           `json:"name"`
	Parent    *int             `json:"parent"`
	Note      *string          `json:"note"`
	MsgFormat *UpdateFormInput `json:"msg_format"`
	Restricts *RestrictsInput  `json:"restricts"`
}

type FindByDomainOrID struct {
	ID     *int    `json:"id"`
	Domain *string `json:"domain"`
}

type Form struct {
	Fields []*FormField `json:"fields"`
}

func (Form) IsRoomFormResult() {}

type FormField struct {
	Key      string    `json:"key"`
	Type     FieldType `json:"type"`
	Optional bool      `json:"optional"`
	Length   *int      `json:"length"`
	Items    []string  `json:"items"`
}

type FormFieldInput struct {
	Key      string    `json:"key"`
	Type     FieldType `json:"type"`
	Optional bool      `json:"optional"`
	Length   *int      `json:"length"`
	Items    []string  `json:"items"`
}

type IntValue struct {
	Value *int `json:"value"`
}

func (IntValue) IsCountMembersResult() {}

type Invite struct {
	Code      string `json:"code"`
	Aliens    *int   `json:"aliens"`
	ExpiresAt *int64 `json:"expires_at"`
}

type InviteInfo struct {
	Unit         *Unit              `json:"unit"`
	Private      bool               `json:"private"`
	CountMembers CountMembersResult `json:"count_members"`
}

func (InviteInfo) IsInviteInfoResult() {}

type Invites struct {
	Invites []*Invite `json:"invites"`
}

func (Invites) IsInvitesResult() {}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Me struct {
	User       *User     `json:"user"`
	Data       *UserData `json:"data"`
	Chats      []*Chat   `json:"chats"`
	OwnedChats []*Chat   `json:"owned_chats"`
}

func (Me) IsMeResult() {}

type Member struct {
	Chat     *Chat      `json:"chat"`
	User     *User      `json:"user"`
	Role     RoleResult `json:"role"`
	Char     *Char      `json:"char"`
	JoinedAt int64      `json:"joined_at"`
	Muted    bool       `json:"muted"`
	Frozen   bool       `json:"frozen"`
}

func (Member) IsMemberResult() {}

type Members struct {
	Members []*Member `json:"members"`
}

func (Members) IsMembersResult() {}

type Message struct {
	ID      int         `json:"id"`
	Room    *Room       `json:"room"`
	ReplyTo *Message    `json:"reply_to"`
	Author  *Unit       `json:"author"`
	Body    string      `json:"body"`
	Type    MessageType `json:"type"`
	Date    int         `json:"date"`
}

func (Message) IsSendMessageToRoomResult() {}
func (Message) IsMessageInfoResult()       {}

type MessageArray struct {
	Messages []*Message `json:"messages"`
}

func (MessageArray) IsRoomMessagesResult() {}

type MessageFilter struct {
	TextFragment *string `json:"text_fragment"`
	AuthorID     *int    `json:"author_id"`
	ChatID       *int    `json:"chat_id"`
	RoomID       *int    `json:"room_id"`
}

type Params struct {
	Limit  *int `json:"limit"`
	Offset *int `json:"offset"`
}

type PermissionHolders struct {
	Roles   *Roles   `json:"roles"`
	Chars   *Chars   `json:"chars"`
	Members *Members `json:"members"`
}

type PermissionHoldersInput struct {
	Roles   []int  `json:"roles"`
	Chars   []Char `json:"chars"`
	Members []int  `json:"members"`
}

type RegisterInput struct {
	Domain   string `json:"domain"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RestrictsInput struct {
	AllowRead  *PermissionHoldersInput `json:"allow_read"`
	AllowWrite *PermissionHoldersInput `json:"allow_write"`
}

type Role struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

func (Role) IsUpdateRoleResult() {}
func (Role) IsUserRoleResult()   {}
func (Role) IsRoleResult()       {}

type Roles struct {
	Roles []*Role `json:"roles"`
}

func (Roles) IsChatRolesResult()     {}
func (Roles) IsRoomWhiteListResult() {}
func (Roles) IsRolesResult()         {}

type Room struct {
	ID        int          `json:"id"`
	Chat      *Chat        `json:"chat"`
	Name      string       `json:"name"`
	ParentID  *int         `json:"parent_id"`
	Note      *string      `json:"note"`
	MsgFormat *Form        `json:"msg_format"`
	Allows    AllowsResult `json:"allows"`
	Messages  []*Message   `json:"messages"`
}

func (Room) IsUpdateRoomResult() {}
func (Room) IsRoomResult()       {}

type Rooms struct {
	Rooms []*Room `json:"rooms"`
}

func (Rooms) IsRoomsResult() {}

type Successful struct {
	Success string `json:"success"`
}

func (Successful) IsJoinByInviteResult()      {}
func (Successful) IsJoinToChatResult()        {}
func (Successful) IsLoginResult()             {}
func (Successful) IsRefreshTokensResult()     {}
func (Successful) IsRegisterResult()          {}
func (Successful) IsSendMessageToRoomResult() {}
func (Successful) IsMutationResult()          {}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (TokenPair) IsLoginResult()         {}
func (TokenPair) IsRefreshTokensResult() {}

type Unit struct {
	ID     int      `json:"id"`
	Domain string   `json:"domain"`
	Name   string   `json:"name"`
	Type   UnitType `json:"type"`
}

func (Unit) IsUnitResult() {}

type UnitArray struct {
	Units []*Unit `json:"units"`
}

func (UnitArray) IsUnitsResult() {}

type UpdateChatInput struct {
	Domain  string `json:"domain"`
	Name    string `json:"name"`
	Private bool   `json:"private"`
}

type UpdateFormInput struct {
	Fields []*FormFieldInput `json:"fields"`
}

type UpdateMeDataInput struct {
	Domain   *string `json:"domain"`
	Name     *string `json:"name"`
	Password *string `json:"password"`
	Email    *string `json:"email"`
}

type UpdateRoleInput struct {
	Name  *string `json:"name"`
	Color *string `json:"color"`
}

type UpdateRoomInput struct {
	Name      *string          `json:"name"`
	ParentID  *int             `json:"parent_id"`
	Note      *string          `json:"note"`
	Restricts *RestrictsInput  `json:"restricts"`
	MsgFormat *UpdateFormInput `json:"msg_format"`
}

type User struct {
	Unit *Unit `json:"unit"`
}

func (User) IsUserResult() {}

type UserData struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (UserData) IsUpdateMeDataResult() {}

type Users struct {
	Users []*User `json:"users"`
}

func (Users) IsUsersResult() {}

type Char string

const (
	CharAdmin Char = "ADMIN"
	CharModer Char = "MODER"
	CharNone  Char = "NONE"
)

var AllChar = []Char{
	CharAdmin,
	CharModer,
	CharNone,
}

func (e Char) IsValid() bool {
	switch e {
	case CharAdmin, CharModer, CharNone:
		return true
	}
	return false
}

func (e Char) String() string {
	return string(e)
}

func (e *Char) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Char(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Char", str)
	}
	return nil
}

func (e Char) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type FieldType string

const (
	FieldTypeEmail    FieldType = "EMAIL"
	FieldTypeDate     FieldType = "DATE"
	FieldTypeLink     FieldType = "LINK"
	FieldTypeText     FieldType = "TEXT"
	FieldTypeNumeric  FieldType = "NUMERIC"
	FieldTypeSelector FieldType = "SELECTOR"
)

var AllFieldType = []FieldType{
	FieldTypeEmail,
	FieldTypeDate,
	FieldTypeLink,
	FieldTypeText,
	FieldTypeNumeric,
	FieldTypeSelector,
}

func (e FieldType) IsValid() bool {
	switch e {
	case FieldTypeEmail, FieldTypeDate, FieldTypeLink, FieldTypeText, FieldTypeNumeric, FieldTypeSelector:
		return true
	}
	return false
}

func (e FieldType) String() string {
	return string(e)
}

func (e *FieldType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = FieldType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid FieldType", str)
	}
	return nil
}

func (e FieldType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type MessageType string

const (
	MessageTypeSystem    MessageType = "SYSTEM"
	MessageTypeUser      MessageType = "USER"
	MessageTypeFormatted MessageType = "FORMATTED"
)

var AllMessageType = []MessageType{
	MessageTypeSystem,
	MessageTypeUser,
	MessageTypeFormatted,
}

func (e MessageType) IsValid() bool {
	switch e {
	case MessageTypeSystem, MessageTypeUser, MessageTypeFormatted:
		return true
	}
	return false
}

func (e MessageType) String() string {
	return string(e)
}

func (e *MessageType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = MessageType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid MessageType", str)
	}
	return nil
}

func (e MessageType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type UnitType string

const (
	UnitTypeChat UnitType = "CHAT"
	UnitTypeUser UnitType = "USER"
)

var AllUnitType = []UnitType{
	UnitTypeChat,
	UnitTypeUser,
}

func (e UnitType) IsValid() bool {
	switch e {
	case UnitTypeChat, UnitTypeUser:
		return true
	}
	return false
}

func (e UnitType) String() string {
	return string(e)
}

func (e *UnitType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = UnitType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid UnitType", str)
	}
	return nil
}

func (e UnitType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
