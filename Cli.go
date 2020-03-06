package main

// 客户端命令控制工具

import (
	"fmt"
	"os"
	"strconv"
)

type Cli struct {
	bc *BlockChain
}

const Usage = `
----Usage:
		addBlock --data Data      "添加区块"
		printChain                "正向打印区块"
		printChainR               "反向打印区块"
		getBalance --address addr "获取账户余额"
		send -f from -t to -a amount -m miner -d data "from向to转账amount比特币，由miner挖矿，并写入data"
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
	case "getBalance":
		if len(args) == 4 && args[2] == "--address" {
			address := args[3]
			self.GetBalance(address)
		} else {
			fmt.Println("参数错误")
			fmt.Printf(Usage)
		}
	case "send":
		fmt.Println("正在转账中，请等待...")
		//  1               5      7         9       10
		// send -f from -t to -a amount -m miner -d data
		from := args[3]
		to := args[5]
		amount, _ := strconv.ParseFloat(args[7], 64)
		miner := args[9]
		data := args[11]
		ok := self.Send(from, to, amount, miner, data)
		if ok {
			fmt.Println("转账成功！谢谢班长")
		}
	default:
		fmt.Printf("Command '%s' not found, did you mean:\n", cmd)
		fmt.Printf(Usage)
	}
}
