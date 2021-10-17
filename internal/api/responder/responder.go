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
		json.NewEncoder(w).Encode(data)
	}
}
