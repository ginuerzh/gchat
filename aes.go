// aes
package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

const key = "abcdefghijklmnopkrstuvwsyz012345"

var (
	block cipher.Block
)

func init() {
	cblock, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}

	block = cblock
}

func Encrypt(src string) ([]byte, error) {
	plaintext := []byte(src)
	if len(plaintext)%aes.BlockSize != 0 {
		b := make([]byte, aes.BlockSize-len(plaintext)%aes.BlockSize+len(plaintext)) // padding
		copy(b, plaintext)
		plaintext = b

		//return nil, errors.New("plaintext is not a multiple of the block size")
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)
	return ciphertext, nil
}

func Decrypt(hexStr string) ([]byte, error) {
	ciphertext, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	// trim suffix padding byte 0
	if index := bytes.IndexByte(ciphertext, 0); index > 0 {
		ciphertext = ciphertext[:index]
	}

	return ciphertext, nil
}
