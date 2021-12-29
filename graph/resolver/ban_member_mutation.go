package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/rules"
)

func (r *mutationResolver) BanMember(ctx context.Context, memberID int) (model.MutationResult, error) {
	node := r.Piper.CreateNode("mutationResolver > BanMember [mid:", memberID, "]")
	defer node.Kill()

	clientID := ctx.Value(rules.UserIDFromToken).(int)
	var (
		defMember models.DefMember
	)

	if node.GetDefMember(memberID, &defMember) ||
		node.CanBan(clientID, defMember.UserID, defMember.ChatID) {
		return node.Err, nil
	}
	if r.Services.Repos.Chats.AddToBanlist(defMember.UserID, defMember.ChatID) != nil ||
		// todo is deprecated
		// todo banlist.userid > memberid
		// chat_namlist > banlist
		// + query banlist
		r.Services.Repos.Chats.RemoveUserFromChat(defMember.UserID, defMember.ChatID) != nil {
		return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), nil
	}
	return resp.Success("пользователь успешно забанен"), nil
}
