#!/bin/bash

. funs.sh

./make-mysql.sh > /dev/null || {
    echo "..."
    exit 1
}

s="cbind"
if [[ "xx" == "x$sx" ]] ; then
    w=`git --work-tree=MySQL --git-dir=MySQL/.git config remote.origin.url`
    if [[ "x$w" == "xhttps://github.com/thoj/Go-MySQL-Client-Library.git" ]] ; then
        s=thoj
    else
        s=philio
    fi
fi

        s=philio

go_tests=`ls src/*_test.go`
#go_tests="src/db_test.go"
go_files="
  src/app.go
  src/appcfg.go
  src/cgi.go
  src/db.go
  src/db_mysql_$s.go
  src/dbmgr.go
  src/err.go
  src/fcgi.go
  src/sm.go
  src/viewmgr.go
"

( build_pack web && build_testmain web ) || ( echo "rc: web: $?" && exit -1 )

## test app for FCGIModel
go_tests=""
go_files="fcgi_test.go"

( build_exe test.fcgi ) || ( echo "rc: test.fcgi: $?" && exit -1 )

go_files="fcgi_page.go"

( build_exe page.fcgi ) || ( echo "rc: page.fcgi: $?" && exit -1 )

exit 0
