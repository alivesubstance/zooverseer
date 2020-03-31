package core

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"github.com/alivesubstance/zooverseer/util"
)

const salt = "8d84b9363adf51458a3e67672176bcfd"

func Encrypt(passphrase string) string {
	block, err := aes.NewCipher([]byte(createHash(salt)))
	util.CheckError(err)

	gcm, err := cipher.NewGCM(block)
	util.CheckError(err)

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	util.CheckError(err)

	ciphertext := gcm.Seal(nonce, nonce, []byte(passphrase), nil)

	return hex.EncodeToString(ciphertext)
}

func Decrypt(cipherText string) string {
	block, err := aes.NewCipher([]byte(createHash(salt)))
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	cipherBytes, err := hex.DecodeString(cipherText)
	util.CheckError(err)

	plaintext, err := gcm.Open(
		nil,
		cipherBytes[:gcm.NonceSize()],
		cipherBytes[gcm.NonceSize():],
		nil,
	)
	if err != nil {
		panic(err.Error())
	}
	return string(plaintext)
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
