package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/iskiy/rabotaua-telegram-bot/pkg/rabotaua"
	"log"
	"strconv"
)

func (b *RabotaUABot) handleViewVacanciesButton(message *tgbotapi.Message) error {
	return b.createParametersMenu(message, ViewVacanciesState)
}

func (b *RabotaUABot) handleViewVacanciesState(message *tgbotapi.Message) error {
	chatID := message.Chat.ID
	if message.Text == cancelButton.Text {
		return b.handleMainMenuCommand(message)
	}
	parameterID, parameters, err := b.getParameterIDFromUser(message)
	log.Printf("Param ID: %d\n", parameterID)
	viewID, err := b.db.InsertVacanciesView(chatID, parameterID, 0)
	if err != nil {
		return err
	}
	idFromUser, err := getIDFromUser(message.Text)
	if err != nil {
		return err
	}
	parametersPage := rabotaua.VacancyParametersPage{VacancyParameters: parameters[idFromUser-1].Params, Count: count}
	searchResult, err := b.client.GetSearchResultFromParametersPage(parametersPage)
	if err != nil {
		return err
	}
	vacancies := searchResult.Vacancy
	if len(vacancies) == 0 {
		err = b.sendMessage(chatID, cfg.Msg.CantFindVacancies, nil)
		if err != nil {
			return err
		}
		return b.handleMainMenuCommand(message)
	}
	err = b.sendMessage(chatID, totalAmountMessage(searchResult.Total), nil)
	if err != nil {
		return err
	}
	err = b.sendMessage(chatID, getVacanciesString(searchResult.Vacancy), generateVacancyViewKeyBoard(strconv.Itoa(viewID)))
	if err != nil {
		return err
	}
	return b.handleMainMenuCommand(message)
}

func totalAmountMessage(total int) string {
	return fmt.Sprintf("Я знайшов %d вакансій, які задовільняють вибрані параметри:", total)
}

func getVacanciesString(vacancies []rabotaua.Vacancy) string {
	var res string
	for _, v := range vacancies {
		res += vacancyString(v)
	}
	return res
}

func vacancyString(v rabotaua.Vacancy) string {
	return fmt.Sprintf("<b>%s</b>, компанія: <b>%s</b>, \n <i>%s</i>\nМісто: %s\n%s\n%s\n\n",
		v.Name, v.CompanyName, v.DateTxt, v.CityName, v.ShortDescription, v.GetURL())
}

func generateVacancyViewKeyBoard(data string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(prev, prevSign+" "+data),
			tgbotapi.NewInlineKeyboardButtonData(next, nextSign+" "+data),
		),
	)
}
