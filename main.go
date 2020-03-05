package main

func main() {
	bc := NewBlockChain("班长")
	cli := Cli{bc}
	cli.Run()
}
