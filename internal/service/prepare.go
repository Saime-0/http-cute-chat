package service

import (
	"github.com/saime-0/http-cute-chat/internal/res"
	"github.com/saime-0/http-cute-chat/pkg/kit"
	"time"
)

func (s *Services) regularSchedule(interval int64) (err error) {
	ready := make(chan int8)

	regFn := func() {
		end := kit.After(interval)

		err = s.prepareScheduleInvites(end)
		if err != nil {
			panic(err)
		}
		err = s.prepareScheduleRegisterSessions(end)
		if err != nil {
			panic(err)
		}
		err = s.prepareScheduleRefreshSessions(end)
		if err != nil {
			panic(err)
		}

		select {
		case ready <- 1:
			println("regularSchedule: сигнал о готовности был услышан") // debug
		default:
			println("regularSchedule: сигнал о готовности никто не услышал") // debug
		}
		println("regularSchedule come true") // debug
	}

	regFn()

	go func() {
		for {
			runAt := kit.After(interval)
			_, err = s.Scheduler.AddTask(
				regFn,
				runAt,
			)
			if err != nil {
				panic(err)
			}
			s.Cache.Set(res.CacheNextRunRegularScheduleAt, runAt)

			println("regularSchedule: начато прослушивание канала") // debug
			<-ready
			println("regularSchedule: получен сигнал о готовности") // debug
		}
	}()

	return nil

}

func (s *Services) prepareScheduleInvites(before int64) (err error) {
	invites, err := s.Repos.Prepares.ScheduleInvites(before)
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
			println("prepareScheduleInvites: удаляю инвайт, тк он уже истек", inv.Code) // debug
			continue
		}

		s.Subix.CreateScheduledInvite(inv.ChatID, inv.Code, inv.Exp)
	}

	return nil
}

func (s *Services) prepareScheduleRegisterSessions(before int64) (err error) {
	sessions, err := s.Repos.Prepares.ScheduleRegisterSessions(before)
	if err != nil {
		return err
	}
	for _, rs := range sessions {

		if rs.Exp <= time.Now().Unix() {
			s.Repos.Users.DeleteRegistrationSession(rs.Email)
			println("prepareScheduleRegisterSessions: удаляю сессию, тк она уже истекла", rs.Email) // debug
			continue
		}
		_, err = s.Scheduler.AddTask(
			func() {
				s.Repos.Users.DeleteRegistrationSession(rs.Email)
				println("prepareScheduleRegisterSessions: спланирвоанное удаление", rs.Email) // debug
			},
			rs.Exp,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Services) prepareScheduleRefreshSessions(before int64) (err error) {
	sessions, err := s.Repos.Prepares.ScheduleRefreshSessions(before)
	if err != nil {
		return err
	}
	for _, rs := range sessions {

		if rs.Exp <= time.Now().Unix() {
			s.Repos.Users.DeleteRefreshSession(rs.ID)
			println("prepareScheduleRefreshSessions: удаляю сессию, тк она уже истекла", rs.ID) // debug
			continue
		}
		_, err = s.Scheduler.AddTask(
			func() {
				s.Repos.Users.DeleteRefreshSession(rs.ID)
				println("prepareScheduleRefreshSessions: спланированное удаление", rs.ID) // debug
			},
			rs.Exp,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
