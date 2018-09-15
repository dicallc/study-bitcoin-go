package block

import (
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
)

// 区块链
type Blockchain struct {
	tip []byte   //最顶层hash
	db  *bolt.DB //BoltDB数据库
}

//返回指定地址能够支配的utxo的交易集合
func (bc *Blockchain) FindUnspentTransactions(address string) []Transaction {
	var unspentTXs []Transaction //
	spentTXOs := make(map[string][]int)
	//迭代所有区块
	bci := bc.Iterator()
	for {
		block := bci.Next()
		//从区块中取出账本迭代UTXO
		for _, tx := range block.Transactions {
			txId := hex.EncodeToString(tx.ID) //交易id转成String
			//去除挖矿奖励
			if tx.IsCoinbase() == false {
				//目的：找到已经消耗的UTXO，把他们放到一个集合里
				//需要两个字段来表示使用过的utxo:a.交易id,b.output的缩影
				for _, input := range tx.TxInputs {
					if input.CanUnlockOutputWith(address) {
						inTxID := hex.EncodeToString(input.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], input.Vout)
					}
				}
			}
			//遍历output 目的找到所有能支配的utxo 这里有个逻辑如果一区块Utxo的Input解锁地址和一个区块Out锁定地址一样了就是已经消费掉了
		Outputs:
			for outIdx, out := range tx.TXOutputs {
				if spentTXOs[txId] != nil {
					for _, spentOut := range spentTXOs[txId] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				if out.CanBeUnlockedWith(address) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

		}
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return unspentTXs
}

//寻找指定地址能够使用的utxo
func (bc *Blockchain) FindUTXO(address string) []TXOutput {
	var Utxos []TXOutput
	//未使用的UTXO
	unspentTransactions := bc.FindUnspentTransactions(address)
	//遍历交易 寻找OutPut是给这个地址的
	for _, tx := range unspentTransactions {
		//遍历output
		for _, out := range tx.TXOutputs {
			//当前地址拥有的utxo
			if out.CanBeUnlockedWith(address) {
				Utxos = append(Utxos, out)
			}
		}
	}
	return Utxos
}

//开采区块
func (bc *Blockchain) AddBlock(transactions []*Transaction) {
	var lastHash []byte //最新一个hash
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	//创造一个新区块
	newBlock := NewBlock(transactions, lastHash)
	//修改"l"的hash
	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}
		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}
		bc.tip = newBlock.Hash

		return nil
	})
}
func (bc *Blockchain) FindSuitableUTXOs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := bc.FindUnspentTransactions(address)
	accumulated := 0
Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.TXOutputs {
			if out.CanBeUnlockedWith(address) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}
	return accumulated, unspentOutputs
}

const dbFile = "blockchian.db"      //定义数据文件名
const blocksBucket = "blocks"       //区块桶
const lastHashKey = "last_hash_key" //区块桶
// 创建区块链,如果已经创建了，就直接返回
func NewBlockchain(address string) *Blockchain {
	if dbExists() == false {
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte("l"))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	bc := Blockchain{tip, db}
	return &bc
}

//获取区块链手柄 数据
func GetBlockChainHandler() *Blockchain {
	if dbExists() == false {
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte("l"))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	bc := Blockchain{tip, db}
	return &bc
}

//迭代器，就是一个对象，它里面包含一个游标，一直向前（后）移动，完成这个容器的遍历

type BlockchainIterator struct {
	currentHash []byte   //当前hash
	db          *bolt.DB //数据库
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}
	return bci
}

//迭代下一个区块
func (i *BlockchainIterator) Next() *Block {
	var block *Block
	err := i.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		encodedBlock := bucket.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)
		return nil
	})
	CheckErr(err)
	i.currentHash = block.PrevBlockHash
	return block
}
func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

//创世块data
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

// 创新一个新的区块数据
func CreateBlockchain(address string) *Blockchain {
	if dbExists() {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
		genesis := NewGenesisBlock(cbtx)
		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic(err)
		}
		err = b.Put(genesis.Hash, genesis.Serialize())
		if err != nil {
			log.Panic(err)
		}
		err = b.Put([]byte("l"), genesis.Hash)
		if err != nil {
			log.Panic(err)
		}
		tip = genesis.Hash
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	bc := Blockchain{tip, db}
	return &bc
}

//开采区块
func (bc *Blockchain) MineBlock(transactions []*Transaction) {
	var lastHash []byte //最新一个hash
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	//创造一个新区块
	newBlock := NewBlock(transactions, lastHash)
	//修改"l"的hash
	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}
		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}
		bc.tip = newBlock.Hash

		return nil
	})
}

//关闭方法
func Close(bc *Blockchain) error {
	return bc.db.Close()
}
func CheckErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
