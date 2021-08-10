package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

const (
	nextSign = "r"
	prevSign = "l"
	next     = "➡️"
	prev     = "⬅️"
	count    = 3
)

// States
const (
	MainMenuState = iota
	ViewVacanciesState
	AddParametersKeywordsState
	AddParametersCityState
	AddParametersScheduleState
	DeleteParametersState
	AddSubscriptionState
	DeleteSubscriptionState
)

const (
	startCommand = "start"
	menuCommand  = "menu"
)

func (b *RabotaUABot) handleCommand(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	switch message.Command() {
	case startCommand:
		err := b.sendMessage(chatID, cfg.Msg.Start, nil)
		if err != nil {
			log.Println(err.Error())
			return
		}
		fallthrough
	case menuCommand:
		err := b.handleMainMenuCommand(message)
		if err != nil {
			log.Println(err.Error())
			return
		}
	default:
		err := b.sendMessage(chatID, cfg.Msg.UndefinedCommand, nil)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}
}

func (b *RabotaUABot) handleMainMenuCommand(message *tgbotapi.Message) error {
	err := b.updateUserState(message.Chat.ID, MainMenuState)
	if err != nil {
		return nil
	}
	return b.sendMessage(message.Chat.ID, cfg.Msg.MainMenu, mainMenuKeyboard)
}

func (b *RabotaUABot) handleMessage(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	state, err := b.db.GetUserState(chatID)
	if err != nil {
		err = b.handleMainMenuCommand(message)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}
	fmt.Printf("Message: %s, state: %d\n", message.Text, state)
	switch state {
	case MainMenuState:
		err = b.handleMainMenuState(message)
	case AddParametersKeywordsState:
		err = b.handleAddParametersKeywordState(message)
	case AddParametersCityState:
		err = b.handleAddParametersCityState(message)
	case AddParametersScheduleState:
		err = b.handleAddParametersScheduleState(message)
	case ViewVacanciesState:
		err = b.handleViewVacanciesState(message)
	case AddSubscriptionState:
		err = b.handleAddSubscriptionsState(message)
	case DeleteParametersState:
		err = b.handleDeleteParametersState(message)
	case DeleteSubscriptionState:
		err = b.handleDeleteSubscriptionsState(message)
	}
	if err != nil {
		log.Println(err.Error())
	}
}

func (b *RabotaUABot) handleMainMenuState(message *tgbotapi.Message) error {
	chatID := message.Chat.ID
	switch message.Text {
	case viewVacanciesButton.Text:
		return b.handleViewVacanciesButton(message)
	case addParametersButton.Text:
		err := b.db.UpdateUserState(chatID, AddParametersKeywordsState)
		if err != nil {
			return err
		}
		return b.sendMessage(chatID, cfg.Msg.EnterPostname, parametersKeyboard)
	case addSubscriptionsButton.Text:
		return b.handleAddSubscriptionsButton(message)
	case deleteParametersButton.Text:
		return b.handleDeleteParametersButton(message)
	case deleteSubscriptionsButton.Text:
		return b.handleDeleteSubscriptionsButton(message)
	default:
		return b.sendMessage(chatID, cfg.Msg.Undefined, nil)
	}
}

func (b *RabotaUABot) updateUserState(chatID int64, state int) error {
	isPresent, err := b.db.IsUserPresent(chatID)
	if err != nil {
		return err
	}
	if !isPresent {
		err = b.db.InsertUser(chatID, state)
		if err != nil {
			return err
		}
	} else {
		err = b.db.UpdateUserState(chatID, state)
		if err != nil {
			return err
		}
	}
	return nil
}
