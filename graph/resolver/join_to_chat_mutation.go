package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/rules"
)

func (r *mutationResolver) JoinToChat(ctx context.Context, chatID int) (model.JoinToChatResult, error) {
	node := r.Piper.CreateNode("mutationResolver > JoinToChat [cid:", chatID, "]")
	defer node.Kill()

	var clientID = ctx.Value(rules.UserIDFromToken).(int)

	if node.ChatExists(chatID) ||
		node.IsNotMember(clientID, chatID) ||
		// todo UserIsNotBanned
		node.ChatIsNotPrivate(chatID) ||
		node.MembersLimit(chatID) ||
		node.ChatsLimit(clientID) {
		return node.Err, nil
	}
	eventReadyMember, err := r.Services.Repos.Chats.AddUserToChat(clientID, chatID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "не удалось присоединиться"), nil
	}

	go r.Services.Subix.NotifyChatMembers(
		chatID,
		eventReadyMember,
	)

	return resp.Success("успешно присоединился к чату"), nil
}
