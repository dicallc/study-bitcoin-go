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
const reindexutxo = "reindexutxo"
const startnode = "startnode"

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
	fmt.Println("  createblockchain -address ADDRESS --创建区块链")
	fmt.Println("  exmple： main createblockchain -address 18g2nCpySpuNv39iUBUz3Df37ojYfJw9rX")

	fmt.Println("  createwallet --创建钱包")
	fmt.Println("  exmple： main createwallet")

	fmt.Println("  getbalance -address ADDRESS - 获取地址余额信息")
	fmt.Println("  exmple： main getbalance  -address 18g2nCpySpuNv39iUBUz3Df37ojYfJw9rX")

	fmt.Println("  listaddresses --输出所有钱包地址")
	fmt.Println("  exmple：main listaddresses")

	fmt.Println("  printchain -- 输出区块信息")
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
	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		fmt.Printf("NODE_ID env. var is not set!")
		os.Exit(1)
	}
	getBalanceCmd := flag.NewFlagSet(getbalance, flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet(createblockchain, flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet(createwallet, flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet(listaddresses, flag.ExitOnError)
	sendCmd := flag.NewFlagSet(send, flag.ExitOnError)
	reindexUTXOCmd := flag.NewFlagSet(reindexutxo, flag.ExitOnError)
	printChainCmd := flag.NewFlagSet(printchain, flag.ExitOnError)
	startNodeCmd := flag.NewFlagSet(startnode, flag.ExitOnError)
	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")
	startNodeMiner := startNodeCmd.String("miner", "", "Enable mining mode and send reward to ADDRESS")
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
	case reindexutxo:
		err := reindexUTXOCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case startnode:
		err := startNodeCmd.Parse(os.Args[2:])
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
	if reindexUTXOCmd.Parsed() {
		cli.reindexUTXO(nodeID)
	}
	if startNodeCmd.Parsed() {
		nodeID := os.Getenv("NODE_ID")
		if nodeID == "" {
			startNodeCmd.Usage()
			os.Exit(1)
		}
		cli.startNode(nodeID, *startNodeMiner)
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}
		cli.send(*sendFrom, *sendTo, *sendAmount)
	}

}
