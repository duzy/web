#!/bin/bash

. funs.sh

./make-mysql.sh > /dev/null || {
    echo "..."
    exit 1
}

go_tests=`ls src/*_test.go`
go_files="
  src/app.go
  src/view.go
  src/cgi.go
  src/sm.go
  src/dbmgr.go
  src/db.go
  src/db_mysql.go
  src/appcfg.go
  src/err.go
"

build_pack web && build_testmain web
