package database

import (
	"fmt"
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
	lastParameterID, err := testStorage.GetLastInsertedParameterForUser(chatID)
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
	err = testStorage.DeleteLastInsertedParameter(chatID)
	if err != nil {
		t.Errorf(err.Error())
	}
	lastParameterID, err = testStorage.GetLastInsertedParameterForUser(chatID)
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
	fmt.Println(params)
}

func TestStorage_InsertSubscription(t *testing.T) {
	var chatID int64 = 1
	parametersID := 1
	err := testStorage.InsertSubscription(chatID, parametersID, time.Now())
	if err != nil {
		t.Errorf(err.Error())
	}
}
