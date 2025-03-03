package session

import (
	"go-orm/log"
	"reflect"
)

const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
)

// CallMethod 如果传入参数 value，则调用values上的hook方法
// 否则调用 s.RefTable()， 即 model上的hook方法
func (s *Session) CallMethod(method string, value interface{}) {
	fm := reflect.ValueOf(s.RefTable()).MethodByName(method)
	if value != nil {
		fm = reflect.ValueOf(value).MethodByName(method)
	}

	param := []reflect.Value{reflect.ValueOf(s)}

	if fm.IsValid() {
		v := fm.Call(param)
		if len(v) > 0 {
			err, ok := v[0].Interface().(error)
			if ok {
				log.Error(err)
			}
		}
	}
	return
}
