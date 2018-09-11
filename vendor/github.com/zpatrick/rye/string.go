package rye

import (
	"net/http"
)

func String(status int, body string) http.Handler {
	return BytesWithHeader(status, []byte(body), http.Header{})
}

func StringWithHeader(status int, body string, header http.Header) http.Handler {
	return BytesWithHeader(status, []byte(body), header)
}
