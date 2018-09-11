package helpers

import (
	"fmt"
	"os"
)

func Exit(code int, a ...interface{}) {
	fmt.Println(a...)
	os.Exit(code)
}

func Exitf(code int, format string, a ...interface{}) {
	Exit(code, fmt.Sprintf(format, a...))
}
