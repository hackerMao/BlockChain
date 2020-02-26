package main

import "fmt"

func main() {
	total := 0.0
	currentReword := 50.0
	blockIntervel := 21.0 //单位：万

	for currentReword > 0 {
		amount := currentReword * blockIntervel
		currentReword *= 0.5
		total += amount
	}
	fmt.Println("比特币总量：", total, "万")
	fmt.Println(4286700/4.7)
}
