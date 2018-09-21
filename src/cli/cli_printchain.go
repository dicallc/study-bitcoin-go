package cli

import (
	"block"
	"fmt"
	"strconv"
)

//打印区块链上所有区块数据
func (cli *CLI) printChain() {
	bc := block.NewBlockchain("")
	defer block.Close(bc)

	bci := bc.Iterator()

	for {
		block := bci.Next()
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)

		fmt.Printf("PoW: %s\n", strconv.FormatBool(block.Validate()))
		fmt.Println()
		//创世块是没有前一个区块的，所以PrevBlockHash的值是没有的
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
