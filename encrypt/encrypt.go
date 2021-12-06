package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	rand1"math/rand"
)


//生成随机字符串
func GetRandomString(length int) string{
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand1.New(rand1.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
// 生成32位MD5
func MD5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}
func wr(str string) {
	f, err := os.Create("aeskey.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = f.WriteString(str)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
}

func main() {
	if len(os.Args) != 2{
		fmt.Println(os.Args[0]+" filename")
	}
	filename := os.Args[1]
	log.Print("File encryption example")
	plaintext, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	key0 := MD5(GetRandomString(16))

	wr(key0)

	// The key should be 16 bytes (AES-128), 24 bytes (AES-192) or
	// 32 bytes (AES-256)

	key, err := ioutil.ReadFile("aeskey.txt")

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Panic(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Panic(err)
	}

	// Never use more than 2^32 random nonces with a given key
	// because of the risk of repeat.
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatal(err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	// Save back to file
	err = ioutil.WriteFile(filename+".cipher", ciphertext, 0777)
	if err != nil {
		log.Panic(err)
	}
}