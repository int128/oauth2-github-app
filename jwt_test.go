package oauth2githubapp

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"strings"
	"testing"
)

func Test_newAppJWT(t *testing.T) {
	key, err := LoadPrivateKey("testdata/oauth2-github-app-test.2021-02-12.private-key.pem")
	if err != nil {
		t.Fatalf("load error: %s", err)
	}
	token, err := newAppJWT("1234567890", key)
	if err != nil {
		t.Fatalf("encode error: %s", err)
	}

	// verify the token
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatalf("invalid token %s", token)
	}
	hashed := sha256.Sum256([]byte(parts[0] + "." + parts[1]))
	sign, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		t.Fatalf("invalid base64 in part[2] = %s", parts[2])
	}
	if err := rsa.VerifyPKCS1v15(&key.PublicKey, crypto.SHA256, hashed[:], sign); err != nil {
		t.Errorf("verify error: %s", err)
	}
}
