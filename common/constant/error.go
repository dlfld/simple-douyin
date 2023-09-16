package constant

import (
	"reflect"
)

const (
	Success      = 0
	ErrUserExist = iota + 4000
	ErrUserNotExist
	ErrLoginFailed
	ErrRegisterFailed
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
	ErrFollowAction
	ErrGetFollowList
	ErrGetFollowerList
	ErrGetFriendList
	ErrBadRequest
	ErrFavoriteAction
	ErrFavoriteList
	ErrCommentAction
	ErrCommentList
	ErrVideoSizeMaxLimit
	ErrVideoTitleLength
	ErrBloomVideo
	ErrBloomComment
	ErrBloomUser
	ErrVideoPublish
	ErrFileCantOpen
	ErrMessageList
	ErrSendMessage
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
	ErrFollowAction:         "关注/取关操作失败",
	ErrGetFollowList:        "获取关注列表失败",
	ErrGetFollowerList:      "获取粉丝列表失败",
	ErrBadRequest:           "请求参数错误",
	ErrGetFriendList:        "获取好友列表失败",
	ErrFavoriteAction:       "点赞操作失败",
	ErrFavoriteList:         "获取点赞列表失败",
	ErrCommentAction:        "评论操作失败",
	ErrCommentList:          "获取评论列表失败",
	ErrBloomVideo:           "该视频不存在",
	ErrBloomComment:         "该评论不存在",
	ErrBloomUser:            "该用户不存在",
	ErrLoginFailed:          "登录失败",
	ErrRegisterFailed:       "注册失败",
	ErrVideoPublish:         "视频发布失败",
	ErrFileCantOpen:         "视频文件无法访问",
	ErrMessageList:          "获取消息列表失败",
	ErrSendMessage:          "发送消息失败",
	ErrVideoSizeMaxLimit:    "视频大小超过限制",
	ErrVideoTitleLength:     "视频标题长度错误",
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
