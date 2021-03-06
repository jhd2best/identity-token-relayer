package rpc

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"identity-token-relayer/config"
	"identity-token-relayer/log"
	"net/http"
)

var srv *http.Server

func StartAndServe() error {
	srv = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Get().RPC.Listen, config.Get().RPC.Port),
		Handler: registerRouter(),
	}

	log.GetLogger().Info("rpc server start", zap.String("address", srv.Addr))

	return srv.ListenAndServe()
}

func Stop(ctx context.Context) {
	_ = srv.Shutdown(ctx)
}
