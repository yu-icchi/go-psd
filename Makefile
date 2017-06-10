bench:
	go test -bench=BenchmarkParse >> psd-benchmark.txt

bench2:
	go test -bench=BenchmarkDecode >> psd-benchmark2.txt

trace:
	go test -c
	GODEBUG=allocfreetrace=1 ./go-psd.test -test.run=none -test.bench=BenchmarkParse -test.benchtime=10ms 2> trace.log

trace2:
	go test -c
	GODEBUG=allocfreetrace=1 ./go-psd.test -test.run=none -test.bench=BenchmarkDecode -test.benchtime=10ms 2> trace2.log
