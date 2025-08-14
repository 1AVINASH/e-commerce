// internal/middleware/json.go
package middleware

import (
	"encoding/json"
	"net/http"
)

func JSON(next func(w http.ResponseWriter, r *http.Request) (interface{}, int)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		data, status := next(w, r)
		w.WriteHeader(status)
		if data != nil {
			json.NewEncoder(w).Encode(data)
		}
	}
}
