package assist

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"golang.org/x/crypto/ssh"
)

func ParseX509PrivateKey(pem []byte) (crypto.PrivateKey, error) {
	//pem = []byte(strings.Replace(string(pem), "\n", "", -1))
	if key, err := x509.ParsePKCS1PrivateKey(pem); err == nil {
		return key, nil
	}
	if key, err := x509.ParsePKCS8PrivateKey(pem); err == nil {
		switch key := key.(type) {
		case *rsa.PrivateKey, *ecdsa.PrivateKey, ed25519.PrivateKey:
			return key, nil
		default:
			return nil, errors.New("tls: found unknown private key type in PKCS#8 wrapping")
		}
	}
	if key, err := x509.ParseECPrivateKey(pem); err == nil {
		return key, nil
	}
	if key, err := ssh.ParseRawPrivateKey(pem); err == nil {
		return key, nil
	}
	return nil, errors.New("tls: failed to parse private key")
}

func ParseX509PublicKey(raw []byte) (crypto.PublicKey, error) {
	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the public key")
	}

	if key, err := x509.ParsePKIXPublicKey(block.Bytes); err == nil {
		return key, nil
	}
	return nil, errors.New("tls: failed to parse public key")
}
