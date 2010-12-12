package main

import (
        //"os"
        //"io"
        "fmt"
        //"flag"
        //"bytes"
        "./_obj/web"
        "./_obj/mysql"
)

func test_mysql() {
        db := mysql.New()
        db.Logging = false
        db.Connect("localhost", "test", "abc", "dusell")
        if db.Errno != 0 {
                fmt.Printf("Error #d %s\n", db.Errno, db.Error)
                goto finish
        }
        defer db.Close()

        db.Query("SET NAMES utf8")
        if db.Errno != 0 {
                fmt.Printf("Error #d %s\n", db.Errno, db.Error)
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
                fmt.Printf("Error #d %s\n", db.Errno, db.Error)
                goto finish
        }

        sql = `
INSERT INTO table_test(b,c) VALUES(?,?)
`
        stmt, _ := db.InitStmt()
        stmt.Prepare(sql)
        if stmt.Errno != 0 {
                fmt.Printf("Error #d %s\n", stmt.Errno, stmt.Error)
                goto finish
        }
        stmt.BindParams("name", "long text should go here...")
        stmt.Execute()
        if stmt.Errno != 0 {
                fmt.Printf("Error #d %s\n", stmt.Errno, stmt.Error)
                goto finish
        }
        stmt.Close()


        sql = `
SELECT a, b, c FROM table_test LIMIT 10
`
        res, _ := db.Query(sql)
        if db.Errno != 0 {
                fmt.Printf("Error #d %s\n", db.Errno, db.Error)
                goto finish
        }

        for {
                //row := res.FetchMap();
                row := res.FetchRow();
                if row == nil {
                        break
                }
                for k, v := range row {
                        fmt.Printf("%s:%v\n", res.Fields[k].Name, v)
                }
                fmt.Printf("\n")
        }

finish:
}

func main() {
        test_mysql()

        fmt.Printf("==================================\n")

        db := web.NewDatabase()
        db.Connect("localhost", "test", "abc", "dusell")
        defer db.Close()

        err := db.Ping()
        if err != nil {
                fmt.Printf("ping-error: %s\n", err)
                goto finish
        }

        sql := `SELECT a, b, c FROM table_test LIMIT 10`
        res, err := db.Query(sql)
        if err != nil {
                fmt.Printf("error: %s\n", err)
                goto finish
        }

        for {
                row := res.FetchRow();
                if row == nil { break }
                for k, v := range row {
                        fmt.Printf("%s:%v\n", res.GetFieldName(k), v)
                }
                fmt.Printf("\n")
        }

        stmt, _ := db.NewStatement()

        sql = `SELECT a, b, c FROM table_test WHERE a<? LIMIT 10`
        err = stmt.Prepare(sql)
        if err != nil {
                fmt.Printf("error: %s\n", err)
                goto finish
        }
        stmt.BindParams(10)
        res, err = stmt.Execute()
        if err != nil {
                fmt.Printf("error: %s\n", err)
                goto finish
        }
        stmt.Close()

        fmt.Printf("==================================\n")
        for {
                row := res.FetchRow();
                if row == nil { break }
                for k, v := range row {
                        fmt.Printf("%s:%v\n", res.GetFieldName(k), v)
                }
                fmt.Printf("\n")
        }
finish:
}
