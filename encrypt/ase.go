package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

var NewAes AesBase

func init() {
	NewAes.Key = []byte("qwerasdfzxcvtyui")
	NewAes.Iv = []byte("qwerasdfzxcvtyui")
	NewAes.Mode = "PKCS7"
}

type AesBase struct {
	Key  []byte
	Iv   []byte
	Mode string
}

/*CBC加密 按照golang标准库的例子代码
不过里面没有填充的部分,所以补上，根据key来决定填充blockSize
*/

// PKCS7Padding 使用PKCS7进行填充，IOS也是7
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}

// aes加密，填充模式由key决定，16位，24,32分别对应AES-128, AES-192, or AES-256.源码好像是写死16了
func (aesBase *AesBase) aesCBCEncrypt(rawData []byte) ([]byte, error) {
	block, err := aes.NewCipher(aesBase.Key)
	if err != nil {
		panic(err)
	}

	//填充原文
	blockSize := block.BlockSize()

	rawData = PKCS7Padding(rawData, blockSize)
	//初始向量IV必须是唯一，但不需要保密
	cipherText := make([]byte, blockSize+len(rawData))
	//block大小 16
	iv := cipherText[:blockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	//block大小和初始向量大小一定要一致
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[blockSize:], rawData)

	return cipherText, nil
}

func (aesBase *AesBase) aesCBCDecrypt(encryptData []byte) ([]byte, error) {
	block, err := aes.NewCipher(aesBase.Key)
	if err != nil {
		panic(err)
	}

	blockSize := block.BlockSize()

	if len(encryptData) < blockSize {
		panic("ciphertext too short")
	}
	iv := encryptData[:blockSize]
	encryptData = encryptData[blockSize:]

	// CBC mode always works in whole blocks.
	if len(encryptData)%blockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(encryptData, encryptData)
	//解填充
	switch aesBase.Mode {
	case "PKCS7":
		encryptData = PKCS7UnPadding(encryptData)
		return encryptData, nil
	case "Zero":
		encryptData = PKCS7UnPadding(encryptData)
		return encryptData, nil
	default:
		panic("Unknown padding mode")
	}
	// encryptData = PKCS7UnPadding(encryptData)
	// return encryptData, nil
}

func (aesBase *AesBase) Encrypt(rawData []byte) (string, error) {
	data, err := NewAes.aesCBCEncrypt(rawData)
	// fmt.Println(data)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func (aesBase *AesBase) Decrypt(rawData string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(rawData)
	if err != nil {
		return "", err
	}
	dnData, err := NewAes.aesCBCDecrypt(data)
	if err != nil {
		return "", err
	}
	return string(dnData), nil
}

func (aesBase *AesBase) EcbDecrypt(data []byte) []byte {
	block, _ := aes.NewCipher(aesBase.Key)
	decrypted := make([]byte, len(data))
	size := block.BlockSize()

	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		block.Decrypt(decrypted[bs:be], data[bs:be])
	}
	switch aesBase.Mode {
	case "PKCS7":
		decrypted = PKCS7UnPadding(decrypted)
		return decrypted
	case "Zero":
		decrypted = PKCS7UnPadding(decrypted)
		return decrypted
	default:
		panic("Unknown padding mode")
	}
}

func (aesBase *AesBase) EcbEncrypt(data []byte) []byte {
	block, _ := aes.NewCipher(aesBase.Key)
	switch aesBase.Mode {
	case "PKCS7":
		data = PKCS7Padding(data, block.BlockSize())
	case "Zero":
		data = PKCS7Padding(data, block.BlockSize())
	default:
		panic("Unknown padding mode")
	}

	encrypted := make([]byte, len(data))
	size := block.BlockSize()

	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		block.Encrypt(encrypted[bs:be], data[bs:be])
	}

	return encrypted
}
