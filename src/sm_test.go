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

func testSaveLoadSession(t *testing.T, cfg *AppConfig_Persister) {
        sid := ""
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
                if !s.save(cfg) { t.Error("Failed session save") }
        }
        {
                s, err := LoadSession(sid, cfg)
                if err != nil { t.Error(err) }
                if s.Id() == "" { t.Error("Failed LoadSession()") }
                if s.Get("test") != "test-value" {
                        t.Error("Session props persist error")
                }
                if s.Get("multiline") != "line1\nline2\nline3\nline4" {
                        t.Error("Session props persist error")
                }
        }
}

func TestSessionPersistFS(t *testing.T) {
        testSaveLoadSession(t, &AppConfig_Persister{ configFS })
}

func TestSessionPersistDB(t *testing.T) {
        testSaveLoadSession(t, &AppConfig_Persister{ configDB })
}
