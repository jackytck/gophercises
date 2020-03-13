package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"
)

// Encrypt encrypts a string with key.
func Encrypt(key, plaintext string) (string, error) {
	hasher := md5.New()
	fmt.Fprint(hasher, key)
	cipherKey := hasher.Sum(nil)
	block, err := aes.NewCipher(cipherKey)
	if err != nil {
		return "", err
	}

	plainByte := []byte(plaintext)
	ciphertext := make([]byte, aes.BlockSize+len(plainByte))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plainByte)

	return fmt.Sprintf("%x\n", ciphertext), nil
}
