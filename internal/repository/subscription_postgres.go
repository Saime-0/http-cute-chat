package repository

import (
	"database/sql"
)

type QueryUserGroup func(objectID int) (users []int, err error)

type SubscribersRepo struct {
	db          *sql.DB
	Members     QueryUserGroup
	RoomReaders QueryUserGroup
}

func NewSubscribersRepo(db *sql.DB) *SubscribersRepo {
	sub := &SubscribersRepo{
		db: db,
	}
	sub.initFuncs()
	return sub
}

func completeUsers(rows *sql.Rows) (users []int, err error) {
	for rows.Next() {
		var id int
		if err = rows.Scan(&id); err != nil {
			return users, err
		}
		users = append(users, id)
	}
	return users, nil
}

func (r *SubscribersRepo) initFuncs() {

	r.Members = func(chatID int) (users []int, err error) {

		rows, err := r.db.Query(`
		SELECT id 
		FROM chat_members
		WHERE chat_id = $1
		`,
			chatID,
		)
		defer rows.Close()
		if err != nil {
			println("Members:", err.Error()) // debug
			return users, err
		}

		users, err = completeUsers(rows)
		if err != nil {
			println("Members(Scan):", err.Error()) // debug
			return users, err
		}

		return users, nil
	}

	r.RoomReaders = func(roomID int) (users []int, err error) {
		rows, err := r.db.Query(`
		SELECT user_id
		FROM chat_members m
		    
		JOIN chats c on m.chat_id = c.id
		JOIN rooms r on m.chat_id = r.chat_id
		LEFT JOIN allows a on r.id = a.room_id
		
		WHERE r.id = $1 
		  	AND (
			    action_type IS NULL 
			    OR action_type = 'READ' 
                    AND (
						group_type = 'ROLES' AND value = m.role_id::VARCHAR 
					    OR group_type = 'CHARS' AND value = m.char::VARCHAR 
					    OR group_type = 'USERS' AND value = m.user_id::VARCHAR
					)
			   	OR m.user_id = c.owner_id
			)
			
		GROUP BY user_id
		`,
			roomID,
		)
		defer rows.Close()
		if err != nil {
			println("Members:", err.Error()) // debug
			return users, err
		}

		users, err = completeUsers(rows)
		if err != nil {
			println("Members(Scan):", err.Error()) // debug
			return users, err
		}

		return users, nil
	}
}
