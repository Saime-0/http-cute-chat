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

type MessageResult interface {
	IsMessageResult()
}

type MessagesResult interface {
	IsMessagesResult()
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

type RoomResult interface {
	IsRoomResult()
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
func (AdvancedError) IsMessageResult()           {}
func (AdvancedError) IsRoomFormResult()          {}
func (AdvancedError) IsMessagesResult()          {}
func (AdvancedError) IsRoomResult()              {}
func (AdvancedError) IsUnitResult()              {}
func (AdvancedError) IsUnitsResult()             {}
func (AdvancedError) IsUserRoleResult()          {}

type Allows struct {
	Room       *Room              `json:"room"`
	AllowRead  *PermissionHolders `json:"allowRead"`
	AllowWrite *PermissionHolders `json:"allowWrite"`
}

func (Allows) IsAllowsResult() {}

type AllowsInput struct {
	AllowRead  *PermissionHoldersInput `json:"allowRead"`
	AllowWrite *PermissionHoldersInput `json:"allowWrite"`
}

type Chars struct {
	Chars []CharType `json:"chars"`
}

type Chat struct {
	Unit         *Unit              `json:"unit"`
	Owner        UserResult         `json:"owner"`
	Rooms        RoomsResult        `json:"rooms"`
	Private      bool               `json:"private"`
	CountMembers CountMembersResult `json:"countMembers"`
	Members      MembersResult      `json:"members"`
	Roles        RolesResult        `json:"roles"`
	Invites      InvitesResult      `json:"invites"`
	Banlist      UsersResult        `json:"banlist"`
	Me           MemberResult       `json:"me"`
}

func (Chat) IsChatResult()         {}
func (Chat) IsJoinByInviteResult() {}
func (Chat) IsJoinToChatResult()   {}
func (Chat) IsUpdateChatResult()   {}

type Chats struct {
	Chats []*Chat `json:"chats"`
}

func (Chats) IsChatsResult() {}

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
	ReplyTo *int   `json:"replyTo"`
	Body    string `json:"body"`
}

type CreateRoleInput struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type CreateRoomInput struct {
	Name   string           `json:"name"`
	Parent *int             `json:"parent"`
	Note   *string          `json:"note"`
	Form   *UpdateFormInput `json:"form"`
	Allows *AllowsInput     `json:"allows"`
}

type FindByDomainOrID struct {
	ID     *int    `json:"id"`
	Domain *string `json:"domain"`
}

type FindMembers struct {
	ChatID   *int      `json:"chatId"`
	MemberID *int      `json:"memberId"`
	RoleID   *int      `json:"roleId"`
	Char     *CharType `json:"char"`
}

type FindMessages struct {
	ChatID       int     `json:"chatId"`
	RoomID       *int    `json:"roomId"`
	AuthorID     *int    `json:"authorId"`
	TextFragment *string `json:"textFragment"`
}

type FindMessagesInRoomByUnionInput struct {
	AfterTime  *int64 `json:"afterTime"`
	BeforeTime *int64 `json:"beforeTime"`
}

type FindRooms struct {
	ChatID   int        `json:"chatId"`
	RoomID   *int       `json:"roomId"`
	ParentID *int       `json:"parentId"`
	IsParent *FetchType `json:"isParent"`
}

type FindUnits struct {
	UnitID       *int      `json:"unitId"`
	UnitDomain   *string   `json:"unitDomain"`
	NameFragment *string   `json:"nameFragment"`
	UnitType     *UnitType `json:"unitType"`
}

type FindUsers struct {
	UserID       *int    `json:"userId"`
	UserDomain   *string `json:"userDomain"`
	NameFragment *string `json:"nameFragment"`
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
	CountMembers CountMembersResult `json:"countMembers"`
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
	OwnedChats []*Chat   `json:"ownedChats"`
}

func (Me) IsMeResult() {}

type Member struct {
	ID       int        `json:"id"`
	Chat     *Chat      `json:"chat"`
	User     *User      `json:"user"`
	Role     RoleResult `json:"role"`
	Char     CharType   `json:"char"`
	JoinedAt int64      `json:"joinedAt"`
	Muted    bool       `json:"muted"`
	Frozen   bool       `json:"frozen"`
}

func (Member) IsMemberResult() {}

type Members struct {
	Members []*Member `json:"members"`
}

func (Members) IsMembersResult() {}

type Message struct {
	ID        int         `json:"id"`
	Room      *Room       `json:"room"`
	ReplyTo   *Message    `json:"replyTo"`
	Author    *Member     `json:"author"`
	Body      string      `json:"body"`
	Type      MessageType `json:"type"`
	CreatedAt int64       `json:"createdAt"`
}

func (Message) IsSendMessageToRoomResult() {}
func (Message) IsMessageResult()           {}

type Messages struct {
	Messages []*Message `json:"messages"`
}

func (Messages) IsMessagesResult() {}

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
	Roles   []int      `json:"roles"`
	Chars   []CharType `json:"chars"`
	Members []int      `json:"members"`
}

type RegisterInput struct {
	Domain   string `json:"domain"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Role struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

func (Role) IsRoleResult()       {}
func (Role) IsUpdateRoleResult() {}
func (Role) IsUserRoleResult()   {}

type Roles struct {
	Roles []*Role `json:"roles"`
}

func (Roles) IsRolesResult()     {}
func (Roles) IsChatRolesResult() {}

type Room struct {
	RoomID   int            `json:"roomId"`
	Chat     *Chat          `json:"chat"`
	Name     string         `json:"name"`
	ParentID *int           `json:"parentId"`
	Note     *string        `json:"note"`
	Form     *Form          `json:"form"`
	Allows   AllowsResult   `json:"allows"`
	Messages MessagesResult `json:"messages"`
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

func (Successful) IsMutationResult()          {}
func (Successful) IsJoinByInviteResult()      {}
func (Successful) IsJoinToChatResult()        {}
func (Successful) IsLoginResult()             {}
func (Successful) IsRefreshTokensResult()     {}
func (Successful) IsRegisterResult()          {}
func (Successful) IsSendMessageToRoomResult() {}

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
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

type Units struct {
	Units []*Unit `json:"units"`
}

func (Units) IsUnitsResult() {}

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

type UpdateRoomAllowsInput struct {
	Allows *AllowsInput `json:"allows"`
}

type UpdateRoomInput struct {
	Name     *string `json:"name"`
	ParentID *int    `json:"parentId"`
	Note     *string `json:"note"`
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

type ActionType string

const (
	ActionTypeRead  ActionType = "READ"
	ActionTypeWrite ActionType = "WRITE"
)

var AllActionType = []ActionType{
	ActionTypeRead,
	ActionTypeWrite,
}

func (e ActionType) IsValid() bool {
	switch e {
	case ActionTypeRead, ActionTypeWrite:
		return true
	}
	return false
}

func (e ActionType) String() string {
	return string(e)
}

func (e *ActionType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ActionType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ActionType", str)
	}
	return nil
}

func (e ActionType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type CharType string

const (
	CharTypeAdmin CharType = "ADMIN"
	CharTypeModer CharType = "MODER"
	CharTypeNone  CharType = "NONE"
)

var AllCharType = []CharType{
	CharTypeAdmin,
	CharTypeModer,
	CharTypeNone,
}

func (e CharType) IsValid() bool {
	switch e {
	case CharTypeAdmin, CharTypeModer, CharTypeNone:
		return true
	}
	return false
}

func (e CharType) String() string {
	return string(e)
}

func (e *CharType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = CharType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid CharType", str)
	}
	return nil
}

func (e CharType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type FetchType string

const (
	FetchTypePositive FetchType = "POSITIVE"
	FetchTypeNeutral  FetchType = "NEUTRAL"
	FetchTypeNegative FetchType = "NEGATIVE"
)

var AllFetchType = []FetchType{
	FetchTypePositive,
	FetchTypeNeutral,
	FetchTypeNegative,
}

func (e FetchType) IsValid() bool {
	switch e {
	case FetchTypePositive, FetchTypeNeutral, FetchTypeNegative:
		return true
	}
	return false
}

func (e FetchType) String() string {
	return string(e)
}

func (e *FetchType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = FetchType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid FetchType", str)
	}
	return nil
}

func (e FetchType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type FieldType string

const (
	FieldTypeEmail   FieldType = "EMAIL"
	FieldTypeDate    FieldType = "DATE"
	FieldTypeLink    FieldType = "LINK"
	FieldTypeText    FieldType = "TEXT"
	FieldTypeNumeric FieldType = "NUMERIC"
)

var AllFieldType = []FieldType{
	FieldTypeEmail,
	FieldTypeDate,
	FieldTypeLink,
	FieldTypeText,
	FieldTypeNumeric,
}

func (e FieldType) IsValid() bool {
	switch e {
	case FieldTypeEmail, FieldTypeDate, FieldTypeLink, FieldTypeText, FieldTypeNumeric:
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

type GroupType string

const (
	GroupTypeRole   GroupType = "ROLE"
	GroupTypeChar   GroupType = "CHAR"
	GroupTypeMember GroupType = "MEMBER"
)

var AllGroupType = []GroupType{
	GroupTypeRole,
	GroupTypeChar,
	GroupTypeMember,
}

func (e GroupType) IsValid() bool {
	switch e {
	case GroupTypeRole, GroupTypeChar, GroupTypeMember:
		return true
	}
	return false
}

func (e GroupType) String() string {
	return string(e)
}

func (e *GroupType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = GroupType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid GroupType", str)
	}
	return nil
}

func (e GroupType) MarshalGQL(w io.Writer) {
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

type ModifyType string

const (
	ModifyTypeAdd    ModifyType = "ADD"
	ModifyTypeReduce ModifyType = "REDUCE"
)

var AllModifyType = []ModifyType{
	ModifyTypeAdd,
	ModifyTypeReduce,
}

func (e ModifyType) IsValid() bool {
	switch e {
	case ModifyTypeAdd, ModifyTypeReduce:
		return true
	}
	return false
}

func (e ModifyType) String() string {
	return string(e)
}

func (e *ModifyType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ModifyType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ModifyType", str)
	}
	return nil
}

func (e ModifyType) MarshalGQL(w io.Writer) {
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
