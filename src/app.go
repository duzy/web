package web

import (
        "os"
        "io"
        "fmt"
        "http"
        "bufio"
        "template"
)

// Make response to a request.
type Handler interface {
        WriteResponse(w io.Writer, app *App)
}

// Indicate a model of a app, e.g. CGIModel, FCGIModel, SCGIModel, etc.
type AppModel interface {
        RequestMethod() string
        PathInfo() string
        QueryString() string
        RequestReader() io.Reader // for reading data like POST messages
        ResponseWriter() io.Writer
}

// Produce a new web.App to talk to a session.
func NewApp(title string, m interface {}) (app *App) {
        am, isAppModel := m.(AppModel)
        if isAppModel {
                app = new(App)
                app.model = am
                app.handlers = make(map[string]Handler)
                app.title = title
                app.pathMatcher = DefaultPathMatcher(false)
        }
        return
}

// Indicate a CGI web session, also holds parameters of a request.
type App struct {
        model AppModel
        handlers map[string]Handler
        title string
        method string
        path string
        pathMatcher PathMatcher
        query map[string][]string
        //Header map[string]string
}

func (app *App) Method() string { return app.method }
func (app *App) Path() string { return app.path }
func (app *App) Query(k string) []string { return app.query[k] }

func (app *App) Handle(url string, h Handler) (prev Handler) {
        prev = app.handlers[url]
        app.handlers[url] = h
        return
}

func (app *App) HandleDefault(h Handler) (prev Handler) {
        prev = app.handlers[""]
        app.handlers[""] = h
        if app.handlers["/"] == nil {
                app.handlers["/"] = h
        }
        return
}

func (app *App) ReturnError(w io.Writer, err os.Error) {
        fmt.Fprintf(w, "error: %v", err)
}

func (app *App) Exec() (err os.Error) {
        /*
        app.method = os.Getenv("REQUEST_METHOD")
        app.path = os.Getenv("PATH_INFO")
        app.query, err = http.ParseQuery(os.Getenv("QUERY_STRING"))
         */
        app.method = app.model.RequestMethod()
        app.path = app.model.PathInfo()
        app.query, err = http.ParseQuery(app.model.QueryString())
        if err != nil { /* TODO: log error */ goto finish }

        w, err := bufio.NewWriterSize(app.model.ResponseWriter(), 1024*10)
        defer w.Flush()

        for k, h := range app.handlers {
                if !app.pathMatcher.PathMatched(k, app.path) { continue }
                h.WriteResponse(w, app)
        }

finish:
        return
}

// Responsible to check two paths to see if they are identical. 
type PathMatcher interface {
        PathMatched(p1 string, p2 string) bool
}

// The default method of matching two path: matched if equal
type DefaultPathMatcher bool

func (m DefaultPathMatcher) PathMatched(p1 string, p2 string) bool {
        m = DefaultPathMatcher(p1 == p2)
        return bool(m)
}

// Produce a new web view(a web.Handler).
func NewView(a interface{}) (h Handler) {
        switch t := a.(type) {
        case string:
                m := Model(&DefaultView{ t, nil })
                h = Handler(&View{ m })
        case Model:
                h = Handler(&View{ t })
        }
        return
}

// Real representation of a web view, it's a web.Handler.
type View struct {
        model Model // this is private field
}

func (v *View) WriteResponse(w io.Writer, app *App) {
        fmt.Fprintf(w, "Content-Type: text/html;\n\n")

        if v.model.GetTemplate() == "" { goto finish }

        t, err := template.ParseFile(v.model.GetTemplate(), nil)
        if err != nil { app.ReturnError(w, err); goto finish }

        err = t.Execute( v.model.MakeFields(app), w )
        if err != nil { app.ReturnError(w, err); goto finish }

finish:
        return
}

// A Model is a true implementation of a web view.
type Model interface {
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

// Use FuncHandler to wrap a func as a web.Handler.
type FuncHandler func(w io.Writer, app *App)

func (f FuncHandler) WriteResponse(w io.Writer, app *App) {
        f(w, app)
}
