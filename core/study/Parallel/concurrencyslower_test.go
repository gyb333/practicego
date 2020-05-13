package Parallel

import "testing"

//go test -bench .
//go test -bench . -cpuprofile prof.cpu
//go tool pprof [binary] prof.cpu对采样文件进行分析。



func BenchmarkSerialSum(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SerialSum()
	}
}

func BenchmarkConcurrentSum(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ConcurrentSum()
	}
}

func BenchmarkChannelSum(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ChannelSum()
	}
}