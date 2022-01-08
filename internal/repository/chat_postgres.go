package repository

import (
	"database/sql"
	"fmt"
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/rules"
	"github.com/saime-0/http-cute-chat/internal/tlog"
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
	if err != nil {
		println("CreateChat:", err.Error()) // debug
	}
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
	tl := tlog.Start("ChatsRepo > Chat [cid:" + strconv.Itoa(chatId) + "]")
	defer tl.Fine()
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
func (r *ChatsRepo) ChatIDByRoleID(roleId int) (chatId int, err error) {
	err = r.db.QueryRow(`
		SELECT chats.id
		FROM chats JOIN roles ON chats.id = roles.chat_id
		WHERE roles.id = $1 
		`,
		roleId,
	).Scan(&chatId)
	if err != nil {
		println("ChatIDByRoleID:", err.Error()) // debug
	}
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
func (r *ChatsRepo) Members(chatId int) (*model.Members, error) {
	tl := tlog.Start("ChatsRepo > Members [cid:" + strconv.Itoa(chatId) + "]")
	defer tl.Fine()
	members := &model.Members{}
	rows, err := r.db.Query(
		`SELECT member.id, units.id, units.domain, units.name, units.type, member.chat_id, member.char, member.joined_at, member.muted, member.frozen
		FROM units INNER JOIN chat_members AS member
		ON units.id = member.user_id
		WHERE member.chat_id = $1`,
		chatId,
	)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
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
			&m.ID,
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
func (r *ChatsRepo) UserIsChatMember(userId int, chatId int) (isMember bool) {
	tl := tlog.Start("ChatsRepo > UserIsChatMember [uid:" + strconv.Itoa(userId) + ",cid:" + strconv.Itoa(chatId) + "]")
	defer tl.Fine()
	err := r.db.QueryRow(`
		SELECT EXISTS(
	        SELECT 1
			FROM chat_members 
			WHERE user_id = $1 AND chat_id = $2
	    )
		`,
		userId,
		chatId,
	).Scan(&isMember)
	if err != nil {
		println("UserIsChatMember:", err.Error()) // debug
		return
	}
	return
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
func (r *ChatsRepo) Invites(chatId int) (*model.Invites, error) {
	invites := &model.Invites{
		Invites: []*model.Invite{},
	}
	rows, err := r.db.Query(
		`SELECT code, aliens, expires_at
		FROM invites
		WHERE chat_id = $1`,
		chatId,
	)
	if err != nil {
		println("Invites:", err.Error()) // debug
		return invites, err
	}
	defer rows.Close()
	for rows.Next() {
		m := &model.Invite{}
		if err = rows.Scan(&m.Code, &m.Aliens, &m.ExpiresAt); err != nil {
			println("Invites:", err.Error()) // debug
			return invites, err
		}
		invites.Invites = append(invites.Invites, m)
	}

	return invites, nil
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
	if err != nil {
		println("DeleteInvite:", err.Error()) // debug
	}
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
	if err != nil {
		println("CreateInvite:", err.Error()) // debug
	}
	return invite, err
}

// InviteIsRelevant
//  alt: InviteIsExists
func (r *ChatsRepo) InviteIsRelevant(code string) (relevant bool) {
	err := r.db.QueryRow(`
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
		println("InviteIsRelevant:", err.Error()) // debug
	}
	if !relevant {
		r.db.Exec(`
			DELETE FROM invites
			WHERE code = $1`,
			code,
		)
	}

	return
}

func (r *ChatsRepo) AddUserByCode(code string, userId int) (err error) {
	err = r.db.QueryRow(`
		WITH l AS (
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
	).Err()
	if err != nil {
		println("AddUserByCode:", err.Error()) // debug
	}
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

func (r *ChatsRepo) Banlist(chatId int) (*model.Users, error) {
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
		chatId,
	)
	if err != nil {
		println("Banlist:", err.Error()) // debug
		return users, err
	}
	defer rows.Close()
	for rows.Next() {
		m := &model.User{}
		if err = rows.Scan(&m.Unit.ID, &m.Unit.Domain, &m.Unit.Name, &m.Unit.Type); err != nil {
			println("Banlist:", err.Error()) // debug
			return users, err
		}
		users.Users = append(users.Users, m)
	}
	if !rows.NextResultSet() {
		println("Banlist:", err.Error()) // debug
		return users, err
	}
	return users, nil
}

func (r *ChatsRepo) MemberRole(memberId int) *model.Role {
	tl := tlog.Start("ChatsRepo > MemberRole [mid:" + strconv.Itoa(memberId) + "]")
	defer tl.Fine()
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
		println("MemberRole:", err.Error()) // debug
	}
	return nil
}

func (r *ChatsRepo) CreateRoleInChat(inp *model.CreateRoleInput) (err error) {
	err = r.db.QueryRow(
		`INSERT INTO roles
		(chat_id, name, color)
		VALUES ($1, $2, $3)
		RETURNING id`,
		inp.ChatID,
		inp.Name,
		inp.Color,
	).Err()
	if err != nil {
		println("CreateRoleInChat:", err.Error()) // debug
	}
	return
}

func (r *ChatsRepo) Roles(chatId int) (*model.Roles, error) {
	roles := &model.Roles{}
	rows, err := r.db.Query(
		`SELECT id, name, color
		FROM roles 
		WHERE chat_id = $1`,
		chatId,
	)
	if err != nil {
		println("Roles:", err.Error()) // debug
		return roles, err
	}
	defer rows.Close()
	for rows.Next() {
		m := &model.Role{}
		if err = rows.Scan(&m.ID, &m.Name, &m.Color); err != nil {
			return roles, err
		}
		roles.Roles = append(roles.Roles, m)
	}

	return roles, nil
}

func (r *ChatsRepo) GetCountChatRoles(chatId int) (count int, err error) {
	err = r.db.QueryRow(
		`SELECT count(*)
		FROM roles 
		WHERE chat_id = $1`,
		chatId,
	).Scan(&count)
	if err != nil {
		println("GetCountChatRoles:", err.Error()) // debug
	}
	return
}

func (r *ChatsRepo) GiveRole(memberId, roleId int) (*model.UpdateMember, error) {
	member := &model.UpdateMember{}
	err := r.db.QueryRow(
		`UPDATE chat_members
		SET role_id = $2
		WHERE id = $1
		RETURNING id, role_id, char, muted, frozen`,
		memberId,
		roleId,
	).Scan(
		&member.ID,
		&member.RoleID,
		&member.Char,
		&member.Muted,
		&member.Frozen,
	)
	if err != nil {
		println("GiveRole:", err.Error()) // debug
	}
	return member, err
}

func (r *ChatsRepo) DeleteRole(roleId int) (err error) {
	err = r.db.QueryRow(
		`DELETE FROM roles
		WHERE id = $1`,
		roleId,
	).Err()
	if err != nil {
		println("DeleteRole:", err.Error()) // debug
	}
	return
}

func (r *ChatsRepo) TakeRole(memberId int) (*model.UpdateMember, error) {
	member := &model.UpdateMember{}
	err := r.db.QueryRow(
		`UPDATE chat_members
		SET role_id = NULL
		WHERE id = $1
		RETURNING id, role_id, char, muted, frozen`,
		memberId,
	).Scan(
		&member.ID,
		&member.RoleID,
		&member.Char,
		&member.Muted,
		&member.Frozen,
	)
	if err != nil {
		println("TakeRole:", err.Error()) // debug
	}
	return member, err
}

func (r *ChatsRepo) HasInvite(chatId int, code string) (has bool) {
	err := r.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 
			FROM invites
			WHERE chat_id = $1 AND code = $2
		)`,
		chatId,
		code,
	).Scan(&has)
	if err != nil {
		println("HasInvite:", err.Error()) // debug
	}
	return
}

func (r *ChatsRepo) RoleExistsByID(chatId, roleId int) (exists bool) {
	err := r.db.QueryRow(
		`SELECT EXISTS(
		SELECT 1 
		FROM roles
		WHERE chat_id = $1 AND id = $2
		)`,
		chatId,
		roleId,
	).Scan(&exists)
	if err != nil {
		println("RoleExistsByID:", err.Error()) // debug
	}
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
	if err != nil {
		println("InviteInfo:", err.Error()) // debug
	}
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
	if err != nil {
		println("Owner:", err.Error()) // debug
	}
	return owner, err
}
func (r *ChatsRepo) UserIsBanned(userId int, chatId int) (banned bool) {
	err := r.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM chat_banlist WHERE user_id = $1 AND chat_id = $2)`,
		userId,
		chatId,
	).Scan(&banned)
	if err != nil {
		println("UserIsBanned:", err.Error()) // debug
	}
	return
}
func (r *ChatsRepo) MemberBy(userId, chatId int) (*model.Member, error) {
	member := &model.Member{
		User: &model.User{
			Unit: &model.Unit{},
		},
		Chat: &model.Chat{
			Unit: &model.Unit{ID: chatId},
		},
	}
	err := r.db.QueryRow(`
		SELECT member.id, units.id, units.domain, units.name, units.type, member.chat_id, member.char, member.joined_at, member.muted, member.frozen
		FROM units 
		JOIN chat_members AS member
		ON units.id = member.user_id
		WHERE member.user_id = $1 AND member.chat_id = $2`,
		userId,
		chatId,
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
	if err != nil {
		println("MemberBy:", err.Error()) // debug
	}
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
		ON units.id = member.user_id
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
	if err != nil {
		println("Member:", err.Error()) // debug
	}
	return member, err
}

func (r *ChatsRepo) Rooms(chatId int) (*model.Rooms, error) {
	rooms := &model.Rooms{
		Rooms: []*model.Room{},
	}
	rows, err := r.db.Query(
		`SELECT id, parent_id, name, note
		FROM rooms
		WHERE chat_id = $1`,
		chatId,
	)
	if err != nil {
		println("Rooms:", err.Error()) // debug
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		m := &model.Room{
			Chat: &model.Chat{
				Unit: &model.Unit{ID: chatId},
			},
		}
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
		println("RolesByArray:", err.Error()) // debug
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

func (r *ChatsRepo) FindMessages(inp *model.FindMessages, params *model.Params, holder *models.AllowHolder) *model.Messages {
	tl := tlog.Start("ChatsRepo > FindMessages [cid:" + strconv.Itoa(inp.ChatID) + "]")
	defer tl.Fine()
	messages := &model.Messages{
		Messages: []*model.Message{},
	}
	if inp.TextFragment != nil {
		*inp.TextFragment = "%" + *inp.TextFragment + "%"
	}

	// language=PostgreSQL
	var rows, err = r.db.Query(`
			SELECT messages.id, reply_to, author, messages.room_id, body, messages.type, created_at
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
							group_type = 'ROLES' AND value = $3::VARCHAR 
							OR group_type = 'CHARS' AND value = $4::VARCHAR  
							OR group_type = 'USERS' AND value = $5::VARCHAR
						)
			        OR owner_id = $5::BIGINT 
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
			LIMIT $8
			OFFSET $9
			`,
		inp.ChatID,
		inp.RoomID,
		holder.RoleID,
		holder.Char,
		holder.UserID,
		inp.AuthorID,
		inp.TextFragment,
		params.Limit,
		params.Offset,
	)
	if err != nil {
		println("FindMessages:", err.Error()) // debug
		return messages
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

func (r *ChatsRepo) ChatIDByMemberID(memberId int) (chatId int, err error) {
	err = r.db.QueryRow(`
		SELECT chat_id
		FROM chat_members
		WHERE id = $1
		`,
		memberId,
	).Scan(&chatId)
	if err != nil {
		println("ChatIDByMemberID:", err.Error()) // debug
	}
	return
}

func (r *ChatsRepo) FindMessage(messageId int) *model.Message {
	panic("Not implemented")
}

func (r *ChatsRepo) MemberIsMuted(memberId int) (muted bool) {
	err := r.db.QueryRow(`
		SELECT EXISTS(
		    SELECT 1
		    FROM chat_members
		    WHERE  id = $1 AND muted = TRUE
		)
		`,
		memberId,
	).Scan(&muted)
	if err != nil {
		println("MemberIsMuted:", err.Error()) // debug
	}
	return
}

func (r *ChatsRepo) FindMembers(inp *model.FindMembers) *model.Members {
	members := &model.Members{
		Members: []*model.Member{},
	}
	rows, err := r.db.Query(`
			SELECT members.id, chat_id, units.id, units.domain, units.name, units.type, role_id, char, joined_at, frozen, muted
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
			AND (
			    $7::BOOLEAN IS NULL 
			    OR frozen = $7 
			)
			`,
		inp.UserID,
		inp.ChatID,
		inp.MemberID,
		inp.Char,
		inp.RoleID,
		inp.Muted,
		inp.Frozen,
	)
	if err != nil {
		println("FindMembers:", err.Error()) // debug
		return members
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
		var (
			_roleid *int
		)
		if err = rows.Scan(&m.ID, &m.Chat.Unit.ID, &m.User.Unit.ID, &m.User.Unit.Domain, &m.User.Unit.Name, &m.User.Unit.Type, &_roleid, &m.Char, &m.JoinedAt, &m.Muted, &m.Frozen); err != nil {
			println("rows.scan:", err.Error()) // debug
			return members
		}
		if _roleid != nil {
			m.Role = &model.Role{
				ID: *_roleid,
			}
		}
		members.Members = append(members.Members, m)
	}

	return members

}

// DemoMembers selectType: 0 is filter by users, 1 - by members and chatid is not count
func (r *ChatsRepo) DemoMembers(chatId, selectType int, ids ...int) [2]*models.DemoMember { // todo selectType to rules.SelectType
	var (
		sqlArr      = kit.IntSQLArray(ids)
		demoMembers [2]*models.DemoMember
	)
	if selectType != 0 && selectType != 1 {
		selectType = 1
	}

	//language=PostgreSQL
	rows, err := r.db.Query(`
		SELECT user_id, chat_members.id, owner_id = user_id as is_owner, char, muted
		FROM chat_members 
		JOIN chats ON chats.id = chat_members.chat_id
		WHERE $2 = 0 AND chats.id = $1 AND user_id IN `+sqlArr+`
		    OR $2 = 1 AND chat_members.id IN `+sqlArr+`
		`,
		chatId,
		selectType,
	)
	if err != nil {
		println("DemoMembers:", err.Error()) // debug
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
			println("rows.scan:", err.Error()) // debug
			return sort()
		}
		demoMembers[i] = m
		i += 1
	}
	return sort()
}

func (r *ChatsRepo) UpdateRole(roleId int, inp *model.UpdateRoleInput) (*model.UpdateRole, error) {
	role := &model.UpdateRole{}
	err := r.db.QueryRow(`
		UPDATE roles
		SET name = COALESCE($2::VARCHAR, name), color = COALESCE($3::VARCHAR, color)
		WHERE id = $1
		RETURNING id, name, color`,
		roleId,
		inp.Name,
		inp.Color,
	).Scan(
		&role.ID,
		&role.Name,
		&role.Color,
	)
	if err != nil {
		println("UpdateRole:", err.Error()) // debug
	}

	return role, err
}
func (r *ChatsRepo) UpdateChat(chatId int, inp *model.UpdateChatInput) (*model.UpdateChat, error) {
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
		chatId,
		inp.Name,
		inp.Domain,
		inp.Private,
	).Scan(
		&chat.ID,
		&chat.Domain,
		&chat.Name,
		&chat.Private,
	)
	if err != nil {
		println("UpdateChat:", err.Error()) // debug
	}
	return chat, err
}

func (r *ChatsRepo) DefMember(memberId int) (defMember models.DefMember, err error) {
	err = r.db.QueryRow(`
		SELECT user_id, chat_id
		FROM chat_members
		WHERE id = $1
		`,
		memberId,
	).Scan(
		&defMember.UserID,
		&defMember.ChatID,
	)
	if err != nil {
		println("DefMember:", err.Error()) // debug
	}
	return
}

func (r *ChatsRepo) FindChats(inp *model.FindChats, params *model.Params) (*model.Chats, error) {
	chats := &model.Chats{
		Chats: []*model.Chat{},
	}
	if inp.NameFragment != nil {
		*inp.NameFragment = "%" + *inp.NameFragment + "%"
	}
	rows, err := r.db.Query(`
		SELECT units.id, units.domain, units.name, units.type, chats.private
		FROM units JOIN chats 
		ON units.id = chats.id 
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
		println("FindChats:", err.Error()) // debug
		return chats, err
	}
	defer rows.Close()
	for rows.Next() {
		m := &model.Chat{
			Unit: &model.Unit{},
		}
		if err = rows.Scan(&m.Unit.ID, &m.Unit.Domain, &m.Unit.Name, &m.Unit.Type, &m.Private); err != nil {
			println("rows:Scan:", err.Error()) // debug
			return chats, err
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
		    muted = COALESCE($4::BOOLEAN, muted), 
		    frozen = COALESCE($5::BOOLEAN, frozen)
		WHERE id = $1
		RETURNING id, role_id, char, muted, frozen`,
		memberID,
		inp.RoleID,
		inp.Char,
		inp.Muted,
		inp.Frozen,
	).Scan(
		&member.ID,
		&member.RoleID,
		&member.Char,
		&member.Muted,
		&member.Frozen,
	)
	if err != nil {
		println("UpdateMember:", err.Error()) // debug
	}
	return member, err
}
