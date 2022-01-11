package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/xiaojun207/gin-boot/boot"
	"log"
	"net/http"
)

type Foo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Page struct {
	PageNum  int `json:"page_num,string"`
	PageSize int `json:"page_size,string"`
}

// 入参可以自动填充到结构体,如果是POST，则将http body数据填充到结构体；
// 返回数据，可以是任意数据类型。如果数据不是boot.ApiResp，则返回数据会被包装为boot.ApiResp的json数据；<br>
// 如果handler执行异常，请求返回会被包装为系统异常
func TestPost1Handler(c *gin.Context, req *Foo) boot.ApiResp {
	log.Println("TestPost1Handler.req:", req.Username)
	return boot.ApiResp{Code: "100200", Msg: "Success", Data: "TestData: " + req.Username}
}

// 空返回值包装测试，返回：{"code":"100200","data":null,"msg":"成功"}
func TestGetEmptyHandler(c *gin.Context, req *Foo) {
	log.Println("TestGetEmptyHandler.req.username:", req.Username, ",password:", req.Password)
}

// 空返回值包装测试，返回：{"code":"100200","data":null,"msg":"成功"}
func TestGetHandler(c *gin.Context) {
	log.Println("TestGetHandler")
}

// 异常全局处理测试，返回：{"code":"100101","data":null,"msg":"TestPost2Handler.TestError"}
func TestPost2Handler(c *gin.Context, req *Foo) {
	log.Println("TestPost2Handler.req.username:", req.Username, ",password:", req.Password)
	panic(errors.New("TestPost2Handler.TestError"))
}

// /testGet?username=admin12&password=1234&page_num=1&page_size=10
// 也可以使用gin方法, GET类型的请求，query参数也可以自动装填到结构体
func TestGet1Handler(c *gin.Context, req *Foo, page Page) interface{} {
	log.Println("TestGet1Handler.req.username:", req.Username, ",password:", req.Password)
	log.Println("TestGet1Handler.page.PageNum:", page.PageNum, ",PageSize:", page.PageSize)
	data := map[string]interface{}{
		"list": []*Foo{req},
		"page": page,
	}
	return data
}

func AuthInterceptor(c *gin.Context, header http.Header) {
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
	router.GET("/testGetEmpty", TestGetEmptyHandler)
	router.GET("/testGet", AuthInterceptor, TestGet1Handler)
	router.GET("/testGet2", TestGetHandler)

	apiRouter := router.Group("/api/")
	apiRouter.GET("/test", TestPost2Handler)
}

func main() {
	boot.Start("8088", "/", webRouter)
}
