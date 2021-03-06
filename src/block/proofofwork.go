package block

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"utils"
)

var (
	maxNonce = math.MaxInt64
)
const targetBits = 20

//工作量证明
type ProofOfWork struct {
	block  *Block //区块体
	target *big.Int //挖矿目标值
}

// 新的工作量证明，并且得到一个难度值
func NewProofOfWork(b *Block) *ProofOfWork {
	//这里将数字1左移256-20=236位得到难度计算值 //创建target,目标值是比特币动态调整的 2000块就调整下
	target := big.NewInt(1)//000000000000000000...01
	target.Lsh(target, uint(256-targetBits))//0x00000100000000

	pow := &ProofOfWork{b, target}

	return pow
}

//将区块体里面的数据转换成一个字节码数组，为下一个区块准备数据
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	//注意一定要将原始数据转换成[]byte，不能直接从字符串转
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			utils.IntToHex(pow.block.Timestamp),
			utils.IntToHex(int64(targetBits)),
			utils.IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

// 运行工作量证明，开始挖矿，找到小于难度目标值的Hash
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)
	for nonce < maxNonce {

		data := pow.prepareData(nonce)

		hash = sha256.Sum256(data)
		fmt.Printf("\r Dig into mine  %x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			fmt.Printf("\nfound hash: %x,nonce:%d\n", hash, nonce)
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")
	return nonce, hash[:]
}
// Validate validates block's PoW
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}