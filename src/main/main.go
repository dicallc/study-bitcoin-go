package main

import (
	"block"
	"cli"
)

func main() {
	//创世块
	bc := block.NewBlockchain()
	cli.Start(bc)
}
