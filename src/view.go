package web

import (
        "io"
        "io/ioutil"
        "bytes"
        //"os"
        //"fmt"
        "template"
)

// Real representation of a web view, it's a web.Handler.
type View struct {
        model ViewModel // this is private field
}

// Make views all Stringer.
type HandlerStringer struct {
        Handler
        app *App
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
func NewView(a interface{}) (h Handler) {
        switch t := a.(type) {
        case ViewModel:
                h = Handler(&View{ t })
        case TemplateString, TemplateFile, string:
                h = NewStandardView(t)
        }
        return
}

func NewStandardView(a interface{}) (h Handler) {
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
        h = Handler(&View{ m })
        return
}

func NewStringer(a interface{}, app *App) (hs *HandlerStringer) {
        switch t := a.(type) {
        case Handler:
                hs = &HandlerStringer{ t, app, }
        case *View:
                hs = &HandlerStringer{ Handler(t), app, }
        case ViewModel:
                hs = &HandlerStringer{ NewView(a), app, }
        }
        return
}

func TemplateStringGetter(s string) TemplateGetter {
        return TemplateGetter((*TemplateString)(&s))
}

func TemplateFileGetter(s string) TemplateGetter {
        return TemplateGetter((*TemplateFile)(&s))
}


func (v *View) HandleSubpath(subpath string, app *App) (handled bool) {
        if sh, ok := v.model.(SubpathHandler); ok {
                handled = sh.HandleSubpath(subpath, app)
        }
        return
}

func (v *View) WriteContent(w io.Writer, app *App) {
        app.SetHeader("Content-Type", "text/html")

        temp := v.model.GetTemplate()
        if temp == "" { goto finish }

        t, err := template.Parse(temp, nil/* TODO: make use of it? */)
        if err != nil { app.ReturnError(w, err); goto finish }

        var f interface{}
        if fm, ok := v.model.(FieldsMaker); ok {
                f = fm.MakeFields(app)
        } else {
                f = v.model
        }

        err = t.Execute( f, w )
        if err != nil { app.ReturnError(w, err); goto finish }

finish:
        return
}

func (vs *HandlerStringer) String() string {
        // TODO: handle with the case of 'app == nil'
        buf := bytes.NewBuffer(make([]byte, 0))
        vs.Handler.WriteContent(buf, vs.app)
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
