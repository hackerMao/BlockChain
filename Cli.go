package main

import (
	"fmt"
	"os"
)

type Cli struct {
	bc *BlockChain
}

const Usage = `
----Usage:
		addBlock --data Data      "添加区块"
		printChain                "正向打印区块"
		printChainR               "反向打印区块"
`

func (self *Cli) Run() {
	// 获取参数
	args := os.Args
	if len(args) < 2 {
		fmt.Printf(Usage)
		return
	}
	cmd := args[1]
	switch cmd {
	case "addBlock":
		fmt.Println("添加区块")
		if len(args) == 4 && args[2] == "--data" {
			// 获取数据
			data := args[3]
			// 添加区块
			self.AddBlock(data)
		} else {
			fmt.Println("参数错误")
			fmt.Println(Usage)
		}
	case "printChain":
		fmt.Println("正向打印区块")
		self.printChain()
	case "printChainR":
		fmt.Println("反向打印区块")
		self.printChain()
	default:
		fmt.Printf("Command '%s' not found, did you mean:\n", cmd)
		fmt.Printf(Usage)
	}
}
