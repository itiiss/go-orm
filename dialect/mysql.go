package dialect

import (
	"fmt"
	"reflect"
	"time"
)

type mysql struct{}

var d Dialect = &mysql{}

func init() {
	RegisterDialect("mysql", d)
}

// DataTypeOf 函数用于将 Go 数据类型映射为 MySQL 数据类型
func (s *mysql) DataTypeOf(typ reflect.Value) string {
	switch typ.Kind() {
	case reflect.Bool:
		return "BOOLEAN"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		return "INT"
	case reflect.Int64, reflect.Uint64:
		return "BIGINT"
	case reflect.Float32, reflect.Float64:
		return "DOUBLE"
	case reflect.String:
		return "VARCHAR(255)" // 对于字符串，通常指定一个长度，这里简单用 255
	case reflect.Array, reflect.Slice:
		return "BLOB"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "DATETIME"
		}
	}
	panic(fmt.Sprintf("invalid sql type %s (%s)", typ.Type().Name(), typ.Kind()))
}

// TableExistSQL 函数用于生成检查 MySQL 中表是否存在的 SQL 语句和参数
func (s *mysql) TableExistSQL(tableName string) (string, []interface{}) {
	args := []interface{}{tableName}
	return "SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?", args
}
