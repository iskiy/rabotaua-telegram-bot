package database

import (
	"crawshaw.io/sqlite"
	"crawshaw.io/sqlite/sqlitex"
	"fmt"
	"time"
)

func (s *Storage) InsertSubscription(chatID int64, parametersID int, subTime time.Time) error {
	return s.execQuery("INSERT INTO subscriptions(chat_id, parameters_id, sub_time) values(?, ?, ?)",
		nil, chatID, parametersID, subTime.Format(time.RFC3339))
}

func (s *Storage) IsSubscriptionPresent(chatID int64, parametersID int) (bool, error) {
	var ids int64 = -1
	fn := func(stmt *sqlite.Stmt) error {
		ids = stmt.ColumnInt64(0)
		return nil
	}
	err := s.execQuery("SELECT chat_id FROM subscriptions WHERE chat_id = ? AND parameters_id = ? ", fn, chatID, parametersID)
	if err != nil {
		return false, err
	}
	if ids > 0 {
		return true, nil
	}
	return false, nil
}

func (s *Storage) GetAllSubscriptions() ([]Subscription, error) {
	query := "SELECT chat_id, parameters_id, sub_time FROM subscriptions"
	return s.getSubscriptionsFromQuery(query)
}

func (s *Storage) GetSubscriptionsFromUser(chatID int64) ([]Subscription, error) {
	query := "SELECT chat_id, parameters_id, sub_time FROM subscriptions WHERE chat_id = " + fmt.Sprintf("%d", chatID)
	return s.getSubscriptionsFromQuery(query)
}

func (s *Storage) getSubscriptionsFromQuery(query string) ([]Subscription, error) {
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

func (s *Storage) UpdateSubscriptionTime(chatID int64, parametersID int, subTime time.Time) error {
	query := "UPDATE subscriptions SET sub_time = ? WHERE chat_id = ? AND parameters_id = ?"
	return s.execQuery(query, nil, subTime.Format(time.RFC3339), chatID, parametersID)
}

func (s *Storage) DeleteSubscription(chatID int64, parametersID int) error {
	query := "DELETE FROM subscriptions WHERE chat_id = ? AND parameters_id = ?"
	return s.execQuery(query, nil, chatID, parametersID)
}
