// aes_test.go
package main

import (
	"fmt"
	"testing"
)

func TestAES(t *testing.T) {
	s := "exampleplaintexts"
	ciphertext, err := Encrypt(s)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("%x", ciphertext))

	text, err := Decrypt(fmt.Sprintf("%x", ciphertext))
	if err != nil {
		t.Fatal(err)
	}
	s = string(text)
	t.Log(s, len(s))
}
