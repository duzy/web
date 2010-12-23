#!/bin/bash

. ../funs.sh

go_tests=`ls src/*_test.go`
go_files="
  src/apis.go
"

name=eBay
build_pack $name
build_testmain $name
