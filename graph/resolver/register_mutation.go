package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"log"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/api/resp"
	"github.com/saime-0/http-cute-chat/internal/api/validator"
	"github.com/saime-0/http-cute-chat/internal/models"
)

func (r *mutationResolver) Register(ctx context.Context, input model.RegisterInput) (model.RegisterResult, error) {
	fmt.Printf("%#v", input)
	switch {
	case !validator.ValidateDomain(input.Domain):
		return resp.ErrInvalidDomain, nil

	case !validator.ValidateName(input.Name):
		return resp.ErrInvalidName, nil

	case !validator.ValidateEmail(input.Email):
		return resp.ErrInvalidEmail, nil

	case !validator.ValidatePassword(input.Password):
		return resp.ErrInvalidPassword, nil
	}
	id, err := r.Services.Repos.Users.CreateUser(&models.CreateUser{
		Domain:   input.Domain,
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		log.Println("err user create")
		return nil, err
	}
	log.Println("New user with id", id)
	return model.Successful{Success: "регистрация прошла успешно"}, nil
}
