package web

import (
        "os"
        "io"
        "fmt"
        "http"
        "bytes"
        "strings"
)

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
        handlers map[string]RequestHandler
        pathMatcher PathMatcher
        requests map[string]*Request
}

func (f FuncHandler) HandleRequest(request *Request, response *Response) (err os.Error) {
        return f(request, response)
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
                err = newError("web.NewApp: Invalid parameter!")
                return
        }

        switch v := m.(type) {
        case AppModel:
                app = new(App)
                err = app.initFromModel(v)
                if err != nil { app = nil }
        }
        return
}

func (app *App) initFromModel(am AppModel) (err os.Error) {
        app.model = am
        app.handlers = make(map[string]RequestHandler)
        app.pathMatcher = DefaultPathMatcher(PathMatchedNothing)
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

func (app *App) Exec() (err os.Error) {
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
                app.requests[id] = request
        }
        request.app = app
        return
}

type noCloseReader struct { io.Reader }
func (c noCloseReader) Close() os.Error { return nil }

func (app *App) ProcessRequest(req *Request) (response *Response, err os.Error) {
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
        n = PathMatchedNothing

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
