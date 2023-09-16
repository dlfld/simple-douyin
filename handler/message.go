package handler

import (
	"net/http"
	"strconv"

	"github.com/douyin/common/constant"
	"github.com/douyin/kitex_gen/message"
	"github.com/gin-gonic/gin"
)

// MessageAction @Summary 消息操作
// @Schemes
// @Description 登录用户对消息的相关操作，目前只支持消息发送
// @Tags 消息接口
// @Accept json
// @Produce json
// @Param token query message.MessageActionRequest true "Message Action Params"
// @Param token query string true "用户鉴权token"
// @Router /douyin/message/action/ [POST]
func MessageAction(c *gin.Context) {

	// 2. 创建发生消息的请求实例
	req := message.NewMessageActionRequest()
	// 3. 前端请求数据绑定到req中
	err := c.ShouldBind(req)
	if err != nil {
		c.JSON(http.StatusOK, constant.NewErrResp(constant.ErrBadRequest))
		return
	}
	if req.Content == "" {
		c.JSON(http.StatusOK, constant.NewErrResp(constant.ErrEmptyMessage))
		return
	}
	if len(req.Content) > 1000 {
		c.JSON(http.StatusOK, constant.NewErrResp(constant.ErrMsgTooLong))
		return
	}
	req.FromUserId = int64(c.GetUint("userID"))
	// 4. 发起RPC调用
	resp, err := rpcCli.messageCli.SendMessage(c, req)
	if err != nil {
		resp = new(message.MessageActionResponse)
		constant.HandlerErr(constant.ErrSendMessage, resp)
	}
	// 5. 处理错误返回响应
	c.JSON(http.StatusOK, resp)
}

// MessageChat @Summary 消息查询
// @Schemes
// @Description 当前登录用户和其他指定用户的聊天消息记录
// @Tags 消息接口
// @Accept json
// @Produce json
// @Param token query message.MessageChatRequest true "Message Action Params"
// @Param token query string true "用户鉴权token"
// @Router /douyin/message/chat/ [GET]
func MessageChat(c *gin.Context) {
	// 2. 创建请求实例
	req := message.NewMessageChatRequest()
	// 3. 前端请求数据绑定到req中
	uid := c.GetUint("userID")
	toUserID, err := strconv.Atoi(c.Query("to_user_id"))
	if err != nil {
		c.JSON(http.StatusOK, constant.NewErrResp(constant.ErrBadRequest))
		return
	}
	timeStr := c.Query("pre_msg_time")
	preMsgTime, err := strconv.Atoi(timeStr)
	if err != nil {
		c.JSON(http.StatusOK, constant.NewErrResp(constant.ErrBadRequest))
		return
	}
	req.FromUserId = int64(uid)
	req.ToUserId = int64(toUserID)
	req.PreMsgTime = int64(preMsgTime)
	// 4. 发起RPC调用
	resp, err := rpcCli.messageCli.MessageList(c, req)
	if err != nil {
		resp = new(message.MessageChatResponse)
		constant.HandlerErr(constant.ErrMessageList, resp)
	}
	// 5. 处理错误返回响应
	c.JSON(http.StatusOK, resp)
}
