package session

import (
	"go-orm/clause"
	"reflect"
)

// Insert 参数是对象指针，可以插入多个
// session.Insert(&user1, &user2)
func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)
	for _, value := range values {
		table := s.Model(value).RefTable()
		// 构造 Insert子语句
		// 如果插入多个对象，会执行多次，但是set的结果是相同的
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		// 从对象中提取出符合schema定义的value
		recordValues = append(recordValues, table.RecordValues(value))
	}

	// 构造Values子语句
	s.clause.Set(clause.VALUES, recordValues...)
	// 调用一次 clause.Build() 按照传入的顺序构造出最终的 SQL 语句
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	// 执行完整的sql获取结果
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Session) Find(values interface{}) error {
	// 利用反射获取value的反射值和元素类型
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType := destSlice.Type().Elem()
	// 通过值和类型，创建新表
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()

	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	// 需要补充其他WHERE，ORDERBY，LIMIT的子语句，需要提前set好
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}

	//遍历查询结果并填充values切片中
	for rows.Next() {
		dest := reflect.New(destType).Elem()
		var values []interface{}
		for _, name := range table.FieldNames {
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}

		err := rows.Scan(values...)
		if err != nil {
			return err
		}
		destSlice.Set(reflect.Append(destSlice, dest))
	}

	return rows.Close()
}
