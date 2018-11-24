package main

import (
        "github.com/duzy/web"
        "os"
        "fmt"
)

type ContentFeed struct {
        COUNTER int
}

type TestPageFeed struct {
        TITLE string
        CONTENT web.View
}

func (feed *ContentFeed) Prepare(request *web.Request) (err os.Error) {
        feed.COUNTER += 1
        return
}

func (feed *TestPageFeed) Prepare(request *web.Request) (err os.Error) {
        //feed.CONTENT.Prepare(request)
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

        contentFeed := &ContentFeed{ COUNTER: 0 }
        content, err := web.NewTemplateFromString("counter: {COUNTER}", contentFeed)
        if err != nil {
                fmt.Printf("error: %v\n", err)
                return
        }

        pageFeed := &TestPageFeed{
        TITLE: "FCGI test page",
        CONTENT: content,
        }
        page, err := web.NewHtmlPage("page.tpl", pageFeed)
        if err != nil {
                fmt.Printf("error: %v\n", err)
                return
        }

        app.HandleDefault(page)
        app.Exec()
}
