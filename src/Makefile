include $(GOROOT)/src/Make.inc

DB_BACKEND := cbind

TARG = web
CGOFILES =
GOFILES = \
  app.go \
  appcfg.go \
  cgi.go \
  db.go \
  db_mysql_$(DB_BACKEND).go \
  dbmgr.go \
  err.go \
  fcgi.go \
  sm.go \
  viewmgr.go \

CGO_CFLAGS = 
CGO_LDFLAGS = -L. -lmysql_wrap

PREREQ += ../../MySQLClient/libmysql_wrap.so

LD_LIBRARY_PATH += $(shell pwd)/../MySQLClient

export LD_LIBRARY_PATH

include $(GOROOT)/src/Make.pkg

../../MySQLClient/libmysql_wrap.so: ../../MySQLClient/Makefile
	cd ../../MySQLClient && gomake
