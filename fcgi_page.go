package main

import (
        "./_obj/web"
        //"io"
        "os"
        "fmt"
)

type ContentFeed struct {
        counter int
}

type TestPageFeed struct {
        title string
        content web.View
}

func (feed *ContentFeed) Prepare(request *web.Request) (err os.Error) {
        feed.counter += 1
        return
}

func (feed *TestPageFeed) Prepare(request *web.Request) (err os.Error) {
        //feed.content.Prepare(request)
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

        contentFeed := &ContentFeed{ 0 }
        content, err := web.NewTemplateFromString("counter: {counter}", contentFeed)
        if err != nil {
                fmt.Printf("error: %v\n", err)
                return
        }

        feed := &TestPageFeed{ "FCGI test page", content }
        page, err := web.NewHtmlPage("page.tpl", feed)
        if err != nil {
                fmt.Printf("error: %v\n", err)
                return
        }

        app.HandleDefault(page)
        app.Exec()
}
