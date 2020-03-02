package main

import (
	"fmt"
	"os"
)

func main() {
	length := len(os.Args)
	fmt.Printf("len(args): %d\n", length)
	for index, value := range os.Args {
		fmt.Printf("args[%d]=%s\n", index, value)
	}
}