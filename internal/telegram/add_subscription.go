package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	database "github.com/iskiy/rabotaua-telegram-bot/internal/database"
	"github.com/iskiy/rabotaua-telegram-bot/pkg/rabotaua"
	"log"
	"sync"
	"time"
)

func (b *RabotaUABot) handleAddSubscriptionsState(message *tgbotapi.Message) error {
	if message.Text == cancelButton.Text {
		return b.handleMainMenuCommand(message)
	}
	chatID := message.Chat.ID
	parametersID, _, err := b.getParameterIDFromMessage(message)
	if err != nil {
		return err
	}
	isPresent, err := b.db.IsSubscriptionPresent(chatID, parametersID)
	if isPresent {
		return b.sendMessage(chatID, "Ти вже підписаний на цей параметр, вибери якийсь інший", cancelKeyboard)
	}
	err = b.db.InsertSubscription(chatID, parametersID, message.Time())
	if err != nil {
		return err
	}
	err = b.sendMessage(chatID, "Чудово, тепер якщо з'являться нові вакансії для параметра, я тобі їх надішлю", nil)
	if err != nil {
		return err
	}
	return b.handleMainMenuCommand(message)
}

func (b *RabotaUABot) handleAddSubscriptionsButton(message *tgbotapi.Message) error {
	return b.createParametersMenu(message, AddSubscriptionState)
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
	err = b.updateUserState(chatID, newUserState)
	if err != nil {
		return err
	}
	text, err := b.generateUserParametersTextMenu(userParams)
	if err != nil {
		return err
	}
	return b.sendMessage(chatID, text, cancelKeyboard)
}

func (b *RabotaUABot) manageSubscriptions() {
	wg := sync.WaitGroup{}
	for {
		subs, err := b.db.GetSubscriptions()
		if err != nil {
			log.Println(err.Error())
			return
		}
		for _, s := range subs {
			wg.Add(1)
			go b.sendSubscriptionMassage(&wg, s)
		}
		wg.Wait()
		time.Sleep(time.Second * 30)
	}
}

func (b *RabotaUABot) sendSubscriptionMassage(wg *sync.WaitGroup, s database.Subscription) {
	defer wg.Done()
	log.Println("search vacancy for sub")
	vacancyParameters, err := b.db.GetParameter(s.ChatID, s.ParametersID)
	if err != nil {
		log.Println(err)
		return
	}
	paramsWithCount := rabotaua.VacancyParametersPage{VacancyParameters: vacancyParameters, Count: 50}
	searchResult, err := b.client.GetSearchResultFromParametersPage(paramsWithCount)
	if err != nil {
		log.Println(err)
		return
	}
	updateTime := time.Now()
	res := getVacanciesAfterSubTime(searchResult.Vacancy, s.SubTime)
	if len(res) == 0 {
		log.Println("len res = 0")
		return
	}
	vacancyParametersStr, err := b.getVacancyParametersString(vacancyParameters)
	if err != nil {
		log.Printf("getVacancyParametersString error %s\n", err.Error())
		return
	}
	for i := 0; i < len(res); i += 5 {
		var msgText string
		to := i + 5
		if to > len(res) {
			to = len(res)
		}
		if i == 0 {
			msgText += fmt.Sprintf("Я знайшов для тебе нові вакансії за твоїми підпискою \n %s: \n", vacancyParametersStr)
		}
		msgText += getVacanciesString(res[i:to])
		err = b.sendMessage(s.ChatID, msgText, nil)
		if err != nil {
			log.Printf("send message error %s\n", err.Error())
			continue
		}
	}
	err = b.db.UpdateSubscriptionTime(s.ChatID, s.ParametersID, updateTime)
	if err != nil {
		log.Printf("UpdateSubscriptionTime error, %s\n", err.Error())
		return
	}
}

func (b *RabotaUABot) getVacancyParametersString(p rabotaua.VacancyParameters) (string, error) {
	cityName, err := b.db.GetCityName(p.CityID)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Посада: %s, місто: %s, вид зайнятості: %s", p.Keywords, cityName, b.schedulesMap[p.ScheduleID]), nil
}

func getVacanciesAfterSubTime(vacancy []rabotaua.Vacancy, subTime time.Time) []rabotaua.Vacancy {
	var res []rabotaua.Vacancy
	for _, v := range vacancy {
		vacancyTime, err := stringToTime(v.Date)
		if err != nil {
			continue
		}
		trueTime := subTime.Add(3 * time.Hour)
		subs := vacancyTime.Sub(trueTime)
		log.Printf("%s - %s =  %s \n", vacancyTime.String(), trueTime, subs.String())
		if subs > 0 {
			res = append(res, v)
		} else {
			break
		}
	}
	return res
}

func stringToTime(stringTime string) (time.Time, error) {
	res, err := time.Parse(time.RFC3339, stringTime+"Z")
	if err != nil {
		return time.Now(), err
	}
	return res, nil
}
