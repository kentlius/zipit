package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

func decryptURL(secureString string, key []byte) (string, error) {
	data, err := base64.RawURLEncoding.DecodeString(secureString)
	if err != nil {
		return "", err
	}

	if len(data) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := data[:aes.BlockSize]
	ciphertext := data[aes.BlockSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return "", fmt.Errorf("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	padding := int(ciphertext[len(ciphertext)-1])
	if padding > len(ciphertext) || padding == 0 {
		return "", fmt.Errorf("invalid padding")
	}

	return string(ciphertext[:len(ciphertext)-padding]), nil
}
