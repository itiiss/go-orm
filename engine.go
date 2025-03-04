package engine

import (
	"database/sql"
	"fmt"
	"go-orm/dialect"
	"go-orm/log"
	"go-orm/session"
	"strings"
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

// 得到a中有的，但是b中没有的字段，a总是较少字段的那一个
func difference(a, b []string) (diff []string) {
	mb := make(map[string]bool)
	// 维护所有b中的字段
	for _, x := range b {
		mb[x] = true
	}
	// 枚举所有a中的字段
	for _, x := range a {
		_, ok := mb[x]
		//如果该字段在b中不存在则添加进diff
		if !ok {
			diff = append(diff, x)
		}
	}
	return
}

func (engine *Engine) Migrate(value interface{}) error {
	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		if !s.Model(value).HasTable() {
			log.Infof("table %s doesn't exist", s.RefTable().Name)
			return nil, s.CreateTable()
		}

		// table是新表的结构，即期望的表结构
		table := s.RefTable()
		// 取出第一条记录
		rows, _ := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1", table.Name)).QueryRows()
		// 获取该记录的所有字段，即当前旧表的结构
		columns, _ := rows.Columns()
		// old A，B
		// new A，C
		// add C，Del B
		// 分别得到增加和减少的字段
		addCols := difference(table.FieldNames, columns)
		delCols := difference(columns, table.FieldNames)
		log.Infof("add cols: %v, delete cols %v", addCols, delCols)

		// 在原表的基础上添加新增的字段，此时包括所有旧字段 + 新表增加的字段
		for _, col := range addCols {
			f := table.GetField(col)
			sql := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table.Name, f.Name, f.Type)
			_, err = s.Raw(sql).Exec()
			if err != nil {
				return
			}
		}
		// current：A，B，C

		// 如果只有新增没有减少，到此就完成了
		if len(delCols) == 0 {
			return
		}

		tmp := "tmp_" + table.Name
		// 取出新表的所有字段，即所有期望的字段，A，C
		fieldStr := strings.Join(table.FieldNames, ",")
		// 创建一个tmp表，从A，B，C 中只选择A，C字段
		s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s from %s", tmp, fieldStr, table.Name))
		// 删除旧表，并把新表改名成旧表
		s.Raw(fmt.Sprintf("DROP TABLE %s", table.Name))
		s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME to %s", tmp, table.Name))
		_, err = s.Exec()
		return
	})

	return err
}
