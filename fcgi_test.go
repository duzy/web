package main

import (
        "./_obj/web"
        "io"
        "fmt"
)

var counter = 0

func hello(w io.Writer, app *web.App) {
        counter += 1
        app.SetHeader("Content-Type", "text/html")
        fmt.Fprintf(w, "<b>test</b>: <small>num=%d</small>\n", counter)

        s := app.Session().Get("test")
        fmt.Fprintf(w, "test: %s\n", s)

        s = fmt.Sprintf("%s%d", s, counter)
        app.Session().Set("test", s)
}

func main() {
        app := web.NewApp(web.NewFCGIModel())
        app.HandleDefault(web.FuncHandler(hello))
        app.Exec()
}
