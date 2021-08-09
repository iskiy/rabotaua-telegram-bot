package database

import (
	"crawshaw.io/sqlite"
	"fmt"
	"github.com/iskiy/rabotaua-telegram-bot/pkg/rabotaua"
)

func (s *Storage) InsertVacanciesView(chatID int64, parametersID int, page int) (int, error) {
	query := "INSERT INTO vacancies_view(parameters_id, chat_id, page) VALUES(?, ?, ?)"
	err := s.execQuery(query, nil, parametersID, chatID, page)
	if err != nil {
		return -1, err
	}
	return s.GetLastInsertedVacancyViewID(chatID, parametersID)
}

func (s *Storage) GetLastInsertedVacancyViewID(chatID int64, parametersID int) (int, error) {
	query := "SELECT MAX(view_id) FROM vacancies_view WHERE chat_id = ? AND parameters_id = ?"
	var res int
	fn := func(stmt *sqlite.Stmt) error {
		res = stmt.ColumnInt(0)
		return nil
	}
	err := s.execQuery(query, fn, chatID, parametersID)
	if err != nil {
		return -1, err
	}
	return res, nil
}

func (s *Storage) GetVacancyParametersPage(vacancyViewID int, count int) (rabotaua.VacancyParametersPage, error) {
	// TODO: COUNT?
	var parametersID = -1
	var keywords string
	var cityIDs int
	var schedulesID int
	var page int
	fn := func(stmt *sqlite.Stmt) error {
		parametersID = stmt.ColumnInt(0)
		keywords = stmt.ColumnText(1)
		cityIDs = stmt.ColumnInt(2)
		schedulesID = stmt.ColumnInt(3)
		page = stmt.ColumnInt(4)
		return nil
	}

	query := "SELECT vv.parameters_id, keywords, city_id, schedule_id, page " +
		"FROM vacancies_view vv INNER JOIN vacancy_parameters p ON vv.chat_id = p.chat_id AND vv.parameters_id = p.parameters_id " +
		"WHERE view_id = ?"
	err := s.execQuery(query, fn, vacancyViewID)
	if err != nil {
		return rabotaua.VacancyParametersPage{}, err
	}
	if parametersID < 0 {
		return rabotaua.VacancyParametersPage{}, fmt.Errorf("can`t find view for %d viewId", vacancyViewID)
	}
	params := rabotaua.VacancyParameters{Keywords: keywords, CityID: cityIDs, ScheduleID: schedulesID}
	return rabotaua.VacancyParametersPage{params, page, count}, nil
}

func (s *Storage) UpdateViewVacanciesPage(viewID int, page int) error {
	query := "UPDATE vacancies_view SET page = ? WHERE view_id = ?"
	return s.execQuery(query, nil, page, viewID)
}
