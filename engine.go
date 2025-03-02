package engine

import (
	"database/sql"
	"go-orm/dialect"
	"go-orm/log"
	"go-orm/session"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return
	}
	// check database connection is alive
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}

	dial, ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorf("dialect %s Not Found", driver)
		return
	}

	e = &Engine{db: db, dialect: dial}
	log.Info("Connect database success")
	return
}

func (e *Engine) Close() {
	err := e.db.Close()
	if err != nil {
		log.Error(err)
	}
	log.Info("Close database success")
}

func (e *Engine) NewSession() *session.Session {
	return session.NewSession(e.db, e.dialect)
}
