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

func (r *mutationResolver) JoinByInvite(ctx context.Context, code string) (model.JoinByInviteResult, error) {
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	pl := piping.NewPipeline(r.Services.Repos)
	var chatID int
	if pl.InviteIsRelevant(code) ||
		pl.GetChatByInvite(code, &chatID) ||
		pl.IsNotMember(clientID, chatID) ||
		// todo UserIsNotBanned
		pl.MembersLimit(chatID) ||
		pl.ChatsLimit(clientID) {
		return pl.Err, nil
	}

	_, err := r.Services.Repos.Chats.AddUserByCode(code, clientID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), nil
	}
	return resp.Success("успешно присоединился к чату"), nil
}
