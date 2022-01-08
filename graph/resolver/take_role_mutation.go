package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/rules"
)

func (r *mutationResolver) TakeRole(ctx context.Context, memberID int) (model.MutationResult, error) {
	node := r.Piper.CreateNode("mutationResolver > TakeRole [mid:", memberID, "]")
	defer node.Kill()

	var (
		clientID       = ctx.Value(rules.UserIDFromToken).(int)
		chatID         int
		clientMemberID int
	)

	if node.ValidID(memberID) ||
		node.GetChatIDByMember(memberID, &chatID) ||
		node.GetMemberBy(clientID, chatID, &clientMemberID) ||
		node.CanTakeRole(clientMemberID, memberID, chatID) {

		return node.Err, nil
	}
	// todo migrate to..nonono..
	eventReadyMember, err := r.Services.Repos.Chats.TakeRole(memberID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "ошибка при попытке забрать роль"), nil
	}
	r.Services.Subix.NotifyChatMembers(
		[]int{chatID},
		eventReadyMember,
	)
	return resp.Success("успешно"), nil
}
