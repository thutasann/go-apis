package main

import (
	"fmt"

	"github.com/thutasann/cryptit/decrypt"
	"github.com/thutasann/cryptit/encrypt"
)

func main() {
	encryptedStr := encrypt.Nimbus("Kodekloud")
	fmt.Println(encryptedStr)
	fmt.Println(decrypt.Nimbus(encryptedStr))
}
