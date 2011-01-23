package web

import (
        "os"
        "io"
        "fmt"
        "http"
        "bytes"
        "strings"
)

const cookieSessionId = "__web_cid" // use a special name for session id

const (
        PathMatchedNothing = 0
        PathMatchedFull = 1 // the paths matched exactly
        PathMatchedParent = 2 // both paths has the same parent: /edit, /edit/1, /edit/2
)

// Responsible to check two paths to see if they are identical.
// TODO: think about paths like '/edit/*'
type PathMatcher interface {
        PathMatch(p1 string, p2 string) (n int, sub string)
}

// The default method of matching two path: matched if equal
type DefaultPathMatcher int

// Make response to a request.
type RequestHandler interface {
        HandleRequest(request *Request, response *Response) (err os.Error)
}

type SubpathHandler interface {
        HandleSubpath(subpath string, request *Request) bool
}

// Use FuncHandler to wrap a func as a web.Handler.
type FuncHandler func(request *Request, response *Response) (err os.Error)

type RequestManager interface {
        GetRequest(id string) (request *Request, err os.Error)
        ProcessRequest(req *Request) (response *Response, err os.Error)
}

// Indicate a model of a app, e.g. CGIModel, FCGIModel, SCGIModel, etc.
type AppModel interface {
        ProcessRequests(rp RequestManager) os.Error
}

type Request struct {
        http.Request

        Path string
        ScriptName string
        QueryString string
        HttpCookie string
        
        app *App
        session *Session
        cookies []*Cookie
        query map[string][]string // parsed query
}

type Response struct {
        http.Response
        BodyWriter io.Writer

        app *App
        cookies []*Cookie
}

type Cookie struct {
        Name string
        Content string // TODO: should be named as Value
        Domain string
        Path string
        Expires string
        Version string
}

// Indicate a CGI web session, also holds parameters of a request.
type App struct {
        model AppModel
        config *AppConfig
        handlers map[string]RequestHandler
        pathMatcher PathMatcher
        requests map[string]*Request
}

func (f FuncHandler) HandleRequest(request *Request, response *Response) (err os.Error) {
        return f(request, response)
}

func (req *Request) Session() *Session { return req.session; }

func (request *Request) initSession() (err os.Error) {
        if request.app == nil {
                err = os.NewError("no associated app")
                return
        }

        appScriptName := request.ScriptName

        if c := request.Cookie(cookieSessionId); c != nil {
                if c.Path == "" {
                        c.Path = appScriptName
                }
        }

        shouldMakeNewSession := true

        if c := request.Cookie(cookieSessionId); c != nil {
                // TODO: check value of c.Name and c.Content
                if request.app.config == nil {
                        err = os.NewError("no app config")
                        return
                }

                sec, e := LoadSession(c.Content, request.app.config.Persister)
                if e == nil {
                        shouldMakeNewSession = false
                        request.session = sec
                } else {
                        // TODO: log errors
                        //panic(err)
                }
        }

        if shouldMakeNewSession {
                request.session = NewSession()
                // FIXME: app.cookies may be not empty since it's
                //        returned by ParseCookies.
                request.cookies = append(request.cookies, &Cookie{
                Name: cookieSessionId,
                Content: request.session.Id(),
                Path: request.ScriptName,
                })
        }

        if request.session == nil { panic("web: no session") }

        // TODO: use hidden form for session state management
        //       for the case that the client did not support
        //       cookies. Also for security reason?
        return
}

func (request *Request) Cookies() []*Cookie { return request.cookies }
func (request *Request) Cookie(k string) (rc *Cookie) {
        for _, c := range request.cookies {
                if c.Name == k {
                        rc = c
                }
        }
        return
}

func (c *Cookie) String() string {
        s := ""
        if c.Name    != "" { s += c.Name + "=" + c.Content }
        if c.Expires != "" { s += "; expires=" + c.Expires }
        if c.Domain  != "" { s += "; domain=" + c.Domain }
        if c.Path    != "" { s += "; path=" + c.Path }
        if c.Version != "" { s += "; version=" + c.Version }
        return s
}

// Parse http cookies header.
// See http://ftp.ics.uci.edu/pub/ietf/http/rfc2109.txt for full spec.
func ParseCookies(s string) (cookies []*Cookie) {
        cookies = make([]*Cookie, 0, 5)
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

// Produce a new web.App to talk to a session.
func NewApp(m interface {}) (app *App, err os.Error) {
        if m == nil {
                err = os.NewError("web.NewApp: Invalid parameter!")
                return
        }

        switch v := m.(type) {
        case AppModel:
                app = new(App)
                err = app.initFromModel(v)
                if err != nil { app = nil }
        case *AppConfig:
                app = new(App)
                err = app.initFromConfig(v)
                if err != nil { app = nil }
        case string: // app config file
                var cfg *AppConfig
                cfg, err = LoadAppConfig(v)
                if err == nil {
                        if cfg == nil {
                                err = os.NewError(fmt.Sprintf("can't load app config '%v'", v))
                                return
                        }
                        // assign default FS persister if nil
                        if cfg.Persister == nil {
                                cfg.Persister = defaultPersisterConfigFS
                        }
                        app = new(App)
                        err = app.initFromConfig(cfg)
                        if err != nil { app = nil }
                }
        }
        return
}

func newAppModelFromAppConfig(cfg *AppConfig) (am AppModel, err os.Error) {
        switch cfg.Model {
        case "CGI": am, err = NewCGIModel()
        }
        return
}

func newAppConfigForAppModel(am AppModel) (config *AppConfig) {
        config = new(AppConfig)
        switch am.(type) {
        case *CGIModel: config.Model = "CGI"
        }
        config.Persister = PersisterConfig(&PersisterConfigFS{ Location: "/tmp/web/sessions", })
        return
}

func (app *App) initFromConfig(cfg *AppConfig) (err os.Error) {
        app.model, err = newAppModelFromAppConfig(cfg)
        if app.model == nil {
                error("initFromConfig: invalid app model '%s'", cfg.Model)
                goto finish
        }

        app.config = cfg

        app.handlers = make(map[string]RequestHandler)
        app.pathMatcher = DefaultPathMatcher(PathMatchedNothing)
        app.requests = make(map[string]*Request)

        // TODO: init database
finish:
        return
}

func (app *App) initFromModel(am AppModel) (err os.Error) {
        app.model = am
        app.handlers = make(map[string]RequestHandler)
        app.pathMatcher = DefaultPathMatcher(PathMatchedNothing)
        app.config = newAppConfigForAppModel(app.model)
        app.requests = make(map[string]*Request)
        return
}

func (app *App) GetModel() AppModel { return app.model }

func (app *App) Handle(url string, h RequestHandler) (prev RequestHandler) {
        prev = app.handlers[url]
        app.handlers[url] = h
        return
}

func (app *App) HandleDefault(h RequestHandler) (prev RequestHandler) {
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

func (res *Response) writeHeader(w io.Writer) (err os.Error) {
        for _, v := range res.cookies {
                fmt.Fprintf(w, "Set-Cookie: %s\n", v.String())
        }

        for k, v := range res.Header {
                fmt.Fprintf(w, "%s: %s\n", k, v)
        }

        fmt.Fprintf(w, "\n")
        return
}

/**
 *  Get database via the name -- specified in the config file like this:
 *    "databases": {
 *        "db-name": {
 *            ...
 *        }
 *    }
 */
func (app *App) GetDatabase(name string) (db Database, err os.Error) {
        dbm := GetDBManager()
        for n, cfg := range app.config.Databases {
                if n == name {
                        db, err = dbm.GetDatabase(cfg)
                        break
                }
        }
        return
}

func (app *App) Exec() (err os.Error) {
        defer func() {
                dbm := GetDBManager()
                dbm.CloseAll() // close all databases
        }()

        err = app.model.ProcessRequests(app)
        if err != nil {
                // ....
        }
        return
}

func (app *App) GetRequest(id string) (request *Request, err os.Error) {
        // TODO: multiple requests management
        request = app.requests[id]
        if request == nil {
                request = &Request{}
                request.initSession()
                app.requests[id] = request
        }
        request.app = app
        return
}

type noCloseReader struct { io.Reader }
func (c noCloseReader) Close() os.Error { return nil }

func (app *App) ProcessRequest(req *Request) (response *Response, err os.Error) {
        defer func() {
                if app.config == nil { error("config: <nil>") }
                if req.session != nil {
                        err = req.session.save(app.config.Persister)
                }
        }()

        req.query, err = http.ParseQuery(req.QueryString)
        if err != nil {
                /* TODO: log error */
                return
        }

        contentBuffer := bytes.NewBuffer(make([]byte, 1024*10))
        contentBuffer.Reset()

        response = &Response{}
        response.app = app
        response.cookies = make([]*Cookie, 0, 8)
        response.Header = make(map[string]string)
        response.Body = noCloseReader{ contentBuffer }
        response.BodyWriter = contentBuffer

        if c := req.Cookie(cookieSessionId); c != nil {
                response.cookies = append(response.cookies, c)
        }

        handled := false
        for k, h := range app.handlers {
                if m, s := app.pathMatcher.PathMatch(k, req.Path); m != PathMatchedNothing {
                        if m == PathMatchedParent {
                                sh, ok := h.(SubpathHandler)
                                if ok && !sh.HandleSubpath(s, req) {
                                        continue
                                }
                        }
                        err = h.HandleRequest(req, response)
                        if err != nil {
                                // TODO: error handling...
                                return
                        }

                        response.ContentLength = int64(contentBuffer.Len())

                        handled = true
                        break // just ignore any other handlers
                }
        }//for (app.handlers)

        if !handled {
                contentBuffer.Reset()
                w := contentBuffer
                fmt.Fprintf(w, "Content-Type: text/html; charset=utf-8\n\n")
                fmt.Fprintf(w, `<html>`)
                fmt.Fprintf(w, `<head>`)
                fmt.Fprintf(w, `<meta http-equiv="content-type" content="text/html;charset=utf-8">`)
                fmt.Fprintf(w, `</head>`)
                fmt.Fprintf(w, `<body>`)
                fmt.Fprintf(w, `<font color="red"><h1>Error: 404</h1></font>`)
                fmt.Fprintf(w, `The requested URL <code>%s%s</code> was not found on this server`, req.ScriptName, req.Path)
                fmt.Fprintf(w, `<body>`)
                fmt.Fprintf(w, `</html>`)
        }

        return
}

func (m DefaultPathMatcher) PathMatch(p1 string, p2 string) (n int, sub string) {
        /* // debug helper code
        defer func() {
                if e := recover(); e != nil {
                        panic(fmt.Sprintf("PathMatch: %s <=> %s, %v", p1, p2, e))
                }
        }()
         */

        n = PathMatchedNothing

        /*
        i := strings.Index(p2, p1) // find p1 in p2
        if i == 0 { // p1 is the prefix of p2
                l := len(p2) - len(p1)
                if l == 0 { // full matched
                        n = PathMatchedFull
                } else if 0 < l { // p1(/edit) <=> p2(/edit/1)
                        n = PathMatchedParent

                        // For cases of: p1(/) <=> p2(/sync),
                        //               p1(/buy/a/) <=> p2(/buy/a/1)
                        if len(p1) == 0 || p1[len(p1)-1:len(p1)] == "/" { sub = "/" }

                        sub += p2[len(p1):len(p1)+l]
                } else if l < 0 {
                        // unmatched
                }
        }
         */
        if len(p1) <= len(p2) {
                pre := p2[0:len(p1)]
                if pre != p1 { /* matched nothing */ return }

                l := len(p2)-len(p1)
                if l == 0 {
                        n = PathMatchedFull
                } else {
                        n = PathMatchedParent

                        // For cases of: p1(/) <=> p2(/sync),
                        //               p1(/buy/a/) <=> p2(/buy/a/1)
                        if len(p1) == 0 || p1[len(p1)-1:len(p1)] == "/" { sub = "/" }

                        sub += p2[len(p1):len(p1)+l]
                }
        }
        return
}
