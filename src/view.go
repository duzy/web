package web

import (
        //"io"
        "io/ioutil"
        "bytes"
        "os"
        //"fmt"
        "template"
)

// Real representation of a web view, it's a web.Handler.
type View struct {
        model ViewModel // this is private field
}

// Make views all Stringer.
type HandlerStringer struct {
        RequestHandler
        response *Response
}

// A ViewModel is a true implementation of a web view.
type ViewModel interface {
        TemplateGetter
}

type TemplateGetter interface { GetTemplate() string }
type FieldsMaker interface { MakeFields(app *App) interface{} }

// Part of web.ViewModel implements GetTemplate(), load template from file.
type TemplateFile string

// Part of web.ViewModel, implements GetTemplate().
type TemplateString string

// Part of web.ViewModel, implements MakeFields().
type StandardFields map[string]interface{}

// The standard ViewModel of a view.
type StandardView struct {
        TemplateGetter
        StandardFields
}

// Produce a new web view(a web.Handler).
func NewView(a interface{}) (h RequestHandler) {
        switch t := a.(type) {
        case ViewModel:
                h = RequestHandler(&View{ t }) // TODO: using ViewRoot to act as a RequestHandler
        case TemplateString, TemplateFile, string:
                h = NewStandardView(t)
        }
        return
}

func NewStandardView(a interface{}) (h RequestHandler) {
        var g TemplateGetter
        switch t := a.(type) {
        case TemplateString:
                g = TemplateGetter(&t)
        case TemplateFile:
                g = TemplateGetter(&t)
        case string:
                g = TemplateFileGetter(t)
        }

        m := ViewModel(&StandardView{
                g,
                make(StandardFields),
        })
        h = RequestHandler(&View{ m })
        return
}

func NewStringer(a interface{}, response *Response) (hs *HandlerStringer) {
        switch t := a.(type) {
        case RequestHandler:
                hs = &HandlerStringer{ t, response, }
        case *View:
                hs = &HandlerStringer{ RequestHandler(t), response, }
        case ViewModel:
                hs = &HandlerStringer{ NewView(a), response, }
        }
        return
}

func TemplateStringGetter(s string) TemplateGetter {
        return TemplateGetter((*TemplateString)(&s))
}

func TemplateFileGetter(s string) TemplateGetter {
        return TemplateGetter((*TemplateFile)(&s))
}


func (v *View) HandleSubpath(subpath string, request *Request) (handled bool) {
        if sh, ok := v.model.(SubpathHandler); ok {
                handled = sh.HandleSubpath(subpath, request)
        }
        return
}

func (v *View) Handle(request *Request, response *Response) (err os.Error) {
        response.Header["Content-Type"] = "text/html"

        temp := v.model.GetTemplate()
        if temp == "" { goto finish }

        t, err := template.Parse(temp, nil/* TODO: make use of it? */)
        if err != nil { response.app.ReturnError(response.content, err); goto finish }

        var f interface{}
        if fm, ok := v.model.(FieldsMaker); ok {
                f = fm.MakeFields(response.app)
        } else {
                f = v.model
        }

        err = t.Execute( f, response.content )
        if err != nil { response.app.ReturnError(response.content, err); goto finish }

finish:
        return
}

func (vs *HandlerStringer) String() string {
        // TODO: handle with the case of 'app == nil'
        buf := bytes.NewBuffer(make([]byte, 0))
        vs.RequestHandler.Handle(buf, vs.response)
        return string(buf.Bytes())
}

func (t *TemplateFile) GetTemplate() (s string) {
        b, err := ioutil.ReadFile(string(*t))
        if err == nil { s = string(b) }
        /* TODO: handle with IO errors */
        return
}

func (t *TemplateString) GetTemplate() string { return string(*t) }

func (sf *StandardFields) MakeFields(app *App) interface{} {
        (*sf)["SCRIPT"] = app.ScriptName()
        (*sf)["title"] = app.config.Title
        return sf
}
