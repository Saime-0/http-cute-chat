package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/rules"
)

func (r *mutationResolver) LeaveFromChat(ctx context.Context, chatID int) (model.MutationResult, error) {
	node := r.Piper.CreateNode("mutationResolver > LeaveFromChat [cid:", chatID, "]")
	defer node.Kill()

	var clientID = ctx.Value(rules.UserIDFromToken).(int)

	if node.ChatExists(chatID) ||
		node.IsMember(clientID, chatID) ||
		node.CanLeaveFromChat(clientID, chatID) {
		return node.Err, nil
	}

	eventReadyMember, err := r.Services.Repos.Chats.RemoveUserFromChat(clientID, chatID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), nil
	}

	go r.Services.Subix.NotifyChatMembers(
		[]int{chatID},
		eventReadyMember,
	)

	return resp.Success("успешно покинул чат"), nil
}
