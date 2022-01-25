package subix

import (
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/models"
)

func (s *Subix) ForceDropScheduledInvite(code string) {
	inv, ok := s.Store.ScheduleInvites[code]
	if ok {
		err := s.sched.DropTask(&inv.Task)
		if err != nil {
			panic(err)
		}
		delete(s.Store.ScheduleInvites, code)
	}
}

func (s *Subix) CreateScheduledInvite(chatID int, code string, exp *int64) {
	if exp == nil {
		return
	}
	task, err := s.sched.AddTask(
		func() {

			_, err := s.repo.Chats.DeleteInvite(code)
			if err != nil {
				//panic(err)
				println("prepareScheduleInvites:", err) // debug
				return
			}
			s.NotifyChatMembers(
				chatID,
				&model.DeleteInvite{
					Reason: model.DeleteInviteReasonExpired,
					Code:   code,
				},
			)
			delete(s.Store.ScheduleInvites, code)
			println("prepareScheduleInvites: удаляю инвайт, тк он уже истек", code) // debug
		},
		*exp,
	)
	if err != nil {
		panic(err)
	}
	s.Store.ScheduleInvites[code] = &models.ScheduleInvite{
		Task: task,
	}
}
