# []Type、[]_Type、_[]Type 的区别

本目录系统性介绍 Go 语言中三种常见切片/指针组合类型：

- `[]Type`：元素为值类型的切片
- `[]*Type`：元素为指针类型的切片
- `*[]Type`：指向切片的指针

---

## 一、底层原理剖析

### 1. 结构与内存布局

- `[]Type`：切片本身 24 字节，元素为值类型，数据连续存储，cache 友好。
- `[]*Type`：切片本身 24 字节，元素为指针（8 字节），指向堆上对象，数据分散，cache miss 多。
- `*[]Type`：8 字节指针，指向切片结构体，数据布局与 `[]Type` 一致，但多一层指针。

### 2. 传递与拷贝

- `[]Type`：函数间传递只拷贝切片头（24 字节），底层数据不拷贝。
- `[]*Type`：同上，切片头 24 字节，指针数组，底层对象不拷贝。
- `*[]Type`：只拷贝 8 字节指针，适合需要修改切片本身的场景。

### 3. 访问与遍历

- `[]Type`：直接访问值，CPU 指令少，cache 命中率高。
- `[]*Type`：访问需解引用，CPU 指令多，cache miss 增多。
- `*[]Type`：遍历与 `[]Type` 一致，区别在于可修改切片本身。

---

## 二、内存占用与 GC 行为

- `[]Type`：数据连续，GC 只追踪一块大内存。
- `[]*Type`：指针数组+大量堆对象，GC 追踪压力大。
- `*[]Type`：本质同 `[]Type`，GC 行为一致。

---

## 三、函数调用链对比

### 1. []Type

```go
func A() []Data { return genData(1e6) }
func B(ds []Data) []Data { return ds }
func C(ds []Data) string { b, _ := json.Marshal(ds); return string(b) }
```

- 传递高效，序列化快，适合小型结构体。

### 2. []\*Type

```go
func AP() []*Data { return genDataPtrs(1e6) }
func BP(ds []*Data) []*Data { return ds }
func CP(ds []*Data) string { b, _ := json.Marshal(ds); return string(b) }
```

- 适合大型结构体或需共享/修改，序列化和遍历性能较低。

### 3. \*[]Type

```go
func APs() *[]Data { ds := genData(1e6); return &ds }
func BPs(ds *[]Data) *[]Data { return ds }
func CPs(ds *[]Data) string { b, _ := json.Marshal(ds); return string(b) }
```

- 适合需要修改切片本身，遍历性能与 `[]Type` 一致。

---

## 四、Benchmark 代码与结果对比

### 1. Benchmark 代码

详见 [main_test.go](main_test.go)，三种基础遍历与三种链式调用均有对应 benchmark：

- `BenchmarkValueSlice`：[]Type 遍历
- `BenchmarkPointerSlice`：[]\*Type 遍历
- `BenchmarkPtrToSlice`：\*[]Type 遍历
- `BenchmarkChainValueSlice`：[]Type 链式调用
- `BenchmarkChainPointerSlice`：[]\*Type 链式调用
- `BenchmarkChainPtrToSlice`：\*[]Type 链式调用

### 2. 运行结果（Apple M2, 1e6 元素，2025-05-28）

#### 基础遍历：

```
BenchmarkValueSlice-8           4094    301527 ns/op      0 B/op         0 allocs/op
BenchmarkPointerSlice-8         1939    612211 ns/op      0 B/op         0 allocs/op
BenchmarkPtrToSlice-8           3618    348847 ns/op      0 B/op         0 allocs/op
```

#### 链式调用：

```
BenchmarkChainValueSlice-8         6   176909785 ns/op   271799080 B/op      13 allocs/op
BenchmarkChainPointerSlice-8       6   187386639 ns/op   503484700 B/op  1000047 allocs/op
BenchmarkChainPtrToSlice-8         7   162671185 ns/op   380446460 B/op      28 allocs/op
```

### 3. 性能原理剖析

- `[]Type`（值切片）
  - 数据连续，遍历时内存访问线性，CPU cache 命中率高，指令流水线友好。
  - 访问元素无需解引用，CPU 指令少。
  - 适合小型结构体和高频遍历场景。
- `[]*Type`（指针切片）
  - 切片存放指针，实际数据分布在堆上，内存不连续。
  - 遍历时需多一次指针解引用，cache miss 增多，CPU 指令数增加。
  - 分配次数多，GC 追踪压力大，序列化时性能下降。
  - 适合大型结构体或需共享/频繁修改的场景。
- `*[]Type`（切片指针）
  - 仅多一层指针，遍历和内存布局与 `[]Type` 一致。
  - 适合需要修改切片本身的场景。

> 综上，`[]Type`/`*[]Type` 性能接近且优于 `[]*Type`，后者在遍历、序列化、GC 等场景下均有明显劣势。

---

## 五、pprof 资源消耗对比

### 1. 内存热点

- `[]Type`/`*[]Type`：主要分配在切片和序列化缓冲区。
- `[]*Type`：分配热点分散，GC 追踪对象多。

### 2. CPU 热点

- `[]Type`/`*[]Type`：encoding/json、内存拷贝。
- `[]*Type`：encoding/json、指针解引用、cache miss。

---

## 六、小型结构体与大型结构体的判断标准

- 小于等于 64 字节为“小型结构体”，推荐 `[]Type` 或 `*[]Type`。
- 大于 64 字节为“大型结构体”，推荐 `[]*Type`。
- 以 cache line（64 字节）为参考，实际可结合 profile 调整。

---

## 七、实际业务建议

- 只读、遍历多、结构体小：优先 `[]Type` 或 `*[]Type`。
- 需共享/频繁修改/结构体大：可用 `[]*Type`。
- 需修改切片本身（如扩容、重置）：用 `*[]Type`。
- 所有链路、benchmark、资源消耗对比均应包含三种写法，main_test.go 可直接运行。

---

> 推荐实际运行 benchmark 和 pprof，结合自身业务场景选择合适的数据结构。
