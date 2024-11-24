package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
)

func SavePublicKey(userID, publicKey string) error {
	dir := "./keys"
	os.MkdirAll(dir, os.ModePerm)
	filePath := filepath.Join(dir, fmt.Sprintf("%s.pem", userID))
	return os.WriteFile(filePath, []byte(publicKey), 0644)
}

func LoadPublicKey(userID string) (*rsa.PublicKey, error) {
	filePath := fmt.Sprintf("./keys/%s.pem", userID)
	publicKeyData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(publicKeyData)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, fmt.Errorf("Failed to decode PEM block containing public key")
	}

	return x509.ParsePKCS1PublicKey(block.Bytes)
}

func GenerateAndEncryptSymmetricKey(userID string) (string, error) {
	publicKey, err := LoadPublicKey(userID)
	if err != nil {
		return "", fmt.Errorf("Failed to load public key: %v", err)
	}

	symmetricKey := make([]byte, 32)
	if _, err := rand.Read(symmetricKey); err != nil {
		return "", fmt.Errorf("Failed to generate symmetric key: %v", err)
	}

	encryptedKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, symmetricKey, nil)
	if err != nil {
		return "", fmt.Errorf("Failed to encrypt symmetric key: %v", err)
	}

	return base64.StdEncoding.EncodeToString(encryptedKey), nil
}
