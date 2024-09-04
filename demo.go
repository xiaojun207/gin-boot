package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/xiaojun207/gin-boot/boot"
	"github.com/xiaojun207/gin-boot/i18n"
	"golang.org/x/text/language"
	"log"
	"net/http"
)

type Foo struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password"`
}

type Page struct {
	boot.BindQuery     // 继承boot.BindQuery的结构体，指定绑定到url参数
	PageNum        int `json:"page_num" form:"page_num"`   // url 中的参数，需要用tag 'form' 标识，才能自动绑定
	PageSize       int `json:"page_size" form:"page_size"` // url 中的参数，需要用tag 'form' 标识，才能自动绑定
}

type QueryHeader struct {
	boot.BindHeader
	Authorization string `header:"authorization"`
}

// TestPost1Handler
// /testPost
// 入参可以自动填充到结构体,如果是POST，则将http body数据填充到结构体；
// 返回数据，可以是任意数据类型。如果数据不是boot.ApiResp，则返回数据会被包装为boot.ApiResp的json数据；<br>
// 如果handler执行异常，请求返回会被包装为系统异常
// 参数page继承boot.BindQuery的结构体，绑定url参数
func TestPost1Handler(c *gin.Context, req *Foo, page Page, header QueryHeader) boot.ApiResp {
	log.Println("TestPost1Handler.req:", req.Username, ",Password:", req.Password)
	log.Println("TestPost1Handler.page.PageNum:", page.PageNum, ",PageSize:", page.PageSize)
	log.Println("TestPost1Handler.Authorization:", header.Authorization)
	return boot.ApiResp{Code: boot.CodeSuccess, Msg: "Success", Data: "TestData: " + req.Username}
}

// TestGetEmptyHandler
/*
/testGetEmpty
空返回值包装测试，返回：{"code":"100200","data":null,"msg":"成功"}
*/
func TestGetEmptyHandler(c *gin.Context, req *Foo) {
	log.Println("TestGetEmptyHandler.req.username:", req.Username, ",password:", req.Password)
}

// TestGetHandler
// 空返回值包装测试，返回：{"code":"100200","data":null,"msg":"成功"}
func TestGetHandler(c *gin.Context, page Page, header QueryHeader) {
	log.Println("TestGetHandler.page.PageNum:", page.PageNum, ",PageSize:", page.PageSize)
	log.Println("TestGetHandler.Authorization:", header.Authorization)
	log.Println("TestGetHandler")
}

// TestPost2Handler
// 异常全局处理测试，返回：{"code":"100101","data":null,"msg":"TestPost2Handler.TestError"}
func TestPost2Handler(c *gin.Context, req *Foo) {
	log.Println("TestPost2Handler.req.username:", req.Username, ",password:", req.Password)
	panic(errors.New("TestPost2Handler.TestError"))
}

// TestGet1Handler
// /testGet?username=admin12&password=1234&page_num=1&page_size=10
// url 中的参数，需要用tag 'form' 标识
// 也可以使用gin方法, GET类型的请求，query参数也可以自动装填到结构体
func TestGet1Handler(c *gin.Context, req *Foo, page Page) interface{} {
	log.Println("TestGet1Handler.req.username:", req.Username, ",password:", req.Password)
	log.Println("TestGet1Handler.page.PageNum:", page.PageNum, ",PageSize:", page.PageSize)
	data := map[string]interface{}{
		"list": []*Foo{req},
		"page": page,
	}
	log.Println("TestGet1Handler.data:", data)
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
	bundle := i18n.NewBundle(
		language.Chinese, //默认中文
		"locals/message.zh-CN.toml",
		"locals/message.en-US.toml",
	)
	router.Use(i18n.Serve(bundle))

	// 静态资源
	router.StaticFile("/", "./views/index.html")
	router.StaticFS("/static/", http.Dir("./views/static/"))

	// 动态API
	router.POST("/testPost", AuthInterceptor, TestPost1Handler)
	router.POST("/testPost2", TestPost2Handler)
	router.GET("/testGetEmpty", TestGetEmptyHandler)
	router.GET("/testGet", TestGetHandler)
	router.GET("/testGet1", AuthInterceptor, TestGet1Handler)

	apiRouter := router.Group("/api/")
	apiRouter.GET("/test", TestPost2Handler)
}

func CustomCode() {
	boot.CodeSuccess = "100200"
	boot.CodeInvalidRequest = "100104"
	boot.CodeServerError = "100100" // 服务错误
}

func main() {
	CustomCode()
	boot.Start("8088", "/", webRouter)
}
