package subix

import (
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/repository"
)

type ID = int

func (s *Subscription) spam(objects []ID, meth repository.QueryUserGroup, body model.EventResult) {
	users, err := meth(objects)
	if err != nil {
		panic(err)
		return
	}
	s.writeToUsers(users, body)
}
func (s *Subscription) NotifyChatMembers(chats []ID, body model.EventResult) {
	s.spam(
		chats,
		s.repo.Subscribers.Members,
		body,
	)
}
func (s *Subscription) NotifyRoomReaders(rooms []ID, body model.EventResult) {
	s.spam(
		rooms,
		s.repo.Subscribers.RoomReaders,
		body,
	)
}
