GOARCH ?= 386
GOROOT ?= /home/duzy/open/go
include $(GOROOT)/src/Make.inc

TARG=mysql
CGOFILES=mysql.go
GOFILES=mysql_com.go
MYSQL_CONFIG=$(shell which mysql_config)
CGO_CFLAGS=$(shell $(MYSQL_CONFIG) --cflags)
CGO_LDFLAGS=$(shell $(MYSQL_CONFIG) --libs_r)

include $(GOROOT)/src/Make.pkg

prereq:
	@test -x "$(MYSQL_CONFIG)" || \
		(echo "Can't find mysql_config in your path."; false)
