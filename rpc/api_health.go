package rpc

import (
	"net/http"
)

func getHealth(r *http.Request) (resp interface{}, err error) {
	return "ok", nil
}
