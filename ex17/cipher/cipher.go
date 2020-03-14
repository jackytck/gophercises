package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

// Encrypt will tak in a key and plaintext and return a hex representation
// of the encrypted value.
func Encrypt(key, plaintext string) (string, error) {
	block, err := newCipherBlock(key)
	if err != nil {
		return "", err
	}

	plainByte := []byte(plaintext)
	ciphertext := make([]byte, aes.BlockSize+len(plainByte))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plainByte)

	return fmt.Sprintf("%x\n", ciphertext), nil
}

// Decrypt will take in a key and a cipherHex (hex representation of
// the ciphertext) and decrypt it.
func Decrypt(key, cipherHex string) (string, error) {
	block, err := newCipherBlock(key)
	if err != nil {
		return "", err
	}
	ciphertext, _ := hex.DecodeString(cipherHex)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("encrypt: cipher too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

func newCipherBlock(key string) (cipher.Block, error) {
	hasher := md5.New()
	fmt.Fprint(hasher, key)
	cipherKey := hasher.Sum(nil)
	return aes.NewCipher(cipherKey)
}
