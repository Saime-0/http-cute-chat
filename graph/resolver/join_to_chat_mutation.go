package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/api/resp"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"github.com/saime-0/http-cute-chat/internal/piping"
)

func (r *mutationResolver) JoinToChat(ctx context.Context, chatID int) (model.JoinToChatResult, error) {
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	pl := piping.NewPipeline(r.Services.Repos)
	if pl.ChatExists(chatID) ||
		pl.IsNotMember(clientID, chatID) ||
		// todo UserIsNotBanned
		pl.ChatIsNotPrivate(chatID) ||
		pl.MembersLimit(chatID) ||
		pl.ChatsLimit(clientID) {
		return pl.Err, nil
	}
	err := r.Services.Repos.Chats.AddUserToChat(clientID, chatID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), nil
	}
	return resp.Success("успешно присоединился к чату"), nil
}
