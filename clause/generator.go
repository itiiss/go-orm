package clause

import (
	"fmt"
	"strings"
)

type generator func(values ...interface{}) (string, []interface{})

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBY] = _orderBy

}

func genBindVars(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ", ")
}

// 第一个参数是table名，之后的参数都是sql变量 vars
func _insert(values ...interface{}) (string, []interface{}) {
	// INSERT INTO $tableName ($fields)
	tableName := values[0].(string)
	field := strings.Join(values[1].([]string), ", ")
	sql := fmt.Sprintf("INSERT INTO %s (%v)", tableName, field)

	return sql, []interface{}{}
}

// 参数都是sql变量 vars
func _values(values ...interface{}) (string, []interface{}) {
	// VALUES ($v1), ($v2), ...
	var bindStr string
	var sql strings.Builder
	var vars []interface{}
	sql.WriteString("VALUES ")

	for i, value := range values {
		// sql的实参
		v := value.([]interface{})
		if bindStr == "" {
			bindStr = genBindVars(len(v))
		}
		// 再加上实参个数的 ?
		sql.WriteString(fmt.Sprintf("(%v)", bindStr))
		// 当前迭代i到了最后，补充 ，
		if i+1 != len(values) {
			sql.WriteString(", ")
		}
		vars = append(vars, v...)
	}
	return sql.String(), vars
}

// 第一个参数表名，第二个参数字段名
func _select(values ...interface{}) (string, []interface{}) {
	// SELECT $fields FROM $tableName
	tableName := values[0].(string)
	field := strings.Join(values[1].([]string), ", ")
	sql := fmt.Sprintf("SELECT %s FROM %s", field, tableName)
	return sql, []interface{}{}
}

// 参数为limit数
func _limit(values ...interface{}) (string, []interface{}) {
	// LIMIT $num
	return "LIMIT ?", values
}

// 第一个参数where条件，之后参数都是vars
func _where(values ...interface{}) (string, []interface{}) {
	sql := fmt.Sprintf("WHERE %s", values[0])
	vars := values[1:]
	return sql, vars
}

// 参数为orderBy 语句
func _orderBy(values ...interface{}) (string, []interface{}) {
	sql := fmt.Sprintf("ORDER BY %s", values[0])
	return sql, []interface{}{}
}
