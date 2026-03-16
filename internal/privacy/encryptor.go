package privacy

import (
"crypto/aes"
"crypto/cipher"
"crypto/rand"
"errors"
"io"

"golang.org/x/crypto/scrypt"
)

const (
saltLen = 32 // bytes — random salt stored with the ciphertext
scryptN = 32768
scryptR = 8
scryptP = 1
keyLen  = 32 // AES-256
)

// deriveKey uses scrypt to derive a 32-byte AES key from passphrase and salt.
// Using scrypt (instead of a bare SHA-256) makes brute-force attacks
// computationally expensive: each attempt requires significant memory and CPU.
func deriveKey(passphrase string, salt []byte) ([]byte, error) {
return scrypt.Key([]byte(passphrase), salt, scryptN, scryptR, scryptP, keyLen)
}

// Encrypt encrypts plaintext with AES-256-GCM.
// The output layout is: salt (32 B) | nonce (12 B) | ciphertext + GCM tag.
// A fresh random salt is generated for every call so that the same passphrase
// and plaintext always produce a different ciphertext.
func Encrypt(plaintext []byte, passphrase string) ([]byte, error) {
salt := make([]byte, saltLen)
if _, err := io.ReadFull(rand.Reader, salt); err != nil {
return nil, err
}
key, err := deriveKey(passphrase, salt)
if err != nil {
return nil, err
}
block, err := aes.NewCipher(key)
if err != nil {
return nil, err
}
gcm, err := cipher.NewGCM(block)
if err != nil {
return nil, err
}
nonce := make([]byte, gcm.NonceSize())
if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
return nil, err
}
out := make([]byte, 0, saltLen+len(nonce)+len(plaintext)+gcm.Overhead())
out = append(out, salt...)
out = append(out, nonce...)
out = gcm.Seal(out, nonce, plaintext, nil)
return out, nil
}

// Decrypt reverses Encrypt.  It expects the salt | nonce | ciphertext layout
// produced by Encrypt.
func Decrypt(ciphertext []byte, passphrase string) ([]byte, error) {
if len(ciphertext) < saltLen {
return nil, errors.New("ciphertext too short")
}
salt, rest := ciphertext[:saltLen], ciphertext[saltLen:]
key, err := deriveKey(passphrase, salt)
if err != nil {
return nil, err
}
block, err := aes.NewCipher(key)
if err != nil {
return nil, err
}
gcm, err := cipher.NewGCM(block)
if err != nil {
return nil, err
}
nonceSize := gcm.NonceSize()
if len(rest) < nonceSize {
return nil, errors.New("ciphertext too short")
}
nonce, ct := rest[:nonceSize], rest[nonceSize:]
return gcm.Open(nil, nonce, ct, nil)
}
