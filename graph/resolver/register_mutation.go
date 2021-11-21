package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/saime-0/http-cute-chat/internal/api/resp"
	"github.com/saime-0/http-cute-chat/internal/api/validator"

	"github.com/saime-0/http-cute-chat/graph/model"
)

func (r *mutationResolver) Register(ctx context.Context, input model.RegisterInput) (model.RegisterResult, error) {
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

	// todo: обращение к бд? вывод?
}
