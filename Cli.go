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
		printChain					"正向打印区块"
		printChainR					"反向打印区块"
		getBalance --address addr	"获取账户余额"
		send -f from -t to -a amount -m miner -d data "from向to转账amount比特币，由miner挖矿，并写入data"
		newWallet					"新建一个钱包（私钥、公钥对）"
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
		fmt.Println("正在执行中，请等待...")
		//  1               5      7         9       10
		// send -f from -t to -a amount -m miner -d data
		from := args[3]
		to := args[5]
		amount, _ := strconv.ParseFloat(args[7], 64)
		miner := args[9]
		data := args[11]
		self.Send(from, to, amount, miner, data)
	case "newWallet":
		fmt.Printf("新建一个钱包：\n")
		self.NewWallet()
	default:
		fmt.Printf("Command '%s' not found, did you mean:\n", cmd)
		fmt.Printf(Usage)
	}
}
