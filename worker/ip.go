package worker

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
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

	lastIPInfo    ip.Info
	lastMessageID int
}

func (w *ipWatcher) Start() error {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
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
	ipInfo, err := ip.Get()
	if err != nil {
		return fmt.Errorf("failed to get ip: %w", err)
	}

	// check if the IP has changed.
	// If it has changed, send a message and record it.
	if ipInfo.IP != w.lastIPInfo.IP {
		w.logger.Info("New IP detected", "original", w.lastIPInfo.IP, "new", ipInfo.IP)
		messageID, err := w.send(ipInfo)
		if err != nil {
			return err
		} else {
			w.lastIPInfo = ipInfo
			w.lastMessageID = messageID
		}
	}

	return nil
}

// send a message
func (w *ipWatcher) send(ipInfo ip.Info) (messageID int, err error) {
	// Message content
	parseMode := models.ParseModeMarkdown
	var sb strings.Builder
	sb.WriteString("`")
	sb.WriteString(fmt.Sprintf("IP:       %s\n", ipInfo.IP))
	sb.WriteString(fmt.Sprintf("City:     %s\n", ipInfo.City))
	sb.WriteString(fmt.Sprintf("Region:   %s\n", ipInfo.Region))
	sb.WriteString(fmt.Sprintf("Country:  %s\n", ipInfo.Country))
	sb.WriteString(fmt.Sprintf("Loc:      %s\n", ipInfo.Loc))
	sb.WriteString(fmt.Sprintf("Org:      %s\n", ipInfo.Org))
	sb.WriteString(fmt.Sprintf("Postal:   %s\n", ipInfo.Postal))
	sb.WriteString(fmt.Sprintf("Timezone: %s\n", ipInfo.Timezone))
	sb.WriteString("`")
	text := sb.String()

	// Request params
	link := new(string)
	*link = fmt.Sprintf("https://maps.apple.com/?ll=%s&z=%f", ipInfo.Loc, 15.0) // Doc: https://developer.apple.com/library/archive/featuredarticles/iPhoneURLScheme_Reference/MapLinks/MapLinks.html
	p := &bot.SendMessageParams{
		ChatID:    w.chatID,
		ParseMode: parseMode,
		Text:      text,
		LinkPreviewOptions: &models.LinkPreviewOptions{
			URL: link,
		},
	}
	// Add reply
	if w.lastMessageID != -1 {
		p.ReplyParameters = &models.ReplyParameters{
			MessageID: w.lastMessageID,
		}
	}

	// Send
	message, err := w.b.SendMessage(context.Background(), p)
	if err != nil {
		return 0, fmt.Errorf("failed to send message: %w", err)
	}

	return message.ID, nil
}
