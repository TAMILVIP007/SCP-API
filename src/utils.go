package src

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func encryptAES(plaintext string) string {
	encyptkey := []byte(Envars.Encyptkey)
	block, _ := aes.NewCipher(encyptkey)
	padding := aes.BlockSize - len(plaintext)%aes.BlockSize
	plaintext += fmt.Sprintf("%c", padding)
	ciphertext := make([]byte, len(plaintext))
	iv := make([]byte, aes.BlockSize)
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext, []byte(plaintext))
	return string(ciphertext)
}

func decryptAES(ciphertext []byte) string {
	encyptkey := []byte(Envars.Encyptkey)
	block, err := aes.NewCipher(encyptkey)
	if err != nil {
		return ""
	}
	iv := make([]byte, aes.BlockSize)
	stream := cipher.NewCTR(block, iv)
	plaintext := make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)
	return string(plaintext)
}

func MatchKey(key string, rights []string) bool {
	for _, v := range rights {
		if v == key {
			return true
		}
	}
	return false
}

func GenToken(userID string, role string) string {
	return encryptAES(userID + ":" + role)
}
