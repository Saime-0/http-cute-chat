package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/utils"
)

func (r *mutationResolver) UpdateMember(ctx context.Context, memberID int, input model.UpdateMemberInput) (model.MutationResult, error) {
	node := r.Piper.CreateNode("mutationResolver > UpdateMember [mid:", memberID, "]")
	defer node.Kill()

	var (
		chatID   int
		clientID = utils.GetAuthDataFromCtx(ctx).UserID
	)

	if node.GetChatIDByMember(memberID, &chatID) ||
		input.RoleID != nil && node.CanGiveRole(clientID, chatID) && node.RoleExists(chatID, *input.RoleID) ||
		input.Char != nil && node.CanGiveChar(clientID, chatID) ||
		input.Muted != nil && node.CanMuteMember(clientID, chatID) {
		return node.Err, nil
	}

	eventReadyMember, err := r.Services.Repos.Chats.UpdateMember(memberID, &input)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "не удалось обновить данные"), nil
	}
	go r.Services.Subix.NotifyChatMembers(
		chatID,
		eventReadyMember,
	)
	return resp.Success("данные участника чата обновлены"), nil
}
