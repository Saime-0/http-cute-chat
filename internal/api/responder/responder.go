package responder

import (
	"encoding/json"
	"net/http"

	"github.com/saime-0/http-cute-chat/internal/api/rules"
)

type Writer struct {
	http.ResponseWriter
	Code int
}

func (w *Writer) WriteHeader(statusCode int) {
	w.Code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func Respond(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(data)
	}
}

type ResponseError struct {
	Error rules.AdvancedError `json:"error"`
}

func Error(w http.ResponseWriter, code int, err *rules.AdvancedError) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(ResponseError{
		Error: *err,
	})

}
