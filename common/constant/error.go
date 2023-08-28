package constant

import (
	"reflect"
)

const (
	Success      = 0
	ErrUserExist = iota + 4000
	ErrUserNotExist
	ErrInvalidPassword
	ErrTokenExpires
	ErrNotLogin
	ErrInvalidToken
	ErrSystemBusy
	ErrUnsupportedOperation
	ErrEmptyMessage
	ErrMsgTooLong
	ErrFeedErr
	ErrPublishList
)

var errMap = map[int32]string{
	Success:                 "success",
	ErrUserExist:            "用户名已存在",
	ErrUserNotExist:         "用户名不存在",
	ErrInvalidPassword:      "用户密码错误",
	ErrTokenExpires:         "token已过期",
	ErrNotLogin:             "未登录",
	ErrInvalidToken:         "无效的Token",
	ErrSystemBusy:           "系统繁忙",
	ErrUnsupportedOperation: "不支持的操作",
	ErrEmptyMessage:         "消息不能为空",
	ErrMsgTooLong:           "消息过长",
	ErrFeedErr:              "视频获取失败",
	ErrPublishList:          "获取已发布视频列表失败",
}

// HandlerErr 只要把对应期望的错误码，以及自己的resp传入，就会设置好StatusCode字段和StatusMsg字段
func HandlerErr(errCode int32, resp interface{}) {
	if resp == nil {
		return
	}
	e := reflect.ValueOf(resp).Elem()
	e.FieldByName("StatusCode").SetInt(int64(errCode))
	msg := errMap[errCode]
	e.FieldByName("StatusMsg").Set(reflect.ValueOf(&msg))
}

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

// NewErrResp 支持new一个错误响应，适用于请求还在路由端时的错误处理
func NewErrResp(errCode int32) Response {
	return Response{
		StatusCode: errCode,
		StatusMsg:  errMap[errCode],
	}
}
