package database

import (
	"crawshaw.io/sqlite"
	"crawshaw.io/sqlite/sqlitex"
	"errors"
	"fmt"
	rabotaua "github.com/iskiy/rabotaua-telegram-bot/pkg/rabotaua"
	"time"
)

type Storage struct {
	dbPool *sqlitex.Pool
}

var StateGetError = errors.New("state get error")

type Subscription struct {
	ChatID       int64
	ParametersID int
	SubTime      time.Time
}

func NewStorage(storagePath string) (*Storage, error) {
	dbPool, err := sqlitex.Open(storagePath, 0, 64)
	if err != nil {
		return &Storage{}, err
	}
	err = initStorageTables(dbPool)
	if err != nil {
		return nil, err
	}
	return &Storage{dbPool}, nil
}

func initStorageTables(dbPool *sqlitex.Pool) error {
	err := execQueryWithoutParameters(dbPool, "CREATE TABLE IF NOT EXISTS users("+
		"chat_id INTEGER PRIMARY KEY,"+
		"user_state INTEGER NOT NULL)")
	if err != nil {
		return err
	}
	err = execQueryWithoutParameters(dbPool, "CREATE TABLE IF NOT EXISTS vacancy_parameters("+
		"parameters_id INTEGER NOT NULL,"+
		"chat_id INTEGER NOT NULL,"+
		"keywords VARCHAR,"+
		"city_id INTEGER,"+
		"schedule_id INTEGER,"+
		"PRIMARY KEY(parameters_id, chat_id)"+
		");")
	if err != nil {
		return err
	}
	err = execQueryWithoutParameters(dbPool, "CREATE TABLE IF NOT EXISTS vacancies_view("+
		"view_id INTEGER PRIMARY KEY AUTOINCREMENT,"+
		"parameters_id INTEGER NOT NULL,"+
		"chat_id INTEGER NOT NULL,"+
		"page INTEGER NOT NULL,"+
		"FOREIGN KEY (parameters_id)"+
		"REFERENCES vacancy_parameters(parameters_id) ON UPDATE CASCADE ON DELETE CASCADE, "+
		"FOREIGN KEY(chat_id) REFERENCES vacancy_parameters(chat_id) ON UPDATE CASCADE ON DELETE CASCADE);")
	if err != nil {
		return err
	}
	err = execQueryWithoutParameters(dbPool, "CREATE TABLE IF NOT EXISTS cities("+
		"city_id INTEGER PRIMARY KEY,"+
		"city_name varchar NOT NULL);")
	if err != nil {
		return err
	}
	err = execQueryWithoutParameters(dbPool, "CREATE TABLE IF NOT EXISTS "+
		"subscriptions(chat_id INTEGER NOT NULL,"+
		"parameters_id INTEGER NOT NULL, "+
		"sub_time VARCHAR NOT NULL, "+
		"PRIMARY KEY(chat_id, parameters_id), "+
		"FOREIGN KEY (parameters_id) REFERENCES vacancy_parameters(parameters_id) ON UPDATE CASCADE ON DELETE CASCADE, "+
		"FOREIGN KEY(chat_id) REFERENCES vacancy_parameters(chat_id) ON UPDATE CASCADE ON DELETE CASCADE);")
	return insertCityDictionary(dbPool)
}

func insertCityDictionary(dbPool *sqlitex.Pool) error {
	conn := dbPool.Get(nil)
	defer dbPool.Put(conn)
	client := rabotaua.NewRabotaClient()
	cities, err := client.GetCities()
	if err != nil {
		return err
	}
	query := "INSERT INTO cities(city_id, city_name) VALUES(?, ?)"
	for _, c := range cities {
		if err := sqlitex.Exec(conn, query, nil, c.ID, c.City); err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) GetCityName(cityID int) (string, error) {
	conn := s.dbPool.Get(nil)
	defer s.dbPool.Put(conn)
	var res string
	fn := func(stmt *sqlite.Stmt) error {
		res = stmt.ColumnText(0)
		return nil
	}
	query := "SELECT city_name FROM cities WHERE city_id = ?"
	err := sqlitex.Exec(conn, query, fn, cityID)
	if err != nil {
		return "", err
	}
	return res, nil
}

func execQueryWithoutParameters(dbPool *sqlitex.Pool, query string) error {
	conn := dbPool.Get(nil)
	defer dbPool.Put(conn)
	err := sqlitex.Exec(conn, query, nil)
	if err != nil {
		return err
	}
	return err
}

func (s *Storage) getIntValueFromQueryForChatID(chatID int64, query string) (int, error) {
	conn := s.dbPool.Get(nil)
	defer s.dbPool.Put(conn)
	var res int
	fn := func(stmt *sqlite.Stmt) error {
		res = stmt.ColumnInt(0)
		return nil
	}
	err := sqlitex.Exec(conn, query, fn, chatID)
	if err != nil {
		return -1, err
	}
	return res, nil
}

func (s *Storage) InsertVacanciesView(chatID int64, parametersID int, page int) (int, error) {
	conn := s.dbPool.Get(nil)
	defer s.dbPool.Put(conn)
	query := "INSERT INTO vacancies_view(parameters_id, chat_id, page) VALUES(?, ?, ?)"
	err := sqlitex.Exec(conn, query, nil, parametersID, chatID, page)
	if err != nil {
		return -1, err
	}
	return s.GetLastInsertedVacancyViewID(chatID, parametersID)
}

func (s *Storage) GetLastInsertedVacancyViewID(chatID int64, parametersID int) (int, error) {
	conn := s.dbPool.Get(nil)
	defer s.dbPool.Put(conn)
	query := "SELECT MAX(view_id) FROM vacancies_view WHERE chat_id = ? AND parameters_id = ?"
	var res int
	fn := func(stmt *sqlite.Stmt) error {
		res = stmt.ColumnInt(0)
		return nil
	}
	err := sqlitex.Exec(conn, query, fn, chatID, parametersID)
	if err != nil {
		return -1, err
	}
	return res, nil
}

func (s *Storage) GetVacancyParametersPage(vacancyViewID int, count int) (rabotaua.VacancyParametersPage, error) {
	conn := s.dbPool.Get(nil)
	defer s.dbPool.Put(conn)

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
	err := sqlitex.Exec(conn, query, fn, vacancyViewID)
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
	conn := s.dbPool.Get(nil)
	defer s.dbPool.Put(conn)
	query := "UPDATE vacancies_view SET page = ? WHERE view_id = ?"
	return sqlitex.Exec(conn, query, nil, page, viewID)
}

func (s *Storage) InsertSubscription(chatID int64, parametersID int, subTime time.Time) error {
	conn := s.dbPool.Get(nil)
	defer s.dbPool.Put(conn)
	return sqlitex.Exec(conn, "INSERT INTO subscriptions(chat_id, parameters_id, sub_time) values(?, ?, ?)",
		nil, chatID, parametersID, subTime.Format(time.RFC3339))
}

func (s *Storage) IsSubscriptionPresent(chatID int64, parametersID int) (bool, error) {
	conn := s.dbPool.Get(nil)
	defer s.dbPool.Put(conn)
	var ids int64 = -1
	fn := func(stmt *sqlite.Stmt) error {
		ids = stmt.ColumnInt64(0)
		return nil
	}
	err := sqlitex.Exec(conn, "SELECT chat_id FROM subscriptions WHERE chat_id = ? AND parameters_id = ? ",
		fn, chatID, parametersID)
	if err != nil {
		return false, err
	}
	if ids > 0 {
		return true, nil
	}
	return false, nil
}

func (s *Storage) GetSubscriptions() ([]Subscription, error) {
	conn := s.dbPool.Get(nil)
	defer s.dbPool.Put(conn)
	var chatIDs []int64
	var parametersIDs []int
	var subTimes []string
	fn := func(stmt *sqlite.Stmt) error {
		chatIDs = append(chatIDs, stmt.ColumnInt64(0))
		parametersIDs = append(parametersIDs, stmt.ColumnInt(1))
		subTimes = append(subTimes, stmt.ColumnText(2))
		return nil
	}
	query := "SELECT chat_id, parameters_id, sub_time FROM subscriptions"
	err := sqlitex.Exec(conn, query, fn)
	if err != nil {
		return []Subscription{}, err
	}
	return generateSubscriptionsFromArgs(chatIDs, parametersIDs, subTimes), nil
}

func generateSubscriptionsFromArgs(chatIDs []int64, parametersIDs []int, subTimes []string) []Subscription {
	resLen := len(parametersIDs)
	res := make([]Subscription, resLen)
	for i := 0; i < resLen; i++ {
		parsedTime, err := time.Parse(time.RFC3339, subTimes[i])
		if err != nil {
			continue
		}
		res[i] = Subscription{ChatID: chatIDs[i], ParametersID: parametersIDs[i], SubTime: parsedTime}
	}
	return res
}

func (s *Storage) UpdateSubscriptionTime(chatID int64, parameters_id int, subTime time.Time) error {
	conn := s.dbPool.Get(nil)
	defer s.dbPool.Put(conn)
	query := "UPDATE subscriptions SET sub_time = ? WHERE chat_id = ? AND parameters_id = ?"
	return sqlitex.Exec(conn, query, nil, subTime.Format(time.RFC3339), chatID, parameters_id)
}
