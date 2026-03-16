package changelog

import (
"crypto/sha256"
"encoding/hex"
"encoding/json"
)

// HashEntry returns a SHA-256 hash of the entry's content.
// The entry's own Hash field is zeroed before hashing so that the result is
// stable whether the entry has already been assigned a hash or not.  This
// makes independent verification of stored entries possible: just call
// HashEntry on a retrieved entry and compare with the stored Hash.
func HashEntry(e *Entry) string {
// Work on a shallow copy so we never mutate the caller's entry.
tmp := *e
tmp.Hash = ""
data, _ := json.Marshal(tmp)
sum := sha256.Sum256(data)
return hex.EncodeToString(sum[:])
}

// ChainHash combines prevHash with the content hash of e to produce the
// tamper-evident chain link.  Because HashEntry zeroes the Hash field,
// ChainHash produces the same value regardless of whether e.Hash is already
// populated or empty.
func ChainHash(prevHash string, e *Entry) string {
combined := prevHash + HashEntry(e)
sum := sha256.Sum256([]byte(combined))
return hex.EncodeToString(sum[:])
}
