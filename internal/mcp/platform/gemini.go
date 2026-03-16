package platform

import "github.com/muah1987/muah-runner/internal/mcp"

func GeminiMCP() mcp.MCPServer {
return mcp.MCPServer{
Name: "gemini-mcp",
URL:  "mcp://google",
Type: "platform",
Tools: []mcp.MCPTool{
{Name: "generate", Description: "Generate content"},
{Name: "embed", Description: "Create embeddings"},
},
}
}
