package mcp

import (
	"github.com/mark3labs/mcp-go/server"
)

// NewServer creates and configures the MCP server instance.
func NewServer() *server.MCPServer {
	return server.NewMCPServer(
		"DevPath MCP",
		"1.0.0",
		server.WithInstructions(
			"DevPath MCP helps you learn programming technologies step by step. "+
				"Use get_today_topic to get a personalized study plan for today, "+
				"generate_practice for hands-on challenges, "+
				"mark_topic_complete to save your progress, "+
				"get_next_topic to see what comes next, "+
				"and review_progress to track your learning journey.",
		),
	)
}
