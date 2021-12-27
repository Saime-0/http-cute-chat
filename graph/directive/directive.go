package directive

import (
	"context"
	"errors"
	"github.com/99designs/gqlgen/graphql"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/rules"
	"reflect"
)

func IsAuth(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	if ctx.Value(rules.UserIDFromToken).(int) == 0 {
		return resp.Error(resp.ErrNoAccess, "клиент не аутентифицирован"), nil
	}
	return next(ctx)
}

func InputUnion(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	inputObj, err := next(ctx)
	if err != nil {
		return inputObj, err
	}

	v := reflect.ValueOf(inputObj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	valueFound := false

	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).IsNil() {
			if valueFound {
				return inputObj, errors.New("only one field of the input union should have a value")
			}

			valueFound = true
		}
	}

	if !valueFound {
		return inputObj, errors.New("one of the input union fields must have a value")
	}

	return inputObj, nil
}
