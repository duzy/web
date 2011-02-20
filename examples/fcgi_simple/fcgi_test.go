package main

import (
        "ds/web"
        //"io"
        "os"
        "fmt"
)

var counter = 0

func hello(request *web.Request, response *web.Response) (err os.Error) {
        counter += 1
        response.Header["Content-Type"] = "text/html"
        fmt.Fprintf(response.BodyWriter, "<b>test</b>: <small>num=%d</small><br/>\n", counter)

        // if request.Session() == nil {
        //         fmt.Fprintf(response.BodyWriter, "request.Session() == nil\n")
        //         return
        // }
        // s := request.Session().Get("test")
        // fmt.Fprintf(response.BodyWriter, "test: %s\n", s)
        // s = fmt.Sprintf("%s%d", s, counter)
        // request.Session().Set("test", s)
        return
}

func main() {
        m, err := web.NewFCGIModel()
        if err != nil {
                fmt.Printf("%v\n", err)
                return
        }

        app, err := web.NewApp(m)
        if err != nil {
                fmt.Printf("%v\n", err)
                return
        }

        if app == nil {
                fmt.Printf("no app\n")
                return
        }

        app.HandleDefault(web.FuncHandler(hello))
        app.Exec()
}
