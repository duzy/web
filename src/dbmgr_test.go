package web

import (
        "testing"
)

func TestDBManager(t *testing.T) {
        var cfg = &AppConfig_Database {
        Host: "localhost",
        User: "test",
        Password: "abc",
        Database: "dusell",
        }

        db, err := GetDBManager().GetDatabase(cfg)
        if db == nil { t.Error(err) }
        defer db.Close()

        err = db.Ping()
        if err != nil { t.Error(err) }

        db2, err := GetDBManager().GetDatabase(cfg)
        if db2 == nil { t.Error(err) }
        defer db2.Close()

        err = db2.Ping()
        if err != nil { t.Error(err) }

        if db != db2 { t.Error("returned two db with one cfg") }

        if rec, ok := db.(*dbrecord); ok {
                if rec.useCount < 2 { t.Error("wrong dbrecord.useCount") }
        } else {
                t.Error("not a dbrecord returned")
        }
}
