package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	myorm "go-orm"
)

const (
	server   = "127.0.0.1"
	port     = "3306"
	user     = "root"
	password = "as951753258"
	database = "orm_test"
)

var source = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True", user, password, server, port, database)

func main() {

	engine, _ := myorm.NewEngine("mysql", source)
	defer engine.Close()

	s := engine.NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name VARCHAR(255));").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name VARCHAR(255));").Exec()

	result, _ := s.Raw("INSERT INTO User(Name) VALUES (?),(?);", "Morris", "Pierre").Exec()
	count, _ := result.RowsAffected()
	fmt.Printf("Exec success, %d affected\n", count)
}
