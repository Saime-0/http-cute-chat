package resp

import "github.com/saime-0/http-cute-chat/graph/model"

//type AdvancedError struct {
//	Code    string `json:"code"`
//	Message string `json:"message"`
//	// http status const
//	// Err     error  `json:"-"`
//}

func Success(msg string) model.Successful {
	return model.Successful{
		Success: msg,
	}
}
func Error(code ErrCode, msg string) *model.AdvancedError {
	return &model.AdvancedError{
		Code:  string(code),
		Error: msg,
	}
}
func ErrorCopy(code ErrCode, msg string) model.AdvancedError {
	return model.AdvancedError{
		Code:  string(code),
		Error: msg,
	}
}

type ErrCode string

const (
	ErrInternalServerError ErrCode = "InternalServerError"
	ErrBadRequest          ErrCode = "BadRequest"
	ErrNotFound            ErrCode = "NotFound"
	ErrNoAccess            ErrCode = "NoAccess"
)

// todo лог запросов с типом результата ответа(если ошибка то полностю ошибку)
