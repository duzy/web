#!/bin/bash

[[ -d MySQL ]] || {
    git clone http://github.com/Philio/GoMySQL.git MySQL ||
    {
        echo "Stoped building MySQL."
        exit -1
    }
}

. funs.sh

go_files="
  MySQL/mysql.go
  MySQL/mysql_const.go
  MySQL/mysql_error.go
  MySQL/mysql_result.go
  MySQL/mysql_packet.go
  MySQL/mysql_statement.go
"

build_pack mysql && build_testmain mysql
[[ "x$?" == "x100" ]] && exit 0 # tells ok
