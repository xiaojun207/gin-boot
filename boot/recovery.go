package boot

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"runtime/debug"
)

// recover错误，转string
func errorToString(r interface{}) (code string, msg string) {
	switch v := r.(type) {
	case WebError:
		return v.Code(), v.Msg()
	case validator.ValidationErrors:
		return "100101", v.Error()
	case error:
		return "100103", v.Error()
	default:
		return "100103", r.(string)
	}
}

// 全局错误处理
func Recovery(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			//打印错误堆栈信息
			log.Printf("panic: %v\n", r)
			if gin.Mode() == gin.DebugMode {
				debug.PrintStack()
			}

			code, msg := errorToString(r)
			Resp(c, code, msg, "")
			//终止后续接口调用，不加的话recover到异常后，还会继续执行接口里后续代码
			c.Abort()
		}
	}()
	//加载完 defer recover，继续后续接口调用
	c.Next()
}
