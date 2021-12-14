package repository

import (
	"database/sql"
	"fmt"
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"github.com/saime-0/http-cute-chat/pkg/kit"
	"strconv"
	"time"

	"github.com/saime-0/http-cute-chat/internal/models"
)

type ChatsRepo struct {
	db *sql.DB
}

func NewChatsRepo(db *sql.DB) *ChatsRepo {
	return &ChatsRepo{
		db: db,
	}
}

func (r *ChatsRepo) CreateChat(ownerId int, chatModel *models.CreateChat) (id int, err error) {
	err = r.db.QueryRow(
		`WITH u AS (
			INSERT INTO units (domain, name, type) 
			VALUES ($1, $2, 'CHAT') 
			RETURNING id
			) 
		INSERT INTO chats (id, owner_id, private) 
		SELECT u.id, $3, $4
		FROM u 
		RETURNING id`,
		chatModel.Domain,
		chatModel.Name,
		ownerId,
		chatModel.Private,
	).Scan(&id)

	return
}

func (r *ChatsRepo) GetChatByID(chatId int) (chat models.Chat, err error) {
	err = r.db.QueryRow(
		`SELECT units.id, units.domain, units.name, chats.private
		FROM units INNER JOIN chats 
		ON units.id = chats.id 
		WHERE units.id = $1`,
		chatId,
	).Scan(
		&chat.Unit.ID,
		&chat.Unit.Domain,
		&chat.Unit.Name,
		&chat.Private,
	)

	return
}
func (r *ChatsRepo) Chat(chatId int) (*model.Chat, error) {
	chat := &model.Chat{
		Unit: &model.Unit{},
	}
	err := r.db.QueryRow(
		`SELECT units.id, units.domain, units.name, units.type, chats.private
		FROM units INNER JOIN chats 
		ON units.id = chats.id 
		WHERE units.id = $1`,
		chatId,
	).Scan(
		&chat.Unit.ID,
		&chat.Unit.Domain,
		&chat.Unit.Name,
		&chat.Unit.Type,
		&chat.Private,
	)

	return chat, err
}
func (r *ChatsRepo) ChatIDByInvite(code string) (chatId int, err error) {
	err = r.db.QueryRow(
		`SELECT chat_id
		FROM invites
		WHERE code = $1`,
		code,
	).Scan(&chatId)
	return
}
func (r *ChatsRepo) CountMembers(chatId int) (count int, err error) {
	err = r.db.QueryRow(
		`SELECT count(*)
		FROM chat_members  
		WHERE chat_id = $1`,
		chatId,
	).Scan(&count)

	return
}

func (r *ChatsRepo) GetChatsByNameFragment(fragment string, limit int, offset int) (chats []models.Chat, err error) {
	rows, err := r.db.Query(
		`SELECT units.id, units.domain, units.name, chats.private
		FROM units INNER JOIN chats 
		ON units.id = chats.id 
		WHERE units.name ILIKE $1 AND chats.private = FALSE
		LIMIT $2
		OFFSET $3`,
		"%"+fragment+"%",
		limit,
		offset,
	)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := models.Chat{}
		if err = rows.Scan(&m.Unit.ID, &m.Unit.Domain, &m.Unit.Name, &m.Private); err != nil {
			return
		}
		if err != nil {
			return
		}
		chats = append(chats, m)
	}

	return

}
func (r *ChatsRepo) Members(chatId int) (*model.Members, error) {
	members := &model.Members{}
	rows, err := r.db.Query(
		`SELECT units.id, units.domain, units.name, units.type, member.chat_id, member.char, member.joined_at, member.muted, member.frozen
		FROM units INNER JOIN chat_members AS member
		ON units.id = member.id
		WHERE member.chat_id = $1`,
		chatId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		m := &model.Member{
			User: &model.User{
				Unit: &model.Unit{},
			},
			Chat: &model.Chat{
				Unit: &model.Unit{},
			},
		}
		if err = rows.Scan(
			&m.User.Unit.ID,
			&m.User.Unit.Domain,
			&m.User.Unit.Name,
			&m.User.Unit.Type,
			&m.Chat.Unit.ID,
			&m.Char,
			&m.JoinedAt,
			&m.Muted,
			&m.Frozen); err != nil {
			return nil, err
		}
		members.Members = append(members.Members, m)
	}

	return members, nil
}

// MembersByArray sorry, method implemented in RoomsRepo!!!
func (r *RoomsRepo) MembersByArray(chatId int, memberIds *[]int) (*model.Members, error) {
	members := &model.Members{}
	// language= PostgreSQL
	query := `SELECT units.id, units.domain, units.name, units.type, member.chat_id, member.char, member.joined_at, member.muted, member.frozen
		FROM units INNER JOIN chat_members AS member
		ON units.id = member.user_id
		WHERE member.user_id IN (` + kit.CommaSeparate(memberIds) + `) AND member.chat_id =` + strconv.Itoa(chatId)
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		m := &model.Member{
			User: &model.User{
				Unit: &model.Unit{},
			},
			Chat: &model.Chat{
				Unit: &model.Unit{},
			},
		}
		if err = rows.Scan(
			&m.User.Unit.ID,
			&m.User.Unit.Domain,
			&m.User.Unit.Name,
			&m.User.Unit.Type,
			&m.Chat.Unit.ID,
			&m.Char,
			&m.JoinedAt,
			&m.Muted,
			&m.Frozen); err != nil {
			return nil, err
		}
		members.Members = append(members.Members, m)
	}
	return members, nil
}
func (r *ChatsRepo) GetChatMember(userId, chatId int) (member models.Member, err error) {
	err = r.db.QueryRow(
		`SELECT units.id, units.domain, units.name, units.type, member.role_id, member.char, member.joined_at
		FROM units INNER JOIN chat_members AS member
		ON units.id = member.id
		WHERE member.user_id = $1 AND member.chat_id = $2`,
		userId,
		chatId,
	).Scan(
		&member.User.Unit.ID,
		&member.User.Unit.Domain,
		&member.User.Unit.Name,
		&member.User.Unit.Type,
		&member.RoleID,
		&member.Char,
		&member.JoinedAt,
	)

	return
}

func (r *ChatsRepo) GetUserChar(userId int, chatId int) (char rules.CharType, err error) {
	err = r.db.QueryRow(
		`SELECT char
		FROM chat_members
		WHERE user_id = $1 AND chat_id = $2`,
		userId,
		chatId,
	).Scan(&char)

	return
}

func (r *ChatsRepo) UserIs(userId int, chatId int, char rules.CharType) (yes bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1
			FROM chat_members
			WHERE user_id = $1 AND chat_id = $2 AND char = $3
    	)`,
		userId,
		chatId,
		char,
	).Scan(&yes)

	return
}

func (r *ChatsRepo) Chars(chatId int, char rules.CharType) (admins models.Members, err error) {
	rows, err := r.db.Query(
		`SELECT units.id, units.domain, units.name, member.role_id, member.char, member.joined_at
		FROM units INNER JOIN chat_members AS member
		ON units.id = member.id
		WHERE member.chat_id = $1 AND member.char = $2`,
		chatId,
		char,
	)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := models.Member{}
		if err = rows.Scan(&m.User.Unit.ID, &m.User.Unit.Domain, &m.User.Unit.Name, &m.RoleID, &m.Char, &m.JoinedAt); err != nil {
			return
		}
		admins.Members = append(admins.Members, m)
	}
	return
}

//func (r *ChatsRepo) UpdateChatData(chatId int, inputModel *models.UpdateChatData) (err error) {
//	if inputModel.Domain != "" {
//		err = r.db.QueryRow(
//			`UPDATE units
//			SET domain = $2
//			WHERE id = $1`,
//			chatId,
//			inputModel.Domain,
//		).Err()
//		if err != nil {
//			return
//		}
//	}
//	if inputModel.Name != "" {
//		err = r.db.QueryRow(
//			`UPDATE units
//			SET name = $2
//			WHERE id = $1`,
//			chatId,
//			inputModel.Name,
//		).Err()
//		if err != nil {
//			return
//		}
//	}
//
//	err = r.db.QueryRow(
//		`UPDATE rooms
//		SET private = $2
//		WHERE id = $1`,
//		chatId,
//		inputModel.Private,
//	).Err()
//
//	return
//}

func (r *ChatsRepo) UserIsChatOwner(userId int, chatId int) bool {
	isOwner := false
	err := r.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM chats WHERE id = $1 AND owner_id = $2)`,
		chatId,
		userId,
	).Scan(&isOwner)
	if err != nil || !isOwner {
		return isOwner
	}
	return isOwner
}
func (r *ChatsRepo) UserIsChatMember(userId int, chatId int) bool {
	isMember := false
	err := r.db.QueryRow(
		`SELECT id
		FROM chat_members 
		WHERE user_id = $1 AND chat_id = $2`,
		userId,
		chatId,
	).Scan(&isMember)
	if err != nil || !isMember {
		return isMember
	}
	return isMember
}
func (r *ChatsRepo) FindMemberBy(userId int, chatId int) *int {
	var memberId *int
	err := r.db.QueryRow(
		`SELECT id
		FROM chat_members 
		WHERE user_id = $1 AND chat_id = $2`,
		userId,
		chatId,
	).Scan(&memberId)
	if err != nil {
		fmt.Println("FindMemberBy: произошла ошибка")
	}
	return memberId
}
func (r *ChatsRepo) AddUserToChat(userId int, chatId int) (err error) {
	err = r.db.QueryRow(
		`INSERT INTO chat_members (user_id, chat_id, role_id, char, joined_at, muted, frozen)
		VALUES ($1, $2, NULL, 'NONE', $3, false, false)`,
		userId,
		chatId,
		time.Now().UTC().Unix(),
	).Err()
	if err != nil {
		return
	}
	return
}

func (r *ChatsRepo) GetCountUserChats(userId int) (count int, err error) {
	err = r.db.QueryRow(
		`SELECT count(*)
		FROM units INNER JOIN chats 
		ON units.id = chats.id 
		WHERE units.id IN (
			SELECT chat_id 
			FROM chat_members
			WHERE user_id = $1
			)`,
		userId,
	).Scan(&count)
	return
}

func (r *ChatsRepo) GetCountRooms(chatId int) (count int, err error) {
	err = r.db.QueryRow(
		`SELECT count(*)
		FROM rooms
		WHERE chat_id = $1`,
		chatId,
	).Scan(&count)
	return
}

func (r *ChatsRepo) ChatExistsByID(chatId int) (exists bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1
			FROM chats
			WHERE id = $1
		)`,
		chatId,
	).Scan(&exists)

	return
}

func (r *ChatsRepo) ChatExistsByDomain(chatDomain string) (exists bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1
			FROM units
			INNER JOIN chats
			ON chats.id = units.id
			WHERE units.domain = $1
		)`,
		chatDomain,
	).Scan(&exists)

	return
}

func (r *ChatsRepo) RemoveUserFromChat(userId int, chatId int) (err error) {
	err = r.db.QueryRow(
		`DELETE FROM chat_members
		WHERE user_id = $1 AND chat_id = $2`,
		userId,
		chatId,
	).Err()

	return
}

func (r *ChatsRepo) GetCountLinks(chatId int) (count int, err error) {
	err = r.db.QueryRow(
		`SELECT count(*)
		FROM invites  
		WHERE chat_id = $1`,
		chatId,
	).Scan(&count)

	return
}
func (r *ChatsRepo) Invites(chatId int) (links models.Invites, err error) {
	rows, err := r.db.Query(
		`SELECT code, aliens, expires_at
		FROM invites
		WHERE chat_id = $1`,
		chatId,
	)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := models.Invite{}
		if err = rows.Scan(&m.Code, &m.Aliens, &m.Exp); err != nil {
			return
		}
		links.Invites = append(links.Invites, m)
	}

	return
}
func (r *ChatsRepo) InviteExistsByCode(code string) (exists bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1
			FROM invites
			WHERE code = $1
			)`,
		code,
	).Scan(&exists)

	return
}
func (r *ChatsRepo) FindInviteLinkByCode(code string) (link models.Invite, err error) {
	err = r.db.QueryRow(
		`SELECT code, aliens, expires_at
		FROM invites  
		WHERE code = $1`,
		code,
	).Scan(
		&link.Code,
		&link.Aliens,
		&link.Exp,
	)

	return
}
func (r *ChatsRepo) DeleteInviteLinkByCode(code string) (err error) {
	err = r.db.QueryRow(
		`DELETE FROM invites
		WHERE code = $1`,
		code,
	).Err()

	return
}

func (r *ChatsRepo) CreateInviteLink(linkModel *models.CreateInvite) (link models.Invite, err error) {
	err = r.db.QueryRow(
		`INSERT INTO invites (chat_id, aliens, expires_at) 
		VALUES ($1, $2, $3)
		RETURNING code, aliens, expires_at`,
		linkModel.ChatID,
		linkModel.Aliens,
		linkModel.Exp,
	).Scan(
		&link.Code,
		&link.Aliens,
		&link.Exp,
	)

	return
}

// InviteIsRelevant
//  alt: InviteIsExists
func (r *ChatsRepo) InviteIsRelevant(code string) (relevant bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1
			FROM invites
			WHERE code = $1 AND aliens > 0 AND expires_at > $2
			)`,
		code,
		time.Now().UTC().Unix(),
	).Scan(&relevant)
	if !relevant {
		r.db.Exec(
			`DELETE FROM invites
			WHERE code = $1`,
			code,
		)
	}

	return
}

func (r *ChatsRepo) AddUserByCode(code string, userId int) (chatId int, err error) {

	err = r.db.QueryRow(
		`WITH l AS (
			UPDATE invites
			SET aliens = aliens - 1
			WHERE code = $1
			RETURNING chat_id
			)
		INSERT INTO chat_members (user_id, chat_id, role_id, char, joined_at, muted, frozen)
		VALUES ($2, l.chat_id, NULL, 'NONE', $3, false, false)`,
		code,
		userId,
		time.Now().UTC().Unix(),
	).Scan(&chatId)

	return
}

func (r *ChatsRepo) ChatIsPrivate(chatId int) (private bool) {
	r.db.QueryRow(
		`SELECT private
		FROM chats
		WHERE id = $1`,
		chatId,
	).Scan(&private)

	return
}

func (r *ChatsRepo) AddToBanlist(userId int, chatId int) (err error) {
	err = r.db.QueryRow(
		`INSERT INTO chat_banlist (chat_id, user_id)
		VALUES ($1, $2)`,
		chatId,
		userId,
	).Err()

	return
}

func (r *ChatsRepo) RemoveFromBanlist(userId int, chatId int) (err error) {
	err = r.db.QueryRow(
		`DELETE FROM chat_banlist
		WHERE chat_id = $1 AND user_id = $2`,
		chatId,
		userId,
	).Err()

	return
}

func (r *ChatsRepo) UserIsBannedInChat(userId int, chatId int) (banned bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1
			FROM chat_banlist
			WHERE chat_id = $1 AND user_id = $2
		)`,
		chatId,
		userId,
	).Scan(&banned)

	return
}

func (r *ChatsRepo) Banlist(chatId int) (users models.Users, err error) {
	rows, err := r.db.Query(
		`SELECT units.id, units.domain, units.name, units.type
		FROM units INNER JOIN users 
		ON units.id = users.id 
		INNER JOIN chat_banlist
		ON units.id = chat_banlist.user_id
		WHERE chat_banlist.chat_id = $1`,
		chatId,
	)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := models.User{}
		if err = rows.Scan(&m.Unit.ID, &m.Unit.Domain, &m.Unit.Name, &m.Unit.Type); err != nil {
			return
		}
		users.Users = append(users.Users, m)
	}
	if !rows.NextResultSet() {
		return
	}
	return
}

func (r *ChatsRepo) MemberRole(memberId int) *model.Role {
	_role := models.RoleReference{}
	err := r.db.QueryRow(`
		SELECT roles.id, roles.name, roles.color
		FROM chat_members 
		LEFT JOIN roles  
		    ON chat_members.role_id = roles.id
		WHERE chat_members.id = $1
		`,
		memberId,
	).Scan(
		&_role.ID,
		&_role.Name,
		&_role.Color,
	)
	if _role.ID != nil {
		return &model.Role{
			ID:    *_role.ID,
			Name:  *_role.Name,
			Color: *_role.Color,
		}
	}
	if err != nil {
		fmt.Println("не найдена роль")
	}
	return nil
}

func (r *ChatsRepo) CreateRoleInChat(chatId int, roleModel *models.CreateRole) (roleId int, err error) {
	err = r.db.QueryRow(
		`INSERT INTO roles
		(chat_id, name, color)
		VALUES ($1, $2, $3)
		RETURNING id`,
		chatId,
		roleModel.Name,
		roleModel.Color,
	).Scan(&roleId)

	return
}

func (r *ChatsRepo) Roles(chatId int) (roles []models.Role, err error) {
	rows, err := r.db.Query(
		`SELECT id, name, color
		FROM roles 
		WHERE chat_id = $1`,
		chatId,
	)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := models.Role{}
		if err = rows.Scan(&m.ID, &m.Name, &m.Color); err != nil {
			return
		}
		roles = append(roles, m)
	}

	return
}

func (r *ChatsRepo) GetCountChatRoles(chatId int) (count int, err error) {
	err = r.db.QueryRow(
		`SELECT count(*)
		FROM roles 
		WHERE chat_id = $1`,
		chatId,
	).Scan(&count)

	return
}

func (r *ChatsRepo) GiveRole(memberId, roleId int) (err error) {
	err = r.db.QueryRow(
		`UPDATE chat_members
		SET role_id = $2
		WHERE id = $1`,
		memberId,
		roleId,
	).Err()

	return
}

//func (r *ChatsRepo) UpdateRoleData(roleId int, inputModel *models.UpdateRole) (err error) {
//	err = r.db.QueryRow(
//		`UPDATE roles
//		SET role_name = $2, color = $3, visible = $4, manage_rooms = $5, room_id =  NULLIF($6, 0), manage_chat = $7, manage_roles = $8, manage_members = $9
//		WHERE id = $1`,
//		roleId,
//		inputModel.RoleName,
//		inputModel.Color,
//		inputModel.Visible,
//		inputModel.ManageRooms,
//		inputModel.RoomID,
//		inputModel.ManageChat,
//		inputModel.ManageRoles,
//		inputModel.ManageMembers,
//	).Err()
//
//	return
//}

func (r *ChatsRepo) DeleteRole(roleId int) (err error) {
	err = r.db.QueryRow(
		`DELETE FROM roles
		WHERE id = $1`,
		roleId,
	).Err()

	return
}

func (r *ChatsRepo) TakeRole(userId int, chatId int) (err error) {
	err = r.db.QueryRow(
		`UPDATE chat_members
		SET role_id = NULL
		WHERE user_id = $1`,
		userId,
		chatId,
	).Err()

	return
}

func (r *ChatsRepo) GetMemberInfo(userId int, chatId int) (user models.MemberInfo, err error) {
	err = r.db.QueryRow(
		`SELECT role_id, joined_at
		FROM chat_members
		WHERE user_id = $1 AND chat_id = $2`,
		userId,
		chatId,
	).Scan(
		&user.RoleID,
		&user.JoinedAt,
	)

	return
}

func (r *ChatsRepo) HasInvite(chatId int, code string) (has bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
		SELECT 1 
		FROM invites
		WHERE chat_id = $1 AND code = $2
		)`,
		chatId,
		code,
	).Scan(&has)

	return
}

func (r *ChatsRepo) RoleExistsByID(chatId, roleId int) (exists bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
		SELECT 1 
		FROM roles
		WHERE chat_id = $1 AND id = $2
		)`,
		chatId,
		roleId,
	).Scan(&exists)

	return
}

func (r *ChatsRepo) InviteInfo(code string) (info *model.InviteInfo, err error) {
	info = &model.InviteInfo{
		Unit: &model.Unit{},
	}
	err = r.db.QueryRow(
		`SELECT units.id, units.domain, units.name, units.type, chats.private
		FROM units 
		INNER JOIN chats
			ON units.id = chats.id
		INNER JOIN invites
			ON units.id = invites.chat_id
		WHERE invites.code = $1`,
		code,
	).Scan(
		&info.Unit.ID,
		&info.Unit.Domain,
		&info.Unit.Name,
		&info.Unit.Type,
		&info.Private,
	)

	return
}

func (r *ChatsRepo) Owner(chatId int) (*model.User, error) {
	owner := &model.User{
		Unit: &model.Unit{},
	}
	err := r.db.QueryRow(
		`SELECT units.id, units.domain, units.name, units.type
		FROM units 
		INNER JOIN users 
			ON units.id = users.id
		INNER JOIN chats 
			ON units.id = chats.owner_id
		WHERE chats.id = $1`,
		chatId,
	).Scan(
		&owner.Unit.ID,
		&owner.Unit.Domain,
		&owner.Unit.Name,
		&owner.Unit.Type,
	)

	return owner, err
}
func (r *ChatsRepo) UserIsBanned(userId int, chatId int) (banned bool) {
	r.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM chat_banlist WHERE user_id = $1 AND chat_id = $2)`,
		userId,
		chatId,
	).Scan(&banned)

	return
}
func (r *ChatsRepo) MemberBy(userId, chatId int) (*model.Member, error) {
	member := &model.Member{
		User: &model.User{
			Unit: &model.Unit{},
		},
		Chat: &model.Chat{
			Unit: &model.Unit{},
		},
	}
	err := r.db.QueryRow(
		`SELECT units.id, units.domain, units.name, units.type, member.chat_id, member.char, member.joined_at, member.muted, member.frozen
		FROM units INNER JOIN chat_members AS member
		ON units.id = member.id
		WHERE member.user_id = $1 AND member.chat_id = $2`,
		userId,
		chatId,
	).Scan(
		&member.User.Unit.ID,
		&member.User.Unit.Domain,
		&member.User.Unit.Name,
		&member.User.Unit.Type,
		&member.Chat.Unit.ID,
		&member.Char,
		&member.JoinedAt,
		&member.Muted,
		&member.Frozen,
	)
	return member, err
}
func (r *ChatsRepo) Member(memberId int) (*model.Member, error) {
	member := &model.Member{
		User: &model.User{
			Unit: &model.Unit{},
		},
		Chat: &model.Chat{
			Unit: &model.Unit{},
		},
	}
	err := r.db.QueryRow(
		`SELECT member.id, units.id, units.domain, units.name, units.type, member.chat_id, member.char, member.joined_at, member.muted, member.frozen
		FROM units INNER JOIN chat_members AS member
		ON units.id = member.id
		WHERE member.id = $1`,
		memberId,
	).Scan(
		&member.ID,
		&member.User.Unit.ID,
		&member.User.Unit.Domain,
		&member.User.Unit.Name,
		&member.User.Unit.Type,
		&member.Chat.Unit.ID,
		&member.Char,
		&member.JoinedAt,
		&member.Muted,
		&member.Frozen,
	)
	return member, err
}

func (r *ChatsRepo) Rooms(chatId int) (*model.Rooms, error) {
	rooms := &model.Rooms{}
	rows, err := r.db.Query(
		`SELECT id, parent_id, name, note
		FROM rooms
		WHERE chat_id = $1`,
		chatId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		m := &model.Room{}
		if err = rows.Scan(&m.RoomID, &m.ParentID, &m.Name, &m.Note); err != nil {
			return nil, err
		}
		rooms.Rooms = append(rooms.Rooms, m)
	}

	return rooms, nil
}

// RolesByArray sorry
func (r *RoomsRepo) RolesByArray(roleIds *[]int) (*model.Roles, error) {
	roles := &model.Roles{}
	query := `SELECT id, name, color
		FROM roles 
		WHERE id IN (` + kit.CommaSeparate(roleIds) + `)`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		m := &model.Role{}
		if err = rows.Scan(&m.ID, &m.Name, &m.Color); err != nil {
			return nil, err
		}
		roles.Roles = append(roles.Roles, m)
	}

	return roles, nil

}

func (r *ChatsRepo) FindMessages(inp *model.FindMessages, holder *models.AllowHolder) *model.Messages {
	messages := &model.Messages{
		Messages: []*model.Message{},
	}
	if inp.TextFragment != nil {
		*inp.TextFragment = "%" + *inp.TextFragment + "%"
	}

	// language=PostgreSQL
	rows, err := r.db.Query(`
		SELECT messages.id, reply_to, author, messages.room_id, body, messages.type, created_at
		FROM messages 
		JOIN allows
			ON messages.room_id = allows.room_id
		WHERE messages.room_id IN (
		    SELECT rooms.id 
		    FROM chats
		    INNER JOIN rooms
		        ON chats.id = rooms.chat_id
		    LEFT JOIN allows 
		    	ON rooms.id = allows.room_id 
		    WHERE chats.id = $1 
			AND (
			    $2::BIGINT IS NULL 
			    OR rooms.id = $2 
			)
	        AND (
	            action_type IS NULL 
				OR action_type = 'READ' 
                AND (
                    group_type = 'ROLES' AND value = $3
                    OR group_type = 'CHARS' AND value = $4
                    OR group_type = 'USERS' AND value = $5
                )
            )
		)
		AND (
		    $6::BIGINT IS NULL 
		    OR author = $6 
		)
		AND (
		    $7::VARCHAR IS NULL 
		    OR body ILIKE $7
		)
		`,
		inp.ChatID,
		inp.RoomID,
		holder.RoleID,
		holder.Char,
		holder.UserID,
		inp.AuthorID,
		inp.TextFragment,
	)
	if err != nil {
		println(err.Error())
		return messages
	}
	defer rows.Close()
	for rows.Next() {
		m := &model.Message{
			Room: &model.Room{},
		}
		var (
			_replid   *int
			_memberId *int
		)
		if err = rows.Scan(&m.ID, &_replid, &_memberId, &m.Room.RoomID, &m.Body, &m.Type, &m.CreatedAt); err != nil {
			println("rows.scan:", err.Error()) // debug
			return messages
		}
		if _replid != nil {
			m.ReplyTo = &model.Message{
				ID: *_replid,
			}
		}
		if _memberId != nil {
			m.Author = &model.Member{
				ID: *_memberId,
			}
		}
		messages.Messages = append(messages.Messages, m)
	}

	return messages
}

func (r *ChatsRepo) ChatIDByMemberID(memberId int) (chatId *int) {
	err := r.db.QueryRow(`
		SELECT chat_id
		FROM chat_members
		WHERE id = $1
		`,
		memberId,
	).Scan(&chatId)
	if err != nil {
		fmt.Println("не найден мембер")
		return
	}
	return
}

func (r *ChatsRepo) FindMessage(messageId int) *models.FindMember {
	panic("Not implemented")
}

func (r *ChatsRepo) MemberIsMuted(memberId int) (muted bool) {
	err := r.db.QueryRow(`
		SELECT EXISTS(
		    SELECT 1
		    FROM chat_members
		    WHERE  id = $1 AND muted = FALSE
		)
		`,
		memberId,
	).Scan(&muted)
	if err != nil {
		println(err.Error())
	}
	return
}

func (r *ChatsRepo) FindMemberID(inp *models.FindMember) {
	panic("Not implemented")
}
