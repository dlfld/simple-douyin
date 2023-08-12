package handler

import (
	"reflect"
)

const errService = "服务繁忙"

func HandlerErr(resp interface{}, err error) {
	e := reflect.ValueOf(resp).Elem()
	code := int64(0)
	msg := "操作成功"
	if err != nil {
		code = 500
		msg = errService
	}
	e.FieldByName("StatusCode").SetInt(code)
	e.FieldByName("StatusMsg").Set(reflect.ValueOf(&msg))
	return
}

//
//func HandlerErr(resp interface{}, err error) {
//
//	return
//}
