package repository

import (
	"database/sql"
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
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
func (r *ChatsRepo) Members(chatId int) (members models.Members, err error) {
	rows, err := r.db.Query(
		`SELECT units.id, units.domain, units.name, units.type, member.role_id, member.char, member.joined_at
		FROM units INNER JOIN chat_members AS member
		ON units.id = member.id
		WHERE member.chat_id = $1
		`,
		chatId,
	)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := models.Member{}
		if err = rows.Scan(&m.User.Unit.ID, &m.User.Unit.Domain, &m.User.Unit.Name, &m.User.Unit.Type, &m.RoleID, &m.Char, &m.JoinedAt); err != nil {
			return
		}
		members.Members = append(members.Members, m)
	}

	return
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
		`SELECT EXISTS(SELECT 1 FROM chat_members WHERE user_id = $1 AND chat_id = $2)`,
		userId,
		chatId,
	).Scan(&isMember)
	if err != nil || !isMember {
		return isMember
	}
	return isMember
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

func (r *ChatsRepo) UserRole(userId int, chatId int) (role models.Role, err error) {
	err = r.db.QueryRow(
		`SELECT id, name, color
		FROM roles 
		WHERE id = (
			SELECT role_id
			FROM chat_members
			WHERE user_id = $1 AND chat_id = $2
		)`,
		userId,
		chatId,
	).Scan(
		&role.ID,
		&role.Name,
		&role.Color,
	)

	return
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

func (r *ChatsRepo) GiveRole(userId, chatId, roleId int) (err error) {
	err = r.db.QueryRow(
		`UPDATE chat_members
		SET role_id = $3
		WHERE user_id = $1 AND chat_id = $2`,
		userId,
		chatId,
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

func (r *ChatsRepo) InviteInfo(code string) (info model.InviteInfo, err error) {
	err = r.db.QueryRow(
		`SELECT units.id, units.domain, units.name, units.type, chats.private
		FROM units INNER JOIN chats
		ON units.id = chats.id
		WHERE units.id IN (
		    SELECT chat_id
		    FROM invites
		    WHERE code = $1
		)`,
	).Scan(
		&info.Unit.ID,
		&info.Unit.Domain,
		&info.Unit.Name,
		&info.Unit.Type,
		&info.Private,
	)

	return
}
