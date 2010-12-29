package web

import (
        "crypto/md5"
        "fmt"
        "time"
        "io"
        "os"
        "strings"
        "bytes"
)

const (
        SQL_CREATE_SESSION_TABLE = `
CREATE TABLE IF NOT EXISTS table_web_sessions(
  sid CHAR(64) PRIMARY KEY,
  props TEXT
)
`
        SQL_INSERT_SESSION_ROW = `
INSERT INTO table_web_sessions(sid,props) VALUES(?,?)
  ON DUPLICATE KEY UPDATE props=VALUES(props)
`
        SQL_SELECT_SESSION_ROW = `
SELECT props FROM table_web_sessions WHERE sid=? LIMIT 1
`
)

type Session struct {
        changed bool
        id string
        props map[string]string
}

// Persist the session between connections.
type SessionPersister interface {
        io.ReadWriteCloser
        // TODO: add LoadSession
        // TODO: add SaveSession or saveSession?
}

type FSSessionPersister struct { file *os.File }
type DBSessionPersister struct {
        sid string
        db Database
        buf *bytes.Buffer
}

var defaultPersisterConfigFS = &PersisterConfigFS{
        "/tmp/web/sessions/",
}

// Make a new session Persister.
// The session id (sid) must be more than 5 chars length.
func NewSessionPersister(sid string, cfg PersisterConfig) (p SessionPersister, err os.Error) {
        if v, ok := cfg.(*PersisterConfigFS); ok {
                p, err = newFSSessionPersister(sid, v)
        } else if v, ok := cfg.(*PersisterConfigDB); ok {
                p, err = newDBSessionPersister(sid, v)
        }
        return
}

func newFSSessionPersister(sid string, cfg *PersisterConfigFS) (p SessionPersister, err os.Error) {
        fs := &FSSessionPersister{}

        const dirLen = 5
        if len(sid) < dirLen { goto finish }

        var n int
        d := cfg.Location

        if d=="" {
                d = defaultPersisterConfigFS.Location
        }

        if d[len(d)-1] != '/' { d += "/" }

        for n=0; n < dirLen; n+=1 {
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

func newDBSessionPersister(sid string, cfg *PersisterConfigDB) (p SessionPersister, err os.Error) {
        dbm := GetDBManager()
        db, err := dbm.GetDatabase(&(cfg.DatabaseConfig))
        if err != nil { /* TODO: error */ goto finish }

        sql := SQL_CREATE_SESSION_TABLE
        _, err = db.Query(sql)
        if err != nil { /* TODO: error */ goto finish }

        //stmt, err := db.NewStatement()
        //if err != nil { /* TODO: error */ goto finish }

        stmt, err := db.Prepare(SQL_SELECT_SESSION_ROW)
        if err != nil { /* TODO: error */ stmt.Close(); goto finish }

        //stmt.BindParams(sid)
        res, err := stmt.Execute(sid)
        if err != nil { /* TODO: error */ stmt.Close(); goto finish }

        stmt.Close()

        props := ""
        row := res.FetchRow()
        if row != nil {
                props = fmt.Sprintf("%s", row[0])
        } else { /* TODO: error */ }

        //fmt.Printf("sid=%s\nprops:\n%s", sid, props)
        
        p = &DBSessionPersister{ sid, db, bytes.NewBufferString(props) }

finish:
        return
}

func (p *FSSessionPersister) Close() os.Error { return p.file.Close() }
func (p *FSSessionPersister) Read(b []byte) (n int, err os.Error) {
        n, err = p.file.Read(b)
        return
}
func (p *FSSessionPersister) Write(b []byte) (n int, err os.Error) {
        n, err = p.file.Write(b)
        return
}

func (p *DBSessionPersister) Read(b []byte) (n int, err os.Error) {
        n, err = p.buf.Read(b)
        return
}
func (p *DBSessionPersister) Write(b []byte) (n int, err os.Error) {
        n, err = p.buf.Write(b)
        return
}
func (p *DBSessionPersister) Close() (err os.Error) {
        stmt, err := p.db.Prepare(SQL_INSERT_SESSION_ROW)
        if err != nil { goto finish }
        _, err = stmt.Execute(p.sid, p.buf.String())
        if err != nil { stmt.Close(); goto finish }
        stmt.Close()
finish:
        err = p.db.Close()
        return
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
        return
}

// TODO: make method of SessionPersister, or avoid 'cfg' parameter
func LoadSession(id string, cfg PersisterConfig) (s *Session, err os.Error) {
        p, err := NewSessionPersister(id, cfg)
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

func propEscape(v string) string {
        s := ""
        for {
                var i int
                if i = strings.Index(v, "\\"); i != -1 {
                        s += v[0:i] + "\\\\"
                        v = v[i+1:len(v)] // slicing
                } else if i = strings.Index(v, "\n"); i != -1 {
                        s += v[0:i] + "\\n"
                        v = v[i+1:len(v)] // slicing
                } else {
                        s += v
                        break
                }
        }
        return s
}

func propUnescape(v string) string {
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
                fmt.Fprintf(w, "%s:%s\n", k, propEscape(v))
        }
        return
}

func ReadSession(r io.Reader) (s *Session, err os.Error) {
        s = new(Session)
        s.props = make(map[string]string)
        n, err := fmt.Fscanf(r, "id:%s", &s.id)
        if n == 1 && err == nil {
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
                                s.props[k] = propUnescape(v)
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

func (s *Session) save(cfg PersisterConfig) (err os.Error) {
        if s.changed {
                var p SessionPersister
                p, err = NewSessionPersister(s.id, cfg)
                if err != nil {
                        fmt.Fprintf(os.Stderr, "error: %s\n", err)
                        goto finish
                }
                if p == nil {
                        fmt.Fprintf(os.Stderr, "error: can't new a session\n")
                        goto finish
                }
                defer p.Close()

                err = WriteSession(p, s)
                if err != nil {
                        fmt.Fprintf(os.Stderr, "error: %s\n", err)
                        goto finish
                }
                s.changed = false
        }

finish:
        return
}
