package repository

import (
	"database/sql"

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
		INSERT INTO chats (id, owner_id) 
		SELECT u.id, $3 
		FROM u 
		RETURNING id`,
		chat_model.Domain,
		chat_model.Name,
		owner_id,
	).Scan(&id)
	if err != nil {
		return
	}
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
	if err != nil {
		return
	}
	return
}

func (r *ChatsRepo) GetChatsByNameFragment(name string, offset int) (chats models.ListChatInfo, err error) {
	rows, err := r.db.Query(
		`SELECT units.id, chats.owner_id, units.domain,units.name
		FROM units INNER JOIN chats 
		ON units.id = chats.id 
		WHERE units.name ILIKE $1
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
		`SELECT units.id, chats.owner_id, units.domain, units.name
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
		`INSERT INTO chat_members (user_id, chat_id)
		VALUES ($1, $2)`,
		user_id,
		chat_id,
	).Err()
	if err != nil {
		return
	}
	return
}

// migrate from UsersRepo
func (r *ChatsRepo) GetChatsOwnedUser(user_id int) (chats models.ListChatInfo, err error) {
	rows, err := r.db.Query(
		`SELECT units.id, chats.owner_id, units.domain,units.name
		FROM units INNER JOIN chats 
		ON units.id = chats.id 
		WHERE units.id IN (
			SELECT chats.id
			FROM chats INNER JOIN chat_members 
			ON chats.owner_id = chat_members.user_id
			WHERE chats.owner_id = $1
			)`,
		user_id,
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

func (r *ChatsRepo) GetChatsInvolvedUser(user_id int) (chats models.ListChatInfo, err error) {
	rows, err := r.db.Query(
		`SELECT units.id, chats.owner_id, units.domain,units.name
		FROM units INNER JOIN chats 
		ON units.id = chats.id 
		WHERE units.id IN (
			SELECT chat_id 
			FROM chat_members
			WHERE user_id = $1
			)`,
		user_id,
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
