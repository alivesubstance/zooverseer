package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const salt = "8d84b9363adf51458a3e67672176bcfd"

func Encrypt(passphrase string) (string, error) {
	if passphrase == "" {
		return "", nil
	}

	decrypted, err := Decrypt(passphrase)
	if err == nil && decrypted != passphrase {
		// passphrase already encrypted. do not encrypt it twice
		return passphrase, nil
	}

	gcm, err := createGCM()
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return "", errors.Wrap(err, "Failed to read nonce")
	}

	cipherText := gcm.Seal(nonce, nonce, []byte(passphrase), nil)

	return hex.EncodeToString(cipherText), nil
}

func DecryptOrPanic(cipherText string) string {
	decrypt, err := Decrypt(cipherText)
	if err != nil {
		log.WithError(err).Panic()
	}

	return decrypt
}

func Decrypt(cipherText string) (string, error) {
	if cipherText == "" {
		return "", nil
	}

	cipherBytes, err := hex.DecodeString(cipherText)
	if err != nil {
		// cipherText have to contains only hexadecimal characters
		// and has even length. otherwise, it failed to decode
		// return original string in this case
		return cipherText, nil
	}

	gcm, err := createGCM()
	if err != nil {
		return "", err
	}

	plaintext, err := gcm.Open(
		nil,
		cipherBytes[:gcm.NonceSize()],
		cipherBytes[gcm.NonceSize():],
		nil,
	)
	if err != nil {
		return "", errors.Wrap(err, "Failed to decrypt string")
	}

	return string(plaintext), nil
}

func createGCM() (cipher.AEAD, error) {
	block, err := aes.NewCipher([]byte(createHash(salt)))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create Cipher")
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create GCM")
	}
	return gcm, nil
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
