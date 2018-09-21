package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const printchain = "printchain"
const send = "send"
const getbalance = "getbalance"
const createblockchain = "createblockchain"
const listaddresses = "listaddresses"
const createwallet = "createwallet"

type CLI struct {
}

func Start() interface{} {
	cl := CLI{}
	cl.run()
	return nil
}

//打印用法
func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  createblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("  createwallet - Generates a new key-pair and saves it into the wallet file")
	fmt.Println("  getbalance -address ADDRESS - Get balance of ADDRESS")
	fmt.Println("  listaddresses - Lists all addresses from the wallet file")
	fmt.Println("  printchain - Print all the blocks of the blockchain")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT - Send AMOUNT of coins from FROM address to TO")
}

//校验参数
func (cli *CLI) validateArgs() {
	//参数的数组少于2
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

// 执行命令方法
func (cli *CLI) run() {
	cli.validateArgs()
	getBalanceCmd := flag.NewFlagSet(getbalance, flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet(createblockchain, flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet(createwallet, flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet(listaddresses, flag.ExitOnError)
	sendCmd := flag.NewFlagSet(send, flag.ExitOnError)
	printChainCmd := flag.NewFlagSet(printchain, flag.ExitOnError)
	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")
	switch os.Args[1] {
	case getbalance:
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case createblockchain:
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case createwallet:
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case listaddresses:
		err := listAddressesCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case printchain:
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case send:
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}
	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress)
	}
	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if createWalletCmd.Parsed() {
		cli.createWallet()
	}
	if listAddressesCmd.Parsed() {
		cli.listAddresses()
	}
	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}
		cli.send(*sendFrom, *sendTo, *sendAmount)
	}

}
