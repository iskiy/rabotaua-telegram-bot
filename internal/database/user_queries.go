package database

import "crawshaw.io/sqlite"

func (s *Storage) InsertUser(chatID int64, state int) error {
	return s.execQuery("INSERT INTO users(chat_id, user_state) values(?, ?)", nil, chatID, state)
}

func (s *Storage) IsUserPresent(chatID int64) (bool, error) {
	var ids int64 = -1
	fn := func(stmt *sqlite.Stmt) error {
		ids = stmt.ColumnInt64(0)
		return nil
	}
	err := s.execQuery("SELECT chat_id FROM users WHERE chat_id = ?",
		fn, chatID)
	if err != nil {
		return false, err
	}
	if ids > 0 {
		return true, nil
	}
	return false, nil
}

func (s *Storage) UpdateUserState(chatID int64, state int) error {
	return s.execQuery("UPDATE users SET user_state = ? WHERE chat_id = ?;", nil, state, chatID)
}

func (s *Storage) GetUserState(chatID int64) (int, error) {
	var states []int
	fn := func(stmt *sqlite.Stmt) error {
		states = append(states, stmt.ColumnInt(0))
		return nil
	}
	err := s.execQuery("SELECT user_state FROM users WHERE chat_id = ?", fn, chatID)
	if err != nil {
		return -1, err
	}
	if len(states) == 0 {
		return -1, StateGetError
	}
	return states[0], err
}
