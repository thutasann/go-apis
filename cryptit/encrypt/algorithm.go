package encrypt

import "strings"

func Nimbus(str string) string {
	var encryptedStr strings.Builder
	for _, c := range str {
		encryptedStr.WriteString(string(rune(c + 3)))
	}
	return encryptedStr.String()
}
