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
	worker, err := newIPWatcher(logger.With("worker", "IP"), botToken, chatId)
	if err != nil {
		panic(err)
	}
	err = worker.Start()
	if err != nil {
		panic(err)
	}
	select {}
}
