package tools

import (
	"context"
	"fmt"

	"github.com/andre-0303/devpath-mcp/internal/service"
	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

func newGetTodayTopicTool() mcp.Tool {
	return mcp.NewTool("get_today_topic",
		mcp.WithDescription("Returns the ideal topic to study today based on your current progress and available time."),
		mcp.WithString("language",
			mcp.Required(),
			mcp.Description("Programming language to study"),
			mcp.Enum("golang"),
		),
		mcp.WithNumber("time_available",
			mcp.Required(),
			mcp.Description("How many minutes you have available to study today (5–480)"),
		),
	)
}

func handleGetTodayTopic(svc *service.Service) mcpserver.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		language, err := req.RequireString("language")
		if err != nil {
			return mcp.NewToolResultError("parameter 'language' is required"), nil
		}

		timeAvail := int(req.GetFloat("time_available", 30))
		if timeAvail < 5 || timeAvail > 480 {
			return mcp.NewToolResultError("time_available must be between 5 and 480 minutes"), nil
		}

		result, err := svc.GetTodayTopic(language, timeAvail)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		fitsLabel := "YES"
		if !result.FitsInTime {
			fitsLabel = fmt.Sprintf("NO (needs %d min — consider a longer session)", result.Topic.EstimatedMinutes)
		}

		practice := ""
		if len(result.Topic.PracticeChallenges) > 0 {
			pr, _ := svc.GeneratePractice(result.Topic.Name, language)
			if pr != nil {
				practice = fmt.Sprintf("\nPractice suggestion: %s", pr.Challenge)
			}
		}

		text := fmt.Sprintf(
			"Today's Topic: %s [%s — %d min]\nReason: %s\nFits in your %d-minute window: %s%s",
			result.Topic.Name,
			result.Topic.Difficulty,
			result.Topic.EstimatedMinutes,
			result.Reason,
			result.TimeAvailable,
			fitsLabel,
			practice,
		)

		return mcp.NewToolResultText(text), nil
	}
}
