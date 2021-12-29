package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/validator"
)

func (r *mutationResolver) Register(ctx context.Context, input model.RegisterInput) (model.RegisterResult, error) {
	switch {
	case !validator.ValidateDomain(input.Domain):
		return resp.Error(resp.ErrBadRequest, "домен не соответствует требованиям"), nil

	case !validator.ValidateName(input.Name):
		return resp.Error(resp.ErrBadRequest, "имя не соответствует требованиям"), nil

	case !validator.ValidateEmail(input.Email):
		return resp.Error(resp.ErrBadRequest, "имеил не соответствует требованиям"), nil

	case !validator.ValidatePassword(input.Password):
		return resp.Error(resp.ErrBadRequest, "пароль не соответствует требованиям"), nil
	}

	_, err := r.Services.Repos.Users.CreateUser(&models.CreateUser{
		Domain:   input.Domain,
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), nil
	}

	return resp.Success("пользователь создан"), nil
}
