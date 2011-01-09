package web

import (
        "fmt"
        //"log"
        "runtime"
        "bytes"
)

// type errorType struct {
//         log
// }

func error(f string, args ...interface{}) {
        calldepth := 5
        _, file, line, ok := runtime.Caller(calldepth)
        if !ok {
                file = "???"
                line = 0
        }

        msg := bytes.NewBuffer(make([]byte, 0, 64))

        fmt.Fprintf(msg, "%s:%d: ", file, line)
        fmt.Fprintf(msg, f, args...)

        panic(msg.String())
}
