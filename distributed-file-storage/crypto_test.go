package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestNewEncryptionKey(t *testing.T) {
	for i := 0; i < 1000; i++ {
		key := NewEncryptionKey()
		fmt.Println(key)
		for i := 0; i < len(key); i++ {
			if key[i] == 0 {
				t.Errorf("0 bytes")
			}
		}
	}
}

func TestCopyEncryptDecrypt(t *testing.T) {
	src := bytes.NewReader([]byte("Foo not bar"))
	dst := new(bytes.Buffer)
	key := NewEncryptionKey()
	_, err := CopyEncrypt(key, src, dst)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("dest string ---> ", dst.String())

	out := new(bytes.Buffer)
	if _, err := CopyDecrypt(key, dst, out); err != nil {
		t.Error(err)
	}

	fmt.Println("out string ---> ", out.String())
}
