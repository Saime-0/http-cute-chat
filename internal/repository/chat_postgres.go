package repository

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"github.com/saime-0/http-cute-chat/graph/model"
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

func (r *ChatsRepo) CreateChat(ownerId int, inp *model.CreateChatInput) (id int, err error) {
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
		inp.Domain,
		inp.Name,
		ownerId,
		inp.Private,
	).Scan(&id)

	return
}

func (r *ChatsRepo) GetChatByID(chatID int) (chat models.Chat, err error) {
	err = r.db.QueryRow(
		`SELECT units.id, units.domain, units.name, chats.private
		FROM units INNER JOIN chats 
		ON units.id = chats.id 
		WHERE units.id = $1`,
		chatID,
	).Scan(
		&chat.Unit.ID,
		&chat.Unit.Domain,
		&chat.Unit.Name,
		&chat.Private,
	)

	return
}
func (r *ChatsRepo) Chat(chatID int) (*model.Chat, error) {
	chat := &model.Chat{
		Unit: new(model.Unit),
	}
	err := r.db.QueryRow(
		`SELECT units.id, units.domain, units.name, units.type, cm.count_value, chats.private
		FROM units 
		JOIN chats ON units.id = chats.id
		JOIN count_members cm on chats.id = cm.chat_id
		WHERE units.id = $1`,
		chatID,
	).Scan(
		&chat.Unit.ID,
		&chat.Unit.Domain,
		&chat.Unit.Name,
		&chat.Unit.Type,
		&chat.CountMembers,
		&chat.Private,
	)

	return chat, err
}
func (r *ChatsRepo) ChatIDByRoleID(roleID int) (chatID int, err error) {
	err = r.db.QueryRow(`
		SELECT chats.id
		FROM chats JOIN roles ON chats.id = roles.chat_id
		WHERE roles.id = $1 
		`,
		roleID,
	).Scan(&chatID)

	return
}
func (r *ChatsRepo) ChatIDByInvite(code string) (chatID int, err error) {
	err = r.db.QueryRow(
		`SELECT chat_id
		FROM invites
		WHERE code = $1`,
		code,
	).Scan(&chatID)
	return
}
func (r *ChatsRepo) CountMembers(chatID int) (count int, err error) {
	err = r.db.QueryRow(
		`SELECT count(*)
		FROM chat_members  
		WHERE chat_id = $1`,
		chatID,
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
func (r *ChatsRepo) Members(chatID int) (*model.Members, error) {
	members := &model.Members{}
	rows, err := r.db.Query(
		`SELECT member.id, units.id, units.domain, units.name, units.type, member.chat_id, member.char, member.joined_at, member.muted
		FROM units INNER JOIN chat_members AS member
		ON units.id = member.user_id
		WHERE member.chat_id = $1`,
		chatID,
	)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		m := &model.Member{
			User: &model.User{
				Unit: new(model.Unit),
			},
			Chat: &model.Chat{
				Unit: new(model.Unit),
			},
		}
		if err = rows.Scan(
			&m.ID,
			&m.User.Unit.ID,
			&m.User.Unit.Domain,
			&m.User.Unit.Name,
			&m.User.Unit.Type,
			&m.Chat.Unit.ID,
			&m.Char,
			&m.JoinedAt,
			&m.Muted); err != nil {
			return nil, err
		}
		members.Members = append(members.Members, m)
	}

	return members, nil
}

// MembersByArray sorry, method implemented in RoomsRepo!!!
func (r *RoomsRepo) MembersByArray(chatID int, memberIDs *[]int) (*model.Members, error) {
	members := &model.Members{}
	// language= PostgreSQL
	query := `SELECT units.id, units.domain, units.name, units.type, member.chat_id, member.char, member.joined_at, member.muted
		FROM units INNER JOIN chat_members AS member
		ON units.id = member.user_id
		WHERE member.user_id IN (` + kit.CommaSeparate(memberIDs) + `) AND member.chat_id =` + strconv.Itoa(chatID)
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		m := &model.Member{
			User: &model.User{
				Unit: new(model.Unit),
			},
			Chat: &model.Chat{
				Unit: new(model.Unit),
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
			&m.Muted); err != nil {
			return nil, err
		}
		members.Members = append(members.Members, m)
	}
	return members, nil
}

func (r *ChatsRepo) UserIsChatOwner(userID int, chatID int) bool {
	isOwner := false
	err := r.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM chats WHERE id = $1 AND owner_id = $2)`,
		chatID,
		userID,
	).Scan(&isOwner)
	if err != nil || !isOwner {
		return isOwner
	}
	return isOwner
}

// deprecated
func (r *ChatsRepo) UserIsChatMember(userID int, chatID int) (isMember bool, err error) {

	err = r.db.QueryRow(`
		SELECT EXISTS(
	        SELECT 1
			FROM chat_members 
			WHERE user_id = $1 AND chat_id = $2
	    )
		`,
		userID,
		chatID,
	).Scan(&isMember)

	return
}
func (r *ChatsRepo) FindMemberBy(userID int, chatID int) *int {
	var memberID *int
	err := r.db.QueryRow(
		`SELECT id
		FROM chat_members 
		WHERE user_id = $1 AND chat_id = $2`,
		userID,
		chatID,
	).Scan(&memberID)
	if err != nil {
		fmt.Println("FindMemberBy: произошла ошибка")
	}
	return memberID
}
func (r *ChatsRepo) AddUserToChat(userID int, chatID int) (*model.CreateMember, error) {
	member := &model.CreateMember{
		Unit: new(model.Unit),
	}
	err := r.db.QueryRow(`		
		WITH member AS (
			    INSERT INTO chat_members (user_id, chat_id)
				VALUES ($1, $2)
				RETURNING id, chat_id, user_id
		)
		SELECT member.id, chat_id, user_id, domain, name, type
		FROM member join units on user_id = units.id`,
		userID,
		chatID,
	).Scan(
		&member.ID,
		&member.ChatID,
		&member.Unit.ID,
		&member.Unit.Domain,
		&member.Unit.Name,
		&member.Unit.Type,
	)

	return member, err
}

func (r *ChatsRepo) GetCountUserChats(userID int) (count int, err error) {
	err = r.db.QueryRow(
		`SELECT count(*)
		FROM units INNER JOIN chats 
		ON units.id = chats.id 
		WHERE units.id IN (
			SELECT chat_id 
			FROM chat_members
			WHERE user_id = $1
			)`,
		userID,
	).Scan(&count)
	return
}

func (r *ChatsRepo) GetCountRooms(chatID int) (count int, err error) {
	err = r.db.QueryRow(
		`SELECT count(*)
		FROM rooms
		WHERE chat_id = $1`,
		chatID,
	).Scan(&count)
	return
}

func (r *ChatsRepo) ChatExistsByID(chatID int) (exists bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1
			FROM chats
			WHERE id = $1
		)`,
		chatID,
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

func (r *ChatsRepo) RemoveUserFromChat(userID int, chatID int) (*model.DeleteMember, error) {
	member := &model.DeleteMember{}
	err := r.db.QueryRow(`
		DELETE FROM chat_members
		WHERE user_id = $1 AND chat_id = $2`,
		userID,
		chatID,
	).Scan(
		&member.ID,
	)

	return member, err
}

func (r *ChatsRepo) GetCountLinks(chatID int) (count int, err error) {
	err = r.db.QueryRow(
		`SELECT count(*)
		FROM invites  
		WHERE chat_id = $1`,
		chatID,
	).Scan(&count)

	return
}
func (r *ChatsRepo) Invites(chatID int) (*model.Invites, error) {
	invites := &model.Invites{
		Invites: []*model.Invite{},
	}
	rows, err := r.db.Query(
		`SELECT code, aliens, expires_at
		FROM invites
		WHERE chat_id = $1`,
		chatID,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		m := &model.Invite{}
		if err = rows.Scan(&m.Code, &m.Aliens, &m.ExpiresAt); err != nil {
			return nil, err
		}
		invites.Invites = append(invites.Invites, m)
	}

	return invites, nil
}

func (r *ChatsRepo) InviteExistsByCode(code string) (exists bool, err error) {
	err = r.db.QueryRow(`
			SELECT EXISTS(
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
func (r *ChatsRepo) DeleteInvite(code string) (*model.DeleteInvite, error) {
	invite := &model.DeleteInvite{
		Reason: model.DeleteInviteReasonByuser,
	}
	err := r.db.QueryRow(`
		DELETE FROM invites
		WHERE code = $1
		RETURNING code`,
		code,
	).Scan(&invite.Code)

	return invite, err
}

func (r *ChatsRepo) CreateInvite(linkModel *model.CreateInviteInput) (*model.CreateInvite, error) {
	invite := &model.CreateInvite{}
	err := r.db.QueryRow(
		`INSERT INTO invites (chat_id, aliens, expires_at) 
		VALUES ($1, $2, unix_utc_now($3::BIGINT))
		RETURNING code, aliens, expires_at`,
		linkModel.ChatID,
		linkModel.Aliens,
		linkModel.Duration,
	).Scan(
		&invite.Code,
		&invite.Aliens,
		&invite.ExpiresAt,
	)

	return invite, err
}

//  alt: InviteIsExists
func (r *ChatsRepo) InviteIsRelevant(code string) (relevant bool, err error) {
	err = r.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM invites
			WHERE code = $1 
			  	AND (
			  	    aliens IS NULL 
			    	OR aliens > 0
				) 
				AND (
					expires_at IS NULL 
					OR expires_at > $2
				)
		)`,
		code,
		time.Now().UTC().Unix(),
	).Scan(&relevant)
	if err != nil {
		return
	}
	if !relevant {
		err = r.db.QueryRow(`
			DELETE FROM invites
			WHERE code = $1`,
			code,
		).Err()
	}

	return
}

func (r *ChatsRepo) AddUserByCode(code string, userID int) (*model.CreateMember, error) {
	member := &model.CreateMember{
		Unit: new(model.Unit),
	}
	err := r.db.QueryRow(`
        WITH inv AS (
			UPDATE invites
			SET aliens = aliens - 1
			WHERE code = $2
			RETURNING chat_id
		), member AS (
			    INSERT INTO chat_members (user_id, chat_id)
				SELECT $1, inv.chat_id
			    FROM inv
				RETURNING id, user_id, chat_id 
		)
		SELECT member.id, chat_id, user_id, domain, name, type
		FROM member join units on user_id = units.id`,
		userID,
		code,
	).Scan(
		&member.ID,
		&member.ChatID,
		&member.Unit.ID,
		&member.Unit.Domain,
		&member.Unit.Name,
		&member.Unit.Type,
	)

	return member, err
}

func (r *ChatsRepo) ChatIsPrivate(chatID int) (private bool) {
	r.db.QueryRow(`
		SELECT private
		FROM chats
		WHERE id = $1`,
		chatID,
	).Scan(&private)

	return
}

func (r *ChatsRepo) BanUserInChat(userID int, chatID int) (*model.DeleteMember, error) {
	member := &model.DeleteMember{}

	err := r.db.QueryRow(`
		WITH m AS (
		    DELETE FROM chat_members
			WHERE user_id = $1 AND chat_id = $2
		    RETURNING id, user_id, chat_id
		),
		ban AS (
			INSERT INTO chat_banlist (user_id, chat_id)
			SELECT user_id, chat_id 
			FROM m
		)
		SELECT m.id FROM m
		`,
		userID,
		chatID,
	).Scan(
		&member.ID,
	)

	return member, err
}

func (r *ChatsRepo) UnbanUserInChat(userID int, chatID int) error {
	err := r.db.QueryRow(`
	    DELETE FROM chat_banlist
		WHERE user_id = $1 AND chat_id = $2
		`,
		userID,
		chatID,
	).Err()

	return err
}

func (r *ChatsRepo) RemoveFromBanlist(userID int, chatID int) (err error) {
	err = r.db.QueryRow(
		`DELETE FROM chat_banlist
		WHERE chat_id = $1 AND user_id = $2`,
		chatID,
		userID,
	).Err()

	return
}

func (r *ChatsRepo) UserIsBannedInChat(userID int, chatID int) (banned bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1
			FROM chat_banlist
			WHERE chat_id = $1 AND user_id = $2
		)`,
		chatID,
		userID,
	).Scan(&banned)

	return
}

func (r *ChatsRepo) Banlist(chatID int) (*model.Users, error) {
	users := &model.Users{
		Users: []*model.User{},
	}
	rows, err := r.db.Query(
		`SELECT units.id, units.domain, units.name, units.type
		FROM units INNER JOIN users 
		ON units.id = users.id 
		INNER JOIN chat_banlist
		ON units.id = chat_banlist.user_id
		WHERE chat_banlist.chat_id = $1`,
		chatID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		m := &model.User{
			Unit: new(model.Unit),
		}
		if err = rows.Scan(&m.Unit.ID, &m.Unit.Domain, &m.Unit.Name, &m.Unit.Type); err != nil {
			return nil, err
		}
		users.Users = append(users.Users, m)
	}

	return users, nil
}

func (r *ChatsRepo) MemberRole(memberID int) (*model.Role, error) {
	role := &model.Role{}

	err := r.db.QueryRow(`
		SELECT coalesce(roles.id, 0), 
		       coalesce(roles.name, ''), 
		       coalesce(roles.color, '')
		FROM (SELECT 1) _
		LEFT JOIN chat_members ON chat_members.id = $1
		LEFT JOIN roles ON chat_members.role_id = roles.id
		`,
		memberID,
	).Scan(
		&role.ID,
		&role.Name,
		&role.Color,
	)
	if err != nil {
		return nil, err
	}
	if role.ID == 0 {
		return nil, nil
	}

	return role, nil
}

func (r *ChatsRepo) CreateRoleInChat(inp *model.CreateRoleInput) (*model.CreateRole, error) {
	role := &model.CreateRole{}
	err := r.db.QueryRow(
		`INSERT INTO roles
		(chat_id, name, color)
		VALUES ($1, $2, $3)
		RETURNING id, chat_id, name, color`,
		inp.ChatID,
		inp.Name,
		inp.Color,
	).Scan(
		&role.ID,
		&role.ChatID,
		&role.Name,
		&role.Color,
	)

	return role, err
}

func (r *ChatsRepo) DeleteRole(roleID int) (*model.DeleteRole, error) {
	role := &model.DeleteRole{}
	err := r.db.QueryRow(`
		DELETE FROM roles
		WHERE id = $1
		RETURNING id`,
		roleID,
	).Scan(
		&role.ID,
	)

	return role, err
}

func (r *ChatsRepo) Roles(chatID int) (*model.Roles, error) {
	roles := &model.Roles{}
	rows, err := r.db.Query(`
		SELECT id, name, color
		FROM roles 
		WHERE chat_id = $1
		`,
		chatID,
	)
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

func (r *ChatsRepo) GetCountChatRoles(chatID int) (count int, err error) {
	err = r.db.QueryRow(`
		SELECT count(*)
		FROM roles 
		WHERE chat_id = $1`,
		chatID,
	).Scan(&count)

	return
}

func (r *ChatsRepo) GiveRole(memberID, roleID int) (*model.UpdateMember, error) {
	member := &model.UpdateMember{}
	err := r.db.QueryRow(
		`UPDATE chat_members
		SET role_id = $2
		WHERE id = $1
		RETURNING id, role_id, char, muted`,
		memberID,
		roleID,
	).Scan(
		&member.ID,
		&member.RoleID,
		&member.Char,
		&member.Muted,
	)

	return member, err
}

func (r *ChatsRepo) TakeRole(memberID int) (*model.UpdateMember, error) {
	member := &model.UpdateMember{}
	err := r.db.QueryRow(`
		WITH x AS (
		    SELECT id
		    FROM chat_members
		    where id = $1 AND role_id is not null 
		)
		UPDATE chat_members
		SET role_id = NULL
		FROM x
		WHERE chat_members.id = x.id
		RETURNING x.id, role_id, char, muted`,
		memberID,
	).Scan(
		&member.ID,
		&member.RoleID,
		&member.Char,
		&member.Muted,
	)

	return member, err
}

func (r *ChatsRepo) TakeChar(memberID int) (*model.UpdateMember, error) {
	member := &model.UpdateMember{}
	err := r.db.QueryRow(`
		WITH x AS (
		    SELECT id
		    FROM chat_members
		    where id = $1 AND char is not null 
		)
		UPDATE chat_members
		SET char = NULL
		FROM x
		WHERE chat_members.id = x.id
		RETURNING x.id, role_id, char, muted`,
		memberID,
	).Scan(
		&member.ID,
		&member.RoleID,
		&member.Char,
		&member.Muted,
	)

	return member, err
}

func (r *ChatsRepo) HasInvite(chatID int, code string) (has bool, err error) {
	err = r.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 
			FROM invites
			WHERE chat_id = $1 AND code = $2
		)`,
		chatID,
		code,
	).Scan(&has)

	return
}

func (r *ChatsRepo) RoleExistsByID(chatID, roleID int) (exists bool, err error) {
	err = r.db.QueryRow(
		`SELECT EXISTS(
		SELECT 1 
		FROM roles
		WHERE chat_id = $1 AND id = $2
		)`,
		chatID,
		roleID,
	).Scan(&exists)

	return
}

func (r *ChatsRepo) InviteInfo(code string) (info *model.InviteInfo, err error) {
	info = &model.InviteInfo{
		Unit: new(model.Unit),
	}
	err = r.db.QueryRow(`
		SELECT units.id, units.domain, units.name, units.type, cm.count_value, chats.private
		FROM units 
		INNER JOIN chats
			ON units.id = chats.id
		JOIN count_members cm ON chats.id = cm.chat_id
		INNER JOIN invites
			ON units.id = invites.chat_id
		WHERE invites.code = $1`,
		code,
	).Scan(
		&info.Unit.ID,
		&info.Unit.Domain,
		&info.Unit.Name,
		&info.Unit.Type,
		&info.CountMembers,
		&info.Private,
	)

	return
}

func (r *ChatsRepo) Owner(chatID int) (*model.User, error) {
	owner := &model.User{
		Unit: new(model.Unit),
	}
	err := r.db.QueryRow(
		`SELECT units.id, units.domain, units.name, units.type
		FROM units 
		INNER JOIN users 
			ON units.id = users.id
		INNER JOIN chats 
			ON units.id = chats.owner_id
		WHERE chats.id = $1`,
		chatID,
	).Scan(
		&owner.Unit.ID,
		&owner.Unit.Domain,
		&owner.Unit.Name,
		&owner.Unit.Type,
	)

	return owner, err
}
func (r *ChatsRepo) UserIsBanned(userID int, chatID int) (banned bool, err error) {
	err = r.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM chat_banlist WHERE user_id = $1 AND chat_id = $2)`,
		userID,
		chatID,
	).Scan(&banned)

	return
}
func (r *ChatsRepo) MemberBy(userID, chatID int) (*model.Member, error) {
	member := &model.Member{
		User: &model.User{
			Unit: new(model.Unit),
		},
		Chat: &model.Chat{
			Unit: &model.Unit{ID: chatID},
		},
	}
	err := r.db.QueryRow(`
		SELECT member.id, units.id, units.domain, units.name, units.type, member.chat_id, member.char, member.joined_at, member.muted
		FROM units 
		JOIN chat_members AS member
		ON units.id = member.user_id
		WHERE member.user_id = $1 AND member.chat_id = $2`,
		userID,
		chatID,
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
	)

	return member, err
}
func (r *ChatsRepo) Member(memberID int) (*model.Member, error) {
	member := &model.Member{
		User: &model.User{
			Unit: new(model.Unit),
		},
		Chat: &model.Chat{
			Unit: new(model.Unit),
		},
	}
	err := r.db.QueryRow(
		`SELECT member.id, units.id, units.domain, units.name, units.type, member.chat_id, member.char, member.joined_at, member.muted
		FROM units INNER JOIN chat_members AS member
		ON units.id = member.user_id
		WHERE member.id = $1`,
		memberID,
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
	)

	return member, err
}

func (r *ChatsRepo) Rooms(chatID int) (*model.Rooms, error) {
	rooms := &model.Rooms{
		Rooms: []*model.Room{},
	}
	rows, err := r.db.Query(`
		SELECT id, parent_id, name, note
		FROM rooms
		WHERE chat_id = $1`,
		chatID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		m := &model.Room{
			Chat: &model.Chat{
				Unit: &model.Unit{ID: chatID},
			},
		}
		if err = rows.Scan(&m.RoomID, &m.ParentID, &m.Name, &m.Note); err != nil {
			return nil, err
		}
		rooms.Rooms = append(rooms.Rooms, m)
	}

	return rooms, nil
}

func (r *RoomsRepo) RolesByArray(roleIDs *[]int) (*model.Roles, error) {
	roles := &model.Roles{}

	rows, err := r.db.Query(`
		SELECT id, name, color
		FROM roles 
		WHERE id = ANY($1)
		`,
		pq.Array(roleIDs),
	)
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

func (r *ChatsRepo) FindMessages(inp *model.FindMessages, params *model.Params, holder *models.AllowHolder) (*model.Messages, error) {

	messages := &model.Messages{
		Messages: []*model.Message{},
	}
	if inp.TextFragment != nil {
		*inp.TextFragment = "%" + *inp.TextFragment + "%"
	}

	var rows, err = r.db.Query(`
		SELECT messages.id, reply_to, user_id, messages.room_id, body, messages.type, created_at
		FROM messages 
		LEFT JOIN allows
			ON messages.room_id = allows.room_id
		WHERE messages.room_id IN (
		    SELECT rooms.id 
		    FROM chats
		    JOIN rooms
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
						group_type = 'ROLE' AND value = $3::VARCHAR
						OR group_type = 'CHAR' AND value = $4::VARCHAR
						OR group_type = 'MEMBER' AND value = $5::VARCHAR
					)
		        OR owner_id = $5::BIGINT 
	        )
		)
		AND (
		    $6::BIGINT IS NULL 
		    OR user_id = $6 
		)
		AND (
		    $7::VARCHAR IS NULL 
		    OR body ILIKE $7
		)
		LIMIT $8
		OFFSET $9
		`,
		inp.ChatID,
		inp.RoomID,
		holder.RoleID,
		holder.Char,
		holder.MemberID,
		inp.UserID,
		inp.TextFragment,
		params.Limit,
		params.Offset,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		m := &model.Message{
			Room: &model.Room{
				Chat: &model.Chat{
					Unit: &model.Unit{ID: inp.ChatID},
				},
			},
		}
		var (
			_replid *int
			_userID *int
		)
		if err = rows.Scan(&m.ID, &_replid, &_userID, &m.Room.RoomID, &m.Body, &m.Type, &m.CreatedAt); err != nil {
			return nil, err
		}
		if _replid != nil {
			m.ReplyTo = &model.Message{
				ID: *_replid,
			}
		}
		if _userID != nil {
			m.User = &model.User{
				Unit: &model.Unit{ID: *_userID},
			}
		}
		messages.Messages = append(messages.Messages, m)
	}

	return messages, nil
}

func (r *ChatsRepo) ChatIDByMemberID(memberID int) (chatID int, err error) {
	err = r.db.QueryRow(`
		select coalesce(
		    (SELECT chat_id
			FROM chat_members
			WHERE id = $1),
		    0
		    ) as chat_id
		`,
		memberID,
	).Scan(&chatID)

	return
}

func (r *ChatsRepo) MemberIsMuted(memberID int) (muted bool, err error) {
	err = r.db.QueryRow(`
		SELECT EXISTS(
		    SELECT 1
		    FROM chat_members
		    WHERE  id = $1 AND muted = TRUE
		)
		`,
		memberID,
	).Scan(&muted)

	return
}

func (r *ChatsRepo) FindMembers(inp *model.FindMembers) (*model.Members, error) {
	members := &model.Members{
		Members: []*model.Member{},
	}
	rows, err := r.db.Query(`
			SELECT members.id, chat_id, units.id, units.domain, units.name, units.type, role_id, char, joined_at, muted
			FROM chat_members as members
			JOIN units ON members.user_id = units.id
			WHERE (
			    $1::BIGINT IS NULL 
			    OR units.id= $1 
			)
			AND (
			    $2::BIGINT IS NULL 
			    OR chat_id = $2 
			)
			AND (
			    $3::BIGINT IS NULL 
			    OR members.id = $3 
			)
			AND (
			    $4::char_type IS NULL 
			    OR char = $4 
			)
			AND (
			    $5::BIGINT IS NULL 
			    OR role_id = $5 
			)
			AND (
			    $6::BOOLEAN IS NULL 
			    OR muted = $6 
			)
			`,
		inp.UserID,
		inp.ChatID,
		inp.MemberID,
		inp.Char,
		inp.RoleID,
		inp.Muted,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		m := &model.Member{
			User: &model.User{
				Unit: new(model.Unit),
			},
			Chat: &model.Chat{
				Unit: new(model.Unit),
			},
		}
		var (
			_roleid *int
		)
		if err = rows.Scan(&m.ID, &m.Chat.Unit.ID, &m.User.Unit.ID, &m.User.Unit.Domain, &m.User.Unit.Name, &m.User.Unit.Type, &_roleid, &m.Char, &m.JoinedAt, &m.Muted); err != nil {
			return nil, err
		}
		if _roleid != nil {
			m.Role = &model.Role{
				ID: *_roleid,
			}
		}
		members.Members = append(members.Members, m)
	}

	return members, nil

}

// DemoMembers
//
// SelectType: 0 is filter by users, 1 - by members and chatid is not count.
//
// Если ids[0] и ids[1] равно, то вернутся array [&DemoMember, nil].
func (r *ChatsRepo) DemoMembers(chatID, selectType int, ids ...int) [2]*models.DemoMember { // todo selectType to rules.SelectType
	var (
		demoMembers [2]*models.DemoMember
	)
	if selectType != 0 && selectType != 1 {
		selectType = 1
	}
	if len(ids) > 2 {
		return demoMembers
	}

	rows, err := r.db.Query(`
		SELECT user_id, chat_members.id, owner_id = user_id as is_owner, char, muted
		FROM chat_members 
		JOIN chats ON chats.id = chat_members.chat_id
		WHERE $2 = 0 AND chats.id = $1 AND user_id = ANY ($3)
		    OR $2 = 1 AND chat_members.id = ANY ($3)
		`,
		chatID,
		selectType,
		pq.Array(ids),
	)
	if err != nil {
		return demoMembers
	}
	defer rows.Close()
	sort := func() [2]*models.DemoMember {
		if demoMembers[1] == nil {
			return demoMembers
		}
		for i, member := range demoMembers {
			if selectType == 0 && member.UserID == ids[0] ||
				selectType == 1 && member.MemberID == ids[0] {
				demoMembers[0], demoMembers[1] = demoMembers[i], demoMembers[1-i]
			}
		}
		return demoMembers
	}
	i := 0
	for rows.Next() {
		m := &models.DemoMember{}
		if err = rows.Scan(&m.UserID, &m.MemberID, &m.IsOwner, &m.Char, &m.Muted); err != nil {
			return sort()
		}
		demoMembers[i] = m
		i += 1
	}
	return sort()
}

func (r *ChatsRepo) UpdateRole(roleID int, inp *model.UpdateRoleInput) (*model.UpdateRole, error) {
	role := &model.UpdateRole{}
	err := r.db.QueryRow(`
		UPDATE roles
		SET name = COALESCE($2::VARCHAR, name), color = COALESCE($3::VARCHAR, color)
		WHERE id = $1
		RETURNING id, name, color`,
		roleID,
		inp.Name,
		inp.Color,
	).Scan(
		&role.ID,
		&role.Name,
		&role.Color,
	)

	return role, err
}
func (r *ChatsRepo) UpdateChat(chatID int, inp *model.UpdateChatInput) (*model.UpdateChat, error) {
	chat := &model.UpdateChat{}
	err := r.db.QueryRow(`
		with c as (
			UPDATE units
			SET 
			    name = COALESCE($2::VARCHAR, name), 
			    domain = COALESCE($3::VARCHAR, domain)
			WHERE id = $1
		    RETURNING domain, name
		)
		UPDATE chats
		SET private = COALESCE($4::BOOLEAN, private)
		FROM c
		WHERE id = $1
		RETURNING id, c.domain, c.name, private
		`,
		chatID,
		inp.Name,
		inp.Domain,
		inp.Private,
	).Scan(
		&chat.ID,
		&chat.Domain,
		&chat.Name,
		&chat.Private,
	)

	return chat, err
}

func (r *ChatsRepo) DefMember(memberID int) (*models.DefMember, error) {
	defMember := new(models.DefMember)
	err := r.db.QueryRow(`
		SELECT coalesce(user_id, 0) userID,
	       	   coalesce(chat_id, 0) chatID
		FROM (SELECT 1) x
		LEFT JOIN chat_members m ON id = $1
		`,
		memberID,
	).Scan(
		&defMember.UserID,
		&defMember.ChatID,
	)
	if err != nil {
		return nil, err
	}
	if defMember.UserID == 0 {
		return nil, nil
	}

	return defMember, nil
}

func (r *ChatsRepo) FindChats(inp *model.FindChats, params *model.Params) (*model.Chats, error) {
	chats := &model.Chats{
		Chats: []*model.Chat{},
	}
	if inp.NameFragment != nil {
		*inp.NameFragment = "%" + *inp.NameFragment + "%"
	}
	rows, err := r.db.Query(`
		SELECT units.id, units.domain, units.name, units.type, cm.count_value, chats.private
		FROM units 
		JOIN chats ON units.id = chats.id 
		JOIN count_members cm on chats.id = cm.chat_id
		WHERE (
		    $1::BIGINT IS NULL 
			OR units.id = $1
		)  
		AND (
		    $2::VARCHAR IS NULL 
			OR domain = $2
		)
		AND (
		    $3::VARCHAR IS NULL 
		    OR name ILIKE $3
		) 
		LIMIT $4
		OFFSET $5
		`,
		inp.ID,
		inp.Domain,
		inp.NameFragment,
		params.Limit,
		params.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		m := &model.Chat{
			Unit: new(model.Unit),
		}
		if err = rows.Scan(&m.Unit.ID, &m.Unit.Domain, &m.Unit.Name, &m.Unit.Type, &m.CountMembers, &m.Private); err != nil {
			return nil, err
		}
		chats.Chats = append(chats.Chats, m)
	}
	return chats, nil
}

func (r *ChatsRepo) UpdateMember(memberID int, inp *model.UpdateMemberInput) (*model.UpdateMember, error) {
	member := &model.UpdateMember{}
	err := r.db.QueryRow(`
		UPDATE chat_members
		SET 
		    role_id = COALESCE($2::BIGINT, role_id), 
		    char = COALESCE($3::char_type, char), 
		    muted = COALESCE($4::BOOLEAN, muted)
		WHERE id = $1
		RETURNING id, role_id, char, muted`,
		memberID,
		inp.RoleID,
		inp.Char,
		inp.Muted,
	).Scan(
		&member.ID,
		&member.RoleID,
		&member.Char,
		&member.Muted,
	)

	return member, err
}

func (r *ChatsRepo) ValidAllows(chatID int, allows *model.AllowsInput) (valid bool, err error) {
	err = r.db.QueryRow(`
		SELECT NOT exists(
		    SELECT 1
		    FROM unnest($2::findallow[]) elem (act,gr,val)
		    LEFT JOIN chat_members cm ON
		        elem.gr = 'MEMBER' AND
		        elem.val::BIGINT = cm.id AND
		        cm.chat_id = $1
		    LEFT JOIN roles r ON
		        elem.gr = 'ROLE' AND
		        elem.val::BIGINT = r.id AND
		        r.chat_id = $1
		    WHERE cm.id IS NULL AND r.id IS NULL AND elem.gr::VARCHAR <> 'CHAR'
		)`,
		chatID,
		pq.Array(allows.Allows),
	).Scan(&valid)

	return valid, err
}

// if noAccessTo = 0 then acces allow to all chats
func (r *ChatsRepo) UserHasAccessToChats(userID int, chats *[]int) (members []*models.SubUser, noAccessTo int, err error) {
	rows, err := r.db.Query(`
	    SELECT cm.id, elem
	    FROM unnest($2::BIGINT[]) elem
	    LEFT JOIN chats c ON c.id = elem
	    LEFT JOIN chat_members cm ON cm.chat_id = c.id AND cm.user_id = $1`,
		userID,
		pq.Array(*chats),
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		member := &models.SubUser{}
		if err = rows.Scan(&member.MemberID, &member.ChatID); err != nil {
			return
		}
		if member.MemberID == nil {
			return nil, *member.ChatID, nil
		}
		members = append(members, member)
	}
	return members, 0, nil
}
