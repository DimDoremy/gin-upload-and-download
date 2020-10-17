package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type MessageJson struct {
	Message  string `json:"message"`
	Filename string `json:"filename"`
}

type MessageQueue struct {
	Message chan interface{}
}

func (mq *MessageQueue) push(message interface{}) {
	mq.Message <- message
}

func (mq MessageQueue) pop() interface{} {
	return <-mq.Message
}

var mq MessageQueue

func MessageSendApi(c *gin.Context) {
	msg := mq.pop()
	c.JSON(http.StatusOK, msg)
}

func MessageReceiveApi(c *gin.Context) {
	var json MessageJson
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求体内容解析失败"})
		return
	}
	mq.push(json)
	c.String(http.StatusOK, "OK")
}
