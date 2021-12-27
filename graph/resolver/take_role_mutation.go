package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/rules"
	"github.com/saime-0/http-cute-chat/internal/tlog"
)

func (r *mutationResolver) TakeRole(ctx context.Context, memberID int) (model.MutationResult, error) {
	tl := tlog.Start("mutationResolver > TakeRole [mid:", memberID, "]")
	defer tl.Fine()

	clientID := ctx.Value(rules.UserIDFromToken).(int)

	node := r.Piper.CreateNode()
	defer node.Kill()

	var (
		chatID         int
		clientMemberID int
	)

	if node.ValidID(memberID) ||
		node.GetChatIDByMember(memberID, &chatID) ||
		node.GetMemberBy(clientID, chatID, &memberID) ||
		node.CanTakeRole(clientMemberID, memberID, chatID) {

		return node.Err, nil
	}
	// todo migrate to
	err := r.Services.Repos.Chats.TakeRole(memberID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), nil
	}
	return resp.Success("успешно"), nil
}
