package session

import (
	"database/sql"
	"go-orm/log"
	"strings"
)

type Session struct {
	db        *sql.DB
	sql       strings.Builder
	sqlValues []interface{}
}

func NewSession(db *sql.DB) *Session {
	return &Session{db: db}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlValues = nil
}

func (s *Session) DB() *sql.DB {
	return s.db
}

func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlValues = append(s.sqlValues, values...)
	return s
}

func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlValues)
	result, err = s.DB().Exec(s.sql.String(), s.sqlValues...)
	if err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlValues)
	return s.DB().QueryRow(s.sql.String(), s.sqlValues...)
}

func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlValues)
	rows, err = s.DB().Query(s.sql.String(), s.sqlValues...)
	if err != nil {
		log.Error(err)
	}
	return
}
