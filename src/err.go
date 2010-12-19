package web

import (
        "fmt"
)

func error(f string, args ...interface{}) {
        panic(fmt.Sprintf(f, args...))
}
