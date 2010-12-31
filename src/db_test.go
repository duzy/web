package web

import (
        "testing"
        "fmt"
)

func TestDatabase(t *testing.T) {
        db := NewDatabase()
        fmt.Print("connect...\n");
        err := db.Connect("localhost", "test", "abc", "dusell")
        if err != nil { t.Error(err); goto finish }
        defer db.Close()

        sql := `SELECT a, b, c FROM table_test LIMIT 10`
        fmt.Print("Query...\n");
        res, err := db.Query(sql)
        if err != nil { t.Error(err); goto finish }

        for {
                fmt.Print("fetch...\n");
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
        fmt.Printf("prepare...\n");
        stmt, err := db.Prepare(sql)
        if err != nil { t.Error(err); goto finish }

        fmt.Printf("execute...\n");
        res, err = stmt.Execute(10)
        if err != nil { t.Error(err); goto finish }

        stmt.Close()

        for {
                fmt.Printf("2fetch...\n")
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
