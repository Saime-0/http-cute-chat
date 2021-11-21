package resp

type AdvancedError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	// http status const
	// Err     error  `json:"-"`
}

var (
	ErrInternalSeverError = &AdvancedError{
		Code:    "InternalServerError",
		Message: "internal server error",
	}
)