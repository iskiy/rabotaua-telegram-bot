package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/iskiy/rabotaua-telegram-bot/internal/database"
	"github.com/iskiy/rabotaua-telegram-bot/pkg/rabotaua"
	"sort"
	"strconv"
	"strings"
)

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

func (b *RabotaUABot) getSubscriptionIDFromUser(message *tgbotapi.Message) (int, error) {
	subsParams, err := b.db.GetUserParametersInSubs(message.Chat.ID)
	if err != nil {
		return -1, err
	}
	id, _, err := b.getIDFromParameters("Введено неправильний номер підписки", subsParams, message)
	return id, err
}

func (b *RabotaUABot) getParameterIDFromUser(message *tgbotapi.Message) (int, []database.UserParameters, error) {
	parameters, err := b.db.GetUserParameters(message.Chat.ID)
	if err != nil {
		return -1, []database.UserParameters{}, err
	}
	return b.getIDFromParameters("Введено неправильний номер параметра", parameters, message)
}

func (b *RabotaUABot) getIDFromParameters(wrongNumText string, params []database.UserParameters, message *tgbotapi.Message) (int, []database.UserParameters, error) {
	IDFromUser, err := getIDFromUser(message.Text)
	if err != nil || IDFromUser > len(params) || IDFromUser <= 0 {
		return -1, []database.UserParameters{}, b.sendMessage(message.Chat.ID, wrongNumText, nil)
	}
	return params[IDFromUser-1].ID, params, nil
}

func (b *RabotaUABot) createParametersMenu(message *tgbotapi.Message, newUserState int) error {
	chatID := message.Chat.ID
	userParams, err := b.db.GetUserParameters(chatID)
	if err != nil {
		return err
	}
	if len(userParams) == 0 {
		return b.sendMessage(chatID, "Йойк, у тебе нема параметрів, додай їх будь ласка", nil)
	}
	return b.printParameters("Введи номер потрібного параметра\n\n", newUserState, userParams, message)
}

func (b *RabotaUABot) printParameters(prefix string, newUserState int, params []database.UserParameters, message *tgbotapi.Message) error {
	err := b.updateUserState(message.Chat.ID, newUserState)
	if err != nil {
		return err
	}
	text, err := b.generateUserParametersTextMenu(prefix, params)
	if err != nil {
		return err
	}
	return b.sendMessage(message.Chat.ID, text, cancelKeyboard)
}

func (b *RabotaUABot) getVacancyParametersString(p rabotaua.VacancyParameters) (string, error) {
	cityName, err := b.db.GetCityName(p.CityID)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Посада: %s, місто: %s, вид зайнятості: %s", p.Keywords, cityName, b.schedulesMap[p.ScheduleID]), nil
}

func getIDFromUser(text string) (int, error) {
	return strconv.Atoi(strings.Split(text, " ")[0])
}

func (b *RabotaUABot) generateUserParametersTextMenu(prefix string, p []database.UserParameters) (string, error) {
	paramsLen := len(p)
	var lines = make([]string, paramsLen)
	for i := 0; i < paramsLen; i++ {
		cityName, err := b.db.GetCityName(p[i].Params.CityID)
		if err != nil {
			return "", err
		}
		lines[i] = fmt.Sprintf("%d Посада: %s, місто: %s, вид зайнятості %s", i+1,
			p[i].Params.Keywords, cityName, b.schedulesMap[p[i].Params.ScheduleID])
	}
	sort.Strings(lines)
	paramsLines := strings.Join(lines, "\n")
	return prefix + paramsLines, nil
}

func (b *RabotaUABot) getSchedulesText() string {
	lines := make([]string, 0, 7)
	for k, v := range b.schedulesMap {
		lines = append(lines, fmt.Sprintf("%d %s \n", k, v))
	}
	sort.Strings(lines)
	return strings.Join(lines, "")
}
