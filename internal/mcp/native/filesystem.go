package native

import "github.com/muah1987/muah-runner/internal/mcp"

func FilesystemMCP() mcp.MCPServer {
return mcp.MCPServer{
Name: "filesystem-mcp",
URL:  "mcp://filesystem",
Type: "native",
Tools: []mcp.MCPTool{
{Name: "read_file", Description: "Read a file"},
{Name: "write_file", Description: "Write a file"},
{Name: "list_directory", Description: "List directory contents"},
{Name: "delete_file", Description: "Delete a file"},
{Name: "move_file", Description: "Move/rename a file"},
},
}
}
