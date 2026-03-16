package memory

// MergeContextLoaded merges two slices of context keys deduplicating.
func MergeContextLoaded(a, b []string) []string {
seen := make(map[string]bool)
var result []string
for _, s := range append(a, b...) {
if !seen[s] {
seen[s] = true
result = append(result, s)
}
}
return result
}

// MergeLessons merges lessons deduplicating.
func MergeLessons(a, b []string) []string {
return MergeContextLoaded(a, b)
}
