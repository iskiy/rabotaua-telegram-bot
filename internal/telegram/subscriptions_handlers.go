package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/iskiy/rabotaua-telegram-bot/internal/database"
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
	parametersID, _, err := b.getParameterIDFromUser(message)
	if err != nil {
		return err
	}
	isPresent, err := b.db.IsSubscriptionPresent(chatID, parametersID)
	if isPresent {
		return b.sendMessage(chatID, cfg.Msg.AlreadySub, cancelKeyboard)
	}
	err = b.db.InsertSubscription(chatID, parametersID, message.Time())
	if err != nil {
		return err
	}
	err = b.sendMessage(chatID, cfg.Msg.AddSubSuccess, nil)
	if err != nil {
		return err
	}
	return b.handleMainMenuCommand(message)
}

func (b *RabotaUABot) handleAddSubscriptionsButton(message *tgbotapi.Message) error {
	return b.createParametersMenu(message, AddSubscriptionState)
}

func (b *RabotaUABot) handleDeleteSubscriptionsButton(message *tgbotapi.Message) error {
	chatID := message.Chat.ID
	userParams, err := b.db.GetUserParametersInSubs(chatID)
	if err != nil {
		return err
	}
	if len(userParams) == 0 {
		return b.sendMessage(chatID, cfg.Msg.EmptySubs, nil)
	}
	return b.printParameters(cfg.Msg.EnterSubsNum, DeleteSubscriptionState, userParams, message)
}

func (b *RabotaUABot) handleDeleteSubscriptionsState(message *tgbotapi.Message) error {
	chatID := message.Chat.ID
	if message.Text == cancelButton.Text {
		return b.handleMainMenuCommand(message)
	}
	id, err := b.getSubscriptionIDFromUser(message)
	if err != nil {
		return err
	}
	err = b.db.DeleteSubscription(chatID, id)
	if err != nil {
		return err
	}
	err = b.sendMessage(chatID, cfg.Msg.DeleteSubSuccess, nil)
	if err != nil {
		return err
	}
	return b.handleMainMenuCommand(message)
}

func (b *RabotaUABot) manageSubscriptions() {
	wg := sync.WaitGroup{}
	for {
		subs, err := b.db.GetAllSubscriptions()
		if err != nil {
			log.Println(err.Error())
			return
		}
		for _, s := range subs {
			wg.Add(1)
			go b.sendSubscriptionMassage(&wg, s)
		}
		wg.Wait()
		time.Sleep(time.Minute * 2)
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
			msgText += fmt.Sprintf(cfg.Msg.FoundSub, vacancyParametersStr)
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
