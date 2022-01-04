package directive

import (
	"context"
	"errors"
	"github.com/99designs/gqlgen/graphql"
	"github.com/saime-0/http-cute-chat/internal/rules"
	"reflect"
)

func IsAuth(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	println("IsAuth directive start!") // debug

	if ctx.Value(rules.UserIDFromToken).(int) == 0 {
		err = errors.New("клиент не аутентифицирован")
		println("IsAuth:", err.Error()) // debug
		return obj, err
	}

	return next(ctx)
}

func InputUnion(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	println("InputUnion directive start!") // debug

	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	valueFound := false

	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).IsNil() {
			if valueFound {
				println("InputUnion:", err.Error()) // debug
				return obj, errors.New("only one field of the input union should have a value")
			}

			valueFound = true
		}
	}

	if !valueFound {
		println("InputUnion:", err.Error()) // debug
		return obj, errors.New("one of the input union fields must have a value")
	}

	return next(ctx)

}

func InputLeastOne(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	println("InputLeastOne directive start!") // debug

	v := reflect.ValueOf(nil)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	valueFound := false

	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).IsNil() {
			valueFound = true
			break
		}
	}

	if !valueFound {
		println("InputLeastOne:", err.Error()) // debug
		return nil, errors.New("one of the input fields must have the value")
	}

	return next(ctx)

}
