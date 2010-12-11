package web

import (
        "crypto/md5"
        "fmt"
        "time"
        "io"
        "os"
)

// Persist the session between connections.
type SessionPersister interface {
        io.ReadWriteCloser
}

// Make a new session Persister.
// The session id (sid) must be more than 5 chars length.
func NewSessionPersister(sid string) (p SessionPersister, err os.Error) {
        fs := &FSSessionPersister{}

        if len(sid) < 5 { goto finish }

        var n int
        d := "/tmp/web/sessions/" // make a configurable base path
        for n=0; n < 5; n+=1 {
                d += sid[n:n+1] + "/"
        }
        err = os.MkdirAll(d, 0700)
        if err != nil { panic(err) }

        d += sid[n:len(sid)]

        f, err := os.Open(d, os.O_RDWR|os.O_CREAT, 0600)
        if err == nil {
                fs.file = f
                p = SessionPersister(fs)
        }

finish:
        return
}

type FSSessionPersister struct { file *os.File }
func (p *FSSessionPersister) Close() os.Error { return p.file.Close() }
func (p *FSSessionPersister) Read(b []byte) (n int, err os.Error) {
        n, err = p.file.Read(b)
        return
}
func (p *FSSessionPersister) Write(b []byte) (n int, err os.Error) {
        n, err = p.file.Write(b)
        return
}

type Session struct {
        id string
        props map[string]string
}

func genSid() (id string) {
        c := md5.New()
        fmt.Fprintf(c, "%v", time.Nanoseconds())
        id = fmt.Sprintf("%x", c.Sum())
        return
}

func NewSession() (s *Session) {
        s = &Session{
        id: genSid(),
        props: make(map[string]string),
        }

        p, err := NewSessionPersister(s.id)
        if err != nil {
                fmt.Fprintf(os.Stderr, "error: %s\n", err)
                goto finish
        }
        defer p.Close()

        err = WriteSession(p, s)
        if err != nil {
                fmt.Fprintf(os.Stderr, "error: %s\n", err)
                goto finish
        }

finish:
        return
}

func LoadSession(id string) (s *Session, err os.Error) {
        p, err := NewSessionPersister(id)
        if err != nil {
                fmt.Fprintf(os.Stderr, "error: %s\n", err)
                goto finish
        }

        defer p.Close()

        s, err = ReadSession(p)
        if err != nil {
                fmt.Fprintf(os.Stderr, "error: %s\n", err)
                goto finish
        }

finish:
        return
}

func WriteSession(w io.Writer, s *Session) (err os.Error) {
        fmt.Fprintf(w, "id:%s\n", s.id)
        for k, v := range s.props {
                fmt.Fprintf(w, "%s:%s\n", k, v)
        }
        return
}

func ReadSession(r io.Reader) (s *Session, err os.Error) {
        s = new(Session)
        n, err := fmt.Fscanf(r, "id:%s\n", &s.id)
        if n == 1 && err == nil {
                for {
                        var k, v string
                        n, err = fmt.Fscanf(r, "%s:%s\n", &k, &v)
                        if n != 2 || err != nil {
                                if err.String() == "Scan: no data for string" {
                                        err = nil
                                }
                                break
                        }
                        s.props[k] = v
                }
        }
        return
}

func (s *Session) Id() string { return s.id }

func (s *Session) Get(k string) string { return s.props[k] }
func (s *Session) Set(k, v string) (prev string) {
        prev = s.props[k]
        s.props[k] = v
        return
}

