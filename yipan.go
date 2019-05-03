package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
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

var block cipher.Block

func LoadConfig() (config onedrive.Config) {
	enc, err := ioutil.ReadFile("config")
	if err != nil {
		panic(err)
	}

	raw := Decrypt(enc)

	c := &onedrive.Config{}
	err = json.Unmarshal(raw, c)
	if err != nil {
		panic(err)
	}
	return *c
}

func Decrypt(enc []byte) (raw []byte) {
	iv := enc[:aes.BlockSize]
	raw = make([]byte, len(enc)-aes.BlockSize)

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(raw, enc[aes.BlockSize:])
	return
}

func Encrypt(raw []byte) (enc []byte) {
	enc = make([]byte, aes.BlockSize+len(raw))
	iv := enc[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(enc[aes.BlockSize:], raw)
	return
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	hexKey := os.Getenv("hexkey")
	if len(hexKey) != 64 {
		panic(fmt.Sprintf("hex key length must be 64: %s", hexKey))
	}

	key, err := hex.DecodeString(hexKey)
	if err != nil {
		panic(err)
	}

	block, err = aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(os.Args) > 1 {
		if os.Args[1] == "enc" {
			raw, err := ioutil.ReadFile("config")
			if err != nil {
				panic(err)
			}

			enc := Encrypt(raw)
			err = ioutil.WriteFile("config", enc, 0644)
			if err != nil {
				panic(err)
			}
			return
		}
		return
	}

	config := LoadConfig()
	onedrive.SetConfig(config)

	//refresh token and save
	config, err = onedrive.Refresh()
	if err != nil {
		panic(err)
	}

	//log.Println(config)

	raw, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}

	enc := Encrypt(raw)
	err = ioutil.WriteFile("config", enc, 0644)
	if err != nil {
		panic(err)
	}

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
