package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/iskiy/rabotaua-telegram-bot/pkg/rabotaua"
	"strings"
)

func (b *RabotaUABot) handleAddParametersKeywordState(message *tgbotapi.Message) error {
	if message.Text == cancelButton.Text {
		return b.handleMainMenuCommand(message)
	}
	chatID := message.Chat.ID
	var keywords string
	if message.Text == cityButton.Text {
		keywords = ""
	} else {
		keywords = strings.TrimSpace(strings.ToLower(message.Text))
	}
	parameters := rabotaua.VacancyParameters{Keywords: keywords, CityID: 0, ScheduleID: 0}
	_, err := b.db.InsertVacancyParameters(message.Chat.ID, parameters)
	if err != nil {
		return err
	}
	err = b.updateUserState(chatID, AddParametersCityState)
	if err != nil {
		return err
	}
	return b.sendMessage(chatID, "Введи назву міста", cityKeyboard)
}

func (b *RabotaUABot) handleAddParametersCityState(message *tgbotapi.Message) error {
	chatID := message.Chat.ID
	switch message.Text {
	case cancelButton.Text:
		return b.handleMainMenuCommand(message)
	case cityButton.Text:
		err := b.updateUserState(chatID, AddParametersScheduleState)
		if err != nil {
			return err
		}
		return b.stepToSchedulesState(message)
	default:
		city, err := b.client.GetCityFromName(message.Text)
		if err != nil {
			if err == rabotaua.CantFindCityError {
				return b.sendMessage(chatID, "Не вдалось знайти це місто в базі даних, спробуй щось інше", nil)
			}
			return err
		}
		if err = b.db.UpdateLastInsertedParameterCityID(chatID, city.ID); err != nil {
			return err
		}
		return b.stepToSchedulesState(message)
	}
}

func (b *RabotaUABot) stepToSchedulesState(message *tgbotapi.Message) error {
	if err := b.updateUserState(message.Chat.ID, AddParametersScheduleState); err != nil {
		return err
	}
	return b.sendMessage(message.Chat.ID, "Тепер вибери номер виду зайнятості:\n"+schedulesText, cityKeyboard)
}

func (b *RabotaUABot) handleAddParametersScheduleState(message *tgbotapi.Message) error {
	chatID := message.Chat.ID
	if message.Text == cancelButton.Text || message.Text == cityButton.Text {
		return b.handleMainMenuCommand(message)
	} else {
		id, err := getIDFromTextMenu(message.Text)
		if err != nil || id > len(b.schedulesMap) || id <= 0 {
			return b.sendMessage(chatID, "Ти ввів неправильний номер виду зайнятості", nil)
		}
		if err := b.db.UpdateLastInsertedParameterScheduleID(chatID, id); err != nil {
			return err
		}
		if err := b.sendMessage(chatID, "Вітаю, параметри були успішно додані", nil); err != nil {
			return err
		}
		return b.handleMainMenuCommand(message)
	}
}
