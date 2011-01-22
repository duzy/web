package main

import (
        "./_obj/web"
        //"io"
        "os"
        "fmt"
)

var counter = 0

func hello(request *web.Request, response *web.Response) (err os.Error) {
        counter += 1
        response.Header["Content-Type"] = "text/html"
        fmt.Fprintf(response.BodyWriter, "<b>test</b>: <small>num=%d</small>\n", counter)

        if request.Session() == nil {
                fmt.Fprintf(response.BodyWriter, "request.Session() == nil\n")
                return
        }

        s := request.Session().Get("test")
        fmt.Fprintf(response.BodyWriter, "test: %s\n", s)

        s = fmt.Sprintf("%s%d", s, counter)
        request.Session().Set("test", s)
        return
}

func main() {
        app, err := web.NewApp(web.NewFCGIModel())
        if err != nil {
                // ...
        }

        app.HandleDefault(web.FuncHandler(hello))
        app.Exec()
}
