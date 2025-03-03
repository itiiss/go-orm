package engine

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

const (
	server   = "127.0.0.1"
	port     = "3306"
	user     = "root"
	password = "as951753258"
	database = "orm_test"
)

var source = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True", user, password, server, port, database)

func OpenDB(t *testing.T) *Engine {
	t.Helper()
	engine, err := NewEngine("mysql", source)
	if err != nil {
		t.Fatal("failed to connect", err)
	}
	return engine
}

func TestNewEngine(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
}
