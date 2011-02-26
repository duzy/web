package web

import (
        "bytes"
        "os"
        "io"
        "io/ioutil"
        "reflect"
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

/**
 Page is a RequestHandler and is designed for handling http requests.
 
 A page is a view and structured by sub-views.

 Usage:

 feed := new(MyFeedOrViewModel)
 page, err := web.NewPage("home.html", feed)
 app.Handle("/home", page)

 */
type Page struct {
        templateView
}

type HtmlPage Page

func NewHtmlPage(fn string, feed interface{}) (page *HtmlPage, err os.Error) {
        p, err := NewPage(fn, feed)
        page = (*HtmlPage)(p)
        return
}

func NewPage(fn string, feed interface{}) (page *Page, err os.Error) {
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

func NewPageFromString(s string, feed interface{}) (page *Page, err os.Error) {
        var temp *template.Template
        temp, err = template.Parse(s, nil/* TODO: make use of it? */)
        if err != nil {
                return
        }

        page = &Page{ templateView{ temp, feed } }
        return
}

func NewHtmlPageFromString(s string, feed interface{}) (page *HtmlPage, err os.Error) {
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
        response.Header.Set("Content-Type", "text/html")
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
func (v nothingPrepareable) Prepare(request *Request) os.Error { return nil }

func NewView(v Renderable) (view View) {
        if _, ok := v.(Renderable); ok {
                
        }

        var vv PrepareableRenderable
        if _, ok := v.(Prepareable); ok {
                vv, _ = v.(PrepareableRenderable)
        } else {
                vv = nothingPrepareable{ v }
        }

        view = viewStringer{ vv }
        return
}

// Template implements View interface.
type templateView struct {
        template *template.Template
        feed interface{}
}

func NewTemplate(fn string, feed interface{}) (view View, err os.Error) {
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

func NewTemplateFromString(s string, feed interface{}) (view View, err os.Error) {
        var temp *template.Template
        temp, err = template.Parse(s, nil/* TODO: make use of it? */)
        if err != nil {
                return
        }

        var tv *templateView
        tv = &templateView{ temp, feed }

        view = NewView(tv)
        return
}

func (v *templateView) Prepare(request *Request) (err os.Error) {
        if p, ok := v.feed.(Prepareable); ok {
                err = p.Prepare(request)
                // TODO: prepare feed fields
        }

        vv := reflect.NewValue(v.feed)
        if p, ok := vv.(*reflect.PtrValue); ok {
                vv = p.Elem()
        }

        if sv, ok := vv.(*reflect.StructValue); ok {
                for i := 0; i < sv.NumField(); i += 1 {
                        f := sv.Field(i)
                        if p, ok := f.Interface().(Prepareable); ok {
                                err = p.Prepare(request)
                        }
                }
        }
        return
}

func (v *templateView) Render(w RenderBuffer) (err os.Error) {
        err = v.template.Execute( w, v.feed )
        return
}
