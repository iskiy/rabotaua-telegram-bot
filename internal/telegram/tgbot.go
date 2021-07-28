package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/iskiy/rabotaua-telegram-bot/internal/database"
	"github.com/iskiy/rabotaua-telegram-bot/pkg/rabotaua"
	"log"
)

type RabotaUABot struct {
	bot          *tgbotapi.BotAPI
	client       *rabotaua.RabotaClient
	db           *database.Storage
	schedulesMap rabotaua.ScheduleMap
}

func NewRabotaUABot(bot *tgbotapi.BotAPI, client *rabotaua.RabotaClient, storage *database.Storage, scheduleMap rabotaua.ScheduleMap) *RabotaUABot {
	return &RabotaUABot{bot: bot, client: client, db: storage, schedulesMap: scheduleMap}
}

func (b *RabotaUABot) Run() error {
	go b.manageSubscriptions()
	err := b.manageUpdates()
	if err != nil {
		return err
	}
	return nil
}

func (b *RabotaUABot) manageUpdates() error {
	schedulesText = b.getSchedulesText()
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := b.bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}
	for update := range updates {
		if update.CallbackQuery != nil {
			err := b.handleCallbackQuery(update.CallbackQuery)
			if err != nil {
				return err
			}
		}
		if update.Message != nil {
			log.Printf("[%s]: %s\n", update.Message.From.UserName, update.Message.Text)
			if update.Message.IsCommand() {
				go b.handleCommand(update.Message)
			} else {
				go b.handleMessage(update.Message)
			}
		}
	}
	return nil
}
