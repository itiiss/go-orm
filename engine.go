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

type TxFunc func(*session.Session) (interface{}, error)

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

// Transaction 将所有的操作放到一个回调函数中，作为入参传递给 engine.Transaction()
// 发生任何错误，自动回滚，如果没有错误发生，则提交
func (engine *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := engine.NewSession()
	if err := s.Begin(); err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = s.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			_ = s.Rollback() // err is non-nil; don't change it
		} else {
			err = s.Commit() // err is nil; if Commit returns error update err
		}
	}()

	return f(s)
}

func (engine *Engine) Close() {
	err := engine.db.Close()
	if err != nil {
		log.Error(err)
	}
	log.Info("Close database success")
}

func (engine *Engine) NewSession() *session.Session {
	return session.NewSession(engine.db, engine.dialect)
}
