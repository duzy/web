include $(GOROOT)/src/Make.inc

DB_BACKEND := cbind

TARG = ds/web
CGOFILES =
GOFILES = \
  app.go \
  cgi.go \
  err.go \
  fcgi.go \
  viewmgr.go \

CGO_CFLAGS = 
CGO_LDFLAGS = -L. -lmysql_wrap

PREREQ += ../../MySQLClient/libmysql_wrap.so

LD_LIBRARY_PATH += $(shell pwd)/../MySQLClient

export LD_LIBRARY_PATH

include $(GOROOT)/src/Make.pkg

../../MySQLClient/libmysql_wrap.so: ../../MySQLClient/Makefile
	cd ../../MySQLClient && gomake
