package database

import (
	"github.com/iskiy/rabotaua-telegram-bot/pkg/rabotaua"
	"testing"
	"time"
)

const memoryStorePath = "file::memory:?cache=shared"

var testStorage *Storage

func TestInitStorage(t *testing.T) {
	storage, err := NewStorage(memoryStorePath)
	if err != nil {
		t.Errorf(err.Error())
	}
	testStorage = storage
}

func TestStorage_VacancyParameters(t *testing.T) {
	var chatID int64 = 1
	params := rabotaua.VacancyParameters{Keywords: "java", CityID: 1, ScheduleID: 1}
	_, err := testStorage.InsertVacancyParameters(chatID, params)
	if err != nil {
		t.Errorf(err.Error())
	}
	parameters, err := testStorage.InsertVacancyParameters(chatID, params)
	if err != nil {
		t.Errorf(err.Error())
	}
	if &parameters == nil || parameters.ID != 2 {
		t.Errorf("parameters result error")
	}
	quantity, err := testStorage.getQuantityParametersForUser(chatID)
	if err != nil {
		t.Errorf(err.Error())
	}
	if quantity != 2 {
		t.Errorf("wrong quantity result")
	}
	lastParameterID, err := testStorage.GetLastInsertedParameterIDForUser(chatID)
	if err != nil {
		t.Errorf(err.Error())
	}
	if lastParameterID != 2 {
		t.Errorf("wrong lastParameterID result")
	}
	userParams, err := testStorage.GetUserParameters(chatID)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(userParams) != 2 {
		t.Errorf("wrong len of userParams")
	}
}

func TestStorage_Updates(t *testing.T) {
	var chatID int64 = 1
	lastParameterID, err := testStorage.GetLastInsertedParameterIDForUser(chatID)
	if err != nil {
		t.Errorf(err.Error())
	}
	err = testStorage.UpdateLastInsertedParameterScheduleID(1, 5)
	if err != nil {
		t.Errorf(err.Error())
	}
	err = testStorage.UpdateLastInsertedParameterCityID(1, 100)
	if err != nil {
		t.Errorf(err.Error())
	}
	lastParameters, err := testStorage.GetParameter(1, lastParameterID)
	if lastParameters.CityID != 100 {
		t.Errorf("wrong updated cityID")
	}
	if lastParameters.ScheduleID != 5 {
		t.Errorf("wrong updated scheduleID")
	}
	parametersFromDB, err := testStorage.GetParameter(1, 1)
	if err != nil {
		t.Errorf(err.Error())
	}
	if parametersFromDB.CityID != 1 {
		t.Errorf("wrong updated cityID")
	}
	err = testStorage.DeleteLastInsertedParameter(chatID)
	if err != nil {
		t.Errorf(err.Error())
	}
	lastParameterID, err = testStorage.GetLastInsertedParameterIDForUser(chatID)
	if err != nil {
		t.Errorf(err.Error())
	}
	if lastParameterID != 1 {
		t.Errorf("wrong last parameter id after delete: %d", lastParameterID)
	}
	if err := testStorage.UpdateParameterCityID(chatID, lastParameterID, 54); err != nil {
		t.Errorf(err.Error())
	}
	if err := testStorage.UpdateParameterScheduleID(chatID, lastParameterID, 5); err != nil {
		t.Errorf(err.Error())
	}
}

func TestStorage_DeleteParameter(t *testing.T) {
	err := testStorage.DeleteParameter(1, 2)
	if err != nil {
		t.Errorf(err.Error())
	}
	parameters, err := testStorage.GetUserParameters(1)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(parameters) != 1 {
		t.Errorf(`wrong len parameters`)
	}
	if parameters[0].Params.Keywords != "java" {
		t.Errorf(`wrong cityID`)
	}
}

func TestStorage_InsertUser(t *testing.T) {
	var chatID int64 = 1
	state := 5
	err := testStorage.InsertUser(chatID, state)
	if err != nil {
		t.Errorf(err.Error())
	}
	insertedState, err := testStorage.GetUserState(chatID)
	if err != nil {
		t.Errorf(err.Error())
	}
	if insertedState != state {
		t.Errorf("%d != %d", insertedState, state)
	}
}

func TestStorage_UpdateUserState(t *testing.T) {
	var chatID int64 = 1
	newState := 1
	err := testStorage.UpdateUserState(chatID, newState)
	if err != nil {
		t.Errorf(err.Error())
	}
	stateFromDB, err := testStorage.GetUserState(chatID)
	if err != nil {
		t.Errorf(err.Error())
	}
	if stateFromDB != newState {
		t.Errorf("%d != %d", stateFromDB, newState)
	}
	isPresent, err := testStorage.IsUserPresent(chatID)
	if err != nil {
		t.Errorf(err.Error())
	}
	if !isPresent {
		t.Errorf("IsUserPresent error")
	}
}

func TestStorage_GetCityName(t *testing.T) {
	cityID := 1
	cityName, err := testStorage.GetCityName(cityID)
	if err != nil {
		t.Errorf(err.Error())
	}
	if cityName != "Київ" {
		t.Errorf("Wrong city name result: %s", cityName)
	}
}

func TestStorage_VacanciesView(t *testing.T) {
	var chatID int64 = 1
	parametersID := 1
	lastID, err := testStorage.InsertVacanciesView(chatID, parametersID, 0)
	if err != nil {
		t.Errorf(err.Error())
	}
	if lastID != 1 {
		t.Errorf("wrong vacancies_view returned ID")
	}
	params, err := testStorage.GetVacancyParametersPage(1, 5)
	if err != nil {
		t.Errorf(err.Error())
	}
	if params.Page != 0 {
		t.Errorf("wrong page")
	}
	testStorage.UpdateViewVacanciesPage(1, 2)
	params, err = testStorage.GetVacancyParametersPage(1, 5)
	if err != nil {
		t.Errorf(err.Error())
	}
	if params.Page != 2 {
		t.Errorf("wrong page")
	}
}

func TestStorage_InsertSubscription(t *testing.T) {
	var chatID int64 = 1
	parametersID := 1
	err := testStorage.InsertSubscription(chatID, parametersID, time.Now())
	if err != nil {
		t.Errorf(err.Error())
	}
	err = testStorage.InsertSubscription(chatID+1, parametersID, time.Now())
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestStorage_IsSubscriptionPresent(t *testing.T) {
	var chatID int64 = 1
	parametersID := 1
	isPresent, err := testStorage.IsSubscriptionPresent(chatID, parametersID)
	if err != nil {
		t.Errorf(err.Error())
	}
	if !isPresent {
		t.Errorf("subscription present in db, but return false")
	}
	isPresent, err = testStorage.IsSubscriptionPresent(5, 1)
	if err != nil {
		t.Errorf(err.Error())
	}
	if isPresent {
		t.Errorf("subscription not present in db, but return true")
	}
}

func TestStorage_GetSubscriptions(t *testing.T) {
	subs, err := testStorage.GetAllSubscriptions()
	var chatID int64 = 1
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(subs) != 2 {
		t.Errorf("wrong len of subs")
	}
	if subs[0].ChatID != 1 || subs[1].ChatID != 2 {
		t.Errorf("wrong subs chatIDs")
	}

	userSubs, err := testStorage.GetSubscriptionsFromUser(chatID)
	if err != nil {
		t.Errorf(err.Error())
	}
	if userSubs[0].ParametersID != subs[0].ParametersID {
		t.Errorf("expected parameter id: %d, received: %d", userSubs[0].ParametersID, subs[0].ParametersID)
	}
	userSubsParams, err := testStorage.GetUserParametersInSubs(chatID)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(userSubsParams) != 1 {
		t.Errorf("wrong len of user subs params, expected: %d, received: %d", 1, len(userSubsParams))
	}
	if userSubsParams[0].ID != userSubs[0].ParametersID {
		t.Errorf("expected parameter id: %d, received: %d", userSubs[0].ParametersID, userSubsParams[0].ID)
	}
}

func TestStorage_UpdateSubscriptionTime(t *testing.T) {
	updateTime := time.Now().Add(time.Hour * 1)
	err := testStorage.UpdateSubscriptionTime(1, 1, updateTime)
	if err != nil {
		t.Errorf(err.Error())
	}
	subs, err := testStorage.GetAllSubscriptions()
	if err != nil {
		t.Errorf(err.Error())
	}
	if subs[0].SubTime.Sub(updateTime) > time.Minute {
		t.Errorf("update time error")
	}
}

func TestStorage_DeleteSubscription(t *testing.T) {
	err := testStorage.DeleteSubscription(1, 1)
	if err != nil {
		t.Errorf(err.Error())
	}
	subs, err := testStorage.GetAllSubscriptions()
	if err != nil {
		t.Errorf(err.Error())
	}
	if subs[0].ChatID != 2 || len(subs) != 1 {
		t.Errorf("delete subs error")
	}
}
