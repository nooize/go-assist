package assist

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
)

func ParsePemCertificate(bytes []byte) (*tls.Certificate, error) {
	var cert tls.Certificate
	pool := append(make([]byte, 0), bytes...)
	for {
		block, rest := pem.Decode(pool)
		if block == nil {
			break
		}
		if block.Type == "CERTIFICATE" {
			cert.Certificate = append(cert.Certificate, block.Bytes)
		} else {
			if pKey, err := ParseX509PrivateKey(block.Bytes); err != nil {
				return nil, fmt.Errorf("fail to parse private key : %s", err.Error())
			} else {
				cert.PrivateKey = pKey
			}
		}
		pool = rest
	}
	if len(cert.Certificate) == 0 {
		return nil, fmt.Errorf("no certificate found")
	} else if cert.PrivateKey == nil {
		return nil, fmt.Errorf("no private key found")
	}
	return &cert, nil
}

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
