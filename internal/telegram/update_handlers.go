package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	database "github.com/iskiy/rabotaua-telegram-bot/internal/database"
	"log"
	"sort"
	"strconv"
	"strings"
)

var (
	schedulesText          string
	viewVacanciesButton    = tgbotapi.NewKeyboardButton("Переглянути вакансії")
	addParametersButton    = tgbotapi.NewKeyboardButton("Додати параметри для пошуку")
	addSubscriptionsButton = tgbotapi.NewKeyboardButton("Додати підписку")
	mainMenuKeyboard       = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			viewVacanciesButton, addParametersButton,
		),
		tgbotapi.NewKeyboardButtonRow(
			addSubscriptionsButton,
		),
	)

	cancelButton   = tgbotapi.NewKeyboardButton("Відмінити")
	cancelKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(cancelButton),
	)

	cityButton   = tgbotapi.NewKeyboardButton("Та мені байдуже")
	cityKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(cityButton, cancelButton),
	)
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
	AddSubscriptionState
)

const (
	startCommand = "start"
	menuCommand  = "menu"
)

func (b *RabotaUABot) handleCommand(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	switch message.Command() {
	case startCommand:
		err := b.sendMessage(chatID, "Привіт, я допоможу тобі знайти роботу твоєї мрії", nil)
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
		err := b.sendMessage(chatID, "Я не знаю що ти від мене хочеш, я таке робити не вмію.", nil)
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
	return b.sendMessage(message.Chat.ID, "Ти в головному меню", mainMenuKeyboard)
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
		return b.sendMessage(chatID, "Введи назву посади", cityKeyboard)
	case addSubscriptionsButton.Text:
		return b.handleAddSubscriptionsButton(message)
	default:
		return b.sendMessage(chatID, "Я не знаю що ти від мене хочеш", nil)
	}
}

func getIDFromTextMenu(text string) (int, error) {
	return strconv.Atoi(strings.Split(text, " ")[0])
}

func (b *RabotaUABot) generateUserParametersTextMenu(p []database.UserParameters) (string, error) {
	paramsLen := len(p)
	var lines = make([]string, paramsLen)
	for i := 0; i < paramsLen; i++ {
		cityName, err := b.db.GetCityName(p[i].Params.CityID)
		if err != nil {
			return "", err
		}
		lines[i] = fmt.Sprintf("%d Посада: %s, місто: %s, вид зайнятості %s", p[i].ID,
			p[i].Params.Keywords, cityName, b.schedulesMap[p[i].Params.ScheduleID])
	}
	sort.Strings(lines)
	paramsLines := strings.Join(lines, "\n")
	return fmt.Sprintf("Введи номер параметру який ти хочеш використати\n\n") + paramsLines, nil
}

func (b *RabotaUABot) getSchedulesText() string {
	lines := make([]string, 0, 7)
	for k, v := range b.schedulesMap {
		lines = append(lines, fmt.Sprintf("%d %s \n", k, v))
	}
	sort.Strings(lines)
	return strings.Join(lines, "")
}

func (b *RabotaUABot) sendMessage(chatID int64, text string, replyMarkup interface{}) error {
	msg := tgbotapi.NewMessage(chatID, text)
	if replyMarkup != nil {
		msg.ReplyMarkup = replyMarkup
	}
	msg.DisableWebPagePreview = true
	msg.ParseMode = "html"
	_, err := b.bot.Send(msg)
	return err
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
