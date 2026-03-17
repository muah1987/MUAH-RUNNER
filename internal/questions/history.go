package questions

import (
"encoding/json"
"fmt"
)

func FormatHistory(answers []Answer) string {
data, _ := json.MarshalIndent(answers, "", "  ")
return fmt.Sprintf("Question History (%d entries):\n%s", len(answers), string(data))
}
