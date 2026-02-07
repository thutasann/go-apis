package decrypt

import "strings"

func Nimbus(str string) string {
	var decryptedStr strings.Builder
	for _, c := range str {
		decryptedStr.WriteString(string(rune(c - 3)))
	}
	return decryptedStr.String()
}
