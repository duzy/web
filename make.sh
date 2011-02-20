#!/bin/bash

. funs.sh

go_tests=`ls src/*_test.go`
#go_tests="src/db_test.go"
go_files="
  app.go
  cgi.go
  err.go
  fcgi.go
  viewmgr.go
"

( build_pack web && build_testmain web ) || ( echo "rc: web: $?" && exit -1 )

## test app for FCGIModel
go_tests=""
go_files="examples/fcgi_simple/fcgi_test.go"

( build_exe test.fcgi ) || ( echo "rc: test.fcgi: $?" && exit -1 )

go_files="examples/fcgi_page/fcgi_page.go"

( build_exe page.fcgi ) || ( echo "rc: page.fcgi: $?" && exit -1 )

exit 0
