package native

import "github.com/muah1987/muah-runner/internal/mcp"

func DockerMCP() mcp.MCPServer {
return mcp.MCPServer{
Name: "docker-mcp",
URL:  "mcp://docker",
Type: "native",
Tools: []mcp.MCPTool{
{Name: "container_run", Description: "Run a container"},
{Name: "container_stop", Description: "Stop a container"},
{Name: "image_build", Description: "Build a Docker image"},
{Name: "container_logs", Description: "Get container logs"},
{Name: "container_exec", Description: "Execute command in container"},
},
}
}
