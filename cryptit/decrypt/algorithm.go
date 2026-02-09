// decrypt package consists of all the decryption algorithms
package decrypt

import "strings"

// decrypt by reducing the ascii code by 3 for each character
func Nimbus(str string) string {
	var decryptedStr strings.Builder
	for _, c := range str {
		decryptedStr.WriteString(string(rune(c - 3)))
	}
	return decryptedStr.String()
}
