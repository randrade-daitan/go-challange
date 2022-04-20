package api

import "net/http"

type Api interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}
