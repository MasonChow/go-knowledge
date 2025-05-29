package main

import (
	"testing"
)

func genSrc(n int) []int {
	ds := make([]int, n)
	for i := range ds {
		ds[i] = i
	}
	return ds
}

// 预设容量
func presetCap(src []int) []int {
	res := make([]int, 0, len(src))
	for _, v := range src {
		if v%2 == 0 {
			res = append(res, v)
		}
	}
	return res
}

// 默认容量
func noCap(src []int) []int {
	res := make([]int, 0)
	for _, v := range src {
		if v%2 == 0 {
			res = append(res, v)
		}
	}
	return res
}

func BenchmarkPresetCap100(b *testing.B) {
	src := genSrc(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = presetCap(src)
	}
}

func BenchmarkNoCap100(b *testing.B) {
	src := genSrc(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = noCap(src)
	}
}

func BenchmarkPresetCap10000(b *testing.B) {
	src := genSrc(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = presetCap(src)
	}
}

func BenchmarkNoCap10000(b *testing.B) {
	src := genSrc(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = noCap(src)
	}
}
