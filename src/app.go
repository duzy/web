package web

import (
        "os"
        "io"
        "fmt"
        "http"
        "bytes"
        "bufio"
)

// Make response to a request.
type Handler interface {
        WriteContent(w io.Writer, app *App) // TODO: returns error
}

// Use FuncHandler to wrap a func as a web.Handler.
type FuncHandler func(w io.Writer, app *App)

func (f FuncHandler) WriteContent(w io.Writer, app *App) {
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
                app.title = title
                app.handlers = make(map[string]Handler)
                app.header = make(map[string]string)
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
        header map[string]string
}

func (app *App) Method() string { return app.method }
func (app *App) Path() string { return app.path }
func (app *App) Query(k string) []string { return app.query[k] }

func (app *App) Header(k string) string { return app.header[k] }
func (app *App) SetHeader(k, v string) (prev string) {
        prev = app.header[k]
        app.header[k] = v
        return
}


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

func (app *App) writeHeader(w io.Writer) (err os.Error) {
        for k, v := range app.header {
                //fmt.Fprintf(os.Stdout, "%s: %s\n", k, v)
                fmt.Fprintf(w, "%s: %s\n", k, v)
        }
        fmt.Fprintf(w, "\n")
        return
}

func (app *App) Exec() (err os.Error) {
        app.method = app.model.RequestMethod()
        app.path = app.model.PathInfo()
        app.query, err = http.ParseQuery(app.model.QueryString())
        if err != nil { /* TODO: log error */ goto finish }

        contentBuffer := bytes.NewBuffer(make([]byte, 1024*10))

        for k, h := range app.handlers {
                if app.pathMatcher.PathMatched(k, app.path) {
                        contentBuffer.Reset()
                        h.WriteContent(contentBuffer, app)

                        w, err := bufio.NewWriterSize(app.model.ResponseWriter(), 1024*10)
                        if err != nil { /*TODO: IO error */ goto finish }

                        defer w.Flush()

                        err = app.writeHeader(w)
                        if err != nil { /*TODO: IO error */ goto finish }

                        w.Write(contentBuffer.Bytes())
                        break
                }
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
