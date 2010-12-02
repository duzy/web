package web

import (
        "os"
        "io"
)

func NewCGIModel() (m AppModel) {
        cgi := new(CGIModel)
        m = AppModel(cgi)
        return
}

// Implements AppModel for CGI web.App.
type CGIModel struct {
}

func (cgi *CGIModel) RequestMethod() string {
        return os.Getenv("REQUEST_METHOD")
}

func (cgi *CGIModel) PathInfo() string {
        return os.Getenv("PATH_INFO")
}

func (cgi *CGIModel) QueryString() string {
        return os.Getenv("QUERY_STRING")
}

func (cgi *CGIModel) ResponseWriter() (w io.Writer) {
        w = os.Stdout
        return
}

func (cgi *CGIModel) RequestReader() (r io.Reader) {
        r = os.Stdin
        return
}
