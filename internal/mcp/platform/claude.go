package platform

import "github.com/muah1987/muah-runner/internal/mcp"

func ClaudeMCP() mcp.MCPServer {
return mcp.MCPServer{
Name: "claude-mcp",
URL:  "mcp://anthropic",
Type: "platform",
Tools: []mcp.MCPTool{
{Name: "complete", Description: "Generate text completion"},
{Name: "analyze", Description: "Analyze content"},
},
}
}
