package boot

import "github.com/gin-gonic/gin"

func _404Handler(c *gin.Context) {
	Resp(c, "100104", "无效的请求", "")
}
