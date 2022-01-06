package subix

import (
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/repository"
)

func (s *Subscription) Spam(objectID int, meth repository.QueryUserGroup, body model.EventResult) {
	users, err := meth(objectID)
	if err != nil {
		panic(err)
		return
	}
	s.writeToUsers(users, body)
}
