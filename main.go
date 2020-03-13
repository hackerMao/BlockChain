package main

func main() {
	bc := NewBlockChain("1CGyBQquLudGjDcPRERPcRCeh8pQJwFGFj")
	cli := Cli{bc}
	cli.Run()
}
