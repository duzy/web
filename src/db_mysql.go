package web

import (
        "./_obj/mysql"
        "fmt"
        "os"
)

type mysqlDatabase struct {
        *mysql.MySQL
}

type mysqlQueryResult struct {
        *mysql.MySQLResult
}

type mysqlStatement struct {
        *mysql.MySQLStatement
}

func NewDatabase() (db Database) {
        p := &mysqlDatabase{ mysql.New() }
        db = Database(p)
        return
}

func formatMySQLError(i interface{}) (err os.Error) {
        switch o := i.(type) {
        case *mysqlDatabase:
                if o.Errno != 0 {
                        err = os.NewError(fmt.Sprintf("[DB:%v] %v",o.Errno,o.Error))
                }
        case *mysqlStatement:
                if o.Errno != 0 {
                        err = os.NewError(fmt.Sprintf("[STMT:%v] %v",o.Errno,o.Error))
                }
        }
        return
}

func (db *mysqlDatabase) Ping() (err os.Error) {
        err = db.MySQL.Ping()
        if err != nil { err = formatMySQLError(db) }
        return
}

func (db *mysqlDatabase) Connect(params ...interface{}) (err os.Error) {
        err = db.MySQL.Connect(params...)
        if err != nil { err = formatMySQLError(db) }
        return
}

func (db *mysqlDatabase) Close() (err os.Error) {
        err = db.MySQL.Close()
        if err != nil { err = formatMySQLError(db) }
        return
}        

func (db *mysqlDatabase) Query(sql string) (res QueryResult, err os.Error) {
        r, err := db.MySQL.Query(sql)
        if err == nil {
                res = QueryResult(&mysqlQueryResult{r})
        } else {
                err = formatMySQLError(db)
        }
        return
}

func (db *mysqlDatabase) MultiQuery(sql string) (res []QueryResult, err os.Error) {
        r, err := db.MySQL.MultiQuery(sql)
        if err == nil {
                res = make([]QueryResult, len(r))
                for i, a := range r {
                        res[i] = QueryResult(&mysqlQueryResult{a})
                }
        } else {
                err = formatMySQLError(db)
        }
        return
}

func (db *mysqlDatabase) NewStatement() (stmt SQLStatement, err os.Error) {
        r, err := db.MySQL.InitStmt()
        if err == nil {
                stmt = SQLStatement(&mysqlStatement{r})
        } else {
                err = formatMySQLError(db)
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
