package native

import "github.com/muah1987/muah-runner/internal/mcp"

func SQLiteMCP() mcp.MCPServer {
return mcp.MCPServer{
Name: "sqlite-mcp",
URL:  "mcp://sqlite",
Type: "native",
Tools: []mcp.MCPTool{
{Name: "query", Description: "Execute a SQL query"},
{Name: "insert", Description: "Insert a row"},
{Name: "create_table", Description: "Create a table"},
},
}
}
