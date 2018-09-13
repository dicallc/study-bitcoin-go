package block

import (
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

// 保存区块数据
func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte
	err := bc.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		if bucket == nil {
			os.Exit(1)
		}
		lastHash = bucket.Get([]byte(lastHashKey))
		return nil
	})
	CheckErr(err)
	//利最后一块hash 挖掘一块新的区块出来
	newBlock := NewBlock(data, lastHash)
	//在挖掘新块之后，我们将其序列化表示保存到数据块中并更新"l"，该密钥现在存储新块的哈希。
	err = bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		err := bucket.Put(newBlock.Hash, newBlock.Serialize())
		CheckErr(err)
		err = bucket.Put([]byte(lastHashKey), newBlock.Hash)
		CheckErr(err)
		bc.tip = newBlock.Hash
		return nil
	})

}

const dbFile = "blockchian.db"      //定义数据文件名
const blocksBucket = "blocks"       //区块桶
const lastHashKey = "last_hash_key" //区块桶
// 创建创世块
func NewBlockchain() *Blockchain {
	//go语言&表示获取存储的内存地址
	var lastHash []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil {
			//没有bucket，要去创建创世块。将数据填写到数据库的bucket中
			fmt.Println("No existing blockchain found. Creating a new one...")
			genesis := NewGenesisBlock()
			//1.创建一个桶
			bucket, err := tx.CreateBucket([]byte(blocksBucket))
			CheckErr(err)
			//2.将创世区块塞入桶中
			err = bucket.Put(genesis.Hash, genesis.Serialize())
			CheckErr(err)
			//3.更新最后的hash
			err = bucket.Put([]byte(lastHashKey), genesis.Hash)
			CheckErr(err)
			//4.返回最后的hash
			lastHash = genesis.Hash
		} else {
			//有bucket，取出最后区块的hash
			lastHash = b.Get([]byte(lastHashKey))
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	bc := Blockchain{lastHash, db}
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

//关闭方法
func Close(bc *Blockchain) error {
	return bc.db.Close()
}
func CheckErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
