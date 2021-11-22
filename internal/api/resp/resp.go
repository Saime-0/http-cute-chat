package resp

import "github.com/saime-0/http-cute-chat/graph/model"

//type AdvancedError struct {
//	Code    string `json:"code"`
//	Message string `json:"message"`
//	// http status const
//	// Err     error  `json:"-"`
//}

var (
	ErrInternalSeverError = &model.AdvancedError{
		Code:    "InternalServerError",
		Message: "internal server error",
	}

	ErrInvalidName = &model.AdvancedError{
		Code:    "InvalidUserName",
		Message: "invalid name",
	}
	ErrInvalidDomain = &model.AdvancedError{
		Code:    "InvalidDomain",
		Message: "invalid domain",
	}
	ErrInvalidEmail = &model.AdvancedError{
		Code:    "InvalidEmail",
		Message: "invalid email",
	}
	ErrInvalidPassword = &model.AdvancedError{
		Code:    "InvalidPassword",
		Message: "invalid password",
	}
	ErrNameFragment = &model.AdvancedError{
		Code:    "NameFragment",
		Message: "invalid name fragment",
	}
)
