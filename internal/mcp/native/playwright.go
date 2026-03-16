package native

import "github.com/muah1987/muah-runner/internal/mcp"

func PlaywrightMCP() mcp.MCPServer {
return mcp.MCPServer{
Name: "playwright-mcp",
URL:  "mcp://playwright",
Type: "native",
Tools: []mcp.MCPTool{
{Name: "screenshot", Description: "Take a screenshot of a URL"},
{Name: "click", Description: "Click an element"},
{Name: "fill_form", Description: "Fill a form"},
{Name: "navigate", Description: "Navigate to a URL"},
{Name: "scrape", Description: "Scrape page content"},
},
}
}
