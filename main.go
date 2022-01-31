package main

import (
	"context"
	"go.uber.org/zap"
	"identity-token-relayer/config"
	"identity-token-relayer/cron"
	"identity-token-relayer/log"
	"identity-token-relayer/model"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	go start()

	// gracefully shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	log.GetLogger().Info("gracefully shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		<-ctx.Done()
		time.Sleep(time.Second)
		os.Exit(0)
	}()

	wg := sync.WaitGroup{}
	groupWait := func(cbs ...func()) {
		for _, cb := range cbs {
			curFunc := cb
			wg.Add(1)
			go func() {
				curFunc()
				wg.Done()
			}()
		}

		wg.Wait()
	}

	groupWait(
		func() {
			cron.StopCron(ctx)
			log.GetLogger().Info("cron stopped")
		},
		func() {
			_ = log.GetLogger().Sync()
		},
		func() {
			log.FlushSentry(ctx)
			log.GetLogger().Info("sentry flushed")
		},
		func() {
			_ = model.CloseDb()
			log.GetLogger().Info("db connection closed")
		},
	)

	log.GetLogger().Info("gracefully shutdown success")
}

func start() {
	// init db
	dbErr := model.InitDb()
	if dbErr != nil {
		log.GetLogger().Fatal("init firebase db failed.", zap.String("error", dbErr.Error()))
	}
	log.GetLogger().Info("init firebase db success")

	// init sentry
	if !config.Get().Debug.DisableSentry {
		log.InitSentry()
	}

	// init cron
	if !config.Get().Debug.DisableCron {
		cron.InitCron()
	}

	// sync projects
	syncErr := model.SyncAllProjects()
	if syncErr != nil {
		log.GetLogger().Error("sync projects failed.", zap.String("error", syncErr.Error()))
	}

	time.Sleep(1*time.Second)
	cron.GetEthTransaction()
}
