package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	pass := "123"
	encrypted, err := Encrypt(pass)
	assert.Nil(t, err)
	assert.Equal(t, encrypted, encrypted)

	againEncrypt, err := Encrypt(encrypted)
	assert.Nil(t, err)
	assert.Equal(t, encrypted, againEncrypt)

	decrypt, err := Decrypt(encrypted)
	assert.Nil(t, err)
	assert.Equal(t, pass, decrypt)
}

func TestDecryptNotEncryptedString(t *testing.T) {
	pass := "123abcABC"
	decrypt, err := Decrypt(pass)
	assert.Nil(t, err)
	assert.Equal(t, pass, decrypt)
}
