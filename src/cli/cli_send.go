package cli

import (
	"block"
	"fmt"
	"log"
	"wallet"
)

//发送
func (cli *CLI) send(from, to string, amount int) {
	if !wallet.ValidateAddress(from) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !wallet.ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}
	bc := block.GetBlockChainHandler()
	defer block.Close(bc)
	//创建新UTXO，放置到新的区块上去
	tx := block.NewUTXOTransaction(from, to, amount, bc)
	//挖矿奖励
	cbTx := block.NewCoinbaseTX(from, "")
	//合并两个UTXO成为[]
	txs := []*block.Transaction{cbTx, tx}
	bc.MineBlock(txs)
	fmt.Println("Success!")
}

//func (cli *CLI) TestSend(from, to string, amount int) {
//	if !wallet.ValidateAddress(from) {
//		log.Panic("ERROR: Sender address is not valid")
//	}
//	if !wallet.ValidateAddress(to) {
//		log.Panic("ERROR: Recipient address is not valid")
//	}
//	bc := block.GetBlockChainHandler()
//	defer block.Close(bc)
//	//创建新UTXO，放置到新的区块上去
//	tx := block.NewUTXOTransaction(from, to, amount, bc)
//	bc.MineBlock([]*block.Transaction{tx})
//	fmt.Println("Success!")
//}
