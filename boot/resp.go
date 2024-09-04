package boot

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaojun207/gin-boot/i18n"
	"net/http"
)

type ApiResp struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func RespWebError(c *gin.Context, err WebError) {
	Resp(c, err.Code(), err.Msg(), "")
}

func RespSuccess(c *gin.Context, data interface{}) {
	Resp(c, CodeSuccess, "成功", data)
}

func Resp(c *gin.Context, code string, msg string, data interface{}) {
	msg = i18n.FormatText(c, &i18n.Message{ID: msg, Other: msg})
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}
