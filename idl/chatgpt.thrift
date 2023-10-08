include "model.thrift"

namespace go gpt

struct gpt_chat_request {
    1:required string msg
}

struct gpt_chat_response{
    1:required string msg
    2:required i32 status
}

service ChatgptService {
    gpt_chat_response GptChat(1:gpt_chat_request req)
}