package web

import (
        "os"
        "io"
        "fmt"
        "http"
        "bytes"
        //"bufio"
        "strings"
)

// Make response to a request.
type Handler interface {
        WriteContent(w io.Writer, app *App) // TODO: returns error
}

type SubpathHandler interface {
        HandleSubpath(subpath string, app *App) bool
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
        ScriptName() string
        HttpCookie() string
        RequestReader() io.Reader // for reading data like POST messages
        ResponseWriter() io.Writer

        HandleErrors() // TODO: think about the error handling
}

type Cookie struct {
        //nv map[string]string
        Name string
        Content string // TODO: should be named as Value
        Domain string
        Path string
        Expires string
        Version string
}

func (c *Cookie) String() string {
        s := ""
        if c.Name    != "" { s += c.Name + "=" + c.Content }
        if c.Expires != "" { s += "; expires=" + c.Expires }
        if c.Domain  != "" { s += "; domain=" + c.Domain }
        if c.Path    != "" { s += "; path=" + c.Path }
        if c.Version != "" { s += "; version=" + c.Version }
        // for n, v := range c.nv {
        //         if s != "" { s += "; " }
        //         s += n + "=" + v;
        // }
        return s
}

// See http://ftp.ics.uci.edu/pub/ietf/http/rfc2109.txt.
func ParseCookies(s string) (cookies []*Cookie) {
        cookies = make([]*Cookie, 0)
        s = strings.TrimSpace(s) // FIXME: needed?
        if ss := strings.Split(s, ";", -1); 0 < len(ss) {
                var c *Cookie
                for _, a := range ss {
                        kv := strings.Split(strings.TrimSpace(a), "=", 2)
                        if len(kv) != 2 { continue }
                        switch kv[0] {
                        case "$Version":
                                // TODO: version checking
                                continue
                        case "$Domain":
                                if c != nil { c.Domain = kv[1] }
                        case "$Path":
                                if c != nil { c.Path = kv[1] }
                        case "$Max-Age":
                                if c != nil { /* TODO: Max-Age selection... */ }
                        default:
                                c = new(Cookie)
                                c.Name = kv[0]
                                c.Content = kv[1]
                                cookies = append(cookies, c)
                                c = nil // reset
                        }
                }
        }
        return
}

// Indicate a CGI web session, also holds parameters of a request.
type App struct {
        model AppModel
        session *Session
        handlers map[string]Handler
        title string
        pathMatcher PathMatcher
        query map[string][]string
        header map[string]string
        cookies []*Cookie
}

// Produce a new web.App to talk to a session.
func NewApp(title string, m interface {}) (app *App) {
        if am, ok := m.(AppModel); ok {
                app = new(App)
                app.model = am

                shouldMakeNewSession := true

                app.cookies = ParseCookies(am.HttpCookie())
                if c := app.Cookie("session"); c != nil {
                        // TODO: check value of c.Name and c.Content
                        sec, err := LoadSession(c.Content)
                        if err == nil {
                                shouldMakeNewSession = false
                                app.session = sec
                        } else {
                                // TODO: log errors
                                //panic(err)
                        }
                }

                if shouldMakeNewSession {
                        app.session = NewSession()
                        app.cookies = append(app.cookies, &Cookie{
                        Name: "session",
                        Content: app.session.Id(),
                        })
                }

                app.title = title
                app.handlers = make(map[string]Handler)
                app.header = make(map[string]string)
                app.pathMatcher = DefaultPathMatcher(PathMatchedNothing)
        }
        return
}

func (app *App) Method() string { return app.model.RequestMethod() }
func (app *App) Path() string { return app.model.PathInfo() }
func (app *App) ScriptName() string { return app.model.ScriptName() }
func (app *App) Query(k string) []string { return app.query[k] }
func (app *App) RequestReader() io.Reader { return app.model.RequestReader() }

// Returns unparsed cookies.
func (app *App) RawCookie() string { return app.model.HttpCookie() }
func (app *App) Cookie(k string) (rc *Cookie) {
        for _, c := range app.cookies {
                if c.Name == k/* TODO: ignore cases? */ {
                        rc = c
                }
        }
        return
}

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

func (app *App) ReturnError(w io.Writer, err interface{}) {
        fmt.Fprintf(w, "Content-Type: text/plain\n\n")
        fmt.Fprintf(w, "error: %v", err)
}

func (app *App) writeHeader(w io.Writer) (err os.Error) {
        for _, v := range app.cookies {
                fmt.Fprintf(w, "Set-Cookie: %s\n", v.String())
        }
        for k, v := range app.header {
                fmt.Fprintf(w, "%s: %s\n", k, v)
        }
        fmt.Fprintf(w, "\n")
        return
}

func (app *App) Exec() (err os.Error) {
        app.query, err = http.ParseQuery(app.model.QueryString())
        if err != nil { /* TODO: log error */ goto finish }

        for k, h := range app.handlers {
                if m, s := app.pathMatcher.PathMatched(k, app.Path()); m != PathMatchedNothing {
                        // TODO: rethink about the buffering performance
                        contentBuffer := bytes.NewBuffer(make([]byte, 1024*10))
                        contentBuffer.Reset()

                        if m == PathMatchedParent {
                                sh, ok := h.(SubpathHandler)
                                if !(ok && sh.HandleSubpath(s, app)) {
                                        continue
                                }
                        }
                        h.WriteContent(contentBuffer, app)

                        headerBuffer := bytes.NewBuffer(make([]byte, 1024))
                        headerBuffer.Reset()
                        err := app.writeHeader(headerBuffer)
                        if err != nil { /*TODO: error */ goto finish }

                        w := app.model.ResponseWriter()

                        /*
                        defer func() {
                                if e := recover(); e != nil {
                                        app.ReturnError(w, e)
                                }
                        }()
                         */
                        
                        w.Write(headerBuffer.Bytes())
                        w.Write(contentBuffer.Bytes())
                        break // just ignore any other handlers
                }
        }//for (app.handlers)

finish:
        return
}

const (
        PathMatchedNothing = 0
        PathMatchedFull = 1 // the paths matched exactly
        PathMatchedParent = 2 // both paths has the same parent: /edit, /edit/1, /edit/2
)

// Responsible to check two paths to see if they are identical.
// TODO: think about paths like '/edit/*'
type PathMatcher interface {
        PathMatched(p1 string, p2 string) (n int, sub string)
}

// The default method of matching two path: matched if equal
type DefaultPathMatcher int

func (m DefaultPathMatcher) PathMatched(p1 string, p2 string) (n int, sub string) {
        n = PathMatchedNothing
        i := strings.Index(p2, p1) // find p1 in p2
        if i == 0 { // p1 is the prefix of p2
                l := len(p2) - len(p1)
                if l == 0 { // full matched
                        n = PathMatchedFull
                } else if 0 < l { // p1(/edit) <=> p2(/edit/1)
                        n = PathMatchedParent
                        sub = p2[len(p1):len(p1)+l]
                } else if l < 0 {
                        // unmatched
                }
        }
        return
}
