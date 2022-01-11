package subix

import "github.com/saime-0/http-cute-chat/internal/models"

type Store struct {
	ScheduleInvites map[string]*models.ScheduleInvite
}

func newStore() *Store {
	return &Store{
		ScheduleInvites: map[string]*models.ScheduleInvite{},
	}
}
