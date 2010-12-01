package main

import (
        "./app"
        "io"
        "fmt"
)

func SayHello(w io.Writer, app *web.App) {
        fmt.Fprintf(w, "hello\n")
}

type SaveSomething struct {
        text string
}

func (this *SaveSomething) WriteResponse(w io.Writer, app *web.App) {
        fmt.Fprintf(w, "Content-Type: text/html\n\n")
        fmt.Fprintf(w, "save a posted message...\n")
}

func main() {
        app := web.NewApp("test app")
        app.Handle("", web.FuncHandler(SayHello))
        app.Handle("/hello", web.FuncHandler(SayHello))
        app.Handle("/view", web.NewView("test.tpl"))
        app.Handle("/save", new(SaveSomething))
        app.Exec()
}
