package cron

import (
	"context"
	"github.com/robfig/cron/v3"
)

var cronClient *cron.Cron

func InitCron() {
	cronClient = cron.New(cron.WithSeconds())

	// add job
	_, _ = cronClient.AddFunc("0 * * * * *", GetEnableProject)
	_, _ = cronClient.AddFunc("10 * * * * *", GetEthTransaction)
	_, _ = cronClient.AddFunc("20,50 * * * * *", SendMappingTransaction)
	_, _ = cronClient.AddFunc("30 * * * * *", CheckMappingTransaction)
	_, _ = cronClient.AddFunc("40 * * * * *", RetryErrorTransaction)

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
