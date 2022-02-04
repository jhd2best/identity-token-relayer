package log

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"identity-token-relayer/config"
	"path/filepath"
	"runtime"
	"time"
)

var logger *zap.Logger

func init() {
	var err error

	sentryHook := zap.Hooks(func(entry zapcore.Entry) error {
		if entry.Level == zapcore.ErrorLevel || entry.Level == zapcore.FatalLevel {
			defer sentry.Flush(2 * time.Second)
			sentry.CaptureMessage(fmt.Sprintf("%s\n%s", entry.Message, entry.Stack))
		}
		return nil
	})

	if config.Get().Debug.Verbose {
		developmentConfig := zap.NewDevelopmentConfig()

		if _, file, _, ok := runtime.Caller(0); ok {
			basePath := filepath.Dir(filepath.Dir(filepath.Dir(file))) + "/"
			developmentConfig.EncoderConfig.EncodeCaller = func(caller zapcore.EntryCaller, encoder zapcore.PrimitiveArrayEncoder) {
				rel, err := filepath.Rel(basePath, caller.File)
				if err != nil {
					encoder.AppendString(caller.FullPath())
				} else {
					encoder.AppendString(fmt.Sprintf("%s:%d", rel, caller.Line))
				}
			}
		}

		// log to file
		if config.Get().Debug.LogPath != "" {
			developmentConfig.OutputPaths = []string{
				config.Get().Debug.LogPath,
			}
		}

		logger, err = developmentConfig.Build(sentryHook)
		if err != nil {
			panic(err)
		}
	} else {
		prodConfig := zap.NewProductionConfig()

		// log to file
		if config.Get().Debug.LogPath != "" {
			prodConfig.OutputPaths = []string{
				config.Get().Debug.LogPath,
			}
		}

		logger, err = prodConfig.Build(sentryHook)
		if err != nil {
			panic(err)
		}
	}
}

func GetLogger() *zap.Logger {
	return logger
}
