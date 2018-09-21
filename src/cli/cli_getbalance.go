package cli

import (
	"block"
	"fmt"
	"log"
	"utils"
	"wallet"
)

//查询余额
func (cli *CLI) getBalance(address string) {
	if !wallet.ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	bc := block.GetBlockChainHandler()
	defer block.Close(bc)
	balance := 0
	pubKeyHash := utils.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	//查询所有未经使用的交易地址
	UTXOS := bc.FindUTXO(pubKeyHash)
	//算出未使用的交易地址的value
	for _, out := range UTXOS {
		balance += out.Value
	}
	fmt.Printf("Balance of '%s': %d\n", address, balance)
}
