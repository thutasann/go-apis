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
	payload := "Foo not bar"
	src := bytes.NewReader([]byte(payload))
	dst := new(bytes.Buffer)
	key := NewEncryptionKey()
	_, err := CopyEncrypt(key, src, dst)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("payload length ---> ", len(payload))
	fmt.Println("dest string ---> ", dst.String())

	out := new(bytes.Buffer)
	nw, err := CopyDecrypt(key, dst, out)
	if err != nil {
		t.Error(err)
	}

	if nw != 16+len(payload) {
		t.Fail()
	}

	if out.String() != payload {
		t.Errorf("encryption failed")
	}

	fmt.Println("out string length ---> ", len(out.String()))
	fmt.Println("out string ---> ", out.String())
}
