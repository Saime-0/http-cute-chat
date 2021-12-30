package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/rules"
)

func (r *mutationResolver) GiveRole(ctx context.Context, memberID int, roleID int) (model.MutationResult, error) {
	node := r.Piper.CreateNode("mutationResolver > GiveRole [mid:", memberID, ",rid:", roleID, "]")
	defer node.Kill()

	var (
		chatID   int
		clientID = ctx.Value(rules.UserIDFromToken).(int)
	)

	// todo MemberExists-> chatid; RoleExists(cid, rid):
	if node.GetChatIDByMember(memberID, &chatID) ||
		node.IsMember(clientID, chatID) ||
		node.CanGiveRole(clientID, chatID) ||
		node.RoleExists(chatID, roleID) {
		return node.Err, nil
	}

	err := r.Services.Repos.Chats.GiveRole(memberID, roleID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), nil
	}
	return resp.Success("роль выдана"), nil
}
