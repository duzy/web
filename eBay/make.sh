#!/bin/bash

if [[ -f ../make.sh ]] && [[ -d ../eBay ]]; then
    [[ -d ../_obj  ]] && {
        [[ -d _obj ]] && rm -rf _obj
        ln -s ../_obj .
    }
    [[ -d ../_test ]] && {
        [[ -d _test ]] && rm -rf _test
        ln -s ../_test .
    }
fi

. ../funs.sh

go_tests=`ls src/*_test.go`
go_files="
  src/urls.go
  src/app.go
  src/types.go
  src/findsvc.go
  src/shopping.go
"

name=eBay
build_pack $name && build_testmain $name
