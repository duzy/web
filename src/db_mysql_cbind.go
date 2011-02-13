package web

import (
        "../MySQLClient/_obj/mysql"
        "fmt"
        "os"
)

/**
 *  Implements Database
 */
type cbindDatabase struct {
        conn *mysql.Connection
}

/**
 *  Implements SQLStatement and QueryResult
 */
type cbindStatement struct {
        stmt *mysql.Statement
}

/**
 *  Implements QueryResult
 */
type cbindQueryResult struct {
        rs *mysql.ResultSet
        fields []mysql.Field
}

func NewDatabase() (db Database) {
        db = Database(&cbindDatabase{})
        return
}

func (db *cbindDatabase) Connect(params ...interface{}) (err os.Error) {
        if db.conn != nil {
                err = os.NewError("Alread connected!")
                return
        }

        if len(params) < 4 {
                err = os.NewError("Wrong number of parameters.")
                return
        }

        host := fmt.Sprintf("%v", params[0])
        user := fmt.Sprintf("%v", params[1])
        pass := fmt.Sprintf("%v", params[2])
        dbnm := fmt.Sprintf("%v", params[3])

        db.conn, err = mysql.Open(host, user, pass, dbnm)
        if err != nil {
                return
        }

        if db.conn == nil {
                err = os.NewError("No database connection.")
                return
        }
        return
}

func (db *cbindDatabase) Close() (err os.Error) {
        if db.conn != nil {
                err = db.conn.Close()
        }
        return
}

func (db *cbindDatabase) Reconnect() (err os.Error) {
        err = os.NewError("unimplemented Reconnect()")
        return
}

func (db *cbindDatabase) Query(sql string) (res QueryResult, err os.Error) {
        if db.conn == nil {
                err = os.NewError("no database connection")
                return
        }

        rs, err := db.conn.Query(sql)
        res = &cbindQueryResult{ rs, nil }
        return
}

func (db *cbindDatabase) Switch(dbnm string) (err os.Error) {
        err = os.NewError("unimplemented Switch(db)")
        return
}

func (db *cbindDatabase) Prepare(sql string) (stmt SQLStatement, err os.Error) {
        if db.conn == nil {
                err = os.NewError("no database connection")
                return
        }

        s, err := db.conn.Prepare(sql)
        if err != nil {
                return
        }

        stmt = &cbindStatement{ s }
        return
}

func (db *cbindDatabase) Escape(s string) string {
        if db.conn == nil {
                //err = os.NewError("no database connection")
                return ""
        }
        //err = os.NewError("unimplemented Switch(db)")
        // TODO:...
        return ""
}


func (s *cbindStatement) Execute(args ...interface{}) (res QueryResult, err os.Error) {
        if s.stmt == nil {
                err = os.NewError("statement not inited")
                return
        }

        err = s.stmt.BindParams(args...)
        if err != nil {
                return
        }

        err = s.stmt.Execute()
        res = QueryResult(s)
        return
}

func (s *cbindStatement) Close() (err os.Error) {
        if s.stmt == nil {
                err = os.NewError("no inited statement")
                return
        }

        err = s.stmt.Close()
        return
}

func (r *cbindStatement) GetFieldCount() uint {
        if r.stmt != nil {
                return r.stmt.FieldCount()
        }
        return 0
}

func (r *cbindStatement) GetFieldName(n int) string {
        if r.stmt != nil && r.stmt.Fields != nil {
                if n < len(r.stmt.Fields) {
                        return r.stmt.Fields[n].Name
                }
        }
        return ""
}

func (r *cbindStatement) GetAffectedRows() uint64 {
        if r.stmt != nil {
                return r.stmt.AffectedRows()
        }
        return 0
}

func (r *cbindStatement) GetInsertId() uint64 {
        if r.stmt != nil {
                return r.stmt.InsertId()
        }
        return 0
}

func (r *cbindStatement) GetRowCount() uint64 {
        if r.stmt != nil {
                return r.stmt.NumRows()
        }
        return 0
}

func (r *cbindStatement) FetchRow() (row []interface{}, err os.Error) {
        if r.stmt != nil {
                row, err = r.stmt.Fetch()
        }
        return
}

func (r *cbindStatement) Free() (err os.Error) { return }

func (r *cbindQueryResult) GetFieldCount() uint {
        if r.rs != nil {
                return r.rs.GetNumFields() // or: r.rs.NumFields
        }
        return 0
}

func (r *cbindQueryResult) GetFieldName(n int) string {
        if r.rs != nil {
                if r.fields == nil {
                        var err os.Error
                        r.fields, err = r.rs.FetchFields()
                        if err != nil {
                                // TODO: error handling
                                return ""
                        }
                }
                return r.fields[n].Name
        }
        return ""
}

func (r *cbindQueryResult) GetAffectedRows() uint64 {
        if r.rs != nil {
                return r.rs.AffectedRows
        }
        return 0
}

func (r *cbindQueryResult) GetInsertId() uint64 {
        if r.rs != nil {
                return r.rs.InsertId
        }
        return 0
}

func (r *cbindQueryResult) GetRowCount() uint64 {
        if r.rs != nil {
                return r.rs.GetNumRows()
        }
        return 0
}

func (r *cbindQueryResult) FetchRow() (row []interface{}, err os.Error) {
        if r.rs == nil {
                err = os.NewError("no result")
                return
        }

        var sa []string
        sa, err = r.rs.FetchRow()
        if err != nil { return }

        row, err = r.rs.ConvertRow(sa)
        return
}

func (r *cbindQueryResult) MoveFirst() {
        if r.rs == nil {
                //err = os.NewError("no result")
                return
        }
        //r.rs.RowSeek()
        // TODO: ...
        return
}

func (r *cbindQueryResult) Free() (err os.Error) {
        if r.rs != nil {
                return r.rs.Free()
        }
        return
}
