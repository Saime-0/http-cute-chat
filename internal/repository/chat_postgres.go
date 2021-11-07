package repository

import (
	"database/sql"
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

func (r *ChatsRepo) CreateChat(owner_id int, chat_model *models.CreateChat) (id int, err error) {
	// todo: add owner to chat members
	err = r.db.QueryRow(
		`WITH u AS (
			INSERT INTO units (domain, name) 
			VALUES ($1, $2) 
			RETURNING id
			) 
		INSERT INTO chats (id, owner_id, private) 
		SELECT u.id, $3, $4
		FROM u 
		RETURNING id`,
		chat_model.Domain,
		chat_model.Name,
		owner_id,
		chat_model.Private,
	).Scan(&id)

	return
}
func (r *ChatsRepo) GetChatByDomain(domain string) (chat models.ChatInfo, err error) {
	err = r.db.QueryRow(
		`SELECT units.id, chats.owner_id, units.domain,units.name
		FROM units INNER JOIN chats 
		ON units.id = chats.id 
		WHERE units.domain = $1`,
		domain,
	).Scan(
		&chat.ID,
		&chat.OwnerID,
		&chat.Domain,
		&chat.Name,
	)
	if err != nil {
		return
	}
	chat.CountMembers, err = r.GetCountChatMembers(chat.ID)
	if err != nil {
		return
	}
	return
}
func (r *ChatsRepo) GetChatByID(chat_id int) (chat models.ChatInfo, err error) {
	err = r.db.QueryRow(
		`SELECT units.id, chats.owner_id, units.domain,units.name
		FROM units INNER JOIN chats 
		ON units.id = chats.id 
		WHERE units.id = $1`,
		chat_id,
	).Scan(
		&chat.ID,
		&chat.OwnerID,
		&chat.Domain,
		&chat.Name,
	)
	if err != nil {
		return
	}
	chat.CountMembers, err = r.GetCountChatMembers(chat_id)
	if err != nil {
		return
	}
	return
}
func (r *ChatsRepo) GetCountChatMembers(chat_id int) (count int, err error) {
	err = r.db.QueryRow(
		`SELECT count(*)
		FROM chat_members  
		WHERE chat_id = $1`,
		chat_id,
	).Scan(&count)

	return
}

func (r *ChatsRepo) GetChatsByNameFragment(name string, offset int) (chats models.ListChatInfo, err error) {
	rows, err := r.db.Query(
		`SELECT units.id, chats.owner_id, units.domain,units.name
		FROM units INNER JOIN chats 
		ON units.id = chats.id 
		WHERE units.name ILIKE $1 AND chats.private = FALSE
		LIMIT 20
		OFFSET $2`,
		"%"+name+"%",
		offset,
	)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := models.ChatInfo{}
		if err = rows.Scan(&m.ID, &m.OwnerID, &m.Domain, &m.Name); err != nil {
			return
		}
		m.CountMembers, err = r.GetCountChatMembers(m.ID)
		if err != nil {
			return
		}
		chats.Chats = append(chats.Chats, m)
	}
	if !rows.NextResultSet() {
		return
	}
	return

}
func (r *ChatsRepo) GetChatMembers(chat_id int) (members models.ListUserInfo, err error) {
	rows, err := r.db.Query(
		`SELECT id, domain, name
		FROM units
		WHERE id IN (
			SELECT user_id 
			FROM chat_members
			WHERE chat_id = $1
			)`,
		chat_id,
	)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := models.UserInfo{}
		if err = rows.Scan(&m.ID, &m.Domain, &m.Name); err != nil {
			return
		}
		members.Users = append(members.Users, m)
	}
	if !rows.NextResultSet() {
		return
	}
	return
}

func (r *ChatsRepo) GetChatDataByID(chat_id int) (chat models.ChatData, err error) {
	err = r.db.QueryRow(
		`SELECT units.id, chats.owner_id, units.domain, units.name, chats.private
		FROM units INNER JOIN chats 
		ON units.id = chats.id 
		WHERE units.id = $1`,
		chat_id,
	).Scan(
		&chat.ID,
		&chat.OwnerID,
		&chat.Domain,
		&chat.Name,
		&chat.Private,
	)

	return
}

func (r *ChatsRepo) UpdateChatData(chat_id int, input_model *models.UpdateChatData) (err error) {
	if input_model.Domain != "" {
		err = r.db.QueryRow(
			`UPDATE units
			SET domain = $2
			WHERE id = $1`,
			chat_id,
			input_model.Domain,
		).Err()
		if err != nil {
			return
		}
	}
	if input_model.Name != "" {
		err = r.db.QueryRow(
			`UPDATE units
			SET name = $2
			WHERE id = $1`,
			chat_id,
			input_model.Name,
		).Err()
		if err != nil {
			return
		}
	}

	err = r.db.QueryRow(
		`UPDATE users
		SET private = $2
		WHERE id = $1`,
		chat_id,
		input_model.Private,
	).Err()

	return
}

func (r *ChatsRepo) UserIsChatOwner(user_id int, chat_id int) bool {
	is_owner := false
	err := r.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM chats WHERE id = $1 AND owner_id = $2)`,
		chat_id,
		user_id,
	).Scan(&is_owner)
	if err != nil || !is_owner {
		return is_owner
	}
	return is_owner
}
func (r *ChatsRepo) UserIsChatMember(user_id int, chat_id int) bool {
	is_member := false
	err := r.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM chat_members WHERE user_id = $1 AND chat_id = $2)`,
		user_id,
		chat_id,
	).Scan(&is_member)
	if err != nil || !is_member {
		return is_member
	}
	return is_member
}
func (r *ChatsRepo) AddUserToChat(user_id int, chat_id int) (err error) {
	err = r.db.QueryRow(
		`INSERT INTO chat_members (user_id, chat_id, joined_at)
		VALUES ($1, $2, $3)`,
		user_id,
		chat_id,
		time.Now().UTC().Unix(),
	).Err()
	if err != nil {
		return
	}
	return
}

// migrate from UsersRepo
func (r *ChatsRepo) GetChatsOwnedUser(user_id int, offset int) (chats models.ListChatInfo, err error) {
	rows, err := r.db.Query(
		`SELECT units.id, chats.owner_id, units.domain,units.name
		FROM units INNER JOIN chats 
		ON units.id = chats.id 
		WHERE units.id IN (
			SELECT chats.id
			FROM chats INNER JOIN chat_members 
			ON chats.owner_id = chat_members.user_id
			WHERE chats.owner_id = $1
			LIMIT 20
			OFFSET $2
			)`,
		user_id,
		offset,
	)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := models.ChatInfo{}
		if err = rows.Scan(&m.ID, &m.OwnerID, &m.Domain, &m.Name); err != nil {
			return
		}
		m.CountMembers, err = r.GetCountChatMembers(m.ID)
		if err != nil {
			return
		}
		chats.Chats = append(chats.Chats, m)
	}
	if !rows.NextResultSet() {
		return
	}
	return
}

func (r *ChatsRepo) GetChatsInvolvedUser(user_id int, offset int) (chats models.ListChatInfo, err error) {
	rows, err := r.db.Query(
		`SELECT units.id, chats.owner_id, units.domain,units.name
		FROM units INNER JOIN chats 
		ON units.id = chats.id 
		WHERE units.id IN (
			SELECT chat_id 
			FROM chat_members
			WHERE user_id = $1
			LIMIT 20
			OFFSET $2
			)`,
		user_id,
		offset,
	)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := models.ChatInfo{}
		if err = rows.Scan(&m.ID, &m.OwnerID, &m.Domain, &m.Name); err != nil {
			return
		}
		m.CountMembers, err = r.GetCountChatMembers(m.ID)
		if err != nil {
			return
		}
		chats.Chats = append(chats.Chats, m)
	}
	if !rows.NextResultSet() {
		return
	}
	return
}

func (r *ChatsRepo) GetCountUserChats(user_id int) (count int, err error) {
	err = r.db.QueryRow(
		`SELECT count(*)
		FROM units INNER JOIN chats 
		ON units.id = chats.id 
		WHERE units.id IN (
			SELECT chat_id 
			FROM chat_members
			WHERE user_id = $1
			)`,
		user_id,
	).Scan(&count)
	return
}

func (r *ChatsRepo) GetCountRooms(chat_id int) (count int, err error) {
	err = r.db.QueryRow(
		`SELECT count(*)
		FROM rooms
		WHERE chat_id = $1`,
		chat_id,
	).Scan(&count)
	return
}

func (r *ChatsRepo) ChatExistsByID(chat_id int) (exists bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1
			FROM chats
			WHERE id = $1
		)`,
		chat_id,
	).Scan(&exists)

	return
}

func (r *ChatsRepo) ChatExistsByDomain(chat_domain string) (exists bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1
			FROM units
			INNER JOIN chats
			ON chats.id = units.id
			WHERE units.domain = $1
		)`,
		chat_domain,
	).Scan(&exists)

	return
}

func (r *ChatsRepo) RemoveUserFromChat(user_id int, chat_id int) (err error) {
	err = r.db.QueryRow(
		`DELETE FROM chat_members
		WHERE user_id = $1 AND chat_id = $2`,
		user_id,
		chat_id,
	).Err()

	return
}

func (r *ChatsRepo) GetCountLinks(chat_id int) (count int, err error) {
	err = r.db.QueryRow(
		`SELECT count(*)
		FROM invite_links  
		WHERE chat_id = $1`,
		chat_id,
	).Scan(&count)

	return
}
func (r *ChatsRepo) GetChatLinks(chat_id int) (links models.InviteLinks, err error) {
	rows, err := r.db.Query(
		`SELECT code, aliens, exp
		FROM invite_links
		WHERE chat_id = $1`,
		chat_id,
	)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := models.InviteLink{}
		if err = rows.Scan(&m.Code, &m.Aliens, &m.Exp); err != nil {
			return
		}
		links.Links = append(links.Links, m)
	}

	return
}
func (r *ChatsRepo) LinkExistsByCode(code string) (exists bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1
			FROM invite_links
			WHERE code = $1
			)`,
		code,
	).Scan(&exists)

	return
}
func (r *ChatsRepo) FindInviteLinkByCode(code string) (link models.InviteLink, err error) {
	err = r.db.QueryRow(
		`SELECT code, aliens, exp
		FROM invite_links  
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
		`DELETE FROM invite_links
		WHERE code = $1`,
		code,
	).Err()

	return
}

func (r *ChatsRepo) CreateInviteLink(link_model *models.CreateInviteLink) (link models.InviteLink, err error) {
	err = r.db.QueryRow(
		`INSERT INTO invite_links (chat_id, aliens, exp) 
		VALUES ($1, $2, $3)
		RETURNING code, aliens, exp`,
		link_model.ChatID,
		link_model.Aliens,
		link_model.Exp,
	).Scan(
		&link.Code,
		&link.Aliens,
		&link.Exp,
	)

	return
}

func (r *ChatsRepo) InviteLinkIsRelevant(code string) (relevant bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1
			FROM invite_links
			WHERE code = $1 AND aliens > 0 AND exp > $2
			)`,
		code,
		time.Now().UTC().Unix(),
	).Scan(&relevant)
	if !relevant {
		r.db.Exec(
			`DELETE FROM invite_links
			WHERE code = $1`,
			code,
		)
	}

	return
}

func (r *ChatsRepo) AddUserByCode(code string, user_id int) (chat_id int, err error) {

	err = r.db.QueryRow(
		`WITH l AS (
			UPDATE invite_links
			SET aliens = aliens - 1
			WHERE code = $1
			RETURNING chat_id
			)
		INSERT INTO chat_members (user_id, chat_id)
		VALUES ($2, l.chat_id)`,
		code,
		user_id,
	).Scan(&chat_id)

	return
}

func (r *ChatsRepo) ChatIsPrivate(chat_id int) (private bool) {
	r.db.QueryRow(
		`SELECT private
		FROM chats
		WHERE chat_id = $1`,
		chat_id,
	).Scan(&private)

	return
}

func (r *ChatsRepo) BanUserInChat(user_id int, chat_id int) (err error) {
	err = r.db.QueryRow(
		`INSERT INTO chat_banlist (chat_id, user_id)
		VALUES ($1, $2)`,
		chat_id,
		user_id,
	).Err()

	return
}

func (r *ChatsRepo) UnbanUserInChat(user_id int, chat_id int) (err error) {
	err = r.db.QueryRow(
		`DELETE FROM chat_banlist
		WHERE chat_id = $1 AND user_id = $2`,
		chat_id,
		user_id,
	).Err()

	return
}

func (r *ChatsRepo) UserIsBannedInChat(user_id int, chat_id int) (banned bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1
			FROM chat_banlist
			WHERE chat_id = $1 AND user_id = $2
		)`,
		chat_id,
		user_id,
	).Scan(&banned)

	return
}

func (r *ChatsRepo) GetChatBanlist(chat_id int) (users models.ListUserInfo, err error) {
	rows, err := r.db.Query(
		`SELECT id, domain, name
		FROM units
		WHERE id IN (
			SELECT user_id 
			FROM chat_banlist
			WHERE chat_id = $1
			)`,
		chat_id,
	)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := models.UserInfo{}
		if err = rows.Scan(&m.ID, &m.Domain, &m.Name); err != nil {
			return
		}
		users.Users = append(users.Users, m)
	}
	if !rows.NextResultSet() {
		return
	}
	return
}

func (r *ChatsRepo) GetUserRoleData(user_id int, chat_id int) (role models.RoleData, err error) {
	err = r.db.QueryRow(
		`SELECT id, role_name, color, visible, manage_rooms, COALESCE(room_id, 0), manage_chat, manage_roles, manage_members
		FROM roles 
		WHERE id = (
			SELECT role_id
			FROM chat_members
			WHERE user_id = $1 AND chat_id = $2
		)`,
		user_id,
		chat_id,
	).Scan(
		&role.ID,
		&role.RoleName,
		&role.Color,
		&role.Visible,
		&role.ManageRooms,
		&role.RoomID,
		&role.ManageChat,
		&role.ManageRoles,
		&role.ManageMembers,
	)

	return
}

func (r *ChatsRepo) GetUserRoleInfo(user_id int, chat_id int) (role models.RoleInfo, err error) {
	err = r.db.QueryRow(
		`SELECT id, role_name, color
		FROM roles 
		WHERE id = (
			SELECT role_id
			FROM chat_members
			WHERE user_id = $1 AND chat_id = $2
		)`,
		user_id,
		chat_id,
	).Scan(
		&role.ID,
		&role.RoleName,
		&role.Color,
	)

	return
}

func (r *ChatsRepo) CreateRoleInChat(chat_id int, role_model *models.CreateRole) (role_id int, err error) {
	err = r.db.QueryRow(
		`INSERT INTO roles
		(chat_id, role_name, color, visible, manage_rooms, room_id, manage_chat, manage_roles, manage_members)
		VALUES ($1, $2, $3, $4, $5, NULLIF($6, 0), $7, $8, $9)
		RETURNING id`,
		chat_id,
		role_model.RoleName,
		role_model.Color,
		role_model.Visible,
		role_model.ManageRooms,
		role_model.RoomID,
		role_model.ManageChat,
		role_model.ManageRoles,
		role_model.ManageMembers,
	).Scan(&role_id)

	return
}

func (r *ChatsRepo) GetChatRolesData(chat_id int) (roles models.ListRolesData, err error) {
	rows, err := r.db.Query(
		`SELECT id, role_name, color, visible, manage_rooms, COALESCE(room_id, 0), manage_chat, manage_roles, manage_members
		FROM roles 
		WHERE chat_id = $1`,
		chat_id,
	)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := models.RoleData{}
		if err = rows.Scan(&m.ID, &m.RoleName, &m.Color, &m.Visible, &m.ManageRooms, &m.RoomID, &m.ManageChat, &m.ManageRoles, &m.ManageMembers); err != nil {
			return
		}
		roles.Roles = append(roles.Roles, m)
	}

	return
}

func (r *ChatsRepo) GetChatRolesInfo(chat_id int) (roles models.ListRolesInfo, err error) {
	rows, err := r.db.Query(
		`SELECT id, role_name, color
		FROM roles 
		WHERE chat_id = $1`,
		chat_id,
	)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := models.RoleInfo{}
		if err = rows.Scan(&m.ID, &m.RoleName, &m.Color); err != nil {
			return
		}
		roles.Roles = append(roles.Roles, m)
	}

	return
}
func (r *ChatsRepo) GetCountChatRoles(chat_id int) (count int, err error) {
	err = r.db.QueryRow(
		`SELECT count(*)
		FROM roles 
		WHERE chat_id = $1`,
		chat_id,
	).Scan(&count)

	return
}

func (r *ChatsRepo) GiveRole(user_id int, role_id int) (err error) {
	err = r.db.QueryRow(
		`UPDATE chat_members
		SET role_id = $2
		WHERE user_id = $1`,
		user_id,
		role_id,
	).Err()

	return
}

func (r *ChatsRepo) RoleExistsByID(role_id int, chat_id int) (exists bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1
			FROM roles
			WHERE role_id = $1 AND chat_id = $2
		)`,
		role_id,
		chat_id,
	).Scan(&exists)

	return
}

func (r *ChatsRepo) UpdateRoleData(role_id int, input_model *models.UpdateRole) (err error) {
	err = r.db.QueryRow(
		`UPDATE roles
		SET role_name = $2, color = $3, visible = $4, manage_rooms = $5, room_id =  NULLIF($6, 0), manage_chat, manage_roles, manage_members
		WHERE id = $1`,
		role_id,
		input_model.RoleName,
		input_model.Color,
		input_model.Visible,
		input_model.ManageRooms,
		input_model.RoomID,
		input_model.ManageChat,
		input_model.ManageRoles,
		input_model.ManageMembers,
	).Err()

	return
}

func (r *ChatsRepo) DeleteRole(role_id int) (err error) {
	err = r.db.QueryRow(
		`DELETE FROM roles
		WHERE id = $1`,
		role_id,
	).Err()

	return
}

func (r *ChatsRepo) TakeRole(user_id int, chat_id int) (err error) {
	err = r.db.QueryRow(
		`UPDATE chat_members
		SET role_id = NULL
		WHERE user_id = $1`,
		user_id,
		chat_id,
	).Err()

	return
}
