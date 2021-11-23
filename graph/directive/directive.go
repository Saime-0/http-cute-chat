package directive

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/saime-0/http-cute-chat/internal/api/resp"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
)

func IsAuth(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	if ctx.Value(rules.UserIDFromToken).(int) == 0 {
		return resp.Error(resp.ErrNoAccess, "клиент не аутентифицирован"), nil
	}
	return next(ctx)
}
