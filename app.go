package web

import (
        "os"
        "io"
        "fmt"
        "http"
        "bufio"
        "template"
)

type Handler interface {
        WriteResponse(w io.Writer, app *App)
}

type App struct {
        handlers map[string]Handler
        title string
        method string
        path string
        query map[string][]string
}

func NewApp(title string) (app *App) {
        app = new(App)
        app.handlers = make(map[string]Handler)
        app.title = title
        return
}

func (app *App) Handle(url string, h Handler) {
        app.handlers[url] = h
}

func (app *App) ReturnError(w io.Writer, err os.Error) {
        fmt.Fprintf(w, "error: %v", err)
}

func (app *App) Exec() (err os.Error) {
        app.method = os.Getenv("REQUEST_METHOD")
        app.path = os.Getenv("PATH_INFO")
        app.query, err = http.ParseQuery(os.Getenv("QUERY_STRING"))
        if err != nil { /* TODO: log error */ goto finish }

        w, err := bufio.NewWriterSize(os.Stdout, 1024*10)
        defer w.Flush()

        for k, h := range app.handlers {
                if !pathMatched(k, app.path) { continue }
                h.WriteResponse(w, app)
        }

finish:
        return
}

func pathMatched(lhs string, rhs string) (matched bool) {
        matched = false
        if lhs == rhs { matched = true }
        return
}

type view struct {
        model Model
}

type Model interface {
        GetTemplate() string
        MakeFields(app *App) (map[string]string)
}

func NewView(a interface{}) (h Handler) {
        v := new(view)
        h = Handler(v)
        switch t := a.(type) {
        case string:
                dv := new(DefaultView)
                dv.Template = t
                v.model = Model(dv)
        case Model:
                v.model = t;
        }
        return
}

func (v *view) WriteResponse(w io.Writer, app *App) {
        fmt.Fprintf(w, "Content-Type: text/html;\n\n")

        if v.model.GetTemplate() == "" { goto finish }

        t, err := template.ParseFile(v.model.GetTemplate(), nil)
        if err != nil { app.ReturnError(w, err); goto finish }

        err = t.Execute( v.model.MakeFields(app), w )
        if err != nil { app.ReturnError(w, err); goto finish }

finish:
        return
}

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

type FuncHandler func(w io.Writer, app *App)

func (f FuncHandler) WriteResponse(w io.Writer, app *App) {
        f(w, app)
}
