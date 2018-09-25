[ripemd160](https://blog.csdn.net/u013397318/article/details/80937583)

[golang利用gob序列化struct对象保存到本地](https://studygolang.com/articles/2888)

### PublicKey转换成address

![mark](http://7xnk07.com1.z0.glb.clouddn.com/blog/180920/H2HIHaE78d.png?imageslim)

### 序列化内容代码

```
	var content bytes.Buffer
	encoder:=gob.NewEncoder(&content)
	err:=encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}
```

#### 增强序列化内容gob

```
gob.Register(elliptic.P256())//注册的目的是为了可以序列化任何类型
```



### 文件操作:ioutil

```
err=ioutil.WriteFile(walletFile,content.Bytes(),0644) //0644是文件权限
```



### 使用脚本解锁区块

之前只是对比一下用户名是否一样，现在必须使用其他值去代替

- 数字签名
- 公钥

### pri-public

就在于小写和大写第一个字母区别



### **make** 和**new** 的区别

make也是用于内存分配的，但是和new不同，它只用于chan、map以及切片的内存创建

而且它返回的类型就是这三个类型本身，而不是他们的指针类型

因为这三种类型就是引用类型，所以就没有必要返回他们的指针了 

### 重构过程

1. 引进钱包概念
2. 修改UTXO的输入输出-主要是重构之前的用户名验证变成数字签名和公钥
3. 正式去修改UTXO-address变成publicHash引发的灾难
4. 签名的加入修改UTXO

### 签名思考图

![mark](http://7xnk07.com1.z0.glb.clouddn.com/blog/180920/e3mEleDfa1.png?imageslim)

签名过程中Block4是新生成的

图蓝色边代表一个区块，而区块包含UTXO，则红色边就是UTXO

- Block4是要生成的区块
- Block3是最新的区块

### 签名有什么用

交易必须签名，这是比特币唯一能够保证不去花别人钱去交易的办法

如果签名无效，交易也是为无效，因此无法添加到区块上

1. 公钥哈希存储在解锁输出中，交易的发件人
2. 公钥哈希存储在新的锁定输出中，交易的收件

### 复习思路

#### 1.创建钱包

1. 获取Wallets

   1.1 创建wallets对象同时把本地内容读取出来

2. 创建一个钱包

3. 保存钱包

4. 打印钱包地址

```go
//创建钱包
func (cli *CLI) createWallet() {
	wallets, _ := wallet.LoadWallets()
	address := wallets.CreateWallet()
	wallets.SaveToFile()
	fmt.Printf("Your new address: %s\n", address)
}
//创建wallets对象同时把本地内容读取出来
func LoadWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)
	err := wallets.LoadFromFile()
	return &wallets, err
}
// 加载钱包,读取本地dat 序列化内容至 Wallets
func (ws *Wallets) LoadFromFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}
	fileContent, err := ioutil.ReadFile(walletFile)
	utils.CheckErr("ioutil.ReadFile", err)
	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}
	ws.Wallets = wallets.Wallets
	return nil
}
//创建一个钱包
func (ws *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := fmt.Sprintf("%s", wallet.GetAddress())
	ws.Wallets[address] = wallet
	return address
}
//保存钱包
func (ws Wallets) SaveToFile() {
	var content bytes.Buffer
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}
	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}
```

##### 钱包的结构

以上就展示了一个创建钱包的主要流程，其中有个流程没做过详细说明

之所以不说就是想捉个击破

就是CreateWallet中，牵扯到钱包的结构了

```go
type Wallet struct {
	/**
	PrivateKey: ECDSA基于椭圆曲线
	使用曲线生成私钥，并从私钥生成公钥
	*/
	PrivateKey ecdsa.PrivateKey //私钥
	PublicKey  []byte           //公钥
}
```

公钥，私钥出来了，我们就看看他是怎么出来的吧

```go
//创建一个新钱包
func NewWallet() *Wallet {
	//公钥私钥生成
	private, public := newKeyPair()
	wallet := Wallet{private, public}
	return &wallet
}
//椭圆算法返回私钥与公钥
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	//实现了P-256的曲线
	curve := elliptic.P256()
	//获取私钥
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	utils.CheckErr("", err)
	//在基于椭圆曲线的算法中，公钥是曲线上的点。因此，公钥是X，Y坐标的组合。在比特币中，这些坐标被连接起来形成一个公钥
	pubKey := append(private.PublicKey.X.Bytes(), 	private.PublicKey.Y.Bytes()...)
	return *private, pubKey
}
```

这里你就记住椭圆曲线算法得到了一个新的公钥私钥，牵扯到密码学就不要深究



##### 获取钱包地址

```go
//得到一个钱包地址

func (w Wallet) GetAddress() []byte {
	//1.使用 RIPEMD160(SHA256(PubKey)) 哈希算法，取公钥并对其哈希两次
	pubKeyHash := HashPubKey(w.PublicKey)
	//2.给哈希加上地址生成算法版本的前缀
	versionedPayload := append([]byte{version}, pubKeyHash...)
	//3.对于第二步生成的结果，使用 SHA256(SHA256(payload)) 再哈希，计算校验和。校验和是结果哈希的前四个字节
	checksum := checksum(versionedPayload)
	//4.将校验和附加到 version+PubKeyHash 的组合中
	fullPayload := append(versionedPayload, checksum...)
	//5.使用 Base58 对 version+PubKeyHash+checksum 组合进行编码
	address := utils.Base58Encode(fullPayload)
	return address
}
// 使用RIPEMD160(SHA256(PubKey))哈希算法得到Hashpubkey
func HashPubKey(pubKey []byte) []byte {
	//1.256
	publicSHA256 := sha256.Sum256(pubKey)
	//2.160
	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	utils.CheckErr("", err)
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)
	return publicRIPEMD160
}
```

这里原理就是 参考本文 PublicKey转换成address 的图

#### 2.获取所有钱包地址

```go
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

//迭代所有钱包地址返回到数组中
func (ws *Wallets) GetAddresses() []string {
	var addresses []string
	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}
	return addresses
}


```



#### 3.创建一个区块链

1. 校验钱包地址
2. 创建一个区块链
3. 关闭数据库

```go
func (cli *CLI) createBlockchain(address string) {
	//校验钱包地址
	if !wallet.ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	//创建一个区块链
	bc := block.CreateBlockchain(address)
	//关闭数据库
	block.Close(bc)
	fmt.Println("Done!")
}

func ValidateAddress(address string) bool {
	pubKeyHash := utils.Base58Decode([]byte(address))
	//倒数两个为Checksum
	actualChecksum := pubKeyHash[len(pubKeyHash)-addressChecksumlen:]
	//第一个是版本号
	version := pubKeyHash[0]
	//剩下全是 pubKeyHash
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumlen]
	//计算checksum是否一致
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}
```

其中校验钱包地址还是根据那个图还原，地址用钱包地址作为区块链地址进行创建

#### 4.send（转账）

这个版本微改了一些东西

比如验证地址正确,UTXO的解锁 加锁变成了 非对称加密的公私钥，数字签名

```go
//发送
func (cli *CLI) send(from, to string, amount int) {
	if !wallet.ValidateAddress(from) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !wallet.ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}
	bc := block.GetBlockChainHandler()
	defer block.Close(bc)
	//创建新UTXO，放置到新的区块上去
	tx := block.NewUTXOTransaction(from, to, amount, bc)
	bc.MineBlock([]*block.Transaction{tx})
	fmt.Println("Success!")
}
```

这里其实没有任何改变,只是在UTXO的输入输出结构变了

```go
type Transaction struct {
	ID        []byte
	TxInputs  []TxInput  //输入
	TXOutputs []TXOutput //输出
}
//输入 收账记录
type TxInput struct {
	Txid      []byte //交易ID的hash
	Vout      int    //所引用Output的索引值
	Signature []byte //数字签名
	PubKey    []byte //钱包里的公钥
}
//输出 付钱记录
type TXOutput struct {
	Value      int    //支付给收款方金额值
	PubKeyHash []byte 
}
```

其中还有几个加锁，解锁的方法比如

##### 输出：

```
//lock只需锁定输出 人话版本：你给别人钱你的把别人地址给设置进入生成一个输出
func (out *TXOutput) Lock(address []byte) {
	//人话：最后的地址 (Base58Decode)解密 ,然后减去后四位的版本信息就得到 pubKeyHash
	pubKeyHash := utils.Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
}

//检查提供的公钥散列是否用于锁定输出
func (out *TXOutput) IsLockWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

//新的交易输出 人话转账给别人是要锁定一下
func NewTxOutput(value int, address string) *TXOutput {
	txo := &TXOutput{value, nil}
	txo.Lock([]byte(address))
	return txo
}
```

这样保证了，输出给别人的钱，只有接受钱的人能够解锁



##### 输入：

```
//判断当前输入是否和某个输出吻合
func (in *TxInput) UseKey(pubKeyHash []byte) bool {
	lockingHash := wallet.HashPubKey(in.PubKey)
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
```

后面这些方法都会用到，回到转账的逻辑中继续

##### NewUTXOTransaction

```go
func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TxInput
	var outputs []TXOutput
	//获取所有钱包地址
	wallets, err := wallet.LoadWallets()
	CheckErr(err)
	//获取付款的钱包
	part_wallet := wallets.GetWallet(from)
	pubKeyHash := wallet.HashPubKey(part_wallet.PublicKey)

	//返回合适的UTXO
	acc, validOutputs := bc.FindSuitableUTXOs(pubKeyHash, amount)
	//判断是否有那么多可花费的币
	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}
	//遍历有效UTXO的合集
	for txid, outIndexs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		CheckErr(err)
		//遍历所有引用UTXO的索引，每一个索引需要创建一个Input
		for _, outindex := range outIndexs {
			input := TxInput{txID, outindex, nil, pubKeyHash}
			inputs = append(inputs, input)
		}
	}
	// Build a list of outputs
	outputs = append(outputs, *NewTxOutput(amount, to))
	if acc > amount {
		outputs = append(outputs, *NewTxOutput(acc-amount, from)) // a change
	}
	tx := Transaction{nil, inputs, outputs}
	tx.SetID()
	//签名
	bc.SignTransaction(&tx, part_wallet.PrivateKey)
	return &tx
}
```

整体思路还是根据from得到钱包，在里面寻找合适UTXO，有合适金额的输出（OutPut）

然后根据输出，创建输入,以便于生成新的UTXO

> 在和上个版本区别：在寻找合适的UTXO中，就把之前的对比地址换成了对比publicHash

当然最重要的是，最后需要签名

```go
func (bc *Blockchain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	preTXs := make(map[string]Transaction)
	//根据传入的UTXO，遍历其输入，找到其中UTXO，这些UTXO基本就是包含其有关输出的
	for _, vin := range tx.TxInputs {
		preTX, err := bc.FindTransaction(vin.Txid)
		CheckErr(err)
		preTXs[hex.EncodeToString(preTX.ID)] = preTX
	}
	//{2222222222,2222-UTXO}
	tx.Sign(privKey, preTXs)
}
//签名过程
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	//判断是否是coinbase交易
	if tx.IsCoinbase() {
		return
	}
	//TX 也是4  根据输入去找有不有对应的输出
	for _, vin := range tx.TxInputs {
		if prevTXs[hex.EncodeToString(vin.Txid)].ID == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}
	//txCopy Block4 UTXO      4444也拷贝一份 prevTXs:Block2 2222
	txCopy := tx.TrimmedCopy()
	//这里for循环就是把制作UTXO的输入 重新做了一遍 防止篡改？
	for inID, vin := range txCopy.TxInputs {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.TxInputs[inID].Signature = nil
		//所消费的Out公钥被引进作为了pubKey
		txCopy.TxInputs[inID].PubKey = prevTx.TXOutputs[vin.Vout].PubKeyHash
		//根据UTXO 数据生成id
		txCopy.ID = txCopy.Hash()
		txCopy.TxInputs[inID].PubKey = nil
		//真正的签名部分 钱包-私钥，UTXO-ID
		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.TxInputs[inID].Signature = signature
	}
}
```

这里如果在Sign方法中看不懂就记住  真正的签名部分 钱包-私钥，UTXO-ID

UTXO-ID又是他自己又从原数据中生成数据，创造了一遍输入，生成的

而这些就是签名，而验证签名又在哪里呢

就在开采区块那加上了

##### Verify 验证UTXO

```
//开采区块
func (bc *Blockchain) MineBlock(transactions []*Transaction) {
	var lastHash []byte //最新一个hash
	for _, tx := range transactions {
		if bc.VerifyTransaction(tx) != true {
			log.Panic("ERROR: Invalid transaction")
		}
	}
....
}
//验证区块
func (bc *Blockchain) VerifyTransaction(tx *Transaction) bool {
	prevTXs := make(map[string]Transaction)
	for _, vin := range tx.TxInputs {
		//根据传入的UTXO，遍历其输入，找到其中UTXO，这些UTXO基本就是包含其有关输出的
		prevTX, err := bc.FindTransaction(vin.Txid)
		CheckErr(err)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}
	return tx.Verify(prevTXs)
}
//验证核心方法
func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}
	for _, vin := range tx.TxInputs {
		if prevTXs[hex.EncodeToString(vin.Txid)].ID == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}
	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()
	for inID, vin := range tx.TxInputs {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.TxInputs[inID].Signature = nil
		txCopy.TxInputs[inID].PubKey = prevTx.TXOutputs[vin.Vout].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.TxInputs[inID].PubKey = nil
		//私钥 ID
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])
		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keyLen / 2)])
		y.SetBytes(vin.PubKey[(keyLen / 2):])
		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		//公钥和id，以及签名数据求证
		if ecdsa.Verify(&rawPubKey, txCopy.ID, &r, &s) == false {
			return false
		}
	}
	return true
}
```

验证主要也只能拿公钥和ID去，这就是非对称加密的使用了