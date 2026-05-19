package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/andre-0303/devpath-mcp/internal/service"
	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

func newGetNextTopicTool() mcp.Tool {
	return mcp.NewTool("get_next_topic",
		mcp.WithDescription("Returns the next topic after the current one in the roadmap."),
		mcp.WithString("language",
			mcp.Required(),
			mcp.Description("Programming language"),
			mcp.Enum("golang"),
		),
		mcp.WithString("current_topic",
			mcp.Required(),
			mcp.Description("Name of the current topic (case-insensitive)"),
		),
	)
}

func handleGetNextTopic(svc *service.Service) mcpserver.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		language, err := req.RequireString("language")
		if err != nil {
			return mcp.NewToolResultError("parameter 'language' is required"), nil
		}
		currentTopic, err := req.RequireString("current_topic")
		if err != nil {
			return mcp.NewToolResultError("parameter 'current_topic' is required"), nil
		}

		result, err := svc.GetNextTopic(language, currentTopic)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		if result.IsComplete {
			text := fmt.Sprintf("'%s' is the last topic in the %s roadmap.\nRoadmap complete! You've mastered all topics.", result.CurrentTopic, language)
			return mcp.NewToolResultText(text), nil
		}

		next := result.NextTopic
		prereqs := "none"
		if len(next.Prerequisites) > 0 {
			prereqs = strings.Join(next.Prerequisites, ", ")
		}

		practice := ""
		if len(next.PracticeChallenges) > 0 {
			pr, _ := svc.GeneratePractice(next.Name, language)
			if pr != nil {
				practice = fmt.Sprintf("\nStart with: %q", pr.Challenge)
			}
		}

		text := fmt.Sprintf(
			"Next topic after '%s' in %s:\n\n  %s [%s — %d min]\n  Prerequisites: %s%s",
			result.CurrentTopic,
			language,
			next.Name,
			next.Difficulty,
			next.EstimatedMinutes,
			prereqs,
			practice,
		)

		return mcp.NewToolResultText(text), nil
	}
}
