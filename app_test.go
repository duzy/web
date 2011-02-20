package web

import (
        "testing"
        "strings"
        "bytes"
        "fmt"
        "os"
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
        buffer := bytes.NewBufferString("")
        reader := bytes.NewBufferString("")
        cgi := &CGIModel{ make(map[string]string), nil, buffer } //NewCGIModel()
        m = &testAppModel{ cgi, buffer, reader }
        m.Setenv("SERVER_PROTOCOL", "HTTP/1.1")
        return
}

func (h *customHandler) HandleRequest(request *Request, response *Response) (err os.Error) {
        w := response.BodyWriter
        fmt.Fprint(w, h.content)
        return
}

func (v *customViewModel) GetTemplate() string {
        return v.template
}

func (v *customViewModel) MakeFields(app *App) (fields interface{}) {
        fields = v
        return
}

func TestFuncHandler(t *testing.T) {
        m := newTestAppModel()
        m.Setenv("PATH_INFO", "/test")
        m.Setenv("REQUEST_URI", "/test")
        m.Setenv("HTTP_COOKIE", "foo=bar")

        a, err := NewApp(AppModel(m))
        if err != nil { t.Error(err); return }

        a.Handle("/test", FuncHandler(func(request *Request, response *Response) (err os.Error) {
                response.Header["Content-Type"] = "text/html"
                if request.HttpCookie != "foo=bar" {
                        t.Errorf("request.HttpCookie not correct: %v", request.HttpCookie)
                }
                if c := request.Cookie("foo"); c != nil {
                        if c.Content != "bar" {
                                t.Errorf("cookie 'foo' not correct: %v", c)
                        }
                } else {
                        t.Errorf("no cookie 'foo' in: %v", request.cookies)
                }
                fmt.Fprint(response.BodyWriter, "test-string")
                //fmt.Print(request)
                return
        }))
        err = a.Exec() // produce the output
        if err != nil {
                t.Errorf("App.Exec: %v", err)
                return
        }

        var n int
        str := m.buffer.String()

        if str == "" {
                t.Error("FuncHandler: no output:\n")
                return
        }

        n = strings.Index(str, "Content-Type: text/html\n")
        if n == -1 {
                t.Error("FuncHandler: no 'Content-Type' in:\n", str)
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
        m.Setenv("REQUEST_URI", "/")
        
        a, err := NewApp(AppModel(m))
        if err != nil { t.Error(err); return }

        h := &customHandler{ "test" }
        a.HandleDefault(h)
        err = a.Exec()
        if err != nil {
                t.Errorf("App.Exec: %v", err)
                return
        }

        str := m.buffer.String()
        n := strings.Index(str, "\ntest")
        if n == -1 { t.Error("customHandler: wrong output:\n", str); return }

        if str[n+1:len(str)] != "test" {
                t.Error("customHandler: expecting 'test'")
                return
        }
}

