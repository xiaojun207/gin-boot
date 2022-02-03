package boot

import (
	"fmt"
	"testing"
)

func TestErrorToString(t *testing.T) {
	e := NewWebError(CodeServerError, "s")
	code, msg := errorToString(e)
	fmt.Println("code:", code, "msg:", msg)
}
