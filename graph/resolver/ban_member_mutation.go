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
		member models.DefMember
	)

	if node.GetDefMember(memberID, &member) ||
		node.CanBan(clientID, member.UserID, member.ChatID) {
		return node.Err, nil
	}
	if r.Services.Repos.Chats.AddToBanlist(member.UserID, member.ChatID) != nil {
		return resp.Error(resp.ErrInternalServerError, "не удалось забанить пользователя"), nil
	}
	eventReadyMember, err := r.Services.Repos.Chats.RemoveUserFromChat(member.UserID, member.ChatID)
	if err != nil {
		panic(err)
	}

	r.Services.Subix.NotifyChatMembers(
		[]int{member.ChatID},
		eventReadyMember,
	)

	return resp.Success("пользователь успешно забанен"), nil
}
