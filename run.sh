# usage: nmgrep pattern file...
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

#importpath=$(gomake -s importpath)
importpath=web
{
	# test functions are named TestFoo
	# the grep -v eliminates methods and other special names
	# that have multiple dots.
	pattern='Test([^a-z].*)?'
	tests=$(nmgrep $pattern _test/$importpath.a $xofile)
	if [ "x$tests" = x ]; then
		echo 'gotest: error: no tests matching '$pattern in _test/$importpath.a $xofile 1>&2
		exit 2
	fi
	# benchmarks are named BenchmarkFoo.
	pattern='Benchmark([^a-z].*)?'
	benchmarks=$(nmgrep $pattern _test/$importpath.a $xofile)

	# package spec
	echo 'package main'
	echo
	# imports
	if echo "$tests" | egrep -v '_test\.' >/dev/null; then
		if [ "$importpath" != "testing" ]; then
			echo 'import "'$importpath'"'
		fi
	fi
	if $havex; then
		echo 'import "./_xtest_"'
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
}>_testmain.go

([[ -d _obj ]] || mkdir -p _obj) \
    && 8g -o _obj/web.8 src/app.go src/view.go src/cgi.go src/sm.go src/db.go src/db_mysql.go src/appcfg.go \
    && gopack grc _obj/web.a _obj/web.8 \
    && 8g -o _obj/test.8 test.go \
    && 8g -o _obj/test_db.8 test_db.go \
    && 8g -o _obj/test_appcfg.8 test_appcfg.go \
    && 8l -o test _obj/test.8 \
    && 8l -o test_db _obj/test_db.8 \
    && 8l -o test_appcfg _obj/test_appcfg.8 \
    && ./test_appcfg $@

#    && gopack grc _obj/web.a _obj/cgi.8 \
