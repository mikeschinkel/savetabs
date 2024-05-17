package tasks

import (
	"context"
	"log/slog"
	"time"
)

var (
	// backgroundTick specifies the number of milliseconds between each background
	// task execution.
	// TODO: Make this a user-configurable setting
	backgroundTick = 1 * time.Minute
)

type Context = context.Context

func BackgroundTasks(ctx Context) {
	ticker := time.NewTicker(backgroundTick)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("Background task stopping...")
			return
		case <-ticker.C:
			err := run(ctx)
			if err != nil {
				slog.Error(err.Error())
			}
		}
	}
}
