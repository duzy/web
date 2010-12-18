#!/bin/bash

prepare() {
    [[ "x$1" == "x" ]] && (echo "'prepare' requires a argument!" && exit -1)
    [[ -d $1 ]] || mkdir -p $1
}

# usage: nmgrep pattern file...
# NOTE: copied from 'gotest'
nmgrep() {
	pat="$1"
	shift
	for i
	do
		# Package symbol "".Foo is pkg.Foo when imported in Go.
		# Figure out pkg.
		case "$i" in
		*.a)
			pkg=$(gopack p $i __.PKGDEF | sed -n 's/^package //p' | sed 's/ .*//' | sed 1q)
			;;
		*)
			pkg=$(sed -n 's/^ .* in package "\(.*\)".*/\1/p' $i | sed 1q)
			;;
		esac
		6nm -s "$i" | egrep ' T .*\.'"$pat"'$' |
                sed 's/.* //; /\..*\./d; s/""\./'"$pkg"'./g'
	done
}

go_files="
  src/app.go
  src/view.go
  src/cgi.go
  src/sm.go
  src/dbmgr.go
  src/db.go
  src/db_mysql.go
  src/appcfg.go
"
go_tests="
`ls src/*_test.go`
"

(prepare _obj && prepare _test)\
    && 8g -o _obj/web.8 $go_files \
    && 8g -o _test/web.8 $go_files $go_tests \
    && gopack grc _obj/web.a _obj/web.8 \
    && gopack grc _test/web.a _test/web.8 \

testmain=_testmain.go
importpath=web
{
	# test functions are named TestFoo
	# the grep -v eliminates methods and other special names
	# that have multiple dots.
	#pattern='Test([^a-z].*)?'
	pattern='Test([A-Z0-9_].*)?'
	tests=$(nmgrep $pattern _test/$importpath.a)
	if [ "x$tests" = x ]; then
		echo 'gotest: error: no tests matching '$pattern in _test/$importpath.a 1>&2
		exit 2
	fi
	# benchmarks are named BenchmarkFoo.
	#pattern='Benchmark([^a-z].*)?'
	pattern='Benchmark([A-Z0-9_].*)?'
	benchmarks=$(nmgrep $pattern _test/$importpath.a)

	# package spec
	echo 'package main'
	echo
	# imports
	if echo "$tests" | egrep -v '_test\.' >/dev/null; then
		if [ "$importpath" != "testing" ]; then
			echo 'import "'$importpath'"'
		fi
	fi
	echo 'import "testing"'
	echo 'import __regexp__ "regexp"' # rename in case tested package is called regexp
	# test array
	echo
	echo 'var tests = []testing.InternalTest{'
	for i in $tests
	do
		echo '	{"'$i'", '$i'},'
	done
	echo '}'
	# benchmark array
	if [ "$benchmarks" = "" ]
	then
		# keep the empty array gofmt-safe.
		# (not an issue for the test array, which is never empty.)
		echo 'var benchmarks = []testing.InternalBenchmark{}'
	else
		echo 'var benchmarks = []testing.InternalBenchmark{'
		for i in $benchmarks
		do
			echo '	{"'$i'", '$i'},'
		done
		echo '}'
	fi
	# body
	echo
	echo 'func main() {'
	echo '	testing.Main(__regexp__.MatchString, tests)'
	echo '	testing.RunBenchmarks(__regexp__.MatchString, benchmarks)'
	echo '}'
}>$testmain

(prepare _test)\
    && 8g -I_test -o _test/main.8 $testmain \
    && 8l -L_test -o testmain _test/main.8 \
