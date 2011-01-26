package web

import (
        "bytes"
        "os"
        "io"
        "io/ioutil"
        "template"
        "fmt"
)

type RenderBuffer struct {
        io.Writer
}

type Renderable interface {
        Render(w RenderBuffer) (err os.Error)
}

type Prepareable interface {
        Prepare(request *Request) (err os.Error)
}

type PrepareableRenderable interface {
        Renderable
        Prepareable
}

type View interface {
        PrepareableRenderable
        fmt.Stringer
}

type TemplateFeed interface {
        PrepareFeed(request *Request) (err os.Error)
}

/**
 Usage:

 feed := new(MyFeedOrViewModel)
 page, err := web.NewPage("home.html", feed)
 app.Handle("/home", page)

 */
type Page struct {
        templateView
}

type HtmlPage Page

func NewHtmlPage(fn string, feed TemplateFeed) (page *HtmlPage, err os.Error) {
        p, err := NewPage(fn, feed)
        page = (*HtmlPage)(p)
        return
}

func NewPage(fn string, feed TemplateFeed) (page *Page, err os.Error) {
        f, err := os.Open(fn, os.O_RDONLY, 0)
        if err != nil {
                return
        }

        s, err := ioutil.ReadAll(f)
        if err != nil {
                return
        }

        page, err = NewPageFromString(string(s), feed)
        return
}

func NewPageFromString(s string, feed TemplateFeed) (page *Page, err os.Error) {
        var temp *template.Template
        temp, err = template.Parse(s, nil/* TODO: make use of it? */)
        if err != nil {
                return
        }
        
        page = &Page{ templateView{ temp, feed } }
        return
}

func NewHtmlPageFromString(s string, feed TemplateFeed) (page *HtmlPage, err os.Error) {
        var pg *Page
        pg, err = NewPageFromString(s, feed)
        if err != nil {
                page = (*HtmlPage)(pg)
        }
        return
}

func (p *Page) HandleRequest(request *Request, response *Response) (err os.Error) {
        buf := RenderBuffer{ response.BodyWriter }

        err = p.templateView.Prepare(request)
        if err != nil {
                return
        }

        err = p.templateView.Render(buf)
        return
}

/*
func (vm *Page) HandleSubpath(subpath string, request *Request) (handled bool) {
        return
}
 */

func (p *HtmlPage) HandleRequest(request *Request, response *Response) (err os.Error) {
        response.Header["Content-Type"] = "text/html"
        err = (*Page)(p).HandleRequest(request, response)
        return
}

type viewStringer struct { PrepareableRenderable }

func (v viewStringer) String() string {
        b := bytes.NewBuffer(make([]byte, 0, 512))
        err := v.Render(RenderBuffer{ b })
        if err != nil {
                return err.String()
        }
        return b.String()
}

type nothingPrepareable struct { Renderable }

func (v *nothingPrepareable) Prepare(request *Request) os.Error { return nil }

func NewView(v Renderable) (view View) {
        if _, ok := v.(Renderable); ok {
                
        }

        if _, ok := v.(Prepareable); ok {
                vv, _ := v.(PrepareableRenderable)
                view = viewStringer{ vv }
        } else {
                vv := &nothingPrepareable{ v }
                view = viewStringer{ vv }
        }

        return
}

// Template implements View interface.
type templateView struct {
        template *template.Template
        feed TemplateFeed
}

func NewTemplate(fn string, feed TemplateFeed) (view View, err os.Error) {
        f, err := os.Open(fn, os.O_RDONLY, 0)
        if err != nil {
                return
        }

        s, err := ioutil.ReadAll(f)
        if err != nil {
                return
        }

        view, err = NewTemplateFromString(string(s), feed)
        return
}

func NewTemplateFromString(s string, feed TemplateFeed) (view View, err os.Error) {
        var temp *template.Template
        temp, err = template.Parse(s, nil/* TODO: make use of it? */)
        if err != nil {
                return
        }

        t := &templateView{ temp, feed }
        view = NewView(t)
        return
}

func (v *templateView) Prepare(request *Request) (err os.Error) {
        err = v.feed.PrepareFeed(request)
        // TODO: prepare feed fields
        return
}

func (v *templateView) Render(w RenderBuffer) (err os.Error) {
        err = v.template.Execute( v.feed, w )
        return
}
