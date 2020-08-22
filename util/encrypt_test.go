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

func TestDecrypt(t *testing.T) {
	encrypted := "c14502b09fd94803793b34167e9ad8f9c14e505395bcf85a4131e7259a549f9658039233a3"
	decrypted, err := Decrypt(encrypted)
	assert.Nil(t, err)
	assert.Equal(t, "z00k33p3r", decrypted)
}
