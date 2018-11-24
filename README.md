# web
A Go web app framework with views for constructing web pages, supports CGI &amp; FCGI.
```go
package main

import (
        "web"
        "fmt"
        "os"
)

var counter = 0

func hello(request *web.Request, response *web.Response) (err os.Error) {
        counter += 1
        response.Header.Set("Content-Type", "text/html")
        fmt.Fprintf(response.BodyWriter, "<b>test</b>: <small>num=%d</small><br/>\n", counter)
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
```
