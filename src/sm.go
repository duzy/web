package web

import (
        "crypto/md5"
        "fmt"
        "time"
        "io"
        "os"
        "strings"
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
        changed bool
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
        changed: true, // mark changed for saving
        id: genSid(),
        props: make(map[string]string),
        }

        //s.save() // check result?
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

func prop_escape(s string) string {
        s = strings.Replace(s, "\\", "\\\\", -1)
        s = strings.Replace(s, "\n", "\\n", -1)
        //fmt.Fprintf(os.Stdout, "s: %s\n", s)
        return s
}

func prop_unescape(v string) string {
        s := ""
        for {
                var i int
                if i = strings.Index(v, "\\"); i == -1 {
                        s += v[0:len(v)]
                        break
                }

                s += v[0:i]
                if len(v) == i+1 {
                        // FIXME: should the last '\' be ignored?
                        break;
                }

                // escape chars, TODO: support C escape chars?
                switch v[i+1:i+2] {
                case "\\": s += "\\"
                case "n": s += "\n"
                }
                v = v[i+2:len(v)] // reset slice
        }
        return s
}

func WriteSession(w io.Writer, s *Session) (err os.Error) {
        fmt.Fprintf(w, "id:%s\n", s.id)
        for k, v := range s.props {
                fmt.Fprintf(w, "%s:%s\n", k, prop_escape(v))
        }
        return
}

func ReadSession(r io.Reader) (s *Session, err os.Error) {
        s = new(Session)
        n, err := fmt.Fscanf(r, "id:%s", &s.id)
        if n == 1 && err == nil {
                s.props = make(map[string]string)
                for {
                        // FIXME: handle with multi-line property
                        var ln, k, v string
                        n, err = fmt.Fscanln(r, &ln)
                        if n != 1 || err != nil {
                                if err != nil && err.String() == "Scan: no data for string" {
                                        err = nil
                                }
                                break
                        }
                        n = strings.Index(ln, ":")
                        if 0 < n {
                                k = ln[0:n]
                                v = ln[n+1:len(ln)]
                                s.props[k] = prop_unescape(v)
                                s.changed = false
                        } else {
                                // FIXME: should break or just ignore?
                        }
                }
        }
        return
}

func (s *Session) Id() string { return s.id }

func (s *Session) Get(k string) string { return s.props[k] }
func (s *Session) Set(k, v string) (prev string) {
        prev = s.props[k]
        if prev != v {
                s.props[k] = v
                s.changed = true
        }
        return
}

func (s *Session) save() (saved bool) {
        saved = false
        if s.changed {
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
                saved = true
                s.changed = false
                //fmt.Fprintf(os.Stdout, "session-saved: %s\n", s.id)
        }

finish:
        return
}
