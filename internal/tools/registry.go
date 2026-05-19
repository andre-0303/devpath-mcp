package tools

import (
	"github.com/andre-0303/devpath-mcp/internal/service"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterAll adds all DevPath tools to the MCP server.
func RegisterAll(s *server.MCPServer, svc *service.Service) {
	s.AddTool(newGetTodayTopicTool(), handleGetTodayTopic(svc))
	s.AddTool(newGetNextTopicTool(), handleGetNextTopic(svc))
	s.AddTool(newGeneratePracticeTool(), handleGeneratePractice(svc))
	s.AddTool(newMarkTopicCompleteTool(), handleMarkTopicComplete(svc))
	s.AddTool(newReviewProgressTool(), handleReviewProgress(svc))
}
