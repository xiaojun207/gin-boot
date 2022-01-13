package boot

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
)

// 获取数据
func loadData(request *http.Request) (body []byte) {
	var err error
	if request.Method == "GET" {
		body, err = queryToJson(request.URL.Query())
	} else {
		body, err = ioutil.ReadAll(request.Body)
	}
	if err != nil {
		log.Println("loadData.err:", err)
	}
	return body
}

// Automatic filling parameters
func AutoFillParams(c *gin.Context, fun interface{}) []reflect.Value {
	funType := reflect.TypeOf(fun)
	//NumIn:返回func类型的参数个数，如果不是函数，将会panic
	values := make([]reflect.Value, funType.NumIn())
	var body []byte
	loadBody := func() {
		if len(body) == 0 {
			body = loadData(c.Request)
		}
	}

	for i := 0; i < funType.NumIn(); i++ {
		paramType := funType.In(i)
		if paramType == PtrGinContextType {
			values[i] = reflect.ValueOf(c)
		} else if paramType == GinContextType {
			values[i] = reflect.ValueOf(*c)
		} else if paramType == HeaderType {
			values[i] = reflect.ValueOf(c.Request.Header)
		} else if paramType == RequestType {
			values[i] = reflect.ValueOf(*c.Request)
		} else if paramType == PtrRequestType {
			values[i] = reflect.ValueOf(c.Request)
		} else if paramType == ResponseType {
			values[i] = reflect.ValueOf(*c.Request.Response)
		} else if paramType == PtrResponseType {
			values[i] = reflect.ValueOf(c.Request.Response)
		} else if paramType.Kind() == reflect.Struct || paramType.Kind() == reflect.Ptr {
			err, pObj := bindParam(paramType, c)
			if err != nil {
				log.Println("AutoFillParams.err:", err, ", Type:", reflect.TypeOf(err))
				panic(err)
			}
			values[i] = reflect.ValueOf(pObj)
		} else if paramType.Kind() == reflect.String {
			loadBody()
			values[i] = reflect.ValueOf(string(body))
		} else {
			log.Println("Param.name:", paramType.Elem().Name())
		}
	}
	return values
}

func bindParam(typ reflect.Type, c *gin.Context) (error, interface{}) {
	shouldBind := func(pObj interface{}) error {
		switch pObj.(type) {
		case IBindHeader:
			return c.ShouldBindHeader(pObj)
		case IBindQuery:
			return c.ShouldBindQuery(pObj)
		case IBindBody:
			return c.ShouldBindBodyWith(pObj, binding.JSON)
		case IBindCookie:
			return cookieBinding.Bind(typ, c.Request, pObj)
		default:
			if c.Request.Method == http.MethodGet {
				return c.ShouldBindQuery(pObj)
			} else {
				return c.ShouldBindBodyWith(pObj, binding.JSON)
			}
		}
	}

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		dst := reflect.New(typ).Elem()
		err := shouldBind(dst.Addr().Interface())
		return err, dst.Addr().Interface()
	} else {
		dst := reflect.New(typ).Elem()
		err := shouldBind(dst.Addr().Interface())
		return err, dst.Interface()
	}
}
