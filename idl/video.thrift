include "model.thrift"

namespace go video

// 视频投稿
struct publish_action_request {
    1:required string token
    2:required binary data
    3:required string title
}

struct publish_action_response {
    1:required i32 status_code
    2:optional string status_msg
}
// 视频流
struct feed_request {
    1:optional i64 latest_time
    2:optional string token
}

struct feed_response {
    1:required i32 status_code
    2:optional string status_msg
    3:required list<model.Video> video_list
    4:optional i64 next_time
}


// 发布列表
struct publish_list_request {
    1:required i64 user_id
    2:required string token
}

struct publish_list_response {
    1:required i32 status_code
    2:optional string status_msg
    3:required list<model.Video> video_list
}

service VideoService {
    feed_response Feed(1:feed_request req)
    publish_action_response PublishAction(1:publish_action_request req)
    publish_list_response PublishList(1:publish_list_request req)
}