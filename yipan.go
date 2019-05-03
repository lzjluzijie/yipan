package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"

	"./onedrive"
)

func LoadConfig() (config onedrive.Config) {
	hexKey := os.Getenv("hexkey")
	if len(hexKey) != 64 {
		panic(fmt.Sprintf("hex key length must be 64: %s", hexKey))
	}

	key, err := hex.DecodeString(hexKey)
	if err != nil {
		panic(err)
	}

	encrypted, err := ioutil.ReadFile("config")
	if err != nil {
		panic(err)
	}

	ciphertext, err := base64.StdEncoding.DecodeString(string(encrypted))
	if err != nil {
		panic(err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	c := &onedrive.Config{}
	err = json.Unmarshal(ciphertext, c)
	if err != nil {
		panic(err)
	}

	return *c
}

func Encrypt() {
	hexKey := os.Getenv("hexkey")
	if len(hexKey) != 64 {
		panic(fmt.Sprintf("hex key length must be 64: %s", hexKey))
	}

	key, err := hex.DecodeString(hexKey)
	if err != nil {
		panic(err)
	}

	plaintext, err := ioutil.ReadFile("config")
	if err != nil {
		panic(err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	err = ioutil.WriteFile("config", []byte(base64.StdEncoding.EncodeToString(ciphertext)), 0777)

}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if len(os.Args) > 1 {
		if os.Args[1] == "enc" {
			Encrypt()
			return
		}
		return
	}

	config := LoadConfig()
	onedrive.SetConfig(config)

	files, err := onedrive.ListChildren("root", "")
	if err != nil {
		log.Println(err.Error())
	}

	redirects, err := os.Create("_redirects")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		_, err = redirects.WriteString(fmt.Sprintf("%s %s\r\n", strings.Replace(url.PathEscape(file.Path), "%2F", "/", -1), file.URL))
		if err != nil {
			panic(err)
		}
	}
}
