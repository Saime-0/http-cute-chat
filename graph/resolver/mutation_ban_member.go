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

func (r *mutationResolver) BanMember(ctx context.Context, memberID int) (model.MutationResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("BanMember", &bson.M{
		"memberID": memberID,
	})
	defer node.MethodTiming()

	clientID := utils.GetAuthDataFromCtx(ctx).UserID
	var (
		member models.DefMember
	)

	if node.GetDefMember(memberID, &member) ||
		node.CanBan(clientID, member.UserID, member.ChatID) {
		return node.GetError(), nil
	}

	eventReadyMember, err := r.Services.Repos.Chats.BanUserInChat(member.UserID, member.ChatID)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return resp.Error(resp.ErrInternalServerError, "не удалось забанить пользователя"), nil
	}

	go r.Subix.NotifyChatMembers(
		member.ChatID,
		eventReadyMember,
	)
	r.Subix.DeleteMember(memberID)

	return resp.Success("пользователь успешно забанен"), nil
}
