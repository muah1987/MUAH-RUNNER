package privacy

import (
"testing"
)

func TestEncryptDecrypt(t *testing.T) {
plaintext := []byte("secret thought: auth token abc123")
passphrase := "test-passphrase"
encrypted, err := Encrypt(plaintext, passphrase)
if err != nil {
t.Fatalf("encrypt error: %v", err)
}
decrypted, err := Decrypt(encrypted, passphrase)
if err != nil {
t.Fatalf("decrypt error: %v", err)
}
if string(decrypted) != string(plaintext) {
t.Errorf("decrypted mismatch: got %s", decrypted)
}
}

func TestEncryptWrongKey(t *testing.T) {
encrypted, _ := Encrypt([]byte("secret"), "right-key")
_, err := Decrypt(encrypted, "wrong-key")
if err == nil {
t.Error("expected error with wrong key")
}
}

func TestClassify(t *testing.T) {
tests := []struct {
thought string
want    Classification
}{
{"checking the auth token", Private},
{"user email address", Redactable},
{"running build script", Public},
}
for _, tt := range tests {
got := Classify(tt.thought)
if got != tt.want {
t.Errorf("Classify(%q) = %s, want %s", tt.thought, got, tt.want)
}
}
}

func TestVault(t *testing.T) {
dir := t.TempDir()
vault := NewVault(dir, "test-pass")
entry := vault.NewEntry("run-001")
entry.Thoughts = []Thought{
{ID: "t1", Content: "secret password check", Classification: Private},
}
if err := vault.Store("run-001", entry); err != nil {
t.Fatalf("store error: %v", err)
}
retrieved, err := vault.Retrieve("run-001")
if err != nil {
t.Fatalf("retrieve error: %v", err)
}
if len(retrieved.Thoughts) != 1 {
t.Errorf("expected 1 thought, got %d", len(retrieved.Thoughts))
}
}
