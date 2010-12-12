package main

import (
        "os"
        "io"
        "fmt"
        "flag"
        "bytes"
        "./_obj/web"
)

var counters map[string]int = make(map[string]int)

func SayHello(w io.Writer, app *web.App) {
        counters["hello"] += 1
        app.SetHeader("Content-Type", "text/plain")
        fmt.Fprintf(w, "hello: %d\n", counters["hello"])

        s := app.Session().Get("test")
        fmt.Fprintf(w, "test: %s\n", s)

        s = fmt.Sprintf("%s%d", s, counters["hello"])
        app.Session().Set("test", s)
}

type SaveSomething struct {
        text string
}

func (this *SaveSomething) WriteContent(w io.Writer, app *web.App) {
        counters["save"] += 1
        app.SetHeader("Content-Type", "text/html")

        r := app.RequestReader()
        msg := bytes.NewBuffer(make([]byte, 0))
        n, err := io.Copy(msg, r)
        if err != nil {
                //...
        }

        fmt.Fprintf(w, "posted message(%d): [%d]%s\n", counters["save"],
                n, string(msg.Bytes()))
}

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
        path := flag.String("path", "", "specify a test PATH_INFO")
        session := flag.String("session", "", "specify a test session id")
        help := flag.Bool("help", false, "show help message")
        flag.Parse()

        if *help {
                printUsage()
                os.Exit(1)
        }

        if *path != "" && os.Getenv("PATH_INFO") == "" {
                os.Setenv("PATH_INFO", *path)
        }

        if *session != "" {
                //fmt.Fprintf(os.Stdout, "-session: %s\n", *session)
                cookie := "__web_cid=" + *session
                os.Setenv("HTTP_COOKIE", cookie)
        }

        hello := web.FuncHandler(SayHello)
        app := web.NewApp("test app", web.NewCGIModel())
        app.HandleDefault(hello) // app.Handle("", hello); app.Handle("/", hello)
        app.Handle("/hello", hello)
        app.Handle("/view", web.NewView("test.tpl"))
        app.Handle("/save", new(SaveSomething))
        app.Exec()
}
