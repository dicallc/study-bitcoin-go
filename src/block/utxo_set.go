package block

import (
	"encoding/hex"
	"github.com/boltdb/bolt"
	"log"
)

//http://zhaocongliang.org/2017/08/31/blockchain-dev/
const utxoBucket = "chainstate"

// UTXOSet represents UTXO set
type UTXOSet struct {
	Blockchain *Blockchain
}

//1.有一个方法，功能：
//遍历整个数据库,读取所有的未花费的UTXO，然后将所有的UTXO存储到数据库
//reset
//去遍历数据库时
//[string]*TxOutputs
//{}
//找到一个公钥哈希的未花费输出，然后用来获取余额
func (u UTXOSet) FindUTXO(pubKeyHash []byte) []TXOutput {
	var UTXOs []TXOutput
	db := u.Blockchain.db

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			outs := DeserializeOutputs(v)

			for _, out := range outs.Outputs {
				if out.IsLockWithKey(pubKeyHash) {
					UTXOs = append(UTXOs, out)
				}
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return UTXOs
}

//当一个新的交易创建的时候。如果找到有所需数量的输出
func (u UTXOSet) FindSpendableOutputs(pubkeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	accumulated := 0
	db := u.Blockchain.db
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			txID := hex.EncodeToString(k)
			outs := DeserializeOutputs(v)
			for outIdx, out := range outs.Outputs {
				if out.IsLockWithKey(pubkeyHash) && accumulated < amount {
					accumulated += out.Value
					unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				}
			}
		}
		return nil
	})
	CheckErr(err)
	return accumulated, unspentOutputs
}
func (u UTXOSet) CountTransactions() int {
	db := u.Blockchain.db
	counter := 0

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			counter++
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return counter
}

//使用UTXO找到未花费输出，然后在数据库中进行存储。这里就是缓存的地方
func (u UTXOSet) Renindex() {
	db := u.Blockchain.db
	bucketName := []byte(utxoBucket)
	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(bucketName)
		if err != nil && err != bolt.ErrBucketNotFound {
			log.Panic(err)
		}
		_, err = tx.CreateBucket(bucketName)
		CheckErr(err)
		return nil
	})
	CheckErr(err)
	utxo := u.Blockchain.FindUTXO()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		for txID, outs := range utxo {
			//key 是UTXO-ID Value是TXOutputs
			key, err := hex.DecodeString(txID)
			CheckErr(err)
			b.Put(key, outs.Serialize())
			CheckErr(err)
		}
		return nil

	})
}

//每生成一个新区块，就更新一下UTXO集
func (u UTXOSet) Update(block *Block) {
	db := u.Blockchain.db
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		//遍历区块中的UTXO
		for _, tx := range block.Transactions {
			if tx.IsCoinbase() == false {
				//1.遍历UTXO中的输入
				for _, vin := range tx.TxInputs {
					updatedOuts := TXOutputs{}
					//从数据库中找寻未花费的输出
					outsBytes := b.Get(vin.Txid)
					outs := DeserializeOutputs(outsBytes)
					//遍历起未花费的输出 与其对比 如果不是就加入其中
					for outIdx, out := range outs.Outputs {
						if outIdx != vin.Vout {
							updatedOuts.Outputs = append(updatedOuts.Outputs, out)
						}
					}
					if len(updatedOuts.Outputs) == 0 {
						err := b.Delete(vin.Txid)
						CheckErr(err)
					} else {
						err := b.Put(vin.Txid, updatedOuts.Serialize())
						CheckErr(err)
					}
				}
			}
			newOutputs := TXOutputs{}
			//遍历UTXO中的输出
			for _, out := range tx.TXOutputs {
				newOutputs.Outputs = append(newOutputs.Outputs, out)
			}
			err := b.Put(tx.ID, newOutputs.Serialize())
			CheckErr(err)
		}
		return nil
	})
	CheckErr(err)
}
