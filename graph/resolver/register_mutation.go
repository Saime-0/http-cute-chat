package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/saime-0/http-cute-chat/internal/api/resp"
	"github.com/saime-0/http-cute-chat/internal/api/responder"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"github.com/saime-0/http-cute-chat/internal/api/validator"
	"net/http"

	"github.com/saime-0/http-cute-chat/graph/model"
)

func (r *mutationResolver) Register(ctx context.Context, input model.RegisterInput) (model.RegisterResult, error) {
	switch {
	case !validator.ValidateDomain(input.Domain):
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidDomain)
		return resp.ErrInvalidDomain, nil

	case !validator.ValidateName(input.Name):
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidName)
		return resp.ErrInvalidName, nil

	case !validator.ValidateEmail(input.Email):
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidEmail)
		return resp.ErrInvalidEmail, nil

	case !validator.ValidatePassword(input.Password):
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidPassword)
		return resp.ErrInvalidPassword, nil
	}

}
