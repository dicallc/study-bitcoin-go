package cli

import (
	"block"
	"fmt"
	"log"
	"wallet"
)

//发送
func (cli *CLI) send(from, to string, amount int, nodeID string, mineNow bool) {
	if !wallet.ValidateAddress(from) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !wallet.ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}
	bc := block.GetBlockChainHandler()
	UTXOSet := block.UTXOSet{bc}
	defer block.Close(bc)
	wallets, err := wallet.LoadWallets(nodeID)
	block.CheckErr(err)
	wallet := wallets.GetWallet(from)
	//创建新UTXO，放置到新的区块上去
	tx := block.NewUTXOTransaction(&wallet, to, amount, &UTXOSet)
	if mineNow {
		//挖矿奖励
		cbTx := block.NewCoinbaseTX(from, "")
		//合并两个UTXO成为[]
		txs := []*block.Transaction{cbTx, tx}
		newBlock := bc.MineBlock(txs)
		UTXOSet.Update(newBlock)
	} else {
		sendTx(block.KnownNodes[0], tx)
	}

	fmt.Println("Success!")
}
