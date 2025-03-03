package session

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go-orm/dialect"
	"os"
	"testing"
)

var (
	TestDB      *sql.DB
	TestDial, _ = dialect.GetDialect("mysql")
)

const (
	server   = "127.0.0.1"
	port     = "3306"
	user     = "root"
	password = "as951753258"
	database = "orm_test"
)

var source = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True", user, password, server, port, database)

func TestMain(m *testing.M) {
	TestDB, _ = sql.Open("mysql", source)
	code := m.Run()
	_ = TestDB.Close()
	os.Exit(code)
}

func NewTestSession() *Session {
	return NewSession(TestDB, TestDial)
}

func TestSession_Exec(t *testing.T) {
	s := NewTestSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	result, _ := s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	if count, err := result.RowsAffected(); err != nil || count != 2 {
		t.Fatal("expect 2, but got", count)
	}
}

func TestSession_QueryRows(t *testing.T) {
	s := NewTestSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	row := s.Raw("SELECT count(*) FROM User").QueryRow()
	var count int
	if err := row.Scan(&count); err != nil || count != 0 {
		t.Fatal("failed to query db", err)
	}
}
