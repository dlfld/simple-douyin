include "model.thrift"

namespace go relation

// 关注操作
struct follow_action_request {
    1:required i64 from_user_id // 当前用户id
    2:required i64 to_user_id
    3:required i32 action_type
}

struct follow_action_response {
    1:required i32 status_code
    2:optional string status_msg
}

// 关注列表
struct following_list_request {
    1:required i64 user_id
    2:required string token
}

struct following_list_response {
    1:required i32 status_code
    2:optional string status_msg
    3:required list<model.User> user_list
}

// 粉丝列表
struct follower_list_request {
    1:required i64 user_id
    2:required string token
}

struct follower_list_response {
    1:required i32 status_code
    2:optional string status_msg
    3:required list<model.User> user_list
}

// 好友列表
struct relation_friend_list_request {
    1:required i64 user_id   // 用户id
}

struct relation_friend_list_response {
    1:required i32 status_code                      //状态码，0- 成功，其他值-失败
    2:optional string status_msg                    //返回状态描述
    3:required list<model.FriendUser> user_list     //好友用户列表，好友指的是相互关注的用户
}



// rpc service
service RelationService {
    follow_action_response FollowAction(1:follow_action_request req)
    following_list_response FollowList(1:following_list_request req)
    follower_list_response FollowerList(1:follower_list_request req)
    relation_friend_list_response FriendList(1:relation_friend_list_request req)
}