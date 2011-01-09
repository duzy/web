package web

import (
        "testing"
)

func TestDBManager(t *testing.T) {
        var cfg = &DatabaseConfig {
        Host: "localhost",
        User: "dusellco_test",
        Password: "abc",
        Database: "dusellco_test",
        }

        db, err := GetDBManager().GetDatabase(cfg)
        if err != nil { t.Error(err); goto finish }
        if db == nil { t.Error("no db obtained"); goto finish }
        defer db.Close()

        //err = db.Ping()
        //if err != nil { t.Error(err) }

        db2, err := GetDBManager().GetDatabase(cfg)
        if db2 == nil { t.Error(err) }
        defer db2.Close()

        //err = db2.Ping()
        //if err != nil { t.Error(err) }

        if db != db2 { t.Error("returned two db with one cfg") }

        if rec, ok := db.(*dbrecord); ok {
                if rec.useCount < 2 { t.Error("wrong dbrecord.useCount") }
        } else {
                t.Error("not a dbrecord returned")
        }

finish:
}
