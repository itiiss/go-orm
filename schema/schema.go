package schema

import (
	"go-orm/dialect"
	"go/ast"
	"reflect"
)

// Field represents a column of table
type Field struct {
	Name string
	Type string
	Tag  string
}

// Schema represents a table of database
type Schema struct {
	Model      interface{}
	Name       string
	Fields     []*Field
	FieldNames []string
	FieldsMap  map[string]*Field
}

func (s *Schema) GetField(name string) *Field {
	return s.FieldsMap[name]
}

// Parse 将任意对象解析成schema实例
func Parse(dest interface{}, d dialect.Dialect) *Schema {
	// TypeOf() 和 ValueOf() 是 reflect 包最常用 2 个方法，分别用来返回入参的类型和值。
	// 因为设计的入参是一个对象的指针，因此需要 reflect.Indirect() 获取指针指向的实例
	// Type()：获取解引用后值的类型。
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:     dest,
		Name:      modelType.Name(), // 获取到结构体的名称作为表名
		FieldsMap: make(map[string]*Field),
	}

	for i := 0; i < modelType.NumField(); i++ {
		// 枚举所有字段，p是字段的反射实例
		p := modelType.Field(i)
		// 只处理非匿名和导出字段
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}

			// 查找字段标签中是否存在 go-orm 标签，如果存在则将其值赋值给 Field 的 Tag 字段
			v, ok := p.Tag.Lookup("go-orm")
			if ok {
				field.Tag = v
			}
			// 更新 Schema对象的Field相关自动
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.FieldsMap[p.Name] = field
		}
	}
	return schema
}
