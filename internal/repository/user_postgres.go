package repository

import (
	"database/sql"
	"fmt"
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

func (r *UsersRepo) CreateUser(userModel *model.RegisterInput) (err error) {
	err = r.db.QueryRow(
		`WITH u AS (
			INSERT INTO units (domain, name, type) 
			VALUES ($1, $2, 'USER') 
			RETURNING id
			) 
		INSERT INTO users (id, app_settings, password, email) 
		SELECT u.id, 'default', $3, $4 
		FROM u 
		RETURNING id`,
		userModel.Domain,
		userModel.Name,
		userModel.Password,
		userModel.Email,
	).Err()
	if err != nil {
		println("CreateUser:", err.Error()) // debug
	}
	return
}

func (r *UsersRepo) UserExistsByInput(inputModel *models.UserInput) (exists bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1
			FROM users
			WHERE email = $1 AND password = $2
			)`,
		inputModel.Email,
		inputModel.Password,
	).Scan(&exists)

	return

}

func (r *UsersRepo) GetUserIdByInput(inputModel *models.UserInput) (id int, err error) {
	err = r.db.QueryRow(
		`SELECT units.id
		FROM units INNER JOIN users 
		ON units.id = users.id 
		WHERE users.email = $1 AND users.password = $2`,
		inputModel.Email,
		inputModel.Password,
	).Scan(&id)
	if err != nil {
		return
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

func (r *UsersRepo) GetUsersByNameFragment(fragment string, offset int) (users models.ListUserInfo, err error) {
	rows, err := r.db.Query(
		`SELECT units.id, units.domain,units.name
		FROM units INNER JOIN users 
		ON units.id = users.id 
		WHERE units.name ILIKE $1
		LIMIT 20
		OFFSET $2`,
		"%"+fragment+"%",
		offset,
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

func (r *UsersRepo) GetUserSettings(userId int) (settings models.UserSettings, err error) {
	err = r.db.QueryRow(
		`SELECT users.app_settings
		FROM units INNER JOIN users 
		ON units.id = users.id 
		WHERE units.id = $1`,
		userId,
	).Scan(
		&settings.AppSettings,
	)
	if err != nil {
		return
	}
	return
}

func (r *UsersRepo) UpdateUserData(userId int, userModel *models.UpdateUserData) (err error) {
	if userModel.Domain != "" {
		err = r.db.QueryRow(
			`UPDATE units
			SET domain = $2
			WHERE id = $1`,
			userId,
			userModel.Domain,
		).Err()
		if err != nil {
			return
		}
	}
	if userModel.Name != "" {
		err = r.db.QueryRow(
			`UPDATE units
			SET name = $2
			WHERE id = $1`,
			userId,
			userModel.Name,
		).Err()
		if err != nil {
			return
		}
	}
	if userModel.Email != "" {
		err = r.db.QueryRow(
			`UPDATE users
			SET email = $2
			WHERE id = $1`,
			userId,
			userModel.Email,
		).Err()
		if err != nil {
			return
		}
	}
	if userModel.Password != "" {
		err = r.db.QueryRow(
			`UPDATE users
			SET password = $2
			WHERE id = $1`,
			userId,
			userModel.Password,
		).Err()
		if err != nil {
			return
		}
	}
	return
}

func (r *UsersRepo) UpdateUserSettings(userId int, settingsModel *models.UpdateUserSettings) error {
	err := r.db.QueryRow(
		`UPDATE users
		SET app_settings = $1
		WHERE id = $2`,
		settingsModel.AppSettings,
		userId,
	).Err()
	if err != nil {
		return err
	}
	return nil
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
		`SELECT units.id, units.domain, units.name, units.type, users.email, users.password 
		FROM units INNER JOIN users
		ON units.id = users.id
		WHERE units.id = $1`,
		usersId,
	).Scan(
		&me.User.Unit.ID,
		&me.User.Unit.Domain,
		&me.User.Unit.Name,
		&me.User.Unit.Type,
		&me.Data.Password,
		&me.Data.Password,
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
	rows, err := r.db.Query(
		`SELECT units.id, units.domain, units.name, units.type, chats.private
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
		fmt.Println(err.Error())
		fmt.Println("пользователи не найдены")
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

func (r UsersRepo) UpdateMe(userId int, inp *model.UpdateMeDataInput) (err error) {
	err = r.db.QueryRow(`
		with chat as (
			UPDATE units
			SET 
			    name = COALESCE($2::VARCHAR, name), 
			    domain = COALESCE($3::VARCHAR, domain)
			WHERE id = $1
		)
		UPDATE users
		SET 
		    password = COALESCE($4::VARCHAR, password),
		    email = COALESCE($4::VARCHAR, email)
		WHERE id = $1

		`,
		userId,
		inp.Name,
		inp.Domain,
		inp.Password,
		inp.Email,
	).Err()
	if err != nil {
		println("UpdateMe:", err.Error()) // debug
		return
	}

	return
}
