package web

import (
        "os"
        "./_obj/mysql"
)

type mysqlDatabase struct {
        *mysql.MySQL
}

func NewDatabase() (db Database) {
        p := &mysqlDatabase{ mysql.New() }
        db = Database(p)
        return
}

func (db *mysqlDatabase) Query(sql string) (res QueryResult, err os.Error) {
        r, err := db.MySQL.Query(sql)
        res = QueryResult(&mysqlQueryResult{r})
        return
}

func (db *mysqlDatabase) MultiQuery(sql string) (res []QueryResult, err os.Error) {
        r, err := db.MySQL.MultiQuery(sql)
        if err == nil {
                res = make([]QueryResult, len(r))
                for i, a := range r {
                        res[i] = QueryResult(&mysqlQueryResult{a})
                }
        }
        return
}

func (db *mysqlDatabase) NewStatement() (stmt SQLStatement, err os.Error) {
        r, err := db.MySQL.InitStmt()
        stmt = SQLStatement(&mysqlStatement{r})
        return
}

type mysqlQueryResult struct {
        *mysql.MySQLResult
}

func (qr *mysqlQueryResult) GetRowCount() uint64 { return qr.MySQLResult.RowCount }
func (qr *mysqlQueryResult) GetFieldCount() uint64 { return qr.MySQLResult.FieldCount }
func (qr *mysqlQueryResult) GetFieldName(n int) string { return qr.MySQLResult.Fields[n].Name }
func (qr *mysqlQueryResult) GetAffectedRows() uint64 { return qr.MySQLResult.AffectedRows }
func (qr *mysqlQueryResult) GetInsertId() uint64 { return qr.MySQLResult.InsertId }
func (qr *mysqlQueryResult) MoveFirst() { qr.MySQLResult.Reset() }

type mysqlStatement struct {
        *mysql.MySQLStatement
}

func (stmt *mysqlStatement) Execute() (res QueryResult, err os.Error) {
        r, err := stmt.MySQLStatement.Execute()
        res = QueryResult(&mysqlQueryResult{r})
        return
}
