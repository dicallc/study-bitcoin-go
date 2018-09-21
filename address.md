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

### 复习

#### 创建钱包

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

#### 创建一个区块链

1. 校验钱包地址
2. 创建一个区块链
3. 关闭数据库

```
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
```

其中校验钱包地址还是根据那个图还原

