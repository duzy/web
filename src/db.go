package web

import (
        "os"
)

type Database interface {
        Ping() (err os.Error)
        Connect(params ...interface{}) (err os.Error)
        Close() (err os.Error)
        Query(sql string) (res QueryResult, err os.Error)
        MultiQuery(sql string) (res []QueryResult, err os.Error)
        NewStatement() (stmt SQLStatement, err os.Error)
        Reconnect() (err os.Error)
        ChangeDatabase(db string) (err os.Error)
}

type QueryResult interface {
        GetFieldCount() uint64
        GetFieldName(n int) string
        GetRowCount() uint64
        GetAffectedRows() uint64
        GetInsertId() uint64
        
        //GetMessage() string

        // Fetch the current row (as an array) and move next
        FetchRow() []interface{}
        MoveFirst()

        // TODO: wrappers for MySQLField, MySQLRow, etc.
}

type SQLStatement interface {
        Prepare(sql string) (err os.Error)
        BindParams(params ...interface{}) (err os.Error)
        Execute() (res QueryResult, err os.Error)
        Reset() (err os.Error)
        Close() (err os.Error)
}

