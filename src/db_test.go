package web

import (
        "testing"
)

func TestDatabase(t *testing.T) {
        db := NewDatabase()
        err := db.Connect("localhost", "test", "abc", "dusell")
        if err != nil { t.Error(err); goto finish }
        defer db.Close()

        sql := `SELECT a, b, c FROM table_test LIMIT 10`
        res, err := db.Query(sql)
        if err != nil { t.Error(err); goto finish }

        for {
                row, err := res.FetchRow();
                if err != nil { t.Error("FetchRow: failed"); goto finish }
                if row == nil { break }
                for k, v := range row {
                        switch res.GetFieldName(k) {
                        case "a":
                        case "b":
                                if v != "name" { t.Error(v); goto finish }
                        case "c":
                                if v != "long text should go here..." { t.Error(v); goto finish }
                        default:
                                t.Error("unknown field: ", res.GetFieldName(k), v)
                        }
                }
        }

        sql = `SELECT a, b, c FROM table_test WHERE a<? LIMIT 10`
        stmt, err := db.Prepare(sql)
        if err != nil { t.Error(err); goto finish }

        res, err = stmt.Execute(10)
        if err != nil { t.Error(err); goto finish }

        stmt.Close()

        for {
                row, err := res.FetchRow();
                if err != nil { t.Error("FetchRow: failed"); goto finish }
                if row == nil { break }
                for k, v := range row {
                        switch res.GetFieldName(k) {
                        case "a":
                        case "b":
                                if v != "name" { t.Error(v) }
                        case "c":
                                if v != "long text should go here..." { t.Error(v) }
                        default:
                                t.Error("unknown field: ", res.GetFieldName(k), v)
                        }
                }
        }

finish:
}
