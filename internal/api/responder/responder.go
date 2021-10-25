package responder

import (
	"encoding/json"
	"net/http"
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
	Error string `json:"error"`
}

func Error(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(ResponseError{
		Error: err.Error(),
	})

}
