package main

import (
	"runtime"
	"testing"
)

type Data struct {
	A int
	B int
	C int
	D int
}

func makeContinuous(n int) []Data {
	return make([]Data, n)
}

func makeHeapPointers(n int) []*Data {
	s := make([]*Data, n)
	for i := 0; i < n; i++ {
		s[i] = &Data{A: i}
	}
	return s
}

func BenchmarkContinuous(b *testing.B) {
	var m0, m1 runtime.MemStats
	for i := 0; i < b.N; i++ {
		runtime.ReadMemStats(&m0)
		_ = makeContinuous(10_000_000)
		runtime.GC()
		runtime.ReadMemStats(&m1)
	}
	b.ReportMetric(float64(m1.NumGC-m0.NumGC), "numGC")
	b.ReportMetric(float64(m1.PauseTotalNs-m0.PauseTotalNs)/1e6, "pauseMs")
}

func BenchmarkHeapPointers(b *testing.B) {
	var m0, m1 runtime.MemStats
	for i := 0; i < b.N; i++ {
		runtime.ReadMemStats(&m0)
		_ = makeHeapPointers(10_000_000)
		runtime.GC()
		runtime.ReadMemStats(&m1)
	}
	b.ReportMetric(float64(m1.NumGC-m0.NumGC), "numGC")
	b.ReportMetric(float64(m1.PauseTotalNs-m0.PauseTotalNs)/1e6, "pauseMs")
}
