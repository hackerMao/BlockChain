package main

func main() {
	bc := NewBlockChain("张三")
	cli := Cli{bc}
	cli.Run()
}
