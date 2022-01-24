package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/resp"
)

func (r *mutationResolver) ConfirmRegistration(ctx context.Context, email string, code string) (model.MutationResult, error) {
	node := r.Piper.CreateNode("mutationResolver > ConfirmRegistration [em:", email, "]")
	defer node.Kill()

	var regi = &models.RegisterData{}
	if node.ValidEmail(email) ||
		node.GetRegistrationSession(email, code, &regi) {
		return node.Err, nil
	}
	fmt.Printf("RegisterData: %#v\n", regi) // debug
	err := r.Services.Repos.Users.CreateUser(regi)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "не удлось создать пользователя"), nil
	}

	fmt.Printf("email: %s\n", email) // debug
	r.Services.Repos.Users.DeleteRegistrationSession(email)

	return resp.Success("пользователь создан"), nil
}
