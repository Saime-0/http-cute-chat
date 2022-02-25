package repository

import (
	"database/sql"
	"github.com/lib/pq"
)

type QueryUserGroup func(objectIDs []int) (users []int, err error)

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

func completeIntArray(rows *sql.Rows) (arr []int, err error) {
	for rows.Next() {
		var id int
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}
		arr = append(arr, id)
	}
	return arr, nil
}

func (r *SubscribersRepo) initFuncs() {

	r.Members = func(chatIDs []int) (users []int, err error) {
		//language=PostgreSQL
		rows, err := r.db.Query(`
			SELECT id 
			FROM chat_members
			WHERE chat_id = ANY ($1)`,
			pq.Array(chatIDs),
		)
		defer rows.Close()
		if err != nil {
			return nil, err
		}

		users, err = completeIntArray(rows)
		if err != nil {
			return nil, err
		}

		return users, nil
	}

	r.RoomReaders = func(roomIDs []int) (members []int, err error) {
		rows, err := r.db.Query(`
			SELECT m.id
			FROM chat_members m

			JOIN chats c on m.chat_id = c.id
			JOIN rooms r on m.chat_id = r.chat_id
			LEFT JOIN allows a on r.id = a.room_id

			WHERE r.id = ANY ($1)
		  		AND (
				    action_type IS NULL 
				    OR action_type = 'READ' 
            	        AND (
							group_type = 'ROLE' AND value = m.role_id::VARCHAR 
						    OR group_type = 'CHAR' AND value = m.char::VARCHAR 
						    OR group_type = 'MEMBER' AND value = m.id::VARCHAR
						)
				   	OR m.user_id = c.owner_id
				)
			GROUP BY m.id`,
			pq.Array(roomIDs),
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		members, err = completeIntArray(rows)
		if err != nil {
			return nil, err
		}

		return members, nil
	}
}
