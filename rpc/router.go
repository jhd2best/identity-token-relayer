package rpc

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func registerRouter() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/health", simpleWrap(getHealth)).Methods("GET")

	return handlers.LoggingHandler(os.Stdout, r)
}
