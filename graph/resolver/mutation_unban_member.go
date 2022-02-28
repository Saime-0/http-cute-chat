package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/cerrors"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mutationResolver) UnbanMember(ctx context.Context, userID int, chatID int) (model.MutationResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("UnbanMember", &bson.M{
		"userID": userID,
		"chatID": chatID,
	})
	defer node.MethodTiming()

	var clientID = utils.GetAuthDataFromCtx(ctx).UserID

	if node.ChatExists(chatID) ||
		node.UserExists(userID) ||
		node.CanUnban(clientID, chatID) ||
		//node.IsNotMember(userID, chatID) ||
		node.IsBanned(userID, chatID) {
		return node.GetError(), nil
	}

	err := r.Services.Repos.Chats.UnbanUserInChat(userID, chatID)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return resp.Error(resp.ErrInternalServerError, "не удалось разбанить пользователя"), nil
	}

	return resp.Success("пользователь успешно разбанен"), nil
}
