#!/bin/bash

. ../funs.sh

go_tests=""
go_files="
  names.go
  HomePage.go
  CPanelPage.go
"
build_pack dusell && build_testmain dusell

go_files="
  main.go
"
build_exe main
[[ "$?" == "100" ]] && exit 0 # tells ok
