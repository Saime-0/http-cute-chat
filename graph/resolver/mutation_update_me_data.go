package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/pkg/errors"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/res"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mutationResolver) UpdateMeData(ctx context.Context, input model.UpdateMeDataInput) (model.MutationResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("UpdateMeData", &bson.M{
		"input": input,
	})
	defer node.MethodTiming()

	clientID := utils.GetAuthDataFromCtx(ctx).UserID

	if input.Name != nil && node.ValidName(*input.Name) ||
		input.Domain != nil && node.ValidDomain(*input.Domain) ||
		input.Password != nil && (node.ValidPassword(*input.Password) || func() bool {
			// the function in the condition is the best possible solution, as I believe. I'm sorry if this makes it difficult to read the code
			var err error
			*input.Password, err = utils.HashPassword(*input.Password, r.Config.PasswordSalt)
			if err != nil {
				node.Healer.Alert(errors.Wrap(err, utils.GetCallerPos()))
				node.SetError(resp.ErrInternalServerError, res.UnexpectedError)
				return true
			}
			return false
		}()) ||
		input.Email != nil && node.ValidDomain(*input.Email) {
		return node.GetError(), nil
	}

	eventReadyUser, err := r.Services.Repos.Users.UpdateMe(clientID, &input)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "не удалось обновить данные"), nil
	}

	if input.Name != nil || input.Domain != nil {
		chats, err := r.Services.Repos.Users.ChatsID(clientID)
		if err != nil {
			return nil, err
		} else {
			go r.Subix.NotifyChats(
				chats,
				eventReadyUser,
			)
		}
	}

	return resp.Success("данные пользователя обновлены"), nil
}
