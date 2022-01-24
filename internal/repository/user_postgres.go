package repository

import (
	"database/sql"
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/models"
)

type UsersRepo struct {
	db *sql.DB
}

func NewUsersRepo(db *sql.DB) *UsersRepo {
	return &UsersRepo{
		db: db,
	}
}

func (r *UsersRepo) CreateUser(userModel *models.RegisterData) (err error) {
	err = r.db.QueryRow(
		`WITH u AS (
			INSERT INTO units (domain, name, type) 
			VALUES ($1, $2, 'USER') 
			RETURNING id
			) 
		INSERT INTO users (id, hashed_password, email) 
		SELECT u.id, $3, $4 
		FROM u 
		RETURNING id`,
		userModel.Domain,
		userModel.Name,
		userModel.HashPassword,
		userModel.Email,
	).Err()
	if err != nil {
		println("CreateUser:", err.Error()) // debug
	}
	return
}

func (r *UsersRepo) User(userID int) (*model.User, error) {
	user := &model.User{
		Unit: &model.Unit{},
	}
	err := r.db.QueryRow(`
		SELECT id, domain, name, type
		FROM units
		WHERE id = $1`,
		userID,
	).Scan(
		&user.Unit.ID,
		&user.Unit.Domain,
		&user.Unit.Name,
		&user.Unit.Type,
	)
	if err != nil {
		println("User:", err.Error()) // debug
	}
	return user, err
}

func (r *UsersRepo) UserExistsByRequisites(inp *models.LoginRequisites) (exists bool) {
	err := r.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM users
			WHERE email = $1 AND hashed_password = $2
		)`,
		inp.Email,
		inp.HashedPasswd,
	).Scan(&exists)
	if err != nil {
		println("UserExistsByRequisites:", err.Error()) // debug
	}

	return

}

func (r *UsersRepo) GetUserIdByRequisites(inp *models.LoginRequisites) (id int, err error) {
	err = r.db.QueryRow(`
		SELECT units.id
		FROM units INNER JOIN users 
		ON units.id = users.id 
		WHERE users.email = $1 AND users.hashed_password = $2`,
		inp.Email,
		inp.HashedPasswd,
	).Scan(&id)
	if err != nil {
		println("GetUserIdByInput:", err.Error()) // debug
	}
	return
}

func (r *UsersRepo) GetUserByDomain(domain string) (user models.UserInfo, err error) {
	err = r.db.QueryRow(
		`SELECT units.id,units.domain,units.name
		FROM units INNER JOIN users 
		ON units.id = users.id 
		WHERE units.domain = $1`,
		domain,
	).Scan(
		&user.ID,
		&user.Domain,
		&user.Name,
	)
	if err != nil {
		return // user, err
	}
	return
}

func (r *UsersRepo) GetUserByID(id int) (user models.UserInfo, err error) {
	err = r.db.QueryRow(
		`SELECT units.id,units.domain,units.name
		FROM units INNER JOIN users 
		ON units.id = users.id 
		WHERE units.id = $1`,
		id,
	).Scan(
		&user.ID,
		&user.Domain,
		&user.Name,
	)
	if err != nil {
		return // user, err
	}
	return
}

func (r *UsersRepo) GetCountUserOwnedChats(userId int) (count int, err error) {
	err = r.db.QueryRow(
		`SELECT count(*)
		FROM chats 
		WHERE owner_id = $1`,
		userId,
	).Scan(&count)
	return
}

func (r *UsersRepo) UserExistsByID(userId int) (exists bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1
			FROM users
			WHERE id = $1
		)`,
		userId,
	).Scan(&exists)

	return
}

func (r *UsersRepo) UserExistsByDomain(userDomain string) (exists bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1
			FROM units
			INNER JOIN users
			ON users.id = units.id
			WHERE units.domain = $1
		)`,
		userDomain,
	).Scan(&exists)

	return
}

func (r *UsersRepo) Me(usersId int) (*model.Me, error) {
	me := &model.Me{
		User: &model.User{
			Unit: &model.Unit{},
		},
		Data: &model.UserData{},
	}
	err := r.db.QueryRow(
		`SELECT units.id, units.domain, units.name, units.type, users.email
		FROM units INNER JOIN users
		ON units.id = users.id
		WHERE units.id = $1`,
		usersId,
	).Scan(
		&me.User.Unit.ID,
		&me.User.Unit.Domain,
		&me.User.Unit.Name,
		&me.User.Unit.Type,
		&me.Data.Email,
	)
	if err != nil {
		println("Me:", err.Error()) // debug
	}

	return me, err
}

func (r *UsersRepo) OwnedChats(userId int) (*model.Chats, error) {
	chats := &model.Chats{
		Chats: []*model.Chat{},
	}
	rows, err := r.db.Query(
		`SELECT units.id, units.domain, units.name, units.type, chats.private
		FROM units INNER JOIN chats 
		ON units.id = chats.id 
		WHERE chats.owner_id = $1`,
		userId,
	)
	if err != nil {
		println("OwnedChats:", err.Error()) // debug
		return chats, err
	}
	defer rows.Close()
	for rows.Next() {
		m := &model.Chat{}
		if err = rows.Scan(&m.Unit.ID, &m.Unit.Domain, &m.Unit.Name, &m.Unit.Type, &m.Private); err != nil {
			println("OwnedChats:", err.Error()) // debug
			return chats, err
		}

		chats.Chats = append(chats.Chats, m)
	}
	if !rows.NextResultSet() {
		println("OwnedChats:", err.Error()) // debug
		return chats, err
	}
	return chats, nil
}

func (r *UsersRepo) Chats(userId int) (*model.Chats, error) {
	chats := &model.Chats{
		Chats: []*model.Chat{},
	}
	rows, err := r.db.Query(`
		SELECT units.id, units.domain, units.name, units.type, chats.private
		FROM units 
		INNER JOIN chats 
			ON units.id = chats.id 
		INNER JOIN chat_members
			ON units.id = chat_members.chat_id
		WHERE chat_members.user_id = $1`,
		userId,
	)
	if err != nil {
		println("Chats:", err.Error()) // debug
		return chats, err
	}
	defer rows.Close()
	for rows.Next() {
		m := &model.Chat{
			Unit: &model.Unit{},
		}
		if err = rows.Scan(&m.Unit.ID, &m.Unit.Domain, &m.Unit.Name, &m.Unit.Type, &m.Private); err != nil {
			println("Chats:", err.Error()) // debug
			return chats, err
		}

		chats.Chats = append(chats.Chats, m)
	}

	return chats, nil
}
func (r *UsersRepo) ChatsID(userId int) ([]int, error) {
	rows, err := r.db.Query(
		`SELECT chat_id
		FROM chat_members
		WHERE user_id = $1`,
		userId,
	)
	if err != nil {
		println("ChatsID:", err.Error()) // debug
		return []int{}, err
	}
	defer rows.Close()

	chats, err := completeIntArray(rows)
	if err != nil {
		println("ChatsID:", err.Error()) // debug
	}

	return chats, nil
}
func (r *UsersRepo) FindUsers(inp *model.FindUsers) *model.Users {
	users := &model.Users{}
	if inp.NameFragment != nil {
		*inp.NameFragment = "%" + *inp.NameFragment + "%"
	}
	rows, err := r.db.Query(`
		SELECT units.id, units.domain, units.name, units.type
		FROM units JOIN users ON units.id = users.id 
		WHERE	($1 IS NULL OR units.id = $1)
			AND ($2 IS NULL OR units.domain = $2)
			AND ($3 IS NULL OR units.name ILIKE $3)
		`,
		inp.ID,
		inp.Domain,
		inp.NameFragment,
	)
	if err != nil {
		println("FindUsers:", err.Error()) // debug
		return users
	}
	defer rows.Close()
	for rows.Next() {
		m := &model.User{
			Unit: &model.Unit{},
		}
		if err = rows.Scan(&m.Unit.ID, &m.Unit.Domain, &m.Unit.Name, &m.Unit.Type); err != nil {
			return users
		}
		users.Users = append(users.Users, m)
	}

	return users
}

func (r UsersRepo) UpdateMe(userId int, inp *model.UpdateMeDataInput) (*model.UpdateUser, error) {
	unit := &model.UpdateUser{}
	err := r.db.QueryRow(`
		WITH u AS (
			UPDATE units
			SET 
			    name = COALESCE($2::VARCHAR, name), 
			    domain = COALESCE($3::VARCHAR, domain)
			WHERE id = $1
		    RETURNING domain, name
		)
		UPDATE users
		SET 
		    hashed_password = COALESCE($4::VARCHAR, hashed_password),
		    email = COALESCE($5::VARCHAR, email)
		FROM u
		WHERE id = $1
		RETURNING id, u.domain, u.name
		`,
		userId,
		inp.Name,
		inp.Domain,
		inp.Password,
		inp.Email,
	).Scan(
		&unit.ID,
		&unit.Domain,
		&unit.Name,
	)
	if err != nil {
		println("UpdateMe:", err.Error()) // debug
	}

	return unit, err
}

func (r UsersRepo) GetRegistrationSession(email, code string) (*models.RegisterData, error) {
	regi := &models.RegisterData{}
	err := r.db.QueryRow(`
	    SELECT domain, name, email, hashed_password
		FROM registration_session
	    WHERE email = $1 AND verify_code = $2
		`,
		email,
		code,
	).Scan(
		&regi.Domain,
		&regi.Name,
		&regi.Email,
		&regi.HashPassword,
	)
	if err != nil {
		println("GetRegistrationSession:", err.Error()) // debug
		return nil, err
	}
	return regi, nil
}

func (r UsersRepo) DeleteRegistrationSession(email string) {
	err := r.db.QueryRow(`
		DELETE FROM registration_session
	    WHERE email = $1
		`,
		email,
	).Err()
	if err != nil {
		println("DeleteRegistrationSession:", err.Error()) // debug
	}
	return
}

func (r *UsersRepo) CreateRegistrationSession(userModel *models.RegisterData, expAt int64) (verifyCode string, err error) {
	err = r.db.QueryRow(`
		INSERT INTO registration_session (domain, name, email, hashed_password, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING verify_code`,
		userModel.Domain,
		userModel.Name,
		userModel.Email,
		userModel.HashPassword,
		expAt,
	).Scan(&verifyCode)
	if err != nil {
		println("CreateRegistrationSession:", err.Error()) // debug
	}
	return
}

func (r *UsersRepo) EmailIsFree(email string) (free bool) {
	err := r.db.QueryRow(`
		SELECT 
		EXISTS (
			SELECT 1 
			FROM users
			WHERE email = $1
		) 
		OR
		EXISTS (
		    SELECT 1
		    FROM registration_session
		    WHERE email = $1
		)`,
		email,
	).Scan(&free)
	if err != nil {
		println("EmailIsFree:", err.Error()) // debug
	}
	return !free
}
