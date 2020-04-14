package yangkai

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

type Aes struct {
	Key string
}

func (this *Aes) AES128(origin string,handle bool) (string,error) {
	if len(this.Key) != 16 {
		return "",errors.New("key 的长度不正确 需要16位数")
	}
	var data string
	var err error
	if handle == true {
		data,err = this.Encrypt(origin,this.Key)
	}else{
		data,err = this.Decrypt(origin,this.Key)
	}

	return data,err
}

func (this *Aes) AES192(origin string,handle bool) (string,error) {
	if len(this.Key) != 24 {
		return "",errors.New("key 的长度不正确 需要24位数")
	}
	var data string
	var err error
	if handle == true {
		data,err = this.Encrypt(origin,this.Key)
	}else{
		data,err = this.Decrypt(origin,this.Key)
	}

	return data,err
}

func (this *Aes) AES256(origin string,handle bool) (string,error) {
	if len(this.Key) != 32 {
		return "",errors.New("key 的长度不正确 需要32位数")
	}
	var data string
	var err error
	if handle == true {
		data,err = this.Encrypt(origin,this.Key)
	}else{
		data,err = this.Decrypt(origin,this.Key)
	}

	return data,err
}


func (this *Aes) Encrypt(orig string, key string) (string,error) {
	// 转成字节数组
	origData := []byte(orig)
	k := []byte(key)
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, err := aes.NewCipher(k)
	if err != nil {
		return "",err
	}
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	origData = this.pkcs7Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	// 创建数组
	cryted := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(cryted, origData)
	return base64.StdEncoding.EncodeToString(cryted),nil
}

func (this *Aes) Decrypt(cryted string, key string) (string,error) {
	// 转成字节数组
	crytedByte, _ := base64.StdEncoding.DecodeString(cryted)
	k := []byte(key)
	// 分组秘钥
	block, err := aes.NewCipher(k)
	if err != nil {
		return "",err
	}
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	// 创建数组
	orig := make([]byte, len(crytedByte))
	// 解密
	blockMode.CryptBlocks(orig, crytedByte)
	// 去补全码
	orig = this.pkcs7UnPadding(orig)
	return string(orig),nil
}
//补码
//AES加密数据块分组长度必须为128bit(byte[16])，密钥长度可以是128bit(byte[16])、192bit(byte[24])、256bit(byte[32])中的任意一个。
func (this *Aes) pkcs7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
//去码
func (this *Aes) pkcs7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

type Rsa struct {

}

// 加密
func (this *Rsa) Encrypt(origData string,publicKey string) (string, error) {
	//解密pem格式的公钥
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return "", errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//加密
	newData,err := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(origData))
	if err != nil {
		return "", err
	}
	return string(newData),err
}

// 解密
func (this *Rsa) Decrypt(ciphertext string,privateKey string) (string, error) {
	//解密
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return "", errors.New("private key error!")
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	// 解密
	newData,err := rsa.DecryptPKCS1v15(rand.Reader, priv, []byte(ciphertext))
	return string(newData),err
}

type Hash struct {

}

func (this *Hash) Sha256(data string) string {
	private := sha256.Sum256([]byte(data))
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%x",private)))
}

func (this *Hash) MD5(data string) string {
	private := md5.Sum([]byte(data))
	return fmt.Sprintf("%x",private)
}