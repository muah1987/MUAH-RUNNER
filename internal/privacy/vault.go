package privacy

import (
"encoding/json"
"fmt"
"os"
"path/filepath"
"time"
)

type Vault struct {
muahDir    string
passphrase string
}

func NewVault(muahDir, passphrase string) *Vault {
return &Vault{muahDir: muahDir, passphrase: passphrase}
}

func (v *Vault) vaultDir() string {
return filepath.Join(v.muahDir, "privacy_cot", "vault")
}

func (v *Vault) Store(runID string, entry *CoTEntry) error {
if err := os.MkdirAll(v.vaultDir(), 0700); err != nil {
return err
}
data, err := json.Marshal(entry)
if err != nil {
return err
}
encrypted, err := Encrypt(data, v.passphrase)
if err != nil {
return err
}
filename := fmt.Sprintf("cot-%s.enc", runID)
return os.WriteFile(filepath.Join(v.vaultDir(), filename), encrypted, 0600)
}

func (v *Vault) Retrieve(runID string) (*CoTEntry, error) {
filename := fmt.Sprintf("cot-%s.enc", runID)
data, err := os.ReadFile(filepath.Join(v.vaultDir(), filename))
if err != nil {
return nil, err
}
decrypted, err := Decrypt(data, v.passphrase)
if err != nil {
return nil, err
}
var entry CoTEntry
if err := json.Unmarshal(decrypted, &entry); err != nil {
return nil, err
}
return &entry, nil
}

func (v *Vault) NewEntry(runID string) *CoTEntry {
return &CoTEntry{
RunID:     runID,
Timestamp: time.Now(),
}
}
