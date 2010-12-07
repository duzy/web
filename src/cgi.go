package web

import (
        "os"
        "io"
)

// Implements AppModel for CGI web.App.
type CGIModel struct {
        overides map[string]string
}

func NewCGIModel() (m AppModel) {
        cgi := &CGIModel{ make(map[string]string) }
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

func (cgi *CGIModel) ResponseWriter() (w io.Writer) {
        w = os.Stdout
        return
}

func (cgi *CGIModel) RequestReader() (r io.Reader) {
        r = os.Stdin
        return
}

