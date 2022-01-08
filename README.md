## GIN-BOOT

### 功能描述
* http请求入参可以自动填充到结构体，如果是POST，则将http body数据填充到结构体；
如果是GET，则将URL query参数自动填充到结构体；
* 返回数据，可以是任意数据类型。如果数据不是boot.ApiResp，则返回数据会被包装为boot.ApiResp的json数据；
* 统一全局异常，请求返回会被包装为系统异常


### 使用demo
```
package main

import (
	"github.com/xiaojun207/gin-boot/boot"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Foo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 入参可以自动填充到结构体,如果是POST，则将http body数据填充到结构体；
// 返回数据，可以是任意数据类型。如果数据不是boot.ApiResp，则返回数据会被包装为boot.ApiResp的json数据；<br>
// 如果handler执行异常，请求返回会被包装为系统异常
func TestPost1Handler(c *gin.Context, req *Foo)  boot.ApiResp {
	log.Println("TestPost1Handler.req:", req.Username)
	return boot.ApiResp{
		Code: "100200",
		Msg:  "Success",
		Data: "TestData: " + req.Username,
	}
}

func TestPost2Handler(c *gin.Context, req *Foo) interface{} {
	log.Println("TestPost2Handler.req:", req.Username)
	return req
}

// 也可以使用gin方法, GET类型的请求，query参数也可以自动装填到结构体
func TestGet1Handler(c *gin.Context, req *Foo)  interface{} {
	log.Println("TestGet1Handler.req.username:", req.Username)
	log.Println("TestGet1Handler.req.password:", req.Password)
	data := map[string]interface{}{
		"list": []*Foo{
			req,
		},
		"page": "pageInfo",
	}
	return data
}


func AuthInterceptor(c *gin.Context) {
	authorization := c.GetHeader("authorization")
	log.Println("AuthInterceptor.authorization:", authorization)
	if authorization == "" {
		log.Println("AuthInterceptor authorization is null")
		boot.Resp(c, "105101", "账户未登录", "")
		c.Abort()
		return
	}

	//TODO getUid By authorization
	id := 100
	c.Set("uid", id)
}


var webRouter = func(router *boot.WebRouter) {
	//router.Use(AuthInterceptor)

	// 静态资源
	router.StaticFile("/", "./views/index.html")
	router.StaticFS("/static/", http.Dir("./views/static/"))

	// 动态API
	router.POST("/testPost", AuthInterceptor, TestPost1Handler)
	router.POST("/testPost2", TestPost2Handler)
	router.GET("/testGet", AuthInterceptor, TestGet1Handler)
}

func main() {
	boot.Start("8088", "/", webRouter)
}

```
