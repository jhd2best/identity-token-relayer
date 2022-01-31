package log

import (
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"
	"identity-token-relayer/config"
	"time"
)

func InitSentry() {
	debugConfig := config.Get().Debug
	err := sentry.Init(sentry.ClientOptions{
		// Either set your DSN here or set the SENTRY_DSN environment variable.
		Dsn: debugConfig.SentryDSN,
		// Either set environment and release here or set the SENTRY_ENVIRONMENT
		// and SENTRY_RELEASE environment variables.
		Environment: "",
		// Release name
		Release: "identity-token-relayer",
		// Enable printing of SDK debug messages.
		Debug: debugConfig.Verbose,
	})
	if err != nil {
		GetLogger().Fatal(fmt.Sprintf("sentry.Init: %s", err))
	}
	GetLogger().Info("load sentry succeed")
}

func FlushSentry(ctx context.Context) {
	if deadline, ok := ctx.Deadline(); ok {
		sentry.Flush(deadline.Sub(time.Now()))
	} else {
		sentry.Flush(3 * time.Second)
	}
}
