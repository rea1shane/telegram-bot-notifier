package worker

import (
	"log/slog"
	"os"
	"testing"
)

var (
	botToken = os.Getenv("TELEGRAM_BOT_KEY")
	chatId   = os.Getenv("TELEGRAM_CHAT_ID")
)

func TestIpWatcher_Start(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	worker, err := newIPWatcher(logger.With("worker", "IP Watcher"), botToken, chatId)
	if err != nil {
		t.Fatalf("failed to new ip worker: %v", err)
	}
	err = worker.Start()
	if err != nil {
		t.Fatalf("failed to start ip worker: %v", err)
	}
	select {}
}
