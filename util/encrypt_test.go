package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	pass := "123"
	encrypted, _ := Encrypt(pass)
	assert.Equal(t, encrypted, encrypted)

	againEncrypt, _ := Encrypt(encrypted)
	assert.Equal(t, encrypted, againEncrypt)

	decrypt, _ := Decrypt(encrypted)
	assert.Equal(t, pass, decrypt)
}
