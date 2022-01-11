package service

import (
	"time"
)

func (s *Services) prepareScheduleInvites() (err error) {
	invites, err := s.Repos.Prepares.ScheduleInvites()
	if err != nil {
		return err
	}
	for _, inv := range invites {
		if inv.Exp == nil {
			continue
		}

		if *inv.Exp <= time.Now().Unix() {
			if _, err := s.Repos.Chats.DeleteInvite(inv.Code); err != nil {
				return err
			}
			continue
		}

		s.Subix.CreateScheduledInvite(inv.ChatID, inv.Code, inv.Exp)
	}

	return nil
}
