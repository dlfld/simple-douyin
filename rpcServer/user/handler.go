package main

import (
	"context"
	"fmt"
	"github.com/douyin/common/mysql"
	"github.com/douyin/kitex_gen/model"
	"github.com/douyin/kitex_gen/user"
	"github.com/douyin/models"
	"github.com/douyin/rpcServer/user/common"
	"github.com/douyin/rpcServer/user/util"
	"golang.org/x/crypto/bcrypt"
	"log"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// UserRegister implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserRegister(ctx context.Context, req *user.UserRegisterRequest) (resp *user.UserRegisterResponse, err error) {
	//获取参数
	username := req.Username
	password := req.Password

	//数据验证密码不为空且小于32为
	if len(password) > 32 || len(password) <= 0 {
		//responseLogin(ctx, -1, "密码必须小于32位且不为空", -1, "")
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
	fmt.Println(username, password)

	//创建用户&&密码加密
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		statusMsg := "加密错误"
		resp = &user.UserRegisterResponse{StatusCode: -1, StatusMsg: &statusMsg, UserId: -1, Token: ""}
		return resp, nil
	}
	newUser := models.User{
		UserName: username,
		Password: string(hasedPassword),
	}

	db, _ := mysql.NewMysqlConn()
	db.Create(&newUser)

	token, err := common.ReleaseToken(newUser)
	if err != nil {
		log.Printf("token generate error: %v", err)
		statusMsg := "token发放错误"
		resp = &user.UserRegisterResponse{StatusCode: -1, StatusMsg: &statusMsg, UserId: -1, Token: ""}
		return resp, nil
	}

	//返回结果
	//responseLogin(ctx, 0, "注册并登录成功", int(newUser.ID), token)
	statusMsg := "注册并登录成功"
	resp = &user.UserRegisterResponse{StatusCode: 0, StatusMsg: &statusMsg, UserId: int64(int(newUser.ID)), Token: token}
	return resp, nil
}

// UserLogin implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserLogin(ctx context.Context, req *user.UserLoginRequest) (resp *user.UserLoginResponse, err error) {
	//获取参数
	username := req.Username
	password := req.Password

	//数据验证
	if len(password) <= 0 || len(username) <= 0 {
		statusMsg := "用户名或密码为空"
		resp.SetStatusCode(-1)
		resp.SetStatusMsg(&statusMsg)
		resp.SetUserId(-1)
		resp.SetToken("")
		return resp, nil
	}
	fmt.Println(username, password)

	//判断用户是否存在
	db, _ := mysql.NewMysqlConn()
	var user models.User
	db.Where("user_name = ?", username).First(&user)

	if user.ID <= 0 {
		//responseLogin(ctx, -1, "用户不存在", -1, "")
		statusMsg := "用户不存在"
		resp.SetStatusCode(-1)
		resp.SetStatusMsg(&statusMsg)
		//resp.SetUserId(-1)
		//resp.SetToken("")
		return resp, nil
	}
	//判断密码是否正确
	if err2 := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err2 != nil {
		//responseLogin(ctx, -1, "密码错误", -1, "")
		statusMsg := "密码错误"
		resp.SetStatusCode(-1)
		resp.SetStatusMsg(&statusMsg)
		//resp.SetUserId(-1)
		//resp.SetToken("")
		return resp, nil
	}

	//发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		log.Printf("token generate error: %v", err)
		statusMsg := "系统错误"
		resp.SetStatusCode(-1)
		resp.SetStatusMsg(&statusMsg)
		//resp.SetUserId(-1)
		//resp.SetToken("")
		return resp, nil
	}

	//返回结果
	statusMsg := "登录成功"
	resp.SetStatusCode(0)
	resp.SetStatusMsg(&statusMsg)
	resp.SetUserId(int64(user.ID))
	resp.SetToken(token)
	return resp, nil
}

//// isSuccess=0为成功，其他值失败
//func responseLogin(ctx *gin.Context, isSuccess int, msg string, userId int, token string) {
//	if userId <= 0 {
//		ctx.JSON(http.StatusOK, gin.H{"status_code": isSuccess, "status_msg": msg, "user_id": nil, "token": ""})
//		return
//	}
//	ctx.JSON(http.StatusOK, gin.H{
//		"status_code": isSuccess,
//		"status_msg":  msg,
//		"user_id":     userId,
//		"token":       token,
//	})
//}

// UserMsg implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserMsg(ctx context.Context, req *user.UserRequest) (resp *user.UserResponse, err error) {
	//获取authorization header（以“Bearer ”开头）
	userId := req.GetUserId()
	tokenString := req.GetToken()

	token, _, err := common.ParseToken(tokenString)
	//useridClaims := claims.UserId
	if err != nil || !token.Valid { //发生错误或者token无效
		statusMsg := "用户签名不符"
		resp.SetStatusCode(-1)
		resp.SetStatusMsg(&statusMsg)
		resp.SetUser(nil)
		return resp, nil
	}

	//通过验证，获取claim中的userId
	db, _ := mysql.NewMysqlConn()
	var user models.User
	db.First(&user, userId)

	//若用户不存在
	if user.ID <= 0 {
		statusMsg := "user权限不足"
		resp.SetStatusCode(-1)
		resp.SetStatusMsg(&statusMsg)
		resp.SetUser(nil)
		return resp, nil
	}

	var u model.User
	u.SetId(int64(user.ID))
	u.SetUserName(user.UserName)
	//用户存在，将user信息写入上下文
	//ctx.Set("user", user)
	statusMsg := "获取用户信息成功"
	resp.SetStatusCode(0)
	resp.SetStatusMsg(&statusMsg)
	resp.SetUser(&u)
	return resp, nil
}
