package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"sync"
)

// TODO: 密钥应从环境变量或密钥管理服务加载，当前硬编码仅用于开发环境
var (
	secretKey     = []byte("hr-onboard-aesgcm-key-32byte!!!!!")
	keyMutex      sync.RWMutex
)

func SetKey(key []byte) error {
	if len(key) != 32 {
		return errors.New("AES-256 密钥长度必须为 32 字节")
	}
	keyMutex.Lock()
	defer keyMutex.Unlock()
	secretKey = make([]byte, 32)
	copy(secretKey, key)
	return nil
}

func getKey() []byte {
	keyMutex.RLock()
	defer keyMutex.RUnlock()
	k := make([]byte, len(secretKey))
	copy(k, secretKey)
	return k
}

func Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}
	key := getKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(encoded string) (string, error) {
	if encoded == "" {
		return "", nil
	}
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	key := getKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("密文长度不足")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
