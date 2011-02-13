#!/bin/bash

which=philio
#which=thoj

[[ -d MySQL ]] || {
    #local u=""
    u=""
    if [[ "x$which" == "xthoj" ]] ; then
        u="https://github.com/thoj/Go-MySQL-Client-Library.git"
    else
        u="http://github.com/Philio/GoMySQL.git"
    fi
    git clone $u MySQL ||
    {
        echo "Stoped building MySQL."
        exit -1
    }
}

. funs.sh

if [[ "x$which" == "xthoj" ]] ; then
go_tests="
  mysql_thoj_test.go
"
go_files="
  MySQL/mysql.go
  MySQL/mysql_const.go
  MySQL/mysql_data.go
  MySQL/mysql_stmt.go
  MySQL/mysql_util.go
"
else
go_tests="
  mysql_philio_test.go
"
go_files="
  MySQL/mysql.go
  MySQL/mysql_const.go
  MySQL/mysql_error.go
  MySQL/mysql_result.go
  MySQL/mysql_packet.go
  MySQL/mysql_statement.go
"
fi

build_pack mysql && build_testmain mysql
[[ "x$?" == "x100" ]] && exit 0 # tells ok

exit 0
