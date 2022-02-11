package config

import (
	"reflect"
)

func (v FromCfgFile) validate() bool {
	s := reflect.ValueOf(v)
	for i := 0; i < s.NumField(); i++ {
		if s.Field(i).IsNil() {
			println(i, "is nil of FromCfgFile validate")
			return false
		}
	}
	return true
}

func (v FromEnv) validate() bool {
	s := reflect.ValueOf(v)
	for i := 0; i < s.NumField(); i++ {
		if s.Field(i).IsZero() {
			println(i, "is zero of FromEnv validate")
			return false
		}
	}
	return true
}
