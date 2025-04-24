package worker

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/rea1shane/telegram-bot-notifier/util/ip"
)

const ipWatcherWorkerName = "IP Watcher"

func newIPWatcher(logger *slog.Logger, botToken string, chatID any) (Worker, error) {
	b, err := bot.New(botToken)
	if err != nil {
		return nil, fmt.Errorf("failed to new bot: %w", err)
	}

	return &ipWatcher{
		logger:        logger,
		b:             b,
		chatID:        chatID,
		lastMessageID: -1,
	}, nil
}

type ipWatcher struct {
	logger *slog.Logger

	b      *bot.Bot
	chatID any

	lastIP        string
	lastMessageID int
}

func (w *ipWatcher) Start() error {
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()

		for {
			begin := time.Now()
			err := w.execute()
			duration := time.Since(begin)
			if err != nil {
				w.logger.Error("Execute failed", "duration_seconds", duration.Seconds(), "err", err)
			} else {
				w.logger.Debug("Execute succeeded", "duration_seconds", duration.Seconds())
			}
			<-ticker.C
		}
	}()
	return nil
}

func (w *ipWatcher) execute() error {
	currentIP, err := ip.Get()
	if err != nil {
		return fmt.Errorf("failed to get ip: %w", err)
	}

	// check if the IP has changed.
	// If it has changed, send a message and record it.
	if currentIP != w.lastIP {
		w.logger.Info("New IP detected", "original", w.lastIP, "new", currentIP)
		messageID, err := w.send(currentIP)
		if err != nil {
			return err
		} else {
			w.lastIP = currentIP
			w.lastMessageID = messageID
		}
	}

	return nil
}

// send a message
func (w *ipWatcher) send(addr string) (messageID int, err error) {
	// Request params
	p := &bot.SendMessageParams{
		ChatID:    w.chatID,
		ParseMode: models.ParseModeMarkdown,
		Text:      fmt.Sprintf("`Current IP: %s`", addr),
	}

	// Add reply
	if w.lastMessageID != -1 {
		p.ReplyParameters = &models.ReplyParameters{
			MessageID: w.lastMessageID,
			Quote:     w.lastIP,
		}
	}

	// Send
	message, err := w.b.SendMessage(context.Background(), p)
	if err != nil {
		return 0, fmt.Errorf("failed to send message: %w", err)
	}

	return message.ID, nil
}
