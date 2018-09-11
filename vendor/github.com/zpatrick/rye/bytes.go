package rye

import (
	"net/http"
)

func Bytes(status int, body []byte) http.Handler {
	return BytesWithHeader(status, body, http.Header{})
}

func BytesWithHeader(status int, body []byte, header http.Header) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setHeader(w, header)
		w.WriteHeader(status)
		if _, err := w.Write(body); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
