package main

import (
	"fmt"
)

func main() {
	fmt.Println("请用 go test -bench=. -benchmem 运行基准测试")
	fmt.Println("建议配合 go test -bench=. -benchmem -cpuprofile=cpu.out -memprofile=mem.out 使用 pprof 分析")
	// Go 1.24+ rand.Seed 已为 no-op，若需恢复旧行为请设置 GODEBUG=randseednop=0
}
