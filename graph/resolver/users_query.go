package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/api/resp"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
)

func (r *queryResolver) Users(ctx context.Context, nameFragment string, params *model.Params) (model.UsersResult, error) {
	if len(nameFragment) > rules.NameMaxLength || len(nameFragment) == 0 {
		return resp.Error(resp.ErrBadRequest, "фрагмент не соответсвует требованиям"), nil
	}
	offset := 0
	if params != nil {
		offset = *params.Offset
	}
	userList, err := r.Services.Repos.Users.GetUsersByNameFragment(nameFragment, offset)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%#v", userList)
	users := model.Users{}
	for _, v := range userList.Users {
		user := model.User{
			Unit: &model.Unit{
				ID:     v.ID,
				Domain: v.Domain,
				Name:   v.Name,
				Type:   model.UnitTypeUser,
			},
		}
		users.Users = append(users.Users, &user)
	}

	return users, nil
}
