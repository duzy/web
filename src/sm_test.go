package web

import (
        "testing"
        "os"
)

var configFS = &PersisterConfigFS{
Location: "/tmp/web-test/PersisterFS",
}

var configDB = &PersisterConfigDB{
        DatabaseConfig {
        Host: "localhost",
        User: "test",
        Password: "abc",
        Database: "dusell",
        },
}

func testSaveLoadSession(t *testing.T, cfg PersisterConfig) (sid string) {
        {
                s := NewSession()
                sid = s.Id()
                if sid == "" { t.Error("Failed NewSession()") }
                v := s.Set("test", "test-value")
                if v != "" { t.Error("Previous test value is not empty") }
                v = s.Get("test")
                if v != "test-value" { t.Error("Set property failed") }
                v = s.Set("multiline", "line1\nline2\nline3\nline4")
                if v != "" { t.Error("Previous 'multiline' value is not empty") }
                v = s.Get("multiline")
                if v != "line1\nline2\nline3\nline4" { t.Error("Set property failed") }
                if err := s.save(cfg); err != nil { t.Error("Failed session save:", err) }
        }
        {
                s, err := LoadSession(sid, cfg)
                if err != nil { t.Error(err); goto finish }
                if s == nil { t.Error("nil session"); goto finish }
                if s.Id() == "" { t.Error("Failed LoadSession()"); goto finish }
                if s.Get("test") != "test-value" {
                        t.Error("Session props persist error")
                        goto finish
                }
                if s.Get("multiline") != "line1\nline2\nline3\nline4" {
                        t.Error("Session props persist error")
                        goto finish
                }
        }
finish:
        return
}

func TestSessionPersistFS(t *testing.T) {
        sid := testSaveLoadSession(t, configFS)
        d := configFS.Location
        d += "/" + sid[0:1]
        d += "/" + sid[1:2]
        d += "/" + sid[2:3]
        d += "/" + sid[3:4]
        d += "/" + sid[4:5]
        d += "/" + sid[5:len(sid)]
        if _, err := os.Stat(d); err != nil {
                t.Error(err)
        }
}

func TestSessionPersistDB(t *testing.T) {
        testSaveLoadSession(t, configDB)
}
