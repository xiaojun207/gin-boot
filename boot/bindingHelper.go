package boot

import (
	"encoding/json"
	"github.com/gin-gonic/gin/binding"
	"github.com/xiaojun207/go-base-utils/utils"
	"net/http"
	"net/url"
	"reflect"
)

type IBindBody interface {
	getBindBodyValue() string
}
type BindBody struct{}

func (e BindBody) getBindBodyValue() {}

// IBindQuery
type IBindQuery interface {
	getBindQueryValue() string
}
type BindQuery struct{}

func (e BindQuery) getBindQueryValue() {}

// IBindHeader
type IBindHeader interface {
	getBindHeaderValue() string
}
type BindHeader struct{}

func (e BindHeader) getBindHeaderValue() {}

// IBindCookie
type IBindCookie interface {
	getBindCookieValue() string
}
type BindCookie struct{}

func (e BindHeader) getBindCookieValue() {}

// CookieBinding
type CookieBinding struct{}

func (e CookieBinding) Name() string {
	return "cookie"
}

func (e CookieBinding) Bind(typ reflect.Type, req *http.Request, obj interface{}) error {
	d, err := cookieToJson(req.Cookies())
	if err != nil {
		return err
	}
	err, obj = utils.NewInterface(typ, d)
	if err != nil {
		return err
	}
	if binding.Validator == nil {
		return nil
	}
	return binding.Validator.ValidateStruct(obj)
}

var cookieBinding = CookieBinding{}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func queryToJson(queryValues url.Values) (res []byte, err error) {
	resMap := map[string]interface{}{}
	for key := range queryValues {
		resMap[key] = queryValues.Get(key)
	}
	res, err = json.Marshal(resMap)
	return
}

///////
func cookieToJson(cookies []*http.Cookie) (res []byte, err error) {
	resMap := map[string]interface{}{}
	for _, c := range cookies {
		resMap[c.Name] = c.Value
	}
	res, err = json.Marshal(resMap)
	return
}
