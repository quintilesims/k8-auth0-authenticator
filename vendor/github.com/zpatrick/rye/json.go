package rye

import (
	"encoding/json"
	"net/http"
)

func JSON(status int, body interface{}) http.Handler {
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	return JSONWithHeader(status, body, header)
}

func JSONWithHeader(status int, body interface{}, header http.Header) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setHeader(w, header)
		w.WriteHeader(status)
		if err := json.NewEncoder(w).Encode(body); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
