package web

import (
        "testing"
        "fmt"
        "os"
)

func TestDatabase(t *testing.T) {
        db := NewDatabase()
        err := db.Connect("localhost", "dusellco_test", "abc", "dusellco_test")
        if err != nil { t.Errorf("Connect: %v", err); return }
        defer db.Close()

        _, err = db.Query(`CREATE TEMPORARY TABLE tt(a VARCHAR(100), b VARCHAR(100), c VARCHAR(100));`)
        if err != nil { t.Errorf("Query: CREATE TABLE: %s", err); return }

        for n := 0; n < 20; n += 1 {
                _, err = db.Query(fmt.Sprintf(`INSERT INTO tt(a,b,c) VALUES("%d","%s","%s");`, n, "name", "long text should go here..."));
                if err != nil { t.Errorf("Query: INSERT: %s", err); return }
        }

        res, err := db.Query(`SELECT a, b, c FROM tt LIMIT 5`)
        if err != nil { t.Errorf("Query: SELECT: %v", err); return }

        for {
                row, err := res.FetchRow();
                if row == nil { break }
                if err != nil { t.Errorf("FetchRow: %v", err); break }
                t.Logf("row: %v\n", row)
                for k, v := range row {
                        switch res.GetFieldName(k) {
                        case "a":
                        case "b":
                                if v != "name" { t.Errorf("field 'b': %v", v); return }
                        case "c":
                                if v != "long text should go here..." { t.Errorf("field 'c': %v", v); return }
                        default:
                                t.Error("unknown field: ", res.GetFieldName(k), v)
                        }
                }
        }

        stmt, err := db.Prepare(`SELECT a, b, c FROM tt WHERE a<? LIMIT 5`)
        if err != nil { t.Errorf("Prepare: %v", err); return }

        res, err = stmt.Execute(10)
        if err != nil { t.Errorf("Execute: %v", err); return }

        defer stmt.Close()

        count := 0
        for {
                row, err := res.FetchRow();
                if err == os.EOF || row == nil {
                        if count != 5 {
                                t.Errorf("FetchRow: not %d rows fetched (%d)", 5, count)
                        }
                        break
                }
                if err != nil { t.Errorf("FetchRow: [STMT] %v", err); return }
                count += 1
                t.Logf("row: [STMT] %v\n", row)
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
        return
}
