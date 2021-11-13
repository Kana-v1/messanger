package symmetricCrypto

import (
	"crypto/rsa"
	"crypto/x509"
	"messanger/internal/logs"

	"github.com/pkg/errors"
)

func DecodePublicKey(publicKey *rsa.PublicKey) []byte {
	return x509.MarshalPKCS1PublicKey(publicKey)
}

func EncodePublicKey(publicKey []byte) (*rsa.PublicKey, error) {
	key, err := x509.ParsePKCS1PublicKey(publicKey)
	if err != nil {
		return nil, errors.Wrap(err, "Can not decode public key")
	}
	return key, nil
}

func DecodePrivateKey(privateKey *rsa.PrivateKey) []byte {
	return x509.MarshalPKCS1PrivateKey(privateKey)
}

func EncodePrivateKey(privateKey []byte) (*rsa.PrivateKey, error) {
	key, err := x509.ParsePKCS1PrivateKey(privateKey)
	if err != nil {
		return nil, errors.Wrap(err, "Can not decode private key")
	}
	return key, nil
}

func GetPublicKeyFromPrivateKey(privateKey []byte) []byte {
	encodedPrivateKey, err := EncodePrivateKey(privateKey)
	if err != nil {
		logs.ErrorLog("cryptoKeys.log", "Can not encode private key from []byte", err)
		return make([]byte, 0)
	}
	return DecodePublicKey(&encodedPrivateKey.PublicKey)
}
