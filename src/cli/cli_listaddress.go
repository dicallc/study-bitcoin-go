package cli

import (
	"fmt"
	"log"
	"wallet"
)

func (cli *CLI) listAddresses() {
	wallets, err := wallet.LoadWallets()
	if err != nil {
		log.Panic(err)
	}
	addresses := wallets.GetAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}
