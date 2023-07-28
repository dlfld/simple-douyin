include "model.thrift"

namespace go message

// 查询消息
struct message_chat_request {
    1:required string token  // 用户鉴权token
    2:required i64 to_user_id // 对方用户id
    3:required i64 pre_msg_time //上次最新消息的时间（新增字段-apk更新中）
}

struct message_chat_response {
    1:required i32 status_code                  // 状态码 0-成功， 其他值-失败
    2:optional string status_msg                // 返回状态描述
    3:required list<model.Message> message_list // 消息列表
}

// 发送消息
struct message_action_request {
    1: required string token    // 用户鉴权token
    2: required i64 to_user_id  // 对方用户id
    3: required i32 action_type // 1-发送消息
    4: required string content  // 消息内容

}

struct message_action_response {
    1: required i32 status_code  // 状态码 0-成功， 其他值-失败
    2: optional string status_msg
}

service MessageService{
    message_chat_response MessageList(1:message_chat_request req)
    message_action_response SendMessage(1:message_action_request req)
}