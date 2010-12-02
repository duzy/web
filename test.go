package main

import (
        "os"
        "io"
        "fmt"
        "flag"
        "./_obj/web"
)

var counters map[string]int = make(map[string]int)

func SayHello(w io.Writer, app *web.App) {
        counters["hello"] += 1
        app.SetHeader("Content-Type", "text/plain")
        fmt.Fprintf(w, "hello: %d\n", counters["hello"])
}

type SaveSomething struct {
        text string
}

func (this *SaveSomething) WriteContent(w io.Writer, app *web.App) {
        counters["save"] += 1
        app.SetHeader("Content-Type", "text/html")
        fmt.Fprintf(w, "save a posted message...(%d)\n", counters["save"])
}

var path = flag.String("path", "", "specify a test PATH_INFO")
var help = flag.Bool("help", false, "show help message")

func printUsage() {
        fmt.Printf(`usage: ./test [-path PATH_INFO]
where the PATH_INFO may be one of:
    /
    /hello
    /view
    /save`)
        fmt.Println()
}

func main() {
        flag.Parse()

        if *help {
                printUsage()
                os.Exit(1)
        }

        if *path != "" && os.Getenv("PATH_INFO") == "" {
                os.Setenv("PATH_INFO", *path)
        }

        hello := web.FuncHandler(SayHello)
        app := web.NewApp("test app", web.NewCGIModel())
        app.HandleDefault(hello) // app.Handle("", hello); app.Handle("/", hello)
        app.Handle("/hello", hello)
        app.Handle("/view", web.NewView("test.tpl"))
        app.Handle("/save", new(SaveSomething))
        app.Exec()
}
