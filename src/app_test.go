package web

import (
        "testing"
        "strings"
        "bytes"
        "fmt"
        "io"
)

type testAppModel struct {
        *CGIModel
        buffer *bytes.Buffer
        reader *bytes.Buffer
}

type customHandler struct {
        content string
}

func newTestAppModel() (m *testAppModel) {
        cgi := &CGIModel{ make(map[string]string) } //NewCGIModel()
        buffer := bytes.NewBufferString("")
        reader := bytes.NewBufferString("")
        m = &testAppModel{ cgi, buffer, reader }
        return
}

func (am *testAppModel) ResponseWriter() (w io.Writer) {
        w = am.buffer
        return
}

func (am *testAppModel) RequestReader() (r io.Reader) {
        r = am.reader
        return
}

func (h *customHandler) WriteContent(w io.Writer, app *App) {
        fmt.Fprint(w, h.content)
}

func TestFuncHandler(t *testing.T) {
        m := newTestAppModel()
        m.Setenv("PATH_INFO", "/test")

        a, err := NewApp(AppModel(m))
        if err != nil { t.Error(err) }

        a.Handle("/test", FuncHandler(func(w io.Writer, app *App) {
                app.SetHeader("Content-Type", "text/html")
                fmt.Fprint(w, "test-string")
        }))
        a.Exec() // produce the output

        var n int
        str := m.buffer.String()

        //fmt.Printf(str)

        n = strings.Index(str, "\nContent-Type: text/html")
        if n == -1 {
                t.Error("FuncHandler: no Content-Type header:\n", str)
        }

        n = strings.Index(str, "\n\n")
        if n == -1 {
                t.Error("FuncHandler: expecting \\n\\n in the output")
        } else {
                str = str[n+2:len(str)]
                if str != "test-string" { t.Error("FuncHandler: wrong output:\n", str) }
        }
}

func TestCustomHandler(t *testing.T) {
        m := newTestAppModel()

        a, err := NewApp(AppModel(m))
        if err != nil { t.Error(err) }

        h := &customHandler{ "test" }
        a.HandleDefault(h)
        a.Exec()

        str := m.buffer.String()
        n := strings.Index(str, "\n\n")
        if n == -1 { t.Error("custom: wrong output\n", str) }

        if str[n+2:len(str)] != "test" {
                t.Error("custom: expecting 'test'")
        }
}

func TestSessionPersistent(t *testing.T) {
        sid := ""
        {
                m := newTestAppModel()
                m.Setenv("PATH_INFO", "/test")

                a, err := NewApp(AppModel(m))
                if err != nil { t.Error(err) }

                a.Handle("/test", FuncHandler(func(w io.Writer, app *App) {
                        app.SetHeader("Content-Type", "text/html")
                        fmt.Fprint(w, "test-string")
                }))
                a.Exec() // produce the output

                str := m.buffer.String()
                n := strings.Index(str, "Set-Cookie:")
                if n == -1 { t.Error("no Set-Cookie for", cookieSessionId, str) }

                ln := strings.Index(str[n:len(str)], "\n")
                if ln == -1 { t.Error("bad output", str) }
                n = strings.Index(str[n:ln], cookieSessionId)
                if n == -1 { t.Error("no cookie",cookieSessionId,"in",str) }

                sid = str[n+len(cookieSessionId)+1:ln]
                if sid=="" { t.Error("empty session id",str) }
        }
        {
                m := newTestAppModel()
                m.Setenv("PATH_INFO", "/test")
                m.Setenv("HTTP_COOKIE", cookieSessionId+"="+sid)

                a, err := NewApp(AppModel(m))
                if err != nil { t.Error(err) }

                a.Handle("/test", FuncHandler(func(w io.Writer, app *App) {
                        app.SetHeader("Content-Type", "text/plain")
                        fmt.Fprint(w, "test-string")
                }))
                a.Exec() // produce the output

                str := m.buffer.String()
                if str=="" { t.Error("empty output") }

                n := strings.Index(str, cookieSessionId+"=")
                if n != -1 { t.Error("session persist failed:\n", str) }

                //fmt.Printf("%s\n")
        }
}

func TestViewTemplate(t *testing.T) {
        m := newTestAppModel()
        m.Setenv("PATH_INFO", "/test")

        a, err := NewApp(AppModel(m))
        if err != nil { t.Error(err) }

        a.config.Title = "test"
        a.Handle("/test", NewView("test.tpl"))
        a.Exec() // produce the output

        str := m.buffer.String()
        if str=="" { t.Error("empty output") }

        n := strings.Index(str, "\n\n")
        if n == -1 {
                t.Error("expecting \\n\\n in the output")
        } else {
                str = str[n+2:len(str)]
                if str != "<b>title</b>: test\n" {
                        t.Error("template: wrong output:\n", str)
                }
        }
}



