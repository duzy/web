package web

import (
        "./_obj/mysql"
        //"fmt"
        "os"
)

type mysqlDatabase struct {
        *mysql.MySQLInstance
}

type mysqlQueryResult struct {
        *mysql.MySQLResponse
}

type mysqlStatement struct {
        *mysql.MySQLStatement
}

func NewDatabase() (db Database) {
        p := &mysqlDatabase{ nil }
        db = Database(p)
        return
}

func (db *mysqlDatabase) Connect(params ...interface{}) (err os.Error) {
        if len(params) != 4/*6*/ {
                err = os.NewError("wrong connection parameters");
                return
        }
        var ok bool
        var netstr, laddrstr, raddrstr, username, password, database string
        netstr = "unix"
        raddrstr = "/var/run/mysqld/mysqld.sock"
        laddrstr, ok = params[0].(string); if !ok { err = os.NewError("not string parameter"); return }
        username, ok = params[1].(string); if !ok { err = os.NewError("not string parameter"); return }
        password, ok = params[2].(string); if !ok { err = os.NewError("not string parameter"); return }
        database, ok = params[3].(string); if !ok { err = os.NewError("not string parameter"); return }
        db.MySQLInstance, err = mysql.Connect(netstr, laddrstr, raddrstr, username, password, database)
        return
}

func (db *mysqlDatabase) Close() (err os.Error) {
        db.MySQLInstance.Quit()
        return
}        

func (db *mysqlDatabase) Switch(s string) (err os.Error) {
        _, err = db.MySQLInstance.Use(s)
        return
}

func (db *mysqlDatabase) Query(sql string) (res QueryResult, err os.Error) {
        r, err := db.MySQLInstance.Query(sql)
        if err == nil {
                res = QueryResult(&mysqlQueryResult{r})
        }
        return
}

func (db *mysqlDatabase) Prepare(sql string) (stmt SQLStatement, err os.Error) {
        r, err := db.MySQLInstance.Prepare(sql)
        if err == nil {
                stmt = SQLStatement(&mysqlStatement{r})
        }
        return
}

func (qr *mysqlQueryResult) GetRowCount() uint64 { return qr.MySQLResult.RowCount }
func (qr *mysqlQueryResult) GetFieldCount() uint64 { return qr.MySQLResult.FieldCount }
func (qr *mysqlQueryResult) GetFieldName(n int) string { return qr.MySQLResult.Fields[n].Name }
func (qr *mysqlQueryResult) GetAffectedRows() uint64 { return qr.MySQLResult.AffectedRows }
func (qr *mysqlQueryResult) GetInsertId() uint64 { return qr.MySQLResult.InsertId }
func (qr *mysqlQueryResult) MoveFirst() { qr.MySQLResult.Reset() }

func (stmt *mysqlStatement) Prepare(sql string) (err os.Error) {
        err = stmt.MySQLStatement.Prepare(sql)
        if err != nil { err = formatMySQLError(stmt) }
        return
}

func (stmt *mysqlStatement) BindParams(params ...interface{}) (err os.Error) {
        err = stmt.MySQLStatement.BindParams(params...)
        if err != nil { err = formatMySQLError(stmt) }
        return
}

func (stmt *mysqlStatement) Execute() (res QueryResult, err os.Error) {
        r, err := stmt.MySQLStatement.Execute()
        if err != nil {
                err = formatMySQLError(stmt)
        } else {
                res = QueryResult(&mysqlQueryResult{r})
        }
        return
}

func (stmt *mysqlStatement) Reset() (err os.Error) {
        err = stmt.MySQLStatement.Reset()
        if err != nil { err = formatMySQLError(stmt) }
        return
}

func (stmt *mysqlStatement) Close() (err os.Error) {
        err = stmt.MySQLStatement.Close()
        if err != nil { err = formatMySQLError(stmt) }
        return
}
