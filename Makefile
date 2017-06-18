bench:
	go test -bench=BenchmarkDecode >> psd-benchmark.txt

trace:
	go test -c
	GODEBUG=allocfreetrace=1 ./go-psd.test -test.run=none -test.bench=BenchmarkDecode -test.benchtime=10ms 2> trace.log

