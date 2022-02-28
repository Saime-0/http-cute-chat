package directive

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/saime-0/http-cute-chat/internal/cerrors"
	"github.com/saime-0/http-cute-chat/internal/utils"
	"reflect"
)

func IsAuth(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {

	if utils.GetAuthDataFromCtx(ctx) == nil {
		err = cerrors.New("не аутентифицирован")
		return obj, err
	}

	return next(ctx)
}

func InputUnion(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	valueFound := false

	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).IsNil() {
			if valueFound {
				return obj, cerrors.New("only one field of the input union should have a value")
			}

			valueFound = true
		}
	}

	if !valueFound {
		return obj, cerrors.New("one of the input union fields must have a value")
	}

	return next(ctx)

}

func InputLeastOne(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {

	input, ok := obj.(map[string]interface{})
	if !ok {
		panic("InputLeastOne: can not convert external map")
	}

	finded := false
	for key, val := range input {
		if key == "find" || key == "input" {
			finded = true
			input = val.(map[string]interface{})
			break
		}
	}
	if !finded {
		panic("InputLeastOne: union input field not found")
	}

	for _, val := range input {
		if fmt.Sprint(val) != "<nil>" {
			return next(ctx)
		}
	}

	return obj, cerrors.New("one of the input fields must have the value")

}
