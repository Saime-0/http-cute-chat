package resolver

import (
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/cerrors"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/res"
	"github.com/saime-0/http-cute-chat/internal/utils"
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
			r.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
			panic(err)
		}
		err = r.prepareScheduleRegisterSessions(end)
		if err != nil {
			r.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
			panic(err)
		}
		err = r.prepareScheduleRefreshSessions(end)
		if err != nil {
			r.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
			panic(err)
		}

		select {
		case ready <- 1:
			//  сигнал о готовности был услышан
		default:
			// сигнал о готовности никто не услышал
		}
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

			// начато прослушивание канала
			<-ready
			// получен сигнал о готовности
		}
	}()

	return nil

}

type scheduleInviteMap map[string]*models.ScheduleInvite

func (r *Resolver) prepareScheduleInvites(before int64) (err error) {
	invites, err := r.Services.Repos.Prepares.ScheduleInvites(before)
	if err != nil {
		return cerrors.Wrap(err, "не удалось подготовить инвайты")
	}

	r.Services.Cache.Set(res.CacheScheduleInvites, &scheduleInviteMap{})

	for _, inv := range invites {
		if inv.Exp == nil {
			continue
		}

		if *inv.Exp <= time.Now().Unix() {
			// удаляю инвайт, тк он уже истек
			if _, err := r.Services.Repos.Chats.DeleteInvite(inv.Code); err != nil {
				return err
			}
			continue
		}

		err = r.CreateScheduledInvite(inv.ChatID, inv.Code, inv.Exp)
		if err != nil {
			return cerrors.Wrap(err, "не удалось подготовить инвайты")
		}
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
			// удаляю сессию, тк она уже истекла
			err := r.Services.Repos.Users.DeleteRegistrationSession(rs.Email)
			if err != nil {
				r.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
			}
			continue
		}
		_, err = r.Services.Scheduler.AddTask(
			func() {
				// спланированное удаление
				err := r.Services.Repos.Users.DeleteRegistrationSession(rs.Email)
				if err != nil {
					r.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
				}
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
			// удаляю сессию, тк она уже истекла
			err := r.Services.Repos.Users.DeleteRefreshSession(rs.ID)
			if err != nil {
				r.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
			}
			continue
		}
		_, err = r.Services.Scheduler.AddTask(
			func() {
				// спланированное удаление
				err := r.Services.Repos.Users.DeleteRefreshSession(rs.ID)
				if err != nil {
					r.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
				}
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
		return cerrors.New("not found scheduled invites map in cache")
	}
	invs, ok := intfMap.(*scheduleInviteMap)
	if !ok {
		return cerrors.New("scheduled invites map has invalid type")
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
		return cerrors.New("exp cannot be equal to nil")
	}
	interfacedMap, ok := r.Services.Cache.Get(res.CacheScheduleInvites)
	if !ok {
		return cerrors.New("not found scheduled invites map in cache")
	}
	invites, ok := interfacedMap.(*scheduleInviteMap)
	if !ok {
		return cerrors.New("scheduled invites map has invalid type")
	}
	task, err := r.Services.Scheduler.AddTask(
		func() {

			_, err := r.Services.Repos.Chats.DeleteInvite(code)
			if err != nil {
				r.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
			}
			r.Subix.NotifyChatMembers(
				chatID,
				&model.DeleteInvite{
					Reason: model.DeleteInviteReasonExpired,
					Code:   code,
				},
			)
			// удаляю инвайт, тк он уже истек
			delete(*invites, code)
		},
		*exp,
	)
	if err != nil {
		r.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return cerrors.Wrap(err, "не удалось запланировать задачу")
	}
	(*invites)[code] = &models.ScheduleInvite{
		Task: task,
	}
	return nil
}
