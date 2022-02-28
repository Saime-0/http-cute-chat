package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/cerrors"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mutationResolver) ConfirmRegistration(ctx context.Context, email string, code string) (model.MutationResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("ConfirmRegistration", &bson.M{
		"email": email,
		"code":  code,
	})
	defer node.MethodTiming()

	var regi = &models.RegisterData{}
	if node.ValidEmail(email) ||
		node.GetRegistrationSession(email, code, &regi) {
		return node.GetError(), nil
	}
	err := r.Services.Repos.Users.CreateUser(regi)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return resp.Error(resp.ErrInternalServerError, "не удлось создать пользователя"), nil
	}

	err = r.Services.Repos.Users.DeleteRegistrationSession(email)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
	}

	return resp.Success("пользователь создан"), nil
}
