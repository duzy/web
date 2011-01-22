#!/bin/bash

. funs.sh

./make-mysql.sh > /dev/null || {
    echo "..."
    exit 1
}

w=`git --work-tree=MySQL --git-dir=MySQL/.git config remote.origin.url`
if [[ "x$w" == "xhttps://github.com/thoj/Go-MySQL-Client-Library.git" ]] ; then
    s=thoj
else
    s=philio
fi

#go_tests=`ls src/*_test.go`
go_tests="src/db_test.go"
go_files="
  src/app.go
  src/cgi.go
  src/fcgi.go
  src/sm.go
  src/dbmgr.go
  src/db.go
  src/db_mysql_$s.go
  src/appcfg.go
  src/err.go
"

#  src/view.go

build_pack web && build_testmain web

## test app for FCGIModel
go_tests=""
go_files="fcgi_test.go"

build_exe test.fcgi

exit 0
