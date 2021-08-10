package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/iskiy/rabotaua-telegram-bot/pkg/rabotaua"
	"log"
	"strings"
)

func (b *RabotaUABot) handleAddParametersKeywordState(message *tgbotapi.Message) error {
	if message.Text == cancelButton.Text {
		return b.handleMainMenuCommand(message)
	}
	chatID := message.Chat.ID
	var keywords string
	if message.Text == whateverButton.Text {
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
	return b.sendMessage(chatID, cfg.Msg.EnterCity, parametersKeyboard)
}

func (b *RabotaUABot) handleAddParametersCityState(message *tgbotapi.Message) error {
	chatID := message.Chat.ID
	switch message.Text {
	case cancelButton.Text:
		err := b.db.DeleteLastInsertedParameter(chatID)
		if err != nil {
			log.Println(err.Error())
		}
		return b.handleMainMenuCommand(message)
	case whateverButton.Text:
		err := b.updateUserState(chatID, AddParametersScheduleState)
		if err != nil {
			return err
		}
		return b.stepToSchedulesState(message)
	default:
		city, err := b.client.GetCityFromName(message.Text)
		if err != nil {
			if err == rabotaua.CantFindCityError {
				return b.sendMessage(chatID, cfg.Msg.UndefinedCity, nil)
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
	return b.sendMessage(message.Chat.ID, cfg.Msg.EnterEmploymentTypeNum+schedulesText, parametersKeyboard)
}

func (b *RabotaUABot) handleAddParametersScheduleState(message *tgbotapi.Message) error {
	chatID := message.Chat.ID
	switch message.Text {
	case cancelButton.Text:
		err := b.db.DeleteLastInsertedParameter(chatID)
		if err != nil {
			log.Println(err)
		}
		fallthrough
	case whateverButton.Text:
		return b.handleMainMenuCommand(message)
	default:
		id, err := getIDFromUser(message.Text)
		if err != nil || id > len(b.schedulesMap) || id <= 0 {
			return b.sendMessage(chatID, cfg.Msg.WrongEmploymentTypeNum, nil)
		}
		if err := b.db.UpdateLastInsertedParameterScheduleID(chatID, id); err != nil {
			return err
		}
		if err := b.sendMessage(chatID, cfg.Msg.AddParamsSuccess, nil); err != nil {
			return err
		}
		return b.handleMainMenuCommand(message)
	}
}

func (b *RabotaUABot) handleDeleteParametersButton(message *tgbotapi.Message) error {
	return b.createParametersMenu(message, DeleteParametersState)
}

func (b *RabotaUABot) handleDeleteParametersState(message *tgbotapi.Message) error {
	chatID := message.Chat.ID
	if message.Text == cancelButton.Text {
		return b.handleMainMenuCommand(message)
	}
	parameterID, _, err := b.getParameterIDFromUser(message)
	if err != nil {
		return err
	}
	err = b.db.DeleteParameter(chatID, parameterID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	err = b.sendMessage(chatID, cfg.Msg.DeleteParameterSuccess, nil)
	if err != nil {
		return err
	}
	return b.handleMainMenuCommand(message)
}
