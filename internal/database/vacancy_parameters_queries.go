package database

import (
	"crawshaw.io/sqlite"
	"github.com/iskiy/rabotaua-telegram-bot/pkg/rabotaua"
)

type UserParameters struct {
	ID     int
	ChatID int64
	Page   int
	Params rabotaua.VacancyParameters
}

func (s *Storage) getIntValueFromQueryForChatID(chatID int64, query string) (int, error) {
	var res int
	fn := func(stmt *sqlite.Stmt) error {
		res = stmt.ColumnInt(0)
		return nil
	}
	err := s.execQuery(query, fn, chatID)
	if err != nil {
		return -1, err
	}
	return res, nil
}

func (s *Storage) InsertVacancyParameters(chatID int64, p rabotaua.VacancyParameters) (UserParameters, error) {
	quantityUserParameters, err := s.getQuantityParametersForUser(chatID)
	if err != nil {
		return UserParameters{}, err
	}
	quantityUserParameters++
	err = s.execQuery("INSERT INTO vacancy_parameters(parameters_id, chat_id, keywords, city_id, schedule_id) values(?, ?, ?, ?, ?);",
		nil, quantityUserParameters, chatID, p.Keywords, p.CityID, p.ScheduleID)
	if err != nil {
		return UserParameters{}, err
	}
	return UserParameters{ID: quantityUserParameters, ChatID: chatID, Params: p}, nil
}

func (s *Storage) UpdateLastInsertedParameterCityID(chatID int64, cityID int) error {
	parameterID, err := s.GetLastInsertedParameterIDForUser(chatID)
	if err != nil {
		return err
	}
	return s.UpdateParameterCityID(chatID, parameterID, cityID)
}

func (s *Storage) UpdateParameterCityID(chatID int64, parameterID int, cityID int) error {
	query := "UPDATE vacancy_parameters SET city_id = ? WHERE chat_id = ? AND parameters_id = ?;"
	return s.execQuery(query, nil, cityID, chatID, parameterID)
}

func (s *Storage) UpdateLastInsertedParameterScheduleID(chatID int64, scheduleID int) error {
	parameterID, err := s.GetLastInsertedParameterIDForUser(chatID)
	if err != nil {
		return err
	}
	return s.UpdateParameterScheduleID(chatID, parameterID, scheduleID)
}

func (s *Storage) UpdateParameterScheduleID(chatID int64, parameterID int, scheduleID int) error {
	query := "UPDATE vacancy_parameters SET schedule_id = ? WHERE chat_id = ? AND parameters_id = ?;"
	return s.execQuery(query, nil, scheduleID, chatID, parameterID)
}

func (s *Storage) DeleteLastInsertedParameter(chatID int64) error {
	lastParameterID, err := s.GetLastInsertedParameterIDForUser(chatID)
	if err != nil {
		return err
	}
	query := "DELETE FROM vacancy_parameters WHERE chat_id = ? AND parameters_id = ?;"
	return s.execQuery(query, nil, chatID, lastParameterID)
}

func (s *Storage) getQuantityParametersForUser(chatID int64) (int, error) {
	query := "SELECT COUNT(*) FROM vacancy_parameters WHERE chat_id = ?"
	return s.getIntValueFromQueryForChatID(chatID, query)
}

func (s *Storage) GetLastInsertedParameterIDForUser(chatID int64) (int, error) {
	query := "SELECT MAX(parameters_id) FROM vacancy_parameters WHERE chat_id = ?"
	return s.getIntValueFromQueryForChatID(chatID, query)
}

func (s *Storage) GetUserParameters(chatID int64) ([]UserParameters, error) {
	query := "SELECT parameters_id, keywords, city_id, schedule_id FROM vacancy_parameters WHERE chat_id = ? ORDER BY parameters_id ASC"
	return s.getUserParametersFromQueryAndChatID(chatID, query)
}

func (s *Storage) GetUserParametersInSubs(chatID int64) ([]UserParameters, error) {
	query := `SELECT subscriptions.parameters_id, keywords, city_id, schedule_id FROM subscriptions INNER JOIN vacancy_parameters v  ON subscriptions.parameters_id = v.parameters_id 
			WHERE subscriptions.chat_id = ? ORDER BY subscriptions.parameters_id ASC`
	return s.getUserParametersFromQueryAndChatID(chatID, query)
}

func (s *Storage) getUserParametersFromQueryAndChatID(chatID int64, query string) ([]UserParameters, error) {
	var parametersIDs []int
	var keywords []string
	var cityIDs []int
	var schedulesIDs []int
	fn := func(stmt *sqlite.Stmt) error {
		parametersIDs = append(parametersIDs, stmt.ColumnInt(0))
		keywords = append(keywords, stmt.ColumnText(1))
		cityIDs = append(cityIDs, stmt.ColumnInt(2))
		schedulesIDs = append(schedulesIDs, stmt.ColumnInt(3))
		return nil
	}
	err := s.execQuery(query, fn, chatID)
	if err != nil {
		return []UserParameters{}, err
	}
	return generateUserParametersFromArgs(chatID, parametersIDs, keywords, cityIDs, schedulesIDs), nil

}

func generateUserParametersFromArgs(chatID int64, parametersIDs []int, keywords []string, cityIDs []int, schedulesIDs []int) []UserParameters {
	resLen := len(parametersIDs)
	res := make([]UserParameters, resLen)
	var vacancyParams rabotaua.VacancyParameters
	for i := 0; i < resLen; i++ {
		vacancyParams = rabotaua.VacancyParameters{Keywords: keywords[i], CityID: cityIDs[i], ScheduleID: schedulesIDs[i]}
		res[i] = UserParameters{ChatID: chatID, ID: parametersIDs[i], Params: vacancyParams}
	}
	return res
}

func (s *Storage) GetParameter(chatID int64, parametersID int) (rabotaua.VacancyParameters, error) {
	var keywords string
	var cityID int
	var schedulesID int
	fn := func(stmt *sqlite.Stmt) error {
		keywords = stmt.ColumnText(0)
		cityID = stmt.ColumnInt(1)
		schedulesID = stmt.ColumnInt(2)
		return nil
	}
	query := "SELECT keywords, city_id, schedule_id FROM vacancy_parameters WHERE chat_id = ? AND parameters_id = ?"
	err := s.execQuery(query, fn, chatID, parametersID)
	if err != nil {
		return rabotaua.VacancyParameters{}, err
	}
	return rabotaua.VacancyParameters{Keywords: keywords, CityID: cityID, ScheduleID: schedulesID}, nil
}

func (s *Storage) DeleteParameter(chatID int64, parametersID int) error {
	query := "DELETE FROM vacancy_parameters WHERE chat_ID = ? AND parameters_id = ?"
	return s.execQuery(query, nil, chatID, parametersID)
}
