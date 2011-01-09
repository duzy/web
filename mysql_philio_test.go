/**
 *  Author: Duzy Chan <code@duzy.info>
 */

package mysql

import (
        "testing"
        "strings"
        "fmt"
)

func TestSimpleCreateInsertSelect(t *testing.T) {
        db := New()
        db.Logging = true
        db.Connect("localhost", "dusellco_test", "abc", "dusellco_test")
        if db.Errno != 0 {
                t.Errorf("Connect: [%d] %v", db.Errno, db.Error)
                return
        }
        defer db.Close()

        db.Query("SET NAMES utf8")
        if db.Errno != 0 {
                t.Errorf("Query: [%d] %v", db.Errno, db.Error)
                return
        }

        sql := `
CREATE TABLE IF NOT EXISTS table_test(
a INT AUTO_INCREMENT PRIMARY KEY,
b VARCHAR(64),
c TEXT
)`
        db.Query(sql)
        if db.Errno != 0 {
                t.Errorf("Query: [%d] %v", db.Errno, db.Error)
                return
        }

        sql = `
INSERT INTO table_test(b,c) VALUES(?,?)
`
        stmt, _ := db.InitStmt()
        stmt.Prepare(sql)
        if stmt.Errno != 0 {
                t.Errorf("Prepare: [%d] %v", stmt.Errno, stmt.Error)
                return
        }
        stmt.BindParams("name", "long text should go here...")
        stmt.Execute()
        if stmt.Errno != 0 {
                t.Errorf("Execute: [%d] %v", stmt.Errno, stmt.Error)
                return
        }
        stmt.Close()


        sql = `
SELECT a, b, c FROM table_test LIMIT 10
`
        res, _ := db.Query(sql)
        if db.Errno != 0 {
                t.Errorf("Query: [%d] %v", db.Errno, db.Error)
                return
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
}

func TestCreateInsertSelectComplex(t *testing.T) {
        sql := `
CREATE TABLE IF NOT EXISTS table_test_categories(
  f0 VARCHAR(32) PRIMARY KEY,
  f1 VARCHAR(128) NOT NULL
);
CREATE TABLE IF NOT EXISTS table_test_items(
  f00 VARCHAR(32) PRIMARY KEY,
  f01 VARCHAR(128) NOT NULL,
  f02 VARCHAR(32) NOT NULL,
  f03 VARCHAR(256),
  f04 VARCHAR(512),
  f05 VARCHAR(256),
  f06 VARCHAR(32),
  f07 VARCHAR(32),
  f08 VARCHAR(512),
  f09 VARCHAR(32),
  f10 VARCHAR(32),
  f11 VARCHAR(256),
  f12 FLOAT,
  f13 CHAR(3),
  f14 VARCHAR(32),
  f15 VARCHAR(256),
  f16 SMALLINT,
  f17 BIT(1),
  f18 BIT(1),
  f19 FLOAT,
  f20 CHAR(3),
  f21 FLOAT,
  f22 CHAR(3),
  f23 INT,
  f24 VARCHAR(256),
  f25 VARCHAR(24),
  f26 VARCHAR(30),
  f27 VARCHAR(30),
  f28 VARCHAR(32),
  f29 BIT(1),
  f30 BIT(1),
  f31 BIT(1),
  f32 BIT(1),
  f33 BIT(1)
)
`
        db := New()
        db.Logging = false
        err := db.Connect("localhost", "dusellco_test", "abc", "dusellco_test")
        if db.Errno != 0 { t.Errorf("Connect: [%d] %s", db.Errno, db.Error); return }
        if err != nil { t.Error(err); return }
        defer db.Close()

        _, err = db.Query(sql)
        if db.Errno != 0 { t.Errorf("Query: [%d] %s", db.Errno, db.Error); return }
        if err != nil { t.Error(err); return }

        sql = `
INSERT INTO table_test_items(
  f00,
  f01,
  f02,
  f03,
  f04,
  f05,
  f06,
  f07,
  f08,
  f09,
  f10,
  f11,
  f12,
  f13,
  f14,
  f15,
  f16,
  f17,
  f18,
  f19,
  f20,
  f21,
  f22,
  f23,
  f24,
  f25,
  f26,
  f27,
  f28,
  f29,
  f30,
  f31,
  f32,
  f33
) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
  ON DUPLICATE KEY UPDATE 
    f01=VALUES(f01),
    f02=VALUES(f02),
    f03=VALUES(f03),
    f04=VALUES(f04),
    f05=VALUES(f05),
    f06=VALUES(f06),
    f07=VALUES(f07),
    f08=VALUES(f08),
    f09=VALUES(f09),
    f10=VALUES(f10),
    f11=VALUES(f11),
    f12=VALUES(f12),
    f13=VALUES(f13),
    f14=VALUES(f14),
    f15=VALUES(f15),
    f16=VALUES(f16),
    f17=VALUES(f17),
    f18=VALUES(f18),
    f19=VALUES(f19),
    f20=VALUES(f20),
    f21=VALUES(f21),
    f22=VALUES(f22),
    f23=VALUES(f23),
    f24=VALUES(f24),
    f25=VALUES(f25),
    f26=VALUES(f26),
    f27=VALUES(f27),
    f28=VALUES(f28),
    f29=VALUES(f29),
    f30=VALUES(f30),
    f31=VALUES(f31),
    f32=VALUES(f32),
    f33=VALUES(f33)
`

        stmt, err := db.InitStmt()
        if err != nil { t.Errorf("InitStmt: [%d] %s", stmt.Errno, stmt.Error); return }

        err = stmt.Prepare(sql)
        if stmt.Errno != 0 {
                if stmt.Errno == 2014 {
                        // [2014] Commands out of sync; you can't run this command now
                        if err = db.Reconnect(); err != nil { t.Errorf("Reconnect: [%d] %s", db.Errno, db.Error); return }
                        err = stmt.Prepare(sql)
                } else {
                        t.Logf("Prepare: [%d] %s", stmt.Errno, stmt.Error);
                        t.Fail();
                        return
                }
        }
        if stmt.Errno != 0 { t.Errorf("Prepare: [%d] %s", stmt.Errno, stmt.Error); return }
        if err != nil { t.Errorf("Prepare: %v", err); return }

        params := make([]interface{}, 34)
        for n := 0; n < 34; n += 1 {
                params[n] = fmt.Sprintf("field %d",n);
        }

        if err = stmt.BindParams(params...); err != nil {
                t.Errorf("BindParams: [%d] %s", stmt.Errno, stmt.Error);
                return
        }
        if _, err = stmt.Execute(); err != nil {
                stmt.Close()
                t.Errorf("Execute: [%d] %s", stmt.Errno, stmt.Error);
                return
        }
        if err = stmt.Close(); err != nil {
                t.Errorf("Close: [%d] %s", stmt.Errno, stmt.Error);
                return
        }

        sql = `
SELECT
  f00,
  f01,
  f02,
  f03,
  f04,
  f05,
  f06,
  f07,
  f08,
  f09,
  f10,
  f11,
  f12,
  f13,
  f14,
  f15,
  f16,
  f17,
  f18,
  f19,
  f20,
  f21,
  f22,
  f23,
  f24,
  f25,
  f26,
  f27,
  f28,
  f29,
  f30,
  f31,
  f32,
  f33
FROM table_test_items
WHERE f00=?
LIMIT 1
`
        /*
        if err = db.Reconnect(); err != nil {
                t.Errorf("Reconnect: [%d] %s", db.Errno, db.Error);
                return
        }
         */

        q := strings.Replace(sql, "?", `"field 0"`, 1)
        res, err := db.Query(q)
        if err != nil { t.Errorf("Query: [%d] %s", stmt.Errno, stmt.Error); return }
        fmt.Printf("%v\n", res)

        stmt, err = db.InitStmt()
        if err != nil { t.Errorf("InitStmt: [%d] %s", stmt.Errno, stmt.Error); return }

        err = stmt.Prepare(sql)
        if stmt.Errno != 0 { t.Errorf("Prepare: [%d] %s", stmt.Errno, stmt.Error); return }
        if err != nil { t.Errorf("Prepare: %v", err); return }

        if err = stmt.BindParams("field 0"); err != nil {
                t.Errorf("BindParams: [%d] %s", stmt.Errno, stmt.Error);
                return
        }

        //fmt.Printf("====\n");
        //db.Logging = true
        if res, err := stmt.Execute(); err != nil {
                stmt.Close()
                t.Errorf("Execute: [%d] %s", stmt.Errno, stmt.Error);
                return
        } else {
                fmt.Printf("%v\n", res)
        }
        //db.Logging = false
        //fmt.Printf("====.\n");

        if err = stmt.Close(); err != nil {
                t.Errorf("Close: [%d] %s", stmt.Errno, stmt.Error);
                return
        }
}
