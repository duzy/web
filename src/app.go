package web

import (
        "os"
        "io"
        "fmt"
        "http"
        //"bytes"
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
type Handler interface {
        WriteContent(w io.Writer, app *App) // TODO: returns error
}

type SubpathHandler interface {
        HandleSubpath(subpath string, app *App) bool
}

// Use FuncHandler to wrap a func as a web.Handler.
type FuncHandler func(w io.Writer, app *App)

type RequestProcessor interface {
        ProcessRequest(req *Request) (err os.Error)
}

// Indicate a model of a app, e.g. CGIModel, FCGIModel, SCGIModel, etc.
type AppModel interface {
        ProcessRequests(rp RequestProcessor) os.Error
        GetRequest() *Request
}

type Request struct {
        http.Request
        session *Session
}

type Response struct {
        http.Response
        query map[string][]string
        cookies []*Cookie
}

type Cookie struct {
        changed bool // 'true' means the cookie must be sent via 'Set-Cookie'
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
        handlers map[string]Handler
        pathMatcher PathMatcher
        requests []*Request
}

func (f FuncHandler) WriteContent(w io.Writer, app *App) {
        f(w, app)
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

// Produce a new web.App to talk to a session.
func NewApp(m interface {}) (app *App, err os.Error) {
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
                } else {
                        //error("LoadAppConfig(%s): %v", v, err)
                }
        }
        return
}

func newAppModelFromAppConfig(cfg *AppConfig) (am AppModel) {
        switch cfg.Model {
        case "CGI": am = NewCGIModel()
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
        app.model = newAppModelFromAppConfig(cfg)
        if app.model == nil {
                error("initFromConfig: invalid app model '%s'", cfg.Model)
                goto finish
        }

        app.config = cfg

        //app.cookies = make([]*Cookie, 0)
        app.handlers = make(map[string]Handler)
        //app.header = make(map[string]string)
        app.pathMatcher = DefaultPathMatcher(PathMatchedNothing)

        //err = app.initSession()

        // TODO: init database
finish:
        return
}

func (app *App) initFromModel(am AppModel) (err os.Error) {
        app.model = am
        //app.cookies = make([]*Cookie, 0)
        app.handlers = make(map[string]Handler)
        //app.header = make(map[string]string)
        app.pathMatcher = DefaultPathMatcher(PathMatchedNothing)
        app.config = newAppConfigForAppModel(app.model)
        //err = app.initSession()
        return
}

/*
func (app *App) initSession() (err os.Error) {
        appScriptName := app.ScriptName()

        cookies := ParseCookies(app.model.HttpCookie())
        for _, c := range cookies {
                switch c.Name {
                case cookieSessionId:
                        if c.Path == "" {
                                c.Path = appScriptName
                        }

                        // FIXME: UA may send more than one session
                        //        cookie, should avoid duplication
                        app.cookies = append(app.cookies, c)
                }
        }

        shouldMakeNewSession := true

        if c := app.Cookie(cookieSessionId); c != nil {
                // TODO: check value of c.Name and c.Content
                sec, e := LoadSession(c.Content, app.config.Persister)
                if e == nil {
                        shouldMakeNewSession = false
                        app.session = sec
                } else {
                        // TODO: log errors
                        //panic(err)
                }
        }

        if shouldMakeNewSession {
                app.session = NewSession()
                // FIXME: app.cookies may be not empty since it's
                //        returned by ParseCookies.
                app.cookies = append(app.cookies, &Cookie{
                changed: true,
                Name: cookieSessionId,
                Content: app.session.Id(),
                Path: app.ScriptName(),
                })
        }

        if app.session == nil { panic("web: no session") }

        // TODO: use hidden form for session state management
        //       for the case that the client did not support
        //       cookies. Also for security reason?
        return
}
 */

func (app *App) GetModel() AppModel { return app.model }
//func (app *App) Session() *Session { return app.session }
/*
func (app *App) Method() string { return app.model.RequestMethod() }
func (app *App) Path() string { return app.model.PathInfo() }
func (app *App) ScriptName() string { return app.model.ScriptName() }
func (app *App) Query(k string) []string { return app.query[k] }
func (app *App) RequestReader() io.Reader { return app.model.RequestReader() }

// Returns unparsed cookies.
func (app *App) RawCookie() string { return app.model.HttpCookie() }
func (app *App) Cookies() []*Cookie { return app.cookies }
func (app *App) Cookie(k string) (rc *Cookie) {
        for _, c := range app.cookies {
                if c.Name == k {
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
 */

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

/*
func (app *App) writeHeader(w io.Writer) (err os.Error) {
        for _, v := range app.cookies {
                // TODO: only Set-Cookie for 'changed' cookies, avoid for
                //       browser returned cookies eg the 'session' cookie.
                if v.changed {
                        fmt.Fprintf(w, "Set-Cookie: %s\n", v.String())
                }
        }

        for k, v := range app.header {
                fmt.Fprintf(w, "%s: %s\n", k, v)
        }

        fmt.Fprintf(w, "\n")
        return
}
 */

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
        err = app.model.ProcessRequests(app)
        return
}

func (app *App) ProcessRequest(req *Request) (err os.Error) {
//         app.query, err = http.ParseQuery(app.model.QueryString())
//         if err != nil {
//                 /* TODO: log error */
//                 goto finish
//         }

//         defer func() {
//                 if app.config == nil { error("config: <nil>") }
//                 err = app.session.save(app.config.Persister)

//                 dbm := GetDBManager()
//                 dbm.CloseAll() // close all databases
//         }()

//         handled := false
//         for k, h := range app.handlers {
//                 if m, s := app.pathMatcher.PathMatch(k, app.Path()); m != PathMatchedNothing {
//                         // TODO: rethink about the buffering performance
//                         contentBuffer := bytes.NewBuffer(make([]byte, 1024*10))
//                         contentBuffer.Reset()

//                         if m == PathMatchedParent {
//                                 sh, ok := h.(SubpathHandler)
//                                 if ok && !sh.HandleSubpath(s, app) {
//                                         continue
//                                 }
//                         }
//                         h.WriteContent(contentBuffer, app)

//                         headerBuffer := bytes.NewBuffer(make([]byte, 1024))
//                         headerBuffer.Reset()
//                         err := app.writeHeader(headerBuffer)
//                         if err != nil { /*TODO: error */ goto finish }

//                         w := app.model.ResponseWriter()

//                         /*
//                         defer func() {
//                                 if e := recover(); e != nil {
//                                         app.ReturnError(w, e)
//                                 }
//                         }()
//                          */
                        
//                         w.Write(headerBuffer.Bytes())
//                         w.Write(contentBuffer.Bytes())
//                         handled = true
//                         break // just ignore any other handlers
//                 }
//         }//for (app.handlers)

//         if !handled {
//                 w := app.model.ResponseWriter()
//                 fmt.Fprintf(w, "Content-Type: text/html; charset=utf-8\n\n")
//                 fmt.Fprintf(w, `<html>`)
//                 fmt.Fprintf(w, `<head>`)
//                 fmt.Fprintf(w, `<meta http-equiv="content-type" content="text/html;charset=utf-8">`)
//                 fmt.Fprintf(w, `</head>`)
//                 fmt.Fprintf(w, `<body>`)
//                 fmt.Fprintf(w, `<font color="red"><h1>Error: 404</h1></font>`)
//                 fmt.Fprintf(w, `The requested URL <code>%s%s</code> was not found on this server`, app.ScriptName(), app.Path())
//                 fmt.Fprintf(w, `<body>`)
//                 fmt.Fprintf(w, `</html>`)
//         }

// finish:
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
