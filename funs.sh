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

print_testmain.go()
{
    importpath=$1

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
}

_build()
{
    local type=$1
    local name=$2

    if [[ "x$type" == "x" ]] ; then
        echo "build type unspecified!"
        return 1
    fi

    if [[ "x$name" == "x" ]] ; then
        echo "usage: build_$type name"
        return 2
    fi

    if [[ "x$go_files" == "x" ]] ; then
        echo "build_$type: variable go_files is empty, pack will not be built"
        return 3
    fi

    (prepare _obj) \
        && 8g $go_incs -o _obj/$name.8 $go_files \
        && {
            [[ "$type" == "pack" ]] && gopack grc _obj/$name.a _obj/$name.8
            [[ "$type" == "exe" ]] && {
                prepare _bin && 8l $go_libs -o _bin/$name _obj/$name.8
            }
        }

    if [[ "x$go_tests" == "x" ]] ; then
        #echo "build_$type: variable go_tests is empty, test pack will not be built."
        return 100
    fi

    (prepare _test)\
        && 8g -o _test/$name.8 $go_files $go_tests \
        && gopack grc _test/$name.a _test/$name.8
}

build_pack()
{
    _build pack $@
}

build_exe()
{
    _build exe $@
}

build_testmain()
{
    local name=$1
    local testmain=_testmain.go

    print_testmain.go $name > $testmain || return 1

    (prepare _test) \
        && 8g -I_test -o _test/main.8 $testmain \
        && 8l -L_test -o testmain _test/main.8 \

}
