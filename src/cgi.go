package web

import (
        "os"
        "io"
        "fmt"
        "runtime"
        "http"
        "net"
)

// Implements AppModel for CGI web.App.
type CGIModel struct {
        overides map[string]string
        request *Request
}

var cgiRequest *Request

// Parse HTTP version: "HTTP/1.2" -> (1, 2, true), copied from http/request.go
func parseHTTPVersion(vers string) (int, int, bool) {
	if len(vers) < 5 || vers[0:5] != "HTTP/" {
		return 0, 0, false
	}
	major, i, ok := atoi(vers, 5)
	if !ok || i >= len(vers) || vers[i] != '.' {
		return 0, 0, false
	}
	var minor int
	minor, i, ok = atoi(vers, i+1)
	if !ok || i != len(vers) {
		return 0, 0, false
	}
	return major, minor, true
}

func initCGIRequest() bool {
        if cgiRequest != nil {
                return true
        }

        request := new(Request)
        ok := false
        //scheme + "://" + os.Getenv("SERVER_NAME") + 
        request.Method = os.Getenv("REQUEST_METHOD")
        request.Proto = os.Getenv("SERVER_PROTOCOL")
        request.ProtoMajor, request.ProtoMinor, ok = parseHTTPVersion(request.Proto)
        if !ok {
                err = os.NewError("malformed HTTP version: "+request.Proto)
                return
        }

        request.RawURL = os.Getenv("REQUEST_URI")
        request.URL, err = http.ParseURL(request.RawURL)
        if err != nil {
                return
        }

        request.Header = make(map[string]string)
        for _, v := os.Environ() {
                if 5 < len(v) && v[0:5] == "HTTP_" {
                        if kv := strings.Split(v); kv != nil {
                                // TODO: convert the uppercase names?
                                request.Header[kv[0][5:]] = kv[1]
                        }
                }
        }

        request.Host = os.Getenv("HTTP_HOST")
        request.Referer = os.Getenv("HTTP_REFERER")
        request.UserAgent = os.Getenv("HTTP_USER_AGENT")

        request.Close = false
        request.Body = os.Stdin
        request.ContentLength, _ = strconv.Atoi64(os.Getenv("HTTP_CONTENT_LENGTH"))

        cgiRequest = request
        return true
}

func NewCGIModel() (m AppModel) /*, err os.Error)*/ {
        if !initCGIRequest() {
                return
        }

        cgi := &CGIModel{
                make(map[string]string),
                cgiRequest,
        }

        m = AppModel(cgi)
        return
}

func (cgi *CGIModel) Getenv(k string) (v string) {
        v = cgi.overides[k]
        if v == "" {
                v = os.Getenv(k)
        }
        return
}

func (cgi *CGIModel) Setenv(k, v string) (prev string) {
        prev = cgi.overides[k]
        if prev == "" {
                prev = os.Getenv(k)
        }
        cgi.overides[k] = v
        return
}


func (cgi *CGIModel) RequestMethod() string {
        return cgi.Getenv("REQUEST_METHOD")
}

func (cgi *CGIModel) PathInfo() string {
        return cgi.Getenv("PATH_INFO")
}

func (cgi *CGIModel) QueryString() string {
        return cgi.Getenv("QUERY_STRING")
}

func (cgi *CGIModel) ScriptName() string {
        return cgi.Getenv("SCRIPT_NAME")
}

func (cgi *CGIModel) HttpCookie() string {
        return cgi.Getenv("HTTP_COOKIE")
}

func (cgi *CGIModel) ResponseWriter() (w io.Writer) {
        w = os.Stdout
        return
}

func (cgi *CGIModel) RequestReader() (r io.Reader) {
        r = os.Stdin
        return
}

func (cgi *CGIModel) ProcessRequests(rp RequestProcessor) (err os.Error) {
        err = rp.ProcessRequest()
        return
}

func (cgi *CGIModel) GetRequest() (req *Request) {
        req = cgi.request
        return
}

func CGIHandleErrors() {
        if err := recover(); err != nil {
                //stack := make([]uintptr, 5)

                stack, file, line, ok := runtime.Caller(5)
                if !ok {
                        file = "???"
                        line = 0
                }

                f := runtime.FuncForPC(stack)

                fmt.Fprintf(os.Stdout, "Content-Type: text/html; charset=utf-8\n\n")
                fmt.Fprintf(os.Stdout, `<font color="red"><b>error</b>:</font> %v<p>`, err)
                fmt.Fprintf(os.Stdout, `%s:%d<br/>`, file, line)
                if f != nil {
                        file, line = f.FileLine(stack)
                        fmt.Fprintf(os.Stdout, `%s:%d: %s<br/>`, file, line, f.Name())
                }
                fmt.Fprintf(os.Stdout, `</p>`)
        }
}
