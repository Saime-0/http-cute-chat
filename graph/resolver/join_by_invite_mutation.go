package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/rules"
)

func (r *mutationResolver) JoinByInvite(ctx context.Context, code string) (model.JoinByInviteResult, error) {
	node := r.Piper.CreateNode("mutationResolver > JoinByInvite [code:", code, "]")
	defer node.Kill()

	var (
		chatID   int
		clientID = ctx.Value(rules.UserIDFromToken).(int)
	)

	if node.InviteIsRelevant(code) ||
		node.GetChatByInvite(code, &chatID) ||
		node.IsNotMember(clientID, chatID) ||
		node.IsNotBanned(clientID, chatID) ||
		node.MembersLimit(chatID) ||
		node.ChatsLimit(clientID) {
		return node.Err, nil
	}

	eventReadyMember, err := r.Services.Repos.Chats.AddUserByCode(code, clientID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "не удалось присоединиться"), nil
	}

	go r.Services.Subix.NotifyChatMembers(
		[]int{chatID},
		eventReadyMember,
	)

	return resp.Success("успешно присоединился к чату"), nil
}
