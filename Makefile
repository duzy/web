# TODO: use gomake for building

LD_LIBRARY_PATH += $(shell pwd)/../MySQLClient

export LD_LIBRARY_PATH

all: remake test

remake: make.sh ; @./$<

test: _bin/test ; @./$<

_bin/test: make.sh ; @./$<
