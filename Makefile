bench:
	go test -bench=BenchmarkParse >> psd.txt

trace:
	go test -c
	GODEBUG=allocfreetrace=1 ./go-psd.test -test.run=none -test.bench=BenchmarkParse -test.benchtime=10ms 2> trace.log
