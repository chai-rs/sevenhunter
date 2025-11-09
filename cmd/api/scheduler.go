package main

import (
	"context"

	"github.com/chai-rs/sevenhunter/internal/repo"
	"github.com/chai-rs/sevenhunter/internal/scheduler"
	logx "github.com/chai-rs/sevenhunter/pkg/logger"
	"github.com/go-co-op/gocron/v2"
)

func startScheduler() (shutdown func() error) {
	sch, err := gocron.NewScheduler()
	if err != nil {
		logx.Fatal().Err(err).Msg("failed to create scheduler")
	}

	// Bind schedulers
	bindUserCountScheduler(sch)

	sch.Start()
	return sch.Shutdown
}

func bindUserCountScheduler(sch gocron.Scheduler) {
	userCountSchduler := scheduler.NewUserCountScheduler(scheduler.UserCountSchedulerOpts{
		UserRepo: repo.NewUserRepo(registry.MongoDB.Database("sevenhunter")),
	})

	_, err := sch.NewJob(
		gocron.CronJob(conf.Scheduler.UserCount, true),
		gocron.NewTask(userCountSchduler.Run, context.Background()),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
		gocron.JobOption(gocron.WithStartImmediately()),
	)

	if err != nil {
		logx.Panic().Err(err).Msg("failed to create user-count job")
	}
}
