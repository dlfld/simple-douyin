package constant

import "errors"

type ErrResponse struct {
	Code int32
	Err  error
}

var (
	ErrUserExist            = ErrResponse{4001, errors.New("用户名已存在")}
	ErrUserNotExist         = ErrResponse{4002, errors.New("用户名不存在")}
	ErrInvalidPassword      = ErrResponse{4003, errors.New("用户密码错误")}
	ErrTokenExpires         = ErrResponse{4004, errors.New("token已过期")}
	ErrNotLogin             = ErrResponse{4005, errors.New("未登录")}
	ErrInvalidToken         = ErrResponse{4007, errors.New("无效的Token")}
	ErrSystemBusy           = ErrResponse{4008, errors.New("系统繁忙")}
	ErrUnsupportedOperation = ErrResponse{4009, errors.New("不支持的操作")}
)

func (err ErrResponse) Error() string {
	return err.Err.Error()
}
