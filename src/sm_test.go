package web

import (
        "testing"
)

var configFS = &AppConfig_PersisterFS{
Location: "/tmp/web/test/PersisterFS",
}

var configDB = &AppConfig_PersisterDB{
        AppConfig_Database {
        Host: "localhost",
        User: "test",
        Password: "abc",
        Database: "dusell",
        },
}

func TestSessionSave(t *testing.T) {
        cfg := &AppConfig_Persister{ configFS }
        sid := ""
        {
                s := NewSession()
                sid = s.Id()
                if sid == "" { t.Error("Failed NewSession()") }
                s.Set("test", "test-value")
                s.save(cfg)
        }
        {
                s, err := LoadSession(sid, cfg)
                if err != nil { t.Error(err) }
                if s.Id() == "" { t.Error("Failed NewSession()") }
                if s.Get("test") != "test-value" {
                        t.Error("Session props persist error")
                }
        }
}
