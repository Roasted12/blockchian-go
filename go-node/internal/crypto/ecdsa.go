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

//
// -----------------------------
// KEY GENERATION
// -----------------------------
//

// GenerateKeyPair creates a new ECDSA private/public key pair
// using the P-256 elliptic curve.
func GenerateKeyPair() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

//
// -----------------------------
// INTERNAL HASHING HELPER
// -----------------------------
//

// hashMessage hashes arbitrary bytes using SHA-256.
// ECDSA always signs hashes, never raw data.
func hashMessage(msg []byte) []byte {
	hash := sha256.Sum256(msg)
	return hash[:]
}

//
// -----------------------------
// SIGNING
// -----------------------------
//

// SignMessage signs canonical transaction bytes using a private key.
// Returns a hex-encoded signature (r || s).
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

//
// -----------------------------
// PUBLIC KEY ENCODING
// -----------------------------
//

// EncodePublicKey serializes an ECDSA public key into hex format.
func EncodePublicKey(pub *ecdsa.PublicKey) string {
	xBytes := pub.X.Bytes()
	yBytes := pub.Y.Bytes()

	combined := append(xBytes, yBytes...)
	return hex.EncodeToString(combined)
}

//
// DecodePublicKey deserializes a hex-encoded public key
// back into an ECDSA public key structure.
//
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

//
// -----------------------------
// SIGNATURE VERIFICATION
// -----------------------------
//

// VerifySignature verifies that data was signed by the owner of the provided public key.
//
// Parameters:
// - data: The canonical bytes that were signed (must match exactly what was signed)
// - signature: Hex-encoded signature (r || s)
// - pubKeyHex: Hex-encoded public key
//
// Returns:
// - true if signature is valid
// - false if signature is invalid
//
// This function avoids importing the chain package, breaking the import cycle.
// The chain package computes canonical bytes and calls this function.
//
func VerifySignature(data []byte, signature, pubKeyHex string) (bool, error) {
	// Hash the data (ECDSA signs hashes, not raw data)
	hashed := hashMessage(data)

	// Decode signature
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

	// Decode public key
	pub, err := DecodePublicKey(pubKeyHex)
	if err != nil {
		return false, err
	}

	// Verify signature
	return ecdsa.Verify(pub, hashed, r, s), nil
}
