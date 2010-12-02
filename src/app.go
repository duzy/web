package web

import (
        "os"
        "io"
        "fmt"
        "http"
        "bufio"
)

// Make response to a request.
type Handler interface {
        WriteResponse(w io.Writer, app *App)
}

// Use FuncHandler to wrap a func as a web.Handler.
type FuncHandler func(w io.Writer, app *App)

func (f FuncHandler) WriteResponse(w io.Writer, app *App) {
        f(w, app)
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
