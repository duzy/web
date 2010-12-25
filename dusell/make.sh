#!/bin/bash

. ../funs.sh

go_tests=""
go_files="
  names.go
  HomePage.go
  CPanelPage.go
"
build_pack dusell #> /dev/null && build_testmain dusell

go_files="
  main.go
"
build_exe main #> /dev/null

exit 0 # tells ok
