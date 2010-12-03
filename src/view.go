package web

import (
        "io"
        //"os"
        //"fmt"
        "template"
)

// Produce a new web view(a web.Handler).
func NewView(a interface{}) (h Handler) {
        switch t := a.(type) {
        case string:
                m := ViewModel(&StandardView{
                        ViewTemplateName{ t },
                        StandardFields{ nil },
                })
                h = Handler(&View{ m })
        case ViewModel:
                h = Handler(&View{ t })
        }
        return
}

// Real representation of a web view, it's a web.Handler.
type View struct {
        model ViewModel // this is private field
}

func (v *View) WriteContent(w io.Writer, app *App) {
        app.SetHeader("Content-Type", "text/html")

        if v.model.GetTemplate() == "" { goto finish }

        t, err := template.ParseFile(v.model.GetTemplate(), nil)
        if err != nil { app.ReturnError(w, err); goto finish }

        err = t.Execute( v.model.MakeFields(app), w )
        if err != nil { app.ReturnError(w, err); goto finish }

finish:
        return
}

// A ViewModel is a true implementation of a web view.
type ViewModel interface {
        GetTemplate() string
        MakeFields(app *App) interface{}
}

type ViewTemplateName struct { Template string }

func (v *ViewTemplateName) GetTemplate() string { return v.Template }

type StandardFields struct { Fields map[string]interface{} }

func (sf *StandardFields) MakeFields(app *App) (m interface{}) {
        if sf.Fields == nil {
                sf.Fields = make(map[string]interface{})
        }
        sf.Fields["title"] = app.title
        m = &sf.Fields
        return 
}

func (sf *StandardFields) SetField(k string, f interface{}) (prev interface{}) {
        prev = sf.Fields[k]
        sf.Fields[k] = f
        return
}

func (sf *StandardFields) GetField(k string) (f interface{}) {
        f = sf.Fields[k]
        return
}

// The standard ViewModel of a view.
type StandardView struct {
        ViewTemplateName
        StandardFields
}
