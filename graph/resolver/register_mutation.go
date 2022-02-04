package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"time"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/res"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/rules"
	"github.com/saime-0/http-cute-chat/internal/utils"
)

func (r *mutationResolver) Register(ctx context.Context, input model.RegisterInput) (model.RegisterResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Register", &bson.M{
		"input": input,
	})
	defer node.MethodTiming()

	if node.ValidRegisterInput(&input) ||
		node.DomainIsFree(input.Domain) ||
		node.EmailIsFree(input.Email) {
		return node.GetError(), nil

	}
	expAt := time.Now().Unix() + rules.LiftimeOfRegistrationSession
	code, err := r.Services.Repos.Users.CreateRegistrationSession(
		&models.RegisterData{
			Domain: input.Domain,
			Name:   input.Name,
			Email:  input.Email,
			HashPassword: func() string {
				hpasswd, err := utils.HashPassword(input.Password, r.Config.PasswordSalt)
				if err != nil {
					panic(err)
				}
				return hpasswd
			}(),
		},
		expAt,
	)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), nil
	}

	err = r.Services.SMTP.Send(
		"код для подтверждения регистрации",
		"Для подтверждения ваших учетных данных используйте код: "+code,
		input.Email,
	)
	if err != nil {
		println("Register:", err.Error()) // debug
		return resp.Error(resp.ErrInternalServerError, "не удалось отправить код подтверждения на указанную почту"), nil
	}

	if runAt, ok := r.Services.Cache.Get(res.CacheNextRunRegularScheduleAt); ok && expAt < runAt.(int64) {
		_, err = r.Services.Scheduler.AddTask(
			func() {
				r.Services.Repos.Users.DeleteRegistrationSession(input.Email)
			},
			expAt,
		)
		if err != nil {
			panic(err)
		}
	}

	return resp.Success("подтвердите регистрацию, код отправлен на указанную почту"), nil
}
