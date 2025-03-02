package session

import (
	"fmt"
	"go-orm/log"
	"go-orm/schema"
	"reflect"
	"strings"
)

// Model 解析传入对象成Schema，保存到refTable中，继续返回s支持链式调用
func (s *Session) Model(value interface{}) *Session {
	// nil or a new model, update refTable
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("Model is not set")
	}
	return s.refTable
}

func (s *Session) CreateTable() error {
	table := s.RefTable()

	var columns []string
	for _, field := range table.Fields {
		columnDef := fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag)
		columns = append(columns, columnDef)
	}

	desc := strings.Join(columns, ",")
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s) ;", table.Name, desc)

	_, err := s.Raw(sql).Exec()
	return err
}

func (s *Session) DropTable() error {
	sql := fmt.Sprintf("DROP TABLE IF EXISTS %s ;", s.RefTable().Name)
	_, err := s.Raw(sql).Exec()
	return err
}

func (s *Session) HasTable() bool {
	// INFORMATION_SCHEMA.TABLES 中查找当前数据库里名为 s.RefTable().Name
	sql, values := s.dialect.TableExistSQL(s.RefTable().Name)
	// 执行该查询sql，存在则返回该表的名称
	row := s.Raw(sql, values...).QueryRow()
	// 将row中保存的表名assign到tmp字符串上
	var tmp string
	_ = row.Scan(&tmp)
	
	return tmp == s.RefTable().Name
}
