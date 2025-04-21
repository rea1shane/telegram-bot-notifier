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

func newIPWatcher(logger *slog.Logger, botToken string, chatID any) (Worker, error) {
	b, err := bot.New(botToken)
	if err != nil {
		return nil, fmt.Errorf("failed to new bot: %w", err)
	}

	return &ipWatcher{
		logger: logger,
		b:      b,
		chatID: chatID,
	}, nil
}

type ipWatcher struct {
	logger *slog.Logger

	b      *bot.Bot
	chatID any

	lastIPInfo  ip.Info
	lastMessage *models.Message
}

func (w *ipWatcher) Start() error {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for {
			err := w.check()
			if err != nil {
				fmt.Printf("failed to check ip: %v\n", err)
			}
			<-ticker.C
		}
	}()
	return nil
}

// check if the IP has changed
func (w *ipWatcher) check() error {
	ipInfo, err := ip.Get()
	if err != nil {
		return fmt.Errorf("failed to get ip: %w", err)
	}

	if w.lastIPInfo.IP == "" || ipInfo.IP != w.lastIPInfo.IP {
		w.logger.Info("New IP detected", "original", w.lastIPInfo.IP, "new", ipInfo.IP)
		return w.update(ipInfo)
	}
	return nil
}

// update IP information and send message
func (w *ipWatcher) update(ipInfo ip.Info) error {
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
	*link = "https://maps.apple.com/?ll=29.8782,121.5494&z=15"
	p := &bot.SendMessageParams{
		ChatID:    w.chatID,
		ParseMode: parseMode,
		Text:      text,
		LinkPreviewOptions: &models.LinkPreviewOptions{
			URL: link,
		},
	}
	// Add reply
	if w.lastMessage != nil {
		p.ReplyParameters = &models.ReplyParameters{
			MessageID: w.lastMessage.ID,
		}
	}

	// Send
	message, err := w.b.SendMessage(context.Background(), p)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	// Update state
	w.lastIPInfo = ipInfo
	w.lastMessage = message

	return nil
}
