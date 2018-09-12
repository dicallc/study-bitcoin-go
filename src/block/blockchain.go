package block

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

// 区块链
type Blockchain struct {
	tip []byte   //最顶层hash
	db  *bolt.DB //BoltDB数据库
}

// 保存区块数据
func (bc *Blockchain) AddBlock(data string) {
	////获取上一个区块
	//prevBlock := bc.Blocks[len(bc.Blocks)-1]
	////创建一个新的区块
	//newBlock := NewBlock(data, prevBlock.Hash)
	////新的区块添加到数组中
	//bc.Blocks = append(bc.Blocks, newBlock)
}

const dbFile = "db/blockchian.db"   //定义数据文件名
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

func CheckErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
