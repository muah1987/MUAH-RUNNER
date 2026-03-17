package sensor

type MCPScanResult struct {
Available []string
Missing   []string
}

func ScanMCPs() MCPScanResult {
result := MCPScanResult{}
tc := DetectToolchain()
if tc.Docker != "" {
result.Available = append(result.Available, "docker-mcp")
} else {
result.Missing = append(result.Missing, "docker-mcp")
}
if tc.Playwright {
result.Available = append(result.Available, "playwright-mcp")
} else {
result.Missing = append(result.Missing, "playwright-mcp")
}
result.Available = append(result.Available, "filesystem-mcp")
return result
}
