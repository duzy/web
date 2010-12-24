#!/bin/bash

. ../funs.sh

go_tests=`ls src/*_test.go`
go_files="
  src/urls.go
  src/apis.go
  src/types.go
  src/findsvc.go
  src/shopping.go
"

name=eBay
build_pack $name
build_testmain $name
