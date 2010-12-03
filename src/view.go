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
                        TemplateName(t),
                        make(StandardFields),
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

type TemplateName string

func (v *TemplateName) GetTemplate() string { return string(*v) }

type StandardFields map[string]interface{}

func (sf *StandardFields) MakeFields(app *App) interface{} {
        (*sf)["title"] = app.title
        return sf
}

// The standard ViewModel of a view.
type StandardView struct {
        TemplateName
        StandardFields
}
