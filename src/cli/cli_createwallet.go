package cli

import (
	"fmt"
	"wallet"
)

//创建钱包
func (cli *CLI) createWallet() {
	wallets, _ := wallet.LoadWallets()
	address := wallets.CreateWallet()
	wallets.SaveToFile()
	fmt.Printf("Your new address: %s\n", address)
}
