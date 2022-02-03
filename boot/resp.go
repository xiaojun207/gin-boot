package boot

import (
	"github.com/gin-gonic/gin"
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

func Resp(c *gin.Context, code string, msg string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

func RespSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": CodeSuccess,
		"msg":  "成功",
		"data": data,
	})
}
