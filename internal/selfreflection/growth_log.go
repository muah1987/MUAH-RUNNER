package selfreflection

import (
"fmt"
"os"
"path/filepath"
)

func AppendGrowthLog(muahDir, entry string) error {
path := filepath.Join(muahDir, "memory", "selfreflection", "growth-log.md")
os.MkdirAll(filepath.Dir(path), 0755)
f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
if err != nil {
return err
}
defer f.Close()
fmt.Fprintln(f, entry)
return nil
}
