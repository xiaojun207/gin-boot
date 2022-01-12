package boot

import (
	"log"
	"testing"
)

type TestQuery2 struct {
	BindQuery
}

func TestQueryParam(t *testing.T) {
	q := TestQuery2{}
	var p interface{}
	p = q
	switch v := p.(type) {
	case error:
		log.Println("v.error", v)
		break
	case BindQuery:
		log.Println("v.GetValue", v)
		break
	default:
		log.Println("TestQueryParam.default")

	}
}
