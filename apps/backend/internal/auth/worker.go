package auth

import (
	"context"
	"log/slog"
	"time"
)

type CleanupWorker struct {
	commandRepo CommandRepository
	interval    time.Duration
	stopChan    chan struct{}
}

func NewCleanupWorker(commandRepo CommandRepository, interval time.Duration) *CleanupWorker {
	return &CleanupWorker{
		commandRepo: commandRepo,
		interval:    interval,
		stopChan:    make(chan struct{}),
	}
}

func (w *CleanupWorker) Start() {
	ticker := time.NewTicker(w.interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				w.cleanup(context.Background())
			case <-w.stopChan:
				ticker.Stop()
				slog.Info("Token cleanup worker stopped")
				return
			}
		}
	}()
	slog.Info("Token cleanup worker started", "interval", w.interval)
}

func (w *CleanupWorker) Stop() {
	close(w.stopChan)
}

func (w *CleanupWorker) cleanup(ctx context.Context) {
	err := w.commandRepo.DeleteExpiredBlacklistedTokens(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to cleanup expired blacklisted tokens", "error", err)
	} else {
		slog.DebugContext(ctx, "Expired blacklisted tokens cleaned up successfully")
	}
}
