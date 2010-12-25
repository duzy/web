#!/bin/bash

. ../funs.sh

go_incs="-I../_obj"
go_libs="-L../_libs"

go_tests=""
go_files="
  names.go
  HomePage.go
  CPanelPage.go
"
build_pack dusell && build_testmain dusell

go_files="
  main.go
  home.go
"
build_exe home

go_files="
  main.go
  cpanel.go
"
build_exe cpanel

[[ "$?" == "100" ]] && exit 0 # tells ok
