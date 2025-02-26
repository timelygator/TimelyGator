package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func JsonOK(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func HttpError(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)
	fmt.Fprintf(w, `{"error": %q}`, err.Error())
}

func HttpErrorString(w http.ResponseWriter, msg string, code int) {
	w.WriteHeader(code)
	fmt.Fprintf(w, `{"error": %q}`, msg)
}
