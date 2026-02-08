package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/big"
)

func GenerateKeyPair() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

func hashMessage(msg []byte) []byte {
	hash := sha256.Sum256(msg)
	return hash[:]
}

func SignMessage(priv *ecdsa.PrivateKey, msg []byte) (string, error) {
	hashed := hashMessage(msg)

	r, s, err := ecdsa.Sign(rand.Reader, priv, hashed)
	if err != nil {
		return "", err
	}

	rBytes := r.Bytes()
	sBytes := s.Bytes()

	signature := append(rBytes, sBytes...)
	return hex.EncodeToString(signature), nil
}

func EncodePublicKey(pub *ecdsa.PublicKey) string {
	xBytes := pub.X.Bytes()
	yBytes := pub.Y.Bytes()

	combined := append(xBytes, yBytes...)
	return hex.EncodeToString(combined)
}

func DecodePublicKey(hexKey string) (*ecdsa.PublicKey, error) {
	bytes, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, err
	}

	if len(bytes)%2 != 0 {
		return nil, errors.New("invalid public key length")
	}

	mid := len(bytes) / 2

	x := new(big.Int).SetBytes(bytes[:mid])
	y := new(big.Int).SetBytes(bytes[mid:])

	return &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}, nil
}

func VerifySignature(data []byte, signature, pubKeyHex string) (bool, error) {
	hashed := hashMessage(data)

	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		return false, err
	}

	if len(sigBytes)%2 != 0 {
		return false, errors.New("invalid signature length")
	}

	mid := len(sigBytes) / 2

	r := new(big.Int).SetBytes(sigBytes[:mid])
	s := new(big.Int).SetBytes(sigBytes[mid:])

	pub, err := DecodePublicKey(pubKeyHex)
	if err != nil {
		return false, err
	}

	return ecdsa.Verify(pub, hashed, r, s), nil
}
