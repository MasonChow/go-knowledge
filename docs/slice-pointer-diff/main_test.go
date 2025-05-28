package main

import (
	"encoding/json"
	"math/rand"
	"testing"
)

type Data struct {
	A int64
	B int64
	C int64
	D int64
}

func genData(n int) []Data {
	ds := make([]Data, n)
	for i := range ds {
		ds[i] = Data{
			A: rand.Int63(),
			B: rand.Int63(),
			C: rand.Int63(),
			D: rand.Int63(),
		}
	}
	return ds
}

func genDataPtrs(n int) []*Data {
	ds := make([]*Data, n)
	for i := range ds {
		d := &Data{
			A: rand.Int63(),
			B: rand.Int63(),
			C: rand.Int63(),
			D: rand.Int63(),
		}
		ds[i] = d
	}
	return ds
}

func BenchmarkValueSlice(b *testing.B) {
	ds := genData(1e6)
	b.ResetTimer()
	var sum int64
	for i := 0; i < b.N; i++ {
		for _, v := range ds {
			sum += v.A
		}
	}
	_ = sum
}

func BenchmarkPointerSlice(b *testing.B) {
	ds := genDataPtrs(1e6)
	b.ResetTimer()
	var sum int64
	for i := 0; i < b.N; i++ {
		for _, v := range ds {
			sum += v.A
		}
	}
	_ = sum
}

func BenchmarkPtrToSlice(b *testing.B) {
	ds := genData(1e6)
	ptr := &ds
	b.ResetTimer()
	var sum int64
	for i := 0; i < b.N; i++ {
		for _, v := range *ptr {
			sum += v.A
		}
	}
	_ = sum
}

// 值切片链路
func A() []Data {
	return genData(1e6)
}
func B(ds []Data) []Data {
	return ds
}
func C(ds []Data) string {
	b, _ := json.Marshal(ds)
	return string(b)
}
func BenchmarkChainValueSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ds := A()
		ds2 := B(ds)
		_ = C(ds2)
	}
}

// 指针切片链路
func AP() []*Data {
	return genDataPtrs(1e6)
}
func BP(ds []*Data) []*Data {
	return ds
}
func CP(ds []*Data) string {
	b, _ := json.Marshal(ds)
	return string(b)
}
func BenchmarkChainPointerSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ds := AP()
		ds2 := BP(ds)
		_ = CP(ds2)
	}
}

// 切片指针链路
func APs() *[]Data {
	ds := genData(1e6)
	return &ds
}
func BPs(ds *[]Data) *[]Data {
	return ds
}
func CPs(ds *[]Data) string {
	b, _ := json.Marshal(ds)
	return string(b)
}
func BenchmarkChainPtrToSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ds := APs()
		ds2 := BPs(ds)
		_ = CPs(ds2)
	}
}
