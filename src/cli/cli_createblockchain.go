package cli

import (
	"block"
	"fmt"
	"log"
	"wallet"
)

//创建一个区块链

func (cli *CLI) createBlockchain(address string) {
	//校验钱包地址
	if !wallet.ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	//创建一个区块链
	bc := block.CreateBlockchain(address)
	//关闭数据库
	block.Close(bc)
	//重置UTXO集
	UTXOSet := block.UTXOSet{bc}
	UTXOSet.Reindex()
	fmt.Println("Done!")
}
