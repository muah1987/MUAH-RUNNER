package platform

import "github.com/muah1987/muah-runner/internal/mcp"

func GitHubMCP() mcp.MCPServer {
return mcp.MCPServer{
Name: "github-mcp",
URL:  "mcp://github",
Type: "platform",
Tools: []mcp.MCPTool{
{Name: "create_pull_request", Description: "Create a GitHub pull request"},
{Name: "list_issues", Description: "List repository issues"},
{Name: "create_issue", Description: "Create a new issue"},
{Name: "get_file_contents", Description: "Get file contents from a repository"},
{Name: "push_files", Description: "Push files to a repository"},
{Name: "list_commits", Description: "List repository commits"},
{Name: "create_branch", Description: "Create a new branch"},
{Name: "merge_pull_request", Description: "Merge a pull request"},
},
}
}
