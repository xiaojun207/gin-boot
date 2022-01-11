package boot

import (
	"fmt"
	"testing"
)

func TestErrorToString(t *testing.T) {
	e := NewWebError("100100", "s")
	code, msg := errorToString(e)
	fmt.Println("code:", code, "msg:", msg)
}
