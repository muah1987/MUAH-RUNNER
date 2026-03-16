package changelog

import (
"crypto/sha256"
"encoding/hex"
"encoding/json"
)

func HashEntry(e *Entry) string {
data, _ := json.Marshal(e)
sum := sha256.Sum256(data)
return hex.EncodeToString(sum[:])
}

func ChainHash(prevHash string, e *Entry) string {
combined := prevHash + HashEntry(e)
sum := sha256.Sum256([]byte(combined))
return hex.EncodeToString(sum[:])
}
