package main

import (
	"context"
	"github.com/douyin/kitex_gen/user"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// UserRegister implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserRegister(ctx context.Context, req *user.UserRegisterRequest) (resp *user.UserRegisterResponse, err error) {
	// TODO: Your code here...
	return
}

// UserLogin implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserLogin(ctx context.Context, req *user.UserLoginRequest) (resp *user.UserLoginResponse, err error) {
	// TODO: Your code here...
	return
}

// UserMsg implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserMsg(ctx context.Context, req *user.UserRequest) (resp *user.UserResponse, err error) {
	// TODO: Your code here...
	return
}
