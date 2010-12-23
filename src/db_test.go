package web

import (
        "testing"
        "./_obj/mysql"
)

func Test_mysql(t *testing.T) {
        db := mysql.New()
        db.Logging = false
        db.Connect("localhost", "test", "abc", "dusell")
        if db.Errno != 0 {
                t.Error("Error", db.Errno, db.Error)
                goto finish
        }
        defer db.Close()

        db.Query("SET NAMES utf8")
        if db.Errno != 0 {
                t.Error("Error", db.Errno, db.Error)
                goto finish
        }

        sql := `
CREATE TABLE IF NOT EXISTS table_test(
a INT AUTO_INCREMENT PRIMARY KEY,
b VARCHAR(64),
c TEXT
)`
        db.Query(sql)
        if db.Errno != 0 {
                t.Error("Error", db.Errno, db.Error)
                goto finish
        }

        sql = `
INSERT INTO table_test(b,c) VALUES(?,?)
`
        stmt, _ := db.InitStmt()
        stmt.Prepare(sql)
        if stmt.Errno != 0 {
                t.Error("Error", stmt.Errno, stmt.Error)
                goto finish
        }
        stmt.BindParams("name", "long text should go here...")
        stmt.Execute()
        if stmt.Errno != 0 {
                t.Error("Error", stmt.Errno, stmt.Error)
                goto finish
        }
        stmt.Close()


        sql = `
SELECT a, b, c FROM table_test LIMIT 10
`
        res, _ := db.Query(sql)
        if db.Errno != 0 {
                t.Error("Error", db.Errno, db.Error)
                goto finish
        }

        for {
                //row := res.FetchMap();
                row := res.FetchRow();
                if row == nil {
                        break
                }
                for k, v := range row {
                        switch res.Fields[k].Name {
                        case "a":
                        case "b":
                                if v != "name" { t.Error(v) }
                        case "c":
                                if v != "long text should go here..." { t.Error(v) }
                        default:
                                t.Error("unknown field: ", res.Fields[k].Name, v)
                        }
                }
        }

finish:
}

func TestDatabase(t *testing.T) {
        db := NewDatabase()
        err := db.Connect("localhost", "test", "abc", "dusell")
        if err != nil { t.Error(err); goto finish }
        defer db.Close()

        err = db.Ping()
        if err != nil { t.Error(err); goto finish }

        sql := `SELECT a, b, c FROM table_test LIMIT 10`
        res, err := db.Query(sql)
        if err != nil { t.Error(err); goto finish }

        for {
                row := res.FetchRow();
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

        stmt, _ := db.NewStatement()

        sql = `SELECT a, b, c FROM table_test WHERE a<? LIMIT 10`
        err = stmt.Prepare(sql)
        if err != nil { t.Error(err); goto finish }

        stmt.BindParams(10)
        res, err = stmt.Execute()
        if err != nil { t.Error(err); goto finish }

        stmt.Close()

        for {
                row := res.FetchRow();
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
