package database

import (
	"crawshaw.io/sqlite"
	"crawshaw.io/sqlite/sqlitex"
	"errors"
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
	err := execQueryWithoutParameters(dbPool, `CREATE TABLE IF NOT EXISTS users(
		chat_id INTEGER PRIMARY KEY,
		user_state INTEGER NOT NULL)`)
	if err != nil {
		return err
	}
	err = execQueryWithoutParameters(dbPool, `CREATE TABLE IF NOT EXISTS vacancy_parameters(
		parameters_id INTEGER NOT NULL, 
		chat_id INTEGER NOT NULL, 
		keywords VARCHAR, 
		city_id INTEGER, 
		schedule_id INTEGER, 
		PRIMARY KEY(parameters_id, chat_id)
		);`)
	if err != nil {
		return err
	}
	err = execQueryWithoutParameters(dbPool, `CREATE TABLE IF NOT EXISTS vacancies_view( 
		view_id INTEGER PRIMARY KEY AUTOINCREMENT, 
		parameters_id INTEGER NOT NULL, 
		chat_id INTEGER NOT NULL, 
		page INTEGER NOT NULL, 
		FOREIGN KEY (parameters_id) 
		REFERENCES vacancy_parameters(parameters_id) ON UPDATE CASCADE ON DELETE CASCADE, 
		FOREIGN KEY(chat_id) REFERENCES vacancy_parameters(chat_id) ON UPDATE CASCADE ON DELETE CASCADE);`)
	if err != nil {
		return err
	}
	err = execQueryWithoutParameters(dbPool, `CREATE TABLE IF NOT EXISTS cities( 
		 city_id INTEGER PRIMARY KEY,  
		 city_name varchar NOT NULL);`)
	if err != nil {
		return err
	}
	err = execQueryWithoutParameters(dbPool, `CREATE TABLE IF NOT EXISTS 
		subscriptions(chat_id INTEGER NOT NULL, 
		parameters_id INTEGER NOT NULL, 
		sub_time VARCHAR NOT NULL, 
		PRIMARY KEY(chat_id, parameters_id), 
		FOREIGN KEY (parameters_id) REFERENCES vacancy_parameters(parameters_id) ON UPDATE CASCADE ON DELETE CASCADE, 
		FOREIGN KEY(chat_id) REFERENCES vacancy_parameters(chat_id) ON UPDATE CASCADE ON DELETE CASCADE);`)
	if err != nil {
		return err
	}
	return insertCityDictionary(dbPool)
}

func (s *Storage) execQuery(query string, resultFn func(stmt *sqlite.Stmt) error, arguments ...interface{}) error {
	conn := s.dbPool.Get(nil)
	defer s.dbPool.Put(conn)
	return sqlitex.Exec(conn, query, resultFn, arguments...)
}

func insertCityDictionary(dbPool *sqlitex.Pool) error {
	conn := dbPool.Get(nil)
	defer dbPool.Put(conn)
	client := rabotaua.NewRabotaClient()
	cities, err := client.GetCities()
	if err != nil {
		return err
	}
	query := "INSERT OR IGNORE INTO cities(city_id, city_name) VALUES(?, ?)"
	for _, c := range cities {
		if err := sqlitex.Exec(conn, query, nil, c.ID, c.City); err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) GetCityName(cityID int) (string, error) {
	var res string
	fn := func(stmt *sqlite.Stmt) error {
		res = stmt.ColumnText(0)
		return nil
	}
	query := "SELECT city_name FROM cities WHERE city_id = ?"
	err := s.execQuery(query, fn, cityID)
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
