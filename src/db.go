package web

import (
        "os"
)

type Database interface {
        Connect(params ...interface{}) (err os.Error)
        Close() (err os.Error)
        Query(sql string) (res QueryResult, err os.Error)
        Switch(db string) (err os.Error)
        Prepare(sql string) (stmt SQLStatement, err os.Error)
}

type SQLStatement interface {
        Execute(args ...interface{}) (res QueryResult, err os.Error)
        Close() (err os.Error)
}

type QueryResult interface {
        GetFieldCount() uint64
        GetFieldName(n int) string
        GetAffectedRows() uint64
        GetInsertId() uint64
        GetRowCount() uint64

        // Fetch the current row (as an array) and move next
        FetchRow() []interface{}
        MoveFirst()
}

