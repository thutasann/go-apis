package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func encodeWithMD5(str string) string {
	var hashCode = md5.Sum([]byte(str))
	return hex.EncodeToString(hashCode[:])
}

func Encode_With_MD5() {
	var password string
	fmt.Scanln(&password)
	fmt.Println("Password encrypted to: ", encodeWithMD5(password))
}
