include $(GOROOT)/src/Make.inc

TARG = ds/web
GOFILES = \
  app.go \
  cgi.go \
  err.go \
  fcgi.go \
  viewmgr.go \

PREREQ += ../MySQLClient/libmysql_wrap.so

LD_LIBRARY_PATH += $(shell pwd)/../MySQLClient

export LD_LIBRARY_PATH

include $(GOROOT)/src/Make.pkg

../MySQLClient/libmysql_wrap.so: ../MySQLClient/Makefile
	cd ../MySQLClient && gomake
