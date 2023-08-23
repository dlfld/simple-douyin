package main

import (
	"context"
	"log"

	"github.com/douyin/common/crud"
	"github.com/douyin/kitex_gen/model"
	"github.com/douyin/kitex_gen/user"
	"github.com/douyin/models"
	"github.com/douyin/rpcServer/user/common"
	"github.com/douyin/rpcServer/user/util"
	"golang.org/x/crypto/bcrypt"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// UserRegister implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserRegister(ctx context.Context, req *user.UserRegisterRequest) (resp *user.UserRegisterResponse, err error) {
	//获取参数
	username := req.GetUsername()
	password := req.GetPassword()

	//密码不为空且小于32位
	if len(password) > 32 || len(password) <= 0 {
		statusMsg := "密码必须小于32位且不为空"
		resp = &user.UserRegisterResponse{StatusCode: -1, StatusMsg: &statusMsg, UserId: -1, Token: ""}
		return resp, nil
	}
	if len(username) > 32 {
		statusMsg := "用户名必须小于32位"
		resp = &user.UserRegisterResponse{StatusCode: -1, StatusMsg: &statusMsg, UserId: -1, Token: ""}
		return resp, nil
	}
	//如果name为空
	if len(username) <= 0 {
		username = util.RandomString(16)
	}

	//判断用户是否存在
	userR, _ := models.GetUserByName(username)
	if userR != nil {
		statusMsg := "username exist"
		resp = &user.UserRegisterResponse{StatusCode: -1, StatusMsg: &statusMsg, UserId: -1, Token: ""}
		return resp, nil
	}

	//创建用户&&密码加密
	encryptPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		statusMsg := "加密错误"
		resp = &user.UserRegisterResponse{StatusCode: -1, StatusMsg: &statusMsg, UserId: -1, Token: ""}
		return resp, nil
	}

	newUser, nerr := models.CreateUser(username, string(encryptPassword))
	if nerr != nil {
		statusMsg := "create user failed"
		resp = &user.UserRegisterResponse{StatusCode: -1, StatusMsg: &statusMsg, UserId: -1, Token: ""}
		return resp, nil
	}

	token, terr := common.ReleaseToken(*newUser)
	if terr != nil {
		statusMsg := "token发放错误"
		resp = &user.UserRegisterResponse{StatusCode: -1, StatusMsg: &statusMsg, UserId: -1, Token: ""}
		return resp, nil
	}

	//返回结果
	statusMsg := "注册并登录成功"
	resp = &user.UserRegisterResponse{StatusCode: 0, StatusMsg: &statusMsg, UserId: int64(int(newUser.ID)), Token: token}
	return resp, nil
}

// UserLogin implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserLogin(ctx context.Context, req *user.UserLoginRequest) (resp *user.UserLoginResponse, err error) {
	//获取参数
	username := req.GetUsername()
	password := req.GetPassword()

	//数据验证
	if len(password) <= 0 || len(username) <= 0 {
		statusMsg := "用户名或密码为空"
		resp = &user.UserLoginResponse{StatusCode: -1, StatusMsg: &statusMsg, UserId: -1, Token: ""}
		return resp, nil
	}

	//判断用户是否存在
	userL, qerr := models.GetUserByName(username)
	if qerr != nil {
		statusMsg := "用户不存在"
		resp = &user.UserLoginResponse{StatusCode: -1, StatusMsg: &statusMsg, UserId: -1, Token: ""}
		return resp, nil
	}
	//判断密码是否正确
	if err2 := bcrypt.CompareHashAndPassword([]byte(userL.Password), []byte(password)); err2 != nil {
		statusMsg := "密码错误"
		resp = &user.UserLoginResponse{StatusCode: -1, StatusMsg: &statusMsg, UserId: -1, Token: ""}
		return resp, nil
	}

	//发放token
	token, err3 := common.ReleaseToken(*userL)
	if err3 != nil {
		log.Printf("token generate error: %v", err)
		statusMsg := "系统错误"
		resp = &user.UserLoginResponse{StatusCode: -1, StatusMsg: &statusMsg, UserId: -1, Token: ""}
		return resp, nil
	}
	crud, _ := crud.NewCachedCRUD()
	crud.CacheUserInfo(userL)
	//返回结果
	statusMsg := "登录成功"
	resp = &user.UserLoginResponse{StatusCode: 0, StatusMsg: &statusMsg, UserId: int64(userL.ID), Token: token}
	return resp, nil
}

// UserMsg implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserMsg(ctx context.Context, req *user.UserRequest) (resp *user.UserResponse, err error) {
	userId := uint(req.GetUserId())
	tokenString := req.GetToken()

	token, claims, err1 := common.ParseToken(tokenString)
	useridClaims := claims.UserId
	//发生错误或者token无效

	if err1 != nil || !token.Valid || useridClaims != userId {
		statusMsg := "用户签名不符"
		resp = &user.UserResponse{StatusCode: -1, StatusMsg: &statusMsg, User: &model.User{}}
		return resp, nil
	}
	crud, _ := crud.NewCachedCRUD()
	// crud.GetUserInfo(strconv.Itoa(int(userId)))
	userI, qerr := models.GetUserByUserId(userId)
	crud.CacheUserInfo(userI)
	if qerr != nil {
		statusMsg := "用户不存在"
		resp = &user.UserResponse{StatusCode: -1, StatusMsg: &statusMsg, User: &model.User{}}
		return resp, nil
	}

	statusMsg := "获取用户信息成功"
	resp = &user.UserResponse{StatusCode: 0, StatusMsg: &statusMsg, User: models.ChangeModel(*userI)}
	return resp, nil
}
