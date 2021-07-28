package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/iskiy/rabotaua-telegram-bot/pkg/rabotaua"
	"strconv"
	"strings"
)

func (b *RabotaUABot) handleCallbackQuery(query *tgbotapi.CallbackQuery) error {
	direction, id, err := getInfoFromCallbackData(query.Data)
	if err != nil {
		return nil
	}
	params, err := b.db.GetVacancyParametersPage(id, count)
	if err != nil {
		return err
	}
	switch direction {
	case prevSign:
		if params.Page <= 0 {
			return nil
		}
		params.Page = params.Page - 1
		return b.updateViewPage(id, params, query)
	case nextSign:
		params.Page = params.Page + 1
		return b.updateViewPage(id, params, query)
	default:
		return nil
	}
}

func (b *RabotaUABot) updateViewPage(viewID int, params rabotaua.VacancyParametersPage, query *tgbotapi.CallbackQuery) error {
	searchResult, err := b.client.GetSearchResultFromParametersPage(params)
	if err != nil {
		return err
	}
	if searchResult.Total <= count*params.Page {
		return nil
	}
	err = b.db.UpdateViewVacanciesPage(viewID, params.Page)
	if err != nil {
		return err
	}
	_, err = b.bot.AnswerCallbackQuery(tgbotapi.NewCallback(query.ID, strconv.Itoa(params.Page+1)))
	if err != nil {
		return err
	}
	return b.updateMessage(viewID, getVacanciesString(searchResult.Vacancy), query)
}

func (b *RabotaUABot) updateMessage(viewID int, text string, query *tgbotapi.CallbackQuery) error {
	chatID := query.Message.Chat.ID
	msgID := query.Message.MessageID
	msg := tgbotapi.NewEditMessageText(chatID, msgID, text)
	msgMark := tgbotapi.NewEditMessageReplyMarkup(chatID, msgID, generateVacancyViewKeyBoard(strconv.Itoa(viewID)))
	msg.ReplyMarkup = msgMark.ReplyMarkup
	msg.ParseMode = "html"
	msg.DisableWebPagePreview = true
	_, err := b.bot.Send(msg)
	return err
}

func getInfoFromCallbackData(data string) (string, int, error) {
	splitData := strings.Split(data, " ")
	if len(splitData) != 2 {
		return "", -1, nil
	}
	direction := splitData[0]
	id, err := strconv.Atoi(splitData[1])
	if err != nil {
		return "", -1, nil
	}
	return direction, id, nil
}
