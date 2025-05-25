package oauth2githubapp

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

func newAppJWT(appID string, key *rsa.PrivateKey) (string, error) {
	h := headerType{
		Algorithm: "RS256",
		Type:      "JWT",
	}
	c := claimsType{
		Iat: time.Now().Unix(),
		Exp: time.Now().Add(3 * time.Minute).Unix(),
		Iss: appID,
	}
	encoded, err := encodeRS256(h, c, key)
	if err != nil {
		return "", fmt.Errorf("jwt encode error: %w", err)
	}
	return encoded, nil
}

type headerType struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

type claimsType struct {
	Iss string `json:"iss"` // GitHub App ID
	Exp int64  `json:"exp"` // the expiration time of the assertion (seconds since Unix epoch)
	Iat int64  `json:"iat"` // the time the assertion was issued (seconds since Unix epoch)
}

func encodeRS256(h headerType, c claimsType, key *rsa.PrivateKey) (string, error) {
	header, err := encodeBase64JSON(&h)
	if err != nil {
		return "", fmt.Errorf("encode header error: %w", err)
	}
	payload, err := encodeBase64JSON(&c)
	if err != nil {
		return "", fmt.Errorf("encode claims error: %w", err)
	}
	sig, err := signPKCS1v15(header+"."+payload, crypto.SHA256, key)
	if err != nil {
		return "", fmt.Errorf("sign error: %w", err)
	}
	return header + "." + payload + "." + sig, nil
}

func encodeBase64JSON(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("marshal error: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func signPKCS1v15(s string, hash crypto.Hash, key *rsa.PrivateKey) (string, error) {
	h := hash.New()
	_, _ = h.Write([]byte(s))
	sig, err := rsa.SignPKCS1v15(rand.Reader, key, hash, h.Sum(nil))
	if err != nil {
		return "", fmt.Errorf("pkcs1 v1.5 error: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(sig), nil
}
