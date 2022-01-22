package directive

import (
	"context"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/saime-0/http-cute-chat/internal/utils"
	"reflect"
)

func IsAuth(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	println("IsAuth directive start!") // debug

	if utils.GetAuthDataFromCtx(ctx) == nil {
		err = errors.New("не аутентифицирован")
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
	fmt.Printf("%#v %T\n", obj, obj)
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
		//input, ok = val.(map[string]interface{})
		//if !ok {
		//	panic(err)
		//}
	}
	if !finded {
		panic("InputLeastOne: union input field not found")
	}

	for _, val := range input {
		if fmt.Sprint(val) != "<nil>" {
			return next(ctx)
		}
	}

	return obj, errors.New("one of the input fields must have the value")

}
