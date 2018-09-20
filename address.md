[ripemd160](https://blog.csdn.net/u013397318/article/details/80937583)

[golang利用gob序列化struct对象保存到本地](https://studygolang.com/articles/2888)

### 公钥私钥转换过程

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