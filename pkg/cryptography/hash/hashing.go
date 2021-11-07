package hash

import (
	"crypto"
)

const DefaultHasheAlgorithm = crypto.SHA512

func Hash(msg []byte) []byte {
	hashedMessage := DefaultHasheAlgorithm.New()
	hashedMessage.Write(msg)
	return hashedMessage.Sum(msg)
}
