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

// customViewModel must be a TemplateGetter and FieldsMaker
type customViewModel struct {
        template string // eg. <title>{title}</title>
        field1 string
        field2 string
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

func (v *customViewModel) GetTemplate() string {
        return v.template
}

func (v *customViewModel) MakeFields(app *App) (fields interface{}) {
        db, err := app.GetDatabase("dusell")
        if err == nil {
                v.field1 = "bold"
                v.field2 = "italic"
                db.Close() // TODO: do something meaningful with db
        } else {
                v.field1 = "ERROR"
                v.field2 = err.String()
        }
        fields = v
        return
}

func TestFuncHandler(t *testing.T) {
        m := newTestAppModel()
        m.Setenv("PATH_INFO", "/test")

        a, err := NewApp(AppModel(m))
        if err != nil { t.Error(err); return }

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
                return
        }

        n = strings.Index(str, "\n\n")
        if n == -1 {
                t.Error("FuncHandler: expecting \\n\\n in the output")
                return
        } else {
                str = str[n+2:len(str)]
                if str != "test-string" { t.Error("FuncHandler: wrong output:\n", str); return }
        }
}

func TestCustomHandler(t *testing.T) {
        m := newTestAppModel()

        a, err := NewApp(AppModel(m))
        if err != nil { t.Error(err); return }

        h := &customHandler{ "test" }
        a.HandleDefault(h)
        a.Exec()

        str := m.buffer.String()
        n := strings.Index(str, "\n\n")
        if n == -1 { t.Error("custom: wrong output\n", str); return }

        if str[n+2:len(str)] != "test" {
                t.Error("custom: expecting 'test'")
                return
        }
}

func TestSessionPersistent(t *testing.T) {
        sid := ""
        {
                m := newTestAppModel()
                m.Setenv("PATH_INFO", "/test")

                a, err := NewApp(AppModel(m))
                if err != nil { t.Error(err); return }

                a.Handle("/test", FuncHandler(func(w io.Writer, app *App) {
                        app.SetHeader("Content-Type", "text/html")
                        fmt.Fprint(w, "test-string")
                        a.Session().Set("test", "test-session");
                }))
                a.Exec() // produce the output

                str := m.buffer.String()
                n := strings.Index(str, "Set-Cookie:")
                if n == -1 { t.Error("no Set-Cookie for", cookieSessionId, str); return }

                ln := strings.Index(str[n:len(str)], "\n")
                if ln == -1 { t.Error("bad output", str); return }
                n = strings.Index(str[n:ln], cookieSessionId)
                if n == -1 { t.Error("no cookie",cookieSessionId,"in",str); return }

                sid = str[n+len(cookieSessionId)+1:ln]
                if sid=="" { t.Error("empty session id",str); return }
        }
        {
                m := newTestAppModel()
                m.Setenv("PATH_INFO", "/test")
                m.Setenv("HTTP_COOKIE", cookieSessionId+"="+sid)

                a, err := NewApp(AppModel(m))
                if err != nil { t.Error(err); return }

                a.Handle("/test", FuncHandler(func(w io.Writer, app *App) {
                        app.SetHeader("Content-Type", "text/plain")
                        fmt.Fprint(w, "test-string")
                        v := a.Session().Get("test")
                        if v != "test-session" { t.Error("session-prop: persist error:", v); return }
                }))
                a.Exec() // produce the output

                str := m.buffer.String()
                if str=="" { t.Error("empty output"); return }

                n := strings.Index(str, cookieSessionId+"=")
                if n != -1 { t.Error("session persist failed:\n", str); return }

                //fmt.Printf("%s\n")
        }
}

func TestViewTemplate(t *testing.T) {
        m := newTestAppModel()
        m.Setenv("PATH_INFO", "/test")

        a, err := NewApp(AppModel(m))
        if err != nil { t.Error(err); return }

        a.config.Title = "test"
        a.Handle("/test", NewView("test.tpl"))
        a.Exec() // produce the output

        str := m.buffer.String()
        if str=="" { t.Error("empty output"); return }

        n := strings.Index(str, "\n\n")
        if n == -1 {
                t.Error("expecting \\n\\n in the output"); return
        } else {
                str = str[n+2:len(str)]
                if str != "<b>title</b>: test\n" {
                        t.Error("template: wrong output:\n", str); return
                }
        }
}

func TestNewAppFromConfig(t *testing.T) {
        a, err := NewApp("test_app.json")
        if err != nil { t.Error(err); return }
        if a.config == nil { t.Error("app not configured"); return }
        if a.config.Title != "test app via json" {
                t.Error("app from json: title not matched:",a.config.Title)
                return
        }
        if a.config.Model != "CGI" {
                t.Error("app from json: model not matched:",a.config.Model)
                return
        }
        if v, ok := a.config.Persister.(*PersisterConfigFS); !ok {
                t.Error("app from json: not FS persister:",v); return
        } else {
                if v.Location != "/tmp/web-test/sessions" {
                        t.Error("app from json: wrong location:",v.Location)
                        return
                }
        }

        var m *testAppModel
        if v, ok := a.model.(*CGIModel); !ok {
                t.Error("app from json: not CGIModel:",a.model); return
        } else {
                // convert the CGIModel into testAppModel
                writer := bytes.NewBufferString("")
                reader := bytes.NewBufferString("")
                m = &testAppModel{ v, writer, reader }
                a.model = AppModel(m)
        }

        a.HandleDefault(NewView("test.tpl"))
        a.Exec() // produce the output

        str := m.buffer.String()
        if str=="" { t.Error("app from json: empty output"); return }

        n := strings.Index(str, "\n\n")
        if n == -1 {
                t.Error("expecting \\n\\n in the output"); return
        } else {
                str = str[n+2:len(str)]
                if str != "<b>title</b>: test app via json\n" {
                        t.Error("template: wrong output:\n", str); return
                }
        }
}

func TestCustomViewModelAndAppGetDatabase(t *testing.T) {
        a, err := NewApp("test_app.json")
        if err != nil { t.Error(err); return }
        if a.config == nil { t.Error("app not configured"); return }

        var m *testAppModel
        if v, ok := a.model.(*CGIModel); !ok {
                t.Error("app from json: not CGIModel:", a.model)
                return
        } else {
                // convert the CGIModel into testAppModel
                writer := bytes.NewBufferString("")
                reader := bytes.NewBufferString("")
                m = &testAppModel{ v, writer, reader }
                a.model = AppModel(m)
        }

        cv := &customViewModel{ }
        cv.template = "<b>{field1}</b>:<i>{field2}</i>"

        a.HandleDefault(NewView(cv))
        a.Exec() // produce the output

        str := m.buffer.String()
        n := strings.Index(str, "\n\n")
        if n == -1 {
                t.Error("expecting \\n\\n in the output"); return
        }

        str = str[n+2:len(str)]
        if str != "<b>bold</b>:<i>italic</i>" {
                t.Error("custom-view and db:", str); return
        }
}

