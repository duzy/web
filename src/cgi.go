package web

import (
        "io"
        "os"
        "fmt"
        "runtime"
        "http"
        "strconv"
        "strings"
        //"net"
) 

// Implements AppModel for CGI web.App.
type CGIModel struct {
        overides map[string]string
        request *Request
        responseWriter io.Writer
}

// Parse HTTP version: "HTTP/1.2" -> (1, 2, true), copied from http/request.go
func parseHTTPVersion(vers string) (int, int, bool) {
	if len(vers) < 5 || vers[0:5] != "HTTP/" {
		return 0, 0, false
	}
        i := strings.Index(vers[5:], ".")
        if i <= 0 {
                return 0, 0, false
        }
	major, err := strconv.Atoi(vers[5:5+i])
	if err != nil {
		return 0, 0, false
	}
	var minor int
	minor, err = strconv.Atoi(vers[5+i+1:])
	if err != nil {
		return 0, 0, false
	}
	return major, minor, true
}

func (cgi *CGIModel) initRequest(rm RequestManager) (err os.Error) {
        if cgi.request != nil {
                return
        }

        request, err := rm.GetRequest("")
        if err != nil {
                return
        }

        //scheme + "://" + cgi.Getenv("SERVER_NAME") + 

        request.Header = make(map[string]string)
        for _, v := range os.Environ() {
                if kv := strings.Split(v, "=", 2); kv != nil {
                        // fmt.Printf("%v, %d\n", kv, len(kv))
                        // TODO: convert the uppercase names?
                        request.Header[kv[0]] = kv[1]
                }
        }

        ok := false
        request.Method = cgi.Getenv("REQUEST_METHOD")
        request.Proto = cgi.Getenv("SERVER_PROTOCOL")
        request.ProtoMajor, request.ProtoMinor, ok = parseHTTPVersion(request.Proto)
        if !ok {
                err = os.NewError("malformed HTTP version: "+request.Proto)
                return
        }

        request.RawURL = cgi.Getenv("REQUEST_URI")
        request.URL, err = http.ParseURL(request.RawURL)
        if err != nil {
                return
        }

        request.Host = cgi.Getenv("HTTP_HOST")
        request.Referer = cgi.Getenv("HTTP_REFERER")
        request.UserAgent = cgi.Getenv("HTTP_USER_AGENT")

        request.Close = false
        request.Body = noCloseReader{ os.Stdin }
        request.ContentLength, _ = strconv.Atoi64(cgi.Getenv("HTTP_CONTENT_LENGTH"))

        request.Path = cgi.Getenv("PATH_INFO")
        request.QueryString = cgi.Getenv("QUERY_STRING")
        request.ScriptName = cgi.Getenv("SCRIPT_NAME")
        request.HttpCookie = cgi.Getenv("HTTP_COOKIE")
        request.cookies = ParseCookies(request.HttpCookie)

        if err = request.initSession(); err != nil {
                return
        }

        cgi.request = request
        return
}

func NewCGIModel() (m AppModel, err os.Error) {
        cgi := &CGIModel{
                make(map[string]string),
                nil, // request will be inited lately
                os.Stdout, // write response to Stdout
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

func (cgi *CGIModel) ProcessRequests(rm RequestManager) (err os.Error) {
        if err = cgi.initRequest(rm); err != nil {
                //fmt.Printf("ProcessRequests: %v\n", err)
                return
        }

        //fmt.Printf("request: %v\n", cgi.request)

        if cgi.request == nil {
                err = os.NewError("bad CGI Request")
                return
        }

        var response *Response
        response, err = rm.ProcessRequest(cgi.request)
        response.writeHeader(cgi.responseWriter)
        io.Copy(cgi.responseWriter, response.Body)
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
