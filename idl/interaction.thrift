include "model.thrift"

namespace go interaction

// 赞操作
struct favorite_action_request {
    1:required string token
    2:required i64 video_id
    3:required i32 action_type
}

struct favorite_action_response {
    1:required i32 status_code
    2:optional string status_msg
}

// 喜欢列表
struct favorite_list_request {
    1:required i64 user_id
    2:required string token
}

struct favorite_list_response {
    1:required i32 status_code
    2:optional string status_msg
    3:required list<model.Video> video_list
}

// 评论操作
struct comment_action_request {
    1:required string token
    2:required i64 video_id
    3:required i32 action_type
    4:optional string comment_text
    5:optional i64 comment_id
}

struct comment_action_response {
    1:required i32 status_code
    2:optional string status_msg
    3:optional model.Comment comment
}

// 视频评论列表
struct comment_list_request {
    1:required string token
    2:required i64 video_id
}

struct comment_list_response {
    1:required i32 status_code
    2:optional string status_msg
    3:required list<model.Comment> comment_list
}

// rpc service
service InteractionService {
    favorite_action_response FavoriteAction(1:favorite_action_request req)
    favorite_list_response FavoriteList(1:favorite_list_request req)
    comment_action_response CommentAction(1:comment_action_request req)
    comment_list_response CommentList(1:comment_list_request req)
}
