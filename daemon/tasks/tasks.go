package tasks

import (
	"context"
	"log/slog"
	"time"
)

const (
	// BackgroundTick specifies the number of milliseconds between each background
	// task execution.
	// TODO: Change to 1 hour(?) after implementation and debugging
	// TODO: Make this a user-configurable setting
	//BackgroundTick = 1 * time.Hour
	//BackgroundTick = 10 * time.Second
	BackgroundTick = 1 * time.Minute
)

type Context = context.Context

type Runner interface {
	Run(Context) error
}

func BackgroundTask(ctx Context, r Runner) {
	ticker := time.NewTicker(BackgroundTick)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("Background task stopping...")
			return
		case <-ticker.C:
			err := r.Run(ctx)
			if err != nil {
				slog.Error(err.Error())
			}
		}
	}
}
