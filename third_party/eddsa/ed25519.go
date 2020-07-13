package eddsa

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"

	"github.com/dgrijalva/jwt-go"
)

type SigningMethodEd25519 struct{}

var (
	SigningMethodEdDSA *SigningMethodEd25519
)

func init() {
	SigningMethodEdDSA = &SigningMethodEd25519{}
	jwt.RegisterSigningMethod(SigningMethodEdDSA.Alg(), func() jwt.SigningMethod {
		return SigningMethodEdDSA
	})
}

func (m *SigningMethodEd25519) Alg() string {
	return "EdDSA"
}

func (m *SigningMethodEd25519) Verify(signingString, signature string, key interface{}) error {
	var err error
	var ed25519Key *ed25519.PublicKey
	var ok bool

	if ed25519Key, ok = key.(*ed25519.PublicKey); !ok {
		return jwt.ErrInvalidKeyType
	}

	if len(*ed25519Key) != ed25519.PublicKeySize {
		return jwt.ErrInvalidKey
	}

	var sig []byte
	if sig, err = jwt.DecodeSegment(signature); err != nil {
		return err
	}

	if !ed25519.Verify(*ed25519Key, []byte(signingString), sig) {
		return jwt.ErrSignatureInvalid
	}

	return nil
}

func (m *SigningMethodEd25519) Sign(signingString string, key interface{}) (string, error) {
	var ed25519Key *ed25519.PrivateKey
	var ok bool

	if ed25519Key, ok = key.(*ed25519.PrivateKey); !ok {
		return "", jwt.ErrInvalidKeyType
	}

	// ed25519.Sign panics if private key not equal to ed25519.PrivateKeySize
	// this allows to avoid recover usage
	if len(*ed25519Key) != ed25519.PrivateKeySize {
		return "", jwt.ErrInvalidKey
	}

	sig := ed25519.Sign(*ed25519Key, []byte(signingString))
	return jwt.EncodeSegment(sig), nil
}

func ParseEdPublicKeyFromPEM(key []byte) (*ed25519.PublicKey, error) {
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, jwt.ErrKeyMustBePEMEncoded
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	epub, ok := pub.(ed25519.PublicKey)
	if !ok {
		return nil, jwt.ErrInvalidKeyType
	}
	return &epub, nil
}

func ParseEdPrivateKeyFromPEM(key []byte) (*ed25519.PrivateKey, error) {
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, jwt.ErrKeyMustBePEMEncoded
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	epriv, ok := priv.(ed25519.PrivateKey)
	if !ok {
		return nil, jwt.ErrInvalidKeyType
	}
	return &epriv, nil
}
