package resolver

import (
	"github.com/pkg/errors"
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/res"
	"github.com/saime-0/http-cute-chat/pkg/kit"
	"time"
)

// да, я знаю что это криндж, не душите

func (r *Resolver) RegularSchedule(interval int64) (err error) {
	ready := make(chan int8)

	regFn := func() {
		end := kit.After(interval)

		err = r.prepareScheduleInvites(end)
		if err != nil {
			panic(err)
		}
		err = r.prepareScheduleRegisterSessions(end)
		if err != nil {
			panic(err)
		}
		err = r.prepareScheduleRefreshSessions(end)
		if err != nil {
			panic(err)
		}

		select {
		case ready <- 1:
			println("RegularSchedule: сигнал о готовности был услышан") // debug
		default:
			println("RegularSchedule: сигнал о готовности никто не услышал") // debug
		}
		println("RegularSchedule come true") // debug
	}

	regFn()

	go func() {
		for {
			runAt := kit.After(interval)
			_, err = r.Services.Scheduler.AddTask(
				regFn,
				runAt,
			)
			if err != nil {
				panic(err)
			}
			r.Services.Cache.Set(res.CacheNextRunRegularScheduleAt, runAt)

			println("RegularSchedule: начато прослушивание канала") // debug
			<-ready
			println("RegularSchedule: получен сигнал о готовности") // debug
		}
	}()

	return nil

}

type scheduleInviteMap map[string]*models.ScheduleInvite

func (r *Resolver) prepareScheduleInvites(before int64) (err error) {
	invites, err := r.Services.Repos.Prepares.ScheduleInvites(before)
	if err != nil {
		return err
	}

	r.Services.Cache.Set(res.CacheScheduleInvites, &scheduleInviteMap{})

	for _, inv := range invites {
		if inv.Exp == nil {
			continue
		}

		if *inv.Exp <= time.Now().Unix() {
			if _, err := r.Services.Repos.Chats.DeleteInvite(inv.Code); err != nil {
				return err
			}
			println("prepareScheduleInvites: удаляю инвайт, тк он уже истек", inv.Code) // debug
			continue
		}

		r.CreateScheduledInvite(inv.ChatID, inv.Code, inv.Exp)
	}

	return nil
}

func (r *Resolver) prepareScheduleRegisterSessions(before int64) (err error) {
	sessions, err := r.Services.Repos.Prepares.ScheduleRegisterSessions(before)
	if err != nil {
		return err
	}
	for _, rs := range sessions {

		if rs.Exp <= time.Now().Unix() {
			r.Services.Repos.Users.DeleteRegistrationSession(rs.Email)
			println("prepareScheduleRegisterSessions: удаляю сессию, тк она уже истекла", rs.Email) // debug
			continue
		}
		_, err = r.Services.Scheduler.AddTask(
			func() {
				r.Services.Repos.Users.DeleteRegistrationSession(rs.Email)
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

func (r *Resolver) prepareScheduleRefreshSessions(before int64) (err error) {
	sessions, err := r.Services.Repos.Prepares.ScheduleRefreshSessions(before)
	if err != nil {
		return err
	}
	for _, rs := range sessions {

		if rs.Exp <= time.Now().Unix() {
			r.Services.Repos.Users.DeleteRefreshSession(rs.ID)
			println("prepareScheduleRefreshSessions: удаляю сессию, тк она уже истекла", rs.ID) // debug
			continue
		}
		_, err = r.Services.Scheduler.AddTask(
			func() {
				r.Services.Repos.Users.DeleteRefreshSession(rs.ID)
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

func (r *Resolver) ForceDropScheduledInvite(code string) error {
	intfMap, ok := r.Services.Cache.Get(res.CacheScheduleInvites)
	if !ok {
		return errors.New("not found scheduled invites map in cache")
	}
	invs, ok := intfMap.(*scheduleInviteMap)
	if !ok {
		return errors.New("scheduled invites map has invalid type")
	}
	inv, ok := (*invs)[code]
	if ok {
		err := r.Services.Scheduler.DropTask(&inv.Task)
		if err != nil {
			panic(err)
		}
		delete(*invs, code)
	}
	return nil
}

func (r *Resolver) CreateScheduledInvite(chatID int, code string, exp *int64) error {
	if exp == nil {
		return errors.New("exp cannot be equal to nil")
	}
	interfacedMap, ok := r.Services.Cache.Get(res.CacheScheduleInvites)
	if !ok {
		return errors.New("not found scheduled invites map in cache")
	}
	invites, ok := interfacedMap.(*scheduleInviteMap)
	if !ok {
		return errors.New("scheduled invites map has invalid type")
	}
	task, err := r.Services.Scheduler.AddTask(
		func() {

			_, err := r.Services.Repos.Chats.DeleteInvite(code)
			if err != nil {
				//return errors.Wrap(err, "не удалось удалить инвайт")
				panic(err)
				//todo log
			}
			r.Subix.NotifyChatMembers(
				chatID,
				&model.DeleteInvite{
					Reason: model.DeleteInviteReasonExpired,
					Code:   code,
				},
			)
			delete(*invites, code)
			println("prepareScheduleInvites: удаляю инвайт, тк он уже истек", code) // debug
		},
		*exp,
	)
	if err != nil {
		return errors.Wrap(err, "не удалось запланировать задачу")
	}
	(*invites)[code] = &models.ScheduleInvite{
		Task: task,
	}
	return nil
}
