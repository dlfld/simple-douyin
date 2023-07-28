include "model.thrift"

namespace go user
// 用户注册
struct user_register_request {
  1:required string username; // 注册用户名，最长32个字符
  2:required string password; // 密码，最长32个字符
}

struct user_register_response {
  1:required i8 status_code; // 状态码，0-成功，其他值-失败
  2:optional string status_msg; // 返回状态描述
  3:required i64 user_id ; // 用户id
  4:required string token; // 用户鉴权token
}
// 用户登录接口
struct user_login_request {
    1:required string username
    2:required string password
}

struct user_login_response {
    1:required i32 status_code
    2:optional string status_msg
    3:required i64 user_id
    4:required string token
}
// 用户信息
struct user_request {
    1:required i64 user_id
    2:required string token
}

struct user_response {
    1:required i32 status_code
    2:optional string status_msg
    3:required model.User user
}
service UserService{
    user_register_response UserRegister(1:user_register_request req)
    user_login_response UserLogin(1:user_login_request req)
    user_response UserMsg(1:user_request req)
}



