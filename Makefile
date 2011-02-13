# TODO: use gomake for building

LD_LIBRARY_PATH += $(shell pwd)/../MySQLClient

export LD_LIBRARY_PATH

test: _bin/test ; @./$<

_bin/test: make.sh ; @./$<
