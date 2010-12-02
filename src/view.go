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
                m := ViewModel(&DefaultView{ t, nil })
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

        //os.Stdout.WriteString(app.Header("Content-Type"))

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
        MakeFields(app *App) (map[string]string)
}

// The default Model of a view.
type DefaultView struct {
        Template string
        Fields map[string]string
}

func (v DefaultView) GetTemplate() (s string) {
        s = v.Template
        return
}

func (v DefaultView) MakeFields(app *App) (m map[string]string) {
        if v.Fields == nil {
                v.Fields = make(map[string]string)
                v.Fields["title"] = app.title
        }
        m = v.Fields
        return 
}
