package subix

import (
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/repository"
)

type ID = int

func (s *Subix) NotifyChatMembers(chat ID, body model.EventResult) {
	s.spam(
		[]ID{chat},
		s.repo.Subscribers.Members,
		body,
	)
}
func (s *Subix) NotifyChats(chats []ID, body model.EventResult) {
	s.spam(
		chats,
		s.repo.Subscribers.Members,
		body,
	)
}
func (s *Subix) NotifyRoomReaders(room ID, body model.EventResult) {
	s.spam(
		[]ID{room},
		s.repo.Subscribers.RoomReaders,
		body,
	)
}

func (s *Subix) spam(objects []ID, meth repository.QueryUserGroup, body interface{}) {
	users, err := meth(objects)
	if err != nil {
		panic(err)
	}

	switch body.(type) {
	case *model.DeleteInvite:
		s.ForceDropScheduledInvite(body.(*model.DeleteInvite).Code)

	case *model.CreateInvite:
		inv := body.(*model.CreateInvite)
		s.CreateScheduledInvite(objects[0], inv.Code, inv.ExpiresAt)

	}

	s.writeToUsers(users, body.(model.EventResult))
}
