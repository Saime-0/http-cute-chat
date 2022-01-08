package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/rules"
)

func (r *mutationResolver) DeleteInvite(ctx context.Context, chatID int, code string) (model.MutationResult, error) {
	node := r.Piper.CreateNode("mutationResolver > DeleteInvite [cid:", chatID, ",code:", code, "]")
	defer node.Kill()

	clientID := ctx.Value(rules.UserIDFromToken).(int)

	if node.ChatExists(chatID) ||
		node.IsMember(clientID, chatID) ||
		node.CanCreateInvite(clientID, chatID) ||
		node.HasInvite(chatID, code) ||
		node.InviteIsRelevant(code) {
		return node.Err, nil
	}

	eventReadyInvite, err := r.Services.Repos.Chats.DeleteInvite(code)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "не удалось удалить инвайт"), nil
	}
	r.Services.Subix.NotifyChatMembers(
		[]int{chatID},
		eventReadyInvite,
	)
	return resp.Success("инвайт успешно удален"), nil
}
