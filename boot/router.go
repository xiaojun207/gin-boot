package boot

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

type WebHandlerFunc interface{}

var (
	WebErrorType      = reflect.TypeOf(NewWebError("100200", ""))
	ApiRespType       = reflect.TypeOf(ApiResp{})
	ErrorType         = reflect.TypeOf(errors.New(""))
	GinContextType    = reflect.TypeOf(gin.Context{})
	PtrGinContextType = reflect.TypeOf(&gin.Context{})
	HeaderType        = reflect.TypeOf(http.Header{})
	RequestType       = reflect.TypeOf(http.Request{})
	PtrRequestType    = reflect.TypeOf(&http.Request{})
	ResponseType      = reflect.TypeOf(http.Response{})
	PtrResponseType   = reflect.TypeOf(&http.Response{})
)

type WebRouter struct {
	Router *gin.RouterGroup
}

func (e *WebRouter) Group(relativePath string, handlers ...gin.HandlerFunc) *WebRouter {
	return &WebRouter{Router: e.Router.Group(relativePath, handlers...)}
}

func (e *WebRouter) Use(middleware ...gin.HandlerFunc) {
	e.Router.Use(middleware...)
}

func (e *WebRouter) StaticFile(relativePath, filepath string) {
	e.Router.StaticFile(relativePath, filepath)
}

func (e *WebRouter) StaticFS(relativePath string, fs http.FileSystem) {
	e.Router.StaticFS(relativePath, fs)
}

func (e *WebRouter) GET(relativePath string, webHandlers ...WebHandlerFunc) {
	e.route(http.MethodGet, relativePath, webHandlers...)
}

func (e *WebRouter) POST(relativePath string, webHandlers ...WebHandlerFunc) {
	e.route(http.MethodPost, relativePath, webHandlers...)
}

func (e *WebRouter) PUT(relativePath string, webHandlers ...WebHandlerFunc) {
	e.route(http.MethodPut, relativePath, webHandlers...)
}

func (e *WebRouter) DELETE(relativePath string, webHandlers ...WebHandlerFunc) {
	e.route(http.MethodDelete, relativePath, webHandlers...)
}

func (e *WebRouter) OPTIONS(relativePath string, webHandlers ...WebHandlerFunc) {
	e.route(http.MethodOptions, relativePath, webHandlers...)
}

func (e *WebRouter) Any(relativePath string, webHandlers ...WebHandlerFunc) {
	e.route(http.MethodGet, relativePath, webHandlers)
	e.route(http.MethodPost, relativePath, webHandlers)
	e.route(http.MethodPut, relativePath, webHandlers)
	e.route(http.MethodPatch, relativePath, webHandlers)
	e.route(http.MethodHead, relativePath, webHandlers)
	e.route(http.MethodOptions, relativePath, webHandlers)
	e.route(http.MethodDelete, relativePath, webHandlers)
	e.route(http.MethodConnect, relativePath, webHandlers)
	e.route(http.MethodTrace, relativePath, webHandlers)
}

func (e *WebRouter) route(method, relativePath string, webHandlers ...WebHandlerFunc) {
	router(e.Router, method, relativePath, webHandlers...)
}

func router(router *gin.RouterGroup, method, relativePath string, webHandlers ...WebHandlerFunc) {
	handlers := make([]gin.HandlerFunc, len(webHandlers))
	for i, webHandler := range webHandlers {
		tmp := webHandler // 这里需要把webHandler赋值到临时变量，否则会被
		lastHandler := i == len(webHandlers)-1
		handlers[i] = func(c *gin.Context) {
			// 这里只能用变量tmp，不能用webHandler
			baseFunc(c, tmp, lastHandler)
		}
	}
	router.Handle(method, relativePath, handlers...)
}

func baseFunc(c *gin.Context, webHandler WebHandlerFunc, lastHandler bool) {
	funType := reflect.TypeOf(webHandler)
	ginFunType := reflect.TypeOf(func(c *gin.Context) {})
	if funType == ginFunType {
		webHandler.(func(*gin.Context))(c)
	} else {
		params := AutoFillParams(c, webHandler)
		returnValues := reflect.ValueOf(webHandler).Call(params)

		if returnValues != nil && len(returnValues) > 0 {
			r := returnValues[0].Interface()

			rType := reflect.TypeOf(r)
			if rType == WebErrorType {
				RespWebError(c, r.(WebError))
			} else if rType == ApiRespType {
				c.JSON(http.StatusOK, r)
			} else if rType == ErrorType {
				Resp(c, "100100", "系统错误", r)
			} else {
				RespSuccess(c, r)
			}
		}
	}

	// 如果最后一个handler结束，并且没有返回数据，则默认成功
	if lastHandler && !c.Writer.Written() {
		RespSuccess(c, nil)
	}
}
