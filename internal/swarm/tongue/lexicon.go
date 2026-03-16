package tongue

// VanadielTongue implements the inter-agent communication language
// Used by Elvaan (Opus) and Tarutaru agents for compressed communication

var Lexicon = map[string]string{
	// actions
	"AOE":  "analyze_and_execute",
	"RT":   "report_task",
	"RS":   "request_support",
	"CMP":  "complete",
	"FAIL": "task_failed",
	"KO":   "context_exhausted",
	// status
	"HP-LOW": "context_critically_low",
	"MP-LOW": "output_tokens_low",
	"FULL":   "resources_full",
	// task types
	"BUILD": "build_task",
	"TEST":  "test_task",
	"SCAN":  "scan_task",
	"DOC":   "documentation_task",
	"SEC":   "security_task",
}

var Grammar = map[string]string{
	"subject_verb_object": "FROM>VERB>TO",
	"status_report":       "FROM:STATUS:HP/MP",
	"task_assign":         "LEADER>TASK>AGENT",
}

func Encode(message string) string {
	encoded := ""
	for _, word := range splitWords(message) {
		if code, ok := reverseLookup(word); ok {
			encoded += code + " "
		} else {
			encoded += word + " "
		}
	}
	return encoded
}

func Decode(encoded string) string {
	decoded := ""
	for _, token := range splitWords(encoded) {
		if meaning, ok := Lexicon[token]; ok {
			decoded += meaning + " "
		} else {
			decoded += token + " "
		}
	}
	return decoded
}

func reverseLookup(meaning string) (string, bool) {
	for code, m := range Lexicon {
		if m == meaning {
			return code, true
		}
	}
	return "", false
}

func splitWords(s string) []string {
	var words []string
	word := ""
	for _, c := range s {
		if c == ' ' || c == '\t' || c == '\n' {
			if word != "" {
				words = append(words, word)
				word = ""
			}
		} else {
			word += string(c)
		}
	}
	if word != "" {
		words = append(words, word)
	}
	return words
}
