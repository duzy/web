package mysql

import (
        "testing"
        "fmt"
        "os"
)

var testStrings = []string{
	"道可道，非常道。", "名可名，非常名。",
	"無名天地之始；", "有名萬物之母。",
	"故常無欲以觀其妙；", "常有欲以觀其徼。",
	"此兩者同出而異名，", "同謂之玄。",
	"玄之又玄，眾妙之門。",
	"test",
	"test2",
	"test3",
	"test4",
	"test5",
}

func TestConnectQuery(t *testing.T) {
        conn, err := Open("localhost", "dusellco_test", "abc", "dusellco_test")
        if err != nil {
                t.Errorf("%s", err)
                return
        }

        defer conn.Close()

        rs, err := conn.Query(`CREATE TEMPORARY TABLE t(i INT, s VARCHAR(100));`)
        if err != nil {
                t.Errorf("%s", err)
                return
        }
        if rs.NumFields != 0 {
                t.Errorf("error: NumFields == %d (!=0)", rs.NumFields)
                return
        }
        if rs.AffectedRows != 0 {
                t.Errorf("error: AffectedRows == %d (!=0)", rs.AffectedRows)
                return
        }
        rs.Free()

        for i := 0; i < len(testStrings); i += 1 {
                rs, err := conn.Query(fmt.Sprintf(`INSERT INTO t(i,s) VALUES(%d, "%s");`, i, testStrings[i]))
                if err != nil {
                        t.Errorf("%s", err)
                        break;
                }
                if rs.NumFields != 0 {
                        t.Errorf("error: INSERT: [%d] NumFields == %d (!=0)", i, rs.NumFields)
                        break;
                }
                if rs.AffectedRows != 1 {
                        t.Errorf("error: INSERT: [%d] AffectedRows == %d (!=1)", i, rs.AffectedRows)
                        break;
                }
                rs.Free()
        }

        rs, err = conn.Query(`SELECT * FROM t ORDER BY i;`)
        if err != nil {
                t.Errorf("%s", err)
                return
        }
        if rs.NumFields != 2 {
                t.Errorf("error: SELECT: NumFields == %d (!=2)", rs.NumFields)
                return
        }
        if n := rs.GetNumFields(); n != 2 {
                t.Errorf("error: SELECT: GetNumFields() == %d (!=2)", n)
                return
        }
        // if rs.AffectedRows != 0 {
        //         t.Errorf("error: SELECT: AffectedRows == %d (!=0)", rs.AffectedRows)
        //         return
        // }

        fields, err := rs.FetchFields()
        if err != nil {
                t.Errorf("%s", err)
                return
        }
        if len(fields) != int(rs.NumFields) {
                t.Errorf("error: len(FetchFields())", len(fields))
                return
        }
        if fields[0].Name != "i" {
                t.Errorf("error: fields[0].Name: %s", fields[0].Name)
                return
        }
        if fields[0].Type != FIELD_TYPE_LONG {
                t.Errorf("error: fields[0].Type: %d", fields[0].Type)
                return
        }
        if fields[0].Table != "t" {
                t.Errorf("error: fields[0].Table: %s", fields[0].Table)
                return
        }
        if fields[0].OriginalTable != "t" {
                t.Errorf("error: fields[0].Table: %s", fields[0].OriginalTable)
                return
        }
        if fields[1].Name != "s" {
                t.Errorf("error: fields[1].Name: %s", fields[1].Name)
                return
        }
        if fields[1].Type != FIELD_TYPE_VAR_STRING {
                t.Errorf("error: fields[1].Type: %d, %d", fields[1].Type, FIELD_TYPE_VAR_STRING)
                return
        }

        for i := 0; ; i += 1 {
                row, err := rs.FetchRow()
                if err != nil {
                        t.Errorf("%s", err)
                        return
                }
                if row == nil {
                        if n := len(testStrings); i != n {
                                t.Errorf("wrong rows: %d != %d", i, n);
                        }
                        break;
                }
                if row[0] != fmt.Sprintf("%d", i) {
                        t.Errorf("row[%d][0]: %s != %d\n", row[0], i)
                        return
                }
                if row[1] != testStrings[i] {
                        t.Errorf("row[%d][1]: %s != %s\n", row[1], testStrings[i])
                        return
                }

                crow, err := rs.ConvertRow(row)
                if err != nil {
                        t.Errorf("%s", err)
                        return
                }
                if len(crow) != len(row) {
                        t.Errorf("ConvertRow: %v <=> %v", crow, row)
                        return
                }
                if n, ok := crow[0].(int); ok {
                        if n != i {
                                t.Errorf("ConvertRow: [%d] %v", n, crow)
                                return
                        }
                } else {
                        t.Errorf("ConvertRow: [0]: %v", crow)
                        return
                }
                if s, ok := crow[1].(string); ok {
                        if s != testStrings[i] {
                                t.Errorf("ConvertRow: [1]: %v", crow)
                                return
                        }
                } else {
                        t.Errorf("ConvertRow: %v", crow)
                        return
                }
        }
}

func TestPreparedStatement(t *testing.T) {
        conn, err := Open("localhost", "dusellco_test", "abc", "dusellco_test")
        if err != nil {
                t.Errorf("%s", err)
                return
        }

        defer conn.Close()

        rs, err := conn.Query(`CREATE TEMPORARY TABLE t(i INT, s VARCHAR(100));`)
        if err != nil {
                t.Errorf("%s", err)
                return
        }
        if rs.NumFields != 0 {
                t.Errorf("error: NumFields == %d (!=0)", rs.NumFields)
                return
        }
        if rs.AffectedRows != 0 {
                t.Errorf("error: AffectedRows == %d (!=0)", rs.AffectedRows)
                return
        }
        rs.Free()

        for i := 0; i < len(testStrings); i += 1 {
                rs, err := conn.Query(fmt.Sprintf(`INSERT INTO t(i,s) VALUES(%d, "%s");`, i, testStrings[i]))
                if err != nil {
                        t.Errorf("%s", err)
                        break;
                }
                if rs.NumFields != 0 {
                        t.Errorf("error: INSERT: [%d] NumFields == %d (!=0)", i, rs.NumFields)
                        break;
                }
                if rs.AffectedRows != 1 {
                        t.Errorf("error: INSERT: [%d] AffectedRows == %d (!=1)", i, rs.AffectedRows)
                        break;
                }
                rs.Free()
        }

        stmt, err := conn.Prepare(`SELECT * FROM t WHERE i >= ?  ORDER BY i;`)
        if err != nil {
                t.Errorf("%s", err)
                return
        }

        if stmt == nil {
                t.Errorf("conn.Prepare failed")
                return
        }

        err = stmt.BindParams(0)
        if err != nil {
                t.Errorf("%s", err)
                stmt.Close()
                return
        }

        err = stmt.Execute()
        if err != nil {
                t.Errorf("%s", err)
                stmt.Close()
                return
        }

        if n := stmt.NumRows(); n != uint64(len(testStrings)) {
                t.Errorf("NumRows: %d", n)
                stmt.Close()
                return
        }

        if n := stmt.AffectedRows(); n != uint64(len(testStrings)) /*0*/ {
                t.Errorf("AffectedRows: %d", n)
                stmt.Close()
                return
        }

        n := 0;
        for {
                row, err := stmt.Fetch()
                if err != nil {
                        if err != os.EOF {
                                t.Errorf("%v", err)
                        }
                        break
                }
                //fmt.Printf("row: %v\n", row)
                if row == nil {
                        t.Errorf("Fetch: row = <nil>")
                        break
                }
                if len(row) != 2 {
                        t.Errorf("Fetch: row %d, wrong columns: %v", n, row)
                }
                if row[0] != n {
                        t.Errorf("Fetch: row %d: [0] %v", n, row)
                }
                if row[1] != testStrings[n] {
                        t.Errorf("Fetch: row %d: [1] %v", n, row)
                }
                n += 1
        }

        if n != len(testStrings) {
                t.Errorf("Fetch: wrong rows: %v != %v", n, len(testStrings))
        }

        stmt.Close()
}
