package session

import (
	"database/sql"
	"go-orm/clause"
	"go-orm/dialect"
	"go-orm/log"
	"go-orm/schema"
	"strings"
)

type Session struct {
	db        *sql.DB
	dialect   dialect.Dialect
	tx        *sql.Tx
	refTable  *schema.Schema
	sql       strings.Builder
	sqlValues []interface{}
	clause    clause.Clause
}

// CommonDB is a minimal function set of db
type CommonDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

var _ CommonDB = (*sql.DB)(nil)
var _ CommonDB = (*sql.Tx)(nil)

func NewSession(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect,
	}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlValues = nil
	s.clause = clause.Clause{}
}

func (s *Session) DB() CommonDB {
	if s.tx != nil {
		return s.tx
	}
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
