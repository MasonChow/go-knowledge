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

type BigData struct {
	A [16]int64  // 128 字节
	B [8]float64 // 64 字节
	C [4]int32   // 16 字节
	D [32]byte   // 32 字节
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

func genBigData(n int) []BigData {
	ds := make([]BigData, n)
	for i := range ds {
		for j := range ds[i].A {
			ds[i].A[j] = rand.Int63()
		}
		for j := range ds[i].B {
			ds[i].B[j] = rand.Float64()
		}
		for j := range ds[i].C {
			ds[i].C[j] = rand.Int31()
		}
		for j := range ds[i].D {
			ds[i].D[j] = byte(rand.Intn(256))
		}
	}
	return ds
}

func genBigDataPtrs(n int) []*BigData {
	ds := make([]*BigData, n)
	for i := range ds {
		d := &BigData{}
		for j := range d.A {
			d.A[j] = rand.Int63()
		}
		for j := range d.B {
			d.B[j] = rand.Float64()
		}
		for j := range d.C {
			d.C[j] = rand.Int31()
		}
		for j := range d.D {
			d.D[j] = byte(rand.Intn(256))
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

func BenchmarkBigValueSlice(b *testing.B) {
	ds := genBigData(1e5)
	b.ResetTimer()
	var sum int64
	for i := 0; i < b.N; i++ {
		for _, v := range ds {
			sum += v.A[0]
		}
	}
	_ = sum
}

func BenchmarkBigPointerSlice(b *testing.B) {
	ds := genBigDataPtrs(1e5)
	b.ResetTimer()
	var sum int64
	for i := 0; i < b.N; i++ {
		for _, v := range ds {
			sum += v.A[0]
		}
	}
	_ = sum
}

func BenchmarkBigPtrToSlice(b *testing.B) {
	ds := genBigData(1e5)
	ptr := &ds
	b.ResetTimer()
	var sum int64
	for i := 0; i < b.N; i++ {
		for _, v := range *ptr {
			sum += v.A[0]
		}
	}
	_ = sum
}

// 大结构体链式调用
func BigA() []BigData {
	return genBigData(1e5)
}
func BigB(ds []BigData) []BigData {
	return ds
}
func BigC(ds []BigData) string {
	b, _ := json.Marshal(ds)
	return string(b)
}
func BenchmarkChainBigValueSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ds := BigA()
		ds2 := BigB(ds)
		_ = BigC(ds2)
	}
}
func BigAP() []*BigData {
	return genBigDataPtrs(1e5)
}
func BigBP(ds []*BigData) []*BigData {
	return ds
}
func BigCP(ds []*BigData) string {
	b, _ := json.Marshal(ds)
	return string(b)
}
func BenchmarkChainBigPointerSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ds := BigAP()
		ds2 := BigBP(ds)
		_ = BigCP(ds2)
	}
}
func BigAPs() *[]BigData {
	ds := genBigData(1e5)
	return &ds
}
func BigBPs(ds *[]BigData) *[]BigData {
	return ds
}
func BigCPs(ds *[]BigData) string {
	b, _ := json.Marshal(ds)
	return string(b)
}
func BenchmarkChainBigPtrToSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ds := BigAPs()
		ds2 := BigBPs(ds)
		_ = BigCPs(ds2)
	}
}
