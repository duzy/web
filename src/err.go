package web

import (
        "fmt"
        //"log"
        "runtime"
        "bytes"
        "os"
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

func newError(msg string) os.Error {
        stack := make([]uintptr, 30)

        n := runtime.Callers(0, stack)
        stack = stack[0:n]

        str := bytes.NewBuffer(make([]byte, 0, 256))

        fmt.Fprintf(str, "%s\n", msg)
        for i := range stack {
                pc := stack[len(stack)-i-1]
                f := runtime.FuncForPC(pc)
                if f != nil {
                        file, line := f.FileLine(pc)
                        fmt.Fprintf(str, "%s:%d: %s\n", file, line, f.Name())
                }
        }
        return os.NewError(str.String())
}
