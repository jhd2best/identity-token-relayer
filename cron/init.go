package cron

import (
	"context"
	"github.com/robfig/cron/v3"
)

var cronClient *cron.Cron

func InitCron() {
	cronClient = cron.New(cron.WithSeconds())

	// migrate transaction
	// _, _ = cronClient.AddFunc("5,35 * * * * *", xchRechargeCheckHandler)

	cronClient.Start()
}

func StopCron(ctx context.Context) {
	if cronClient == nil {
		return
	}

	select {
	case <-cronClient.Stop().Done():
		return
	case <-ctx.Done():
		return
	}
}
