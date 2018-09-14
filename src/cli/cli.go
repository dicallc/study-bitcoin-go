package cli

import (
	"block"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

type CLI struct {
}

func (cli *CLI) createBlockchain(address string) {
	bc := block.CreateBlockchain(address)
	block.Close(bc) //关闭数据库
	fmt.Println("Done!")
}

func Start(bc *block.Blockchain) interface{} {
	cl := CLI{}
	cl.run()
	return nil
}

//查询余额
func (cli *CLI) getBalance(address string) {
	bc := block.NewBlockchain(address)
	defer block.Close(bc)
	balance := 0
	//查询所有未经使用的交易地址
	UTXOS := bc.FindUTXO(address)
	//算出未使用的交易地址的value
	for _, out := range UTXOS {
		balance += out.Value
	}
	fmt.Printf("Balance of '%s': %d\n", address, balance)
}

////添加区块数据
//func (cli *CLI) addBlock(data string) {
//	cli.bc.AddBlock(data)
//	fmt.Println("Success!")
//}
//
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

//打印用法
func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  addblock -data BLOCK_DATA - add a block to the blockchain")
	fmt.Println("  printchain - print all the blocks of the blockchain")
}

//校验参数
func (cli *CLI) validateArgs() {
	//参数的数组少于2
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

//发送
func (cli *CLI) send(from, to string, amount int) {
	bc := block.NewBlockchain(from)
	defer block.Close(bc)

	tx := block.NewUTXOTransaction(from, to, amount, bc)
	bc.MineBlock([]*block.Transaction{tx})
	fmt.Println("Success!")
}

const addblock = "addblock"
const printchain = "printchain"
const send = "send"

// 执行命令方法
func (cli *CLI) run() {
	cli.validateArgs()
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	addBlockCmd := flag.NewFlagSet(addblock, flag.ExitOnError)
	sendCmd := flag.NewFlagSet(send, flag.ExitOnError)
	printChainCmd := flag.NewFlagSet(printchain, flag.ExitOnError)
	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")
	switch os.Args[1] {
	case addblock:
		err := addBlockCmd.Parse(os.Args[2:])
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
