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
- **原因**：切片元素为值类型，数据在内存中连续存储，遍历时 cache 命中率高，CPU 指令少，序列化时无需多次解引用，整体性能最佳。

### 2. []\*Type

```go
func AP() []*Data { return genDataPtrs(1e6) }
func BP(ds []*Data) []*Data { return ds }
func CP(ds []*Data) string { b, _ := json.Marshal(ds); return string(b) }
```

- 适合大型结构体或需共享/修改，序列化和遍历性能较低。
- **原因**：切片元素为指针，遍历和序列化时需多次解引用，CPU 指令增多，且指针指向的对象分布在堆上，内存分散，cache miss 增多，GC 追踪压力大，导致整体性能下降。

### 3. \*[]Type

```go
func APs() *[]Data { ds := genData(1e6); return &ds }
func BPs(ds *[]Data) *[]Data { return ds }
func CPs(ds *[]Data) string { b, _ := json.Marshal(ds); return string(b) }
```

- 适合需要修改切片本身，遍历性能与 `[]Type` 一致。
- **原因**：本质上还是值类型切片，只是多了一层指针，遍历和序列化时数据依然连续，cache 友好，性能与 `[]Type` 基本一致，但可修改切片本身（如扩容、重置）。

---

## 四、Benchmark 代码与结果对比

### 1. Benchmark 代码

详见 [main_test.go](main_test.go)，三种基础遍历与三种链式调用均有对应 benchmark，且分别覆盖小结构体（Data）和大结构体（BigData）：

- `BenchmarkValueSlice`/`BenchmarkBigValueSlice`：[]Type 遍历
- `BenchmarkPointerSlice`/`BenchmarkBigPointerSlice`：[]\*Type 遍历
- `BenchmarkPtrToSlice`/`BenchmarkBigPtrToSlice`：\*[]Type 遍历
- `BenchmarkChainValueSlice`/`BenchmarkChainBigValueSlice`：[]Type 链式调用
- `BenchmarkChainPointerSlice`/`BenchmarkChainBigPointerSlice`：[]\*Type 链式调用
- `BenchmarkChainPtrToSlice`/`BenchmarkChainBigPtrToSlice`：\*[]Type 链式调用

### 2. 运行结果（Apple M2, 2025-05-28）

#### 小结构体遍历：

```
BenchmarkValueSlice-8           4094    301527 ns/op      0 B/op         0 allocs/op
BenchmarkPointerSlice-8         1939    612211 ns/op      0 B/op         0 allocs/op
BenchmarkPtrToSlice-8           3618    348847 ns/op      0 B/op         0 allocs/op
```

| 类型     | 执行次数 | 平均耗时（约） | 内存占用 | 分配次数 | 说明               |
| -------- | -------- | -------------- | -------- | -------- | ------------------ |
| []Type   | 4094     | 302 μs         | 0        | 0        | 普通切片，遍历快   |
| []\*Type | 1939     | 612 μs         | 0        | 0        | 指针切片，遍历慢   |
| \*[]Type | 3618     | 349 μs         | 0        | 0        | 切片指针，遍历较快 |

#### 小结构体链式调用：

```
BenchmarkChainValueSlice-8         6   176909785 ns/op   271799080 B/op      13 allocs/op
BenchmarkChainPointerSlice-8       6   187386639 ns/op   503484700 B/op  1000047 allocs/op
BenchmarkChainPtrToSlice-8         7   162671185 ns/op   380446460 B/op      28 allocs/op
```

| 类型     | 执行次数 | 平均耗时（约） | 内存占用（约） | 分配次数（约） | 说明                             |
| -------- | -------- | -------------- | -------------- | -------------- | -------------------------------- |
| []Type   | 6        | 177 ms         | 259 MB         | 13             | 普通切片，内存占用低，分配少     |
| []\*Type | 6        | 187 ms         | 480 MB         | 100 万         | 指针切片，内存占用高，分配极多   |
| \*[]Type | 7        | 163 ms         | 363 MB         | 28             | 切片指针，内存占用中等，分配较少 |

#### 大结构体遍历：

```
BenchmarkBigValueSlice-8         35137     33003 ns/op      0 B/op         0 allocs/op
BenchmarkBigPointerSlice-8        4675    291991 ns/op      0 B/op         0 allocs/op
BenchmarkBigPtrToSlice-8         39066     31747 ns/op      0 B/op         0 allocs/op
```

| 类型     | 执行次数 | 平均耗时（约） | 内存占用 | 分配次数 | 说明               |
| -------- | -------- | -------------- | -------- | -------- | ------------------ |
| []Type   | 35,137   | 33 μs          | 0        | 0        | 普通切片，遍历极快 |
| []\*Type | 4,675    | 292 μs         | 0        | 0        | 指针切片，遍历慢   |
| \*[]Type | 39,066   | 32 μs          | 0        | 0        | 切片指针，遍历极快 |

#### 大结构体链式调用：

```
BenchmarkChainBigValueSlice-8         5   222286700 ns/op   235427011 B/op      35 allocs/op
BenchmarkChainBigPointerSlice-8       5   219274592 ns/op   263071192 B/op  100044 allocs/op
BenchmarkChainBigPtrToSlice-8         5   221266492 ns/op   289114571 B/op      50 allocs/op
```

| 类型     | 执行次数 | 平均耗时（约） | 内存占用（约） | 分配次数（约） | 说明                       |
| -------- | -------- | -------------- | -------------- | -------------- | -------------------------- |
| []Type   | 5        | 222 ms         | 224 MB         | 35             | 普通切片，分配较少         |
| []\*Type | 5        | 219 ms         | 251 MB         | 10 万          | 指针切片，分配次数明显增多 |
| \*[]Type | 5        | 221 ms         | 276 MB         | 50             | 切片指针，分配适中         |

### 3. 性能原理剖析与结论

- **小结构体场景**：
  - `[]Type`/`*[]Type` 遍历和链式调用均明显优于 `[]*Type`，cache 友好，CPU 指令少。
  - `[]*Type` 分配次数多，GC 压力大，序列化和遍历性能劣势明显。
- **大结构体场景**：
  - `[]*Type` 在链式调用时内存分配和 GC 压力优势明显（分配次数远低于小结构体），但遍历性能仍不如 `[]Type`/`*[]Type`。
  - `[]Type`/`*[]Type` 适合高频遍历、只读场景，`[]*Type` 适合需频繁修改/共享大对象。
- **函数调用链建议**：
  - 只读/遍历多/结构体小：优先 `[]Type` 或 `*[]Type`。
  - 结构体大且需共享/频繁修改：可用 `[]*Type`，但需权衡遍历和 GC 性能。
  - 需修改切片本身（如扩容、重置）：用 `*[]Type`。

> 所有链路、benchmark、资源消耗对比均已覆盖三种写法和大小结构体，main_test.go 可直接运行。

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

### 如何快速判断结构体大小

1. **字段类型估算**：
   - 常见类型字节数：int64/uint64/float64 为 8 字节，int32/float32 为 4 字节，指针为 8 字节（64 位系统）。
   - 结构体总大小 ≈ 各字段类型字节数之和（注意内存对齐可能略大）。
2. **代码辅助法**：
   - 可用 `unsafe.Sizeof` 快速获取结构体实际大小：
     ```go
     fmt.Println(unsafe.Sizeof(YourStruct{}))
     ```
   - 也可借助 IDE 悬浮提示、GoLand/VSCode 插件等。
3. **经验法则**：
   - 结构体字段总数较多、包含数组/嵌套 struct/大字符串/切片等，通常为“大型结构体”。
   - 仅有少量基础类型字段，通常为“小型结构体”。

> 建议开发者在设计数据结构时，关注结构体大小，必要时用 `unsafe.Sizeof` 明确验证。

---

## 七、实际业务建议

- 只读、遍历多、结构体小：优先 `[]Type` 或 `*[]Type`。
- 需共享/频繁修改/结构体大：可用 `[]*Type`。
- 需修改切片本身（如扩容、重置）：用 `*[]Type`。
- 所有链路、benchmark、资源消耗对比均应包含三种写法，main_test.go 可直接运行。

---

> 推荐实际运行 benchmark 和 pprof，结合自身业务场景选择合适的数据结构。
