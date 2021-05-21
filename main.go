package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
)

const (
	keyFile       = "aes.key"
	encryptedfile = "aes.enc"
)

var abc = []byte("1234567812345678")

func rKey(filename string) ([]byte, error) {
	key, err := ioutil.ReadFile(filename)
	if err != nil {
		return key, err
	}
	block, _ := pem.Decode(key)
	return block.Bytes, nil
}

func cKey() []byte {
	genkey := make([]byte, 16)
	_, err := rand.Read(genkey)
	if err != nil {
		log.Fatalf("failed to read key: %s", err)
	}
	return genkey
}

func sKey(filename string, key []byte) {
	block := &pem.Block{
		Type:  "AES KEY",
		Bytes: key,
	}
	err := ioutil.WriteFile(filename, pem.EncodeToMemory(block), 9854)
	if err != nil {
		log.Fatalf("Failed tio save the key %s: %s", filename, err)
	}
}

func aesKey() []byte {
	file := fmt.Sprintf(keyFile)
	key, err := rKey(file)
	if err != nil {
		log.Println("Create a new AES KEY")
		key = cKey()
		sKey(file, key)
	}
	return key
}

func createCipher() cipher.Block {
	c, err := aes.NewCipher(aesKey())
	if err != nil {
		log.Fatalf("failed to create aes  %s", err)
	}
	return c
}

func encryption(plainText string) {
	bytes := []byte(plainText)
	blockCipher := createCipher()
	stream := cipher.NewCTR(blockCipher, abc)
	stream.XORKeyStream(bytes, bytes)
	err := ioutil.WriteFile(fmt.Sprintf(encryptedfile), bytes, 0644)
	if err != nil {
		log.Fatalf("writing encryption file %s", err)
	} else {
		fmt.Printf("Message encrypted: %s\n\n", encryptedfile)
	}
}
func decryption() []byte {
	bytes, err := ioutil.ReadFile(fmt.Sprintf(encryptedfile))
	if err != nil {
		log.Fatalf("Reading encrypted file %s", err)
	}
	blockCipher := createCipher()
	stream := cipher.NewCTR(blockCipher, abc)
	stream.XORKeyStream(bytes, bytes)
	return bytes
}

func main() {
	var plainText = "AES is now used Worldwide"
	encryption(plainText)
	fmt.Printf("Decrypted Msg : %s", decryption())
}
