package directive

import (
	"context"
	"errors"
	"github.com/99designs/gqlgen/graphql"
	"github.com/saime-0/http-cute-chat/internal/rules"
	"reflect"
)

func IsAuth(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	inputObj, err := next(ctx)
	println("IsAuth directive start!") // debug

	if err != nil {
		println("IsAuth:", err.Error()) // debug
		return inputObj, err
	}

	if ctx.Value(rules.UserIDFromToken).(int) == 0 {
		err = errors.New("клиент не аутентифицирован")
		println("IsAuth:", err.Error()) // debug
		return inputObj, err
	}
	return inputObj, nil
}

func InputUnion(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	inputObj, err := next(ctx)
	if err != nil {
		println("InputUnion:", err.Error()) // debug
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
				println("InputUnion:", err.Error()) // debug
				return inputObj, errors.New("only one field of the input union should have a value")
			}

			valueFound = true
		}
	}

	if !valueFound {
		println("InputUnion:", err.Error()) // debug
		return inputObj, errors.New("one of the input union fields must have a value")
	}

	return inputObj, nil
}

func InputLeastOne(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	inputObj, err := next(ctx)
	if err != nil {
		println("InputLeastOne:", err.Error()) // debug
		return inputObj, err
	}

	v := reflect.ValueOf(inputObj)
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
		return inputObj, errors.New("one of the input fields must have the value")
	}

	return inputObj, nil
}
