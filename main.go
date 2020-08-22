package main

import (
	"fmt"
)

// Usage: 谁 200 谁结婚 收礼 结婚 2020-01-01
func main() {
	cfg := loadConfig()
	initDB(cfg.DB)

	fmt.Printf("成功导入%d条记录!\n", parseFile(cfg.File.In))
	fmt.Printf("成功导出所有的%d条记录!\n", exportFile(cfg.File.Out))
}
