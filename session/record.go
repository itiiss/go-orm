package session

import (
	"errors"
	"go-orm/clause"
	"reflect"
)

// Insert 参数是对象指针，可以插入多个
// session.Insert(&user1, &user2)
func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)
	for _, value := range values {

		// hooks： 执行 value 对象上挂载的 BeforeInsert方法
		s.CallMethod(BeforeInsert, value)

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

	s.CallMethod(AfterInsert, nil)
	return result.RowsAffected()
}

func (s *Session) Find(values interface{}) error {
	s.CallMethod(BeforeQuery, nil)
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

		s.CallMethod(AfterQuery, dest.Addr().Interface())
		destSlice.Set(reflect.Append(destSlice, dest))
	}

	return rows.Close()
}

// Update 接受 2 种入参，平铺开来的键值对和 map 类型的键值对
func (s *Session) Update(kv ...interface{}) (int64, error) {
	s.CallMethod(BeforeUpdate, nil)
	// 判断传入参数的类型
	m, ok := kv[0].(map[string]interface{})
	// 因为 generator 接受的参数是 map 类型的键值对，如果是不是 map 类型，则会自动转换
	if !ok {
		m = make(map[string]interface{})
		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
	}

	s.clause.Set(clause.UPDATE, s.RefTable().Name, m)
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterUpdate, nil)
	return result.RowsAffected()
}

func (s *Session) Delete() (int64, error) {
	s.CallMethod(BeforeDelete, nil)

	s.clause.Set(clause.DELETE, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}

	s.CallMethod(AfterDelete, nil)
	return result.RowsAffected()
}

func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)
	row := s.Raw(sql, vars...).QueryRow()

	var tmp int64
	err := row.Scan(&tmp)
	if err != nil {
		return 0, err
	}
	return tmp, nil
}

func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	return s
}

func (s *Session) Where(desc string, args ...interface{}) *Session {
	var vars []interface{}
	s.clause.Set(clause.WHERE, append(append(vars, desc), args...)...)
	return s
}

func (s *Session) OrderBy(desc string) *Session {
	s.clause.Set(clause.ORDERBY, desc)
	return s
}

// First 根据传入的类型，利用反射构造切片
// 调用 Limit(1) 限制返回的行数，调用 Find 方法获取到查询结果
func (s *Session) First(values interface{}) error {
	dest := reflect.Indirect(reflect.ValueOf(values))
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()

	err := s.Limit(1).Find(destSlice.Addr().Interface())
	if err != nil {
		return err
	}

	if destSlice.Len() == 0 {
		return errors.New("not found")
	}
	dest.Set(destSlice.Index(0))
	return nil
}
