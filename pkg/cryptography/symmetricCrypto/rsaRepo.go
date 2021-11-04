package symmetricCrypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"messanger/internal/logs"

	"github.com/pkg/errors"
)

type CryptoKeys struct {
	*rsa.PrivateKey
}

func GenerateKeys() *CryptoKeys {
	keys, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		logs.ErrorLog("cryptoError.log", "Can not generate crypro key", err)
		return nil
	}
	return &CryptoKeys{keys}
}

func EncryptMessage(msg []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	encryptedBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, msg, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Can not encrypt message")
	}
	return encryptedBytes, nil
}

func DecryptMessage(msg []byte, privateKey *CryptoKeys) ([]byte, error) {
	decrytpedMessage, err := privateKey.Decrypt(nil, msg, &rsa.OAEPOptions{Hash: crypto.SHA256})
	if err != nil {
		return nil, errors.Wrap(err, "Can not decrypt message")
	}
	return decrytpedMessage, nil
}
