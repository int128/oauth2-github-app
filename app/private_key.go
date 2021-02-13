package app

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

// LoadPrivateKey reads the private key file of GitHub App
func LoadPrivateKey(name string) (*rsa.PrivateKey, error) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}
	return ParsePrivateKey(b)
}

// ParsePrivateKey parses the private key of GitHub App
func ParsePrivateKey(b []byte) (*rsa.PrivateKey, error) {
	d, _ := pem.Decode(b)
	if d == nil {
		return nil, fmt.Errorf("no pem block found")
	}
	k, err := x509.ParsePKCS1PrivateKey(d.Bytes)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}
	return k, nil
}
