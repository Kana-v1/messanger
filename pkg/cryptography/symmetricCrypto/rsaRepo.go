package symmetricCrypto

import (
	"crypto/rand"
	"crypto/rsa"
	"messanger/internal/logs"
	"messanger/pkg/cryptography/hash"

	"github.com/pkg/errors"
)

func GenerateKeys() *rsa.PrivateKey {
	keys, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		logs.ErrorLog("cryptoError.log", "Can not generate crypro key", err)
		return nil
	}
	return keys
}

func EncryptMessage(msg []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	encryptedBytes, err := rsa.EncryptOAEP(hash.DefaultHasheAlgorithm.New(), rand.Reader, publicKey, msg, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Can not encrypt message")
	}
	return encryptedBytes, nil
}

func DecryptMessage(msg []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	decrytpedMessage, err := privateKey.Decrypt(nil, msg, &rsa.OAEPOptions{Hash: hash.DefaultHasheAlgorithm})
	if err != nil {
		return nil, errors.Wrap(err, "Can not decrypt message")
	}
	return decrytpedMessage, nil
}
