package tools

import (
	"context"
	"fmt"

	"github.com/andre-0303/devpath-mcp/internal/service"
	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

func newGeneratePracticeTool() mcp.Tool {
	return mcp.NewTool("generate_practice",
		mcp.WithDescription("Generates a daily-rotated practice challenge for a topic to reinforce learning."),
		mcp.WithString("topic",
			mcp.Required(),
			mcp.Description("Name of the topic to practice (case-insensitive)"),
		),
		mcp.WithString("language",
			mcp.Description("Programming language (default: golang)"),
			mcp.Enum("golang"),
		),
	)
}

func handleGeneratePractice(svc *service.Service) mcpserver.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		topic, err := req.RequireString("topic")
		if err != nil {
			return mcp.NewToolResultError("parameter 'topic' is required"), nil
		}

		language := req.GetString("language", "golang")

		result, err := svc.GeneratePractice(topic, language)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		hint := ""
		if result.Hint != "" {
			hint = fmt.Sprintf("\nHint: %s", result.Hint)
		}

		text := fmt.Sprintf(
			"Practice Challenge — %s [%s]\n\nChallenge: %s%s",
			result.Topic.Name,
			result.Topic.Difficulty,
			result.Challenge,
			hint,
		)

		return mcp.NewToolResultText(text), nil
	}
}
