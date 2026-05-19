package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/andre-0303/devpath-mcp/internal/service"
	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

func newReviewProgressTool() mcp.Tool {
	return mcp.NewTool("review_progress",
		mcp.WithDescription("Shows your full learning progress: completed topics, next recommended, and weak areas."),
		mcp.WithString("language",
			mcp.Required(),
			mcp.Description("Programming language"),
			mcp.Enum("golang"),
		),
	)
}

func handleReviewProgress(svc *service.Service) mcpserver.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		language, err := req.RequireString("language")
		if err != nil {
			return mcp.NewToolResultError("parameter 'language' is required"), nil
		}

		result, err := svc.ReviewProgress(language)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		separator := strings.Repeat("═", 36)

		if result.CompletedCount == 0 {
			nextLine := ""
			if result.NextRecommended != nil {
				nextLine = fmt.Sprintf("  %s [%s — %d min]", result.NextRecommended.Name, result.NextRecommended.Difficulty, result.NextRecommended.EstimatedMinutes)
			}
			text := fmt.Sprintf(
				"Progress Report — %s\n%s\n\nCompleted: 0/%d topics (0%%)\n\nNo topics completed yet. Start with:\n%s\n\nTip: Use get_today_topic to get a personalized study plan for today.",
				language, separator, result.TotalTopics, nextLine,
			)
			return mcp.NewToolResultText(text), nil
		}

		// Completed topics — format in rows of 3.
		var completedLines []string
		for i := 0; i < len(result.CompletedTopics); i += 3 {
			end := i + 3
			if end > len(result.CompletedTopics) {
				end = len(result.CompletedTopics)
			}
			row := result.CompletedTopics[i:end]
			parts := make([]string, len(row))
			for j, t := range row {
				parts[j] = fmt.Sprintf("  ✓ %-18s", t)
			}
			completedLines = append(completedLines, strings.Join(parts, ""))
		}

		nextSection := "  None — roadmap complete!"
		if result.NextRecommended != nil {
			prereqStatus := "none"
			if len(result.NextRecommended.Prerequisites) > 0 {
				prereqStatus = strings.Join(result.NextRecommended.Prerequisites, " ✓, ") + " ✓"
			}
			nextSection = fmt.Sprintf(
				"  %s [%s — %d min]\n  Prerequisites: %s",
				result.NextRecommended.Name,
				result.NextRecommended.Difficulty,
				result.NextRecommended.EstimatedMinutes,
				prereqStatus,
			)
		}

		weakSection := "  None — great job staying on track!"
		if len(result.WeakAreas) > 0 {
			lines := make([]string, len(result.WeakAreas))
			for i, wa := range result.WeakAreas {
				lines[i] = fmt.Sprintf("  %d. %-20s — blocks %d downstream topic(s)", i+1, wa.TopicName, wa.DependentCount)
			}
			weakSection = strings.Join(lines, "\n")
		}

		text := fmt.Sprintf(
			"Progress Report — %s\n%s\n\nCompleted: %d/%d topics (%.1f%%)\n\nCompleted topics:\n%s\n\nNext recommended:\n%s\n\nWeak areas (uncompleted blockers):\n%s",
			language,
			separator,
			result.CompletedCount,
			result.TotalTopics,
			result.PercentDone,
			strings.Join(completedLines, "\n"),
			nextSection,
			weakSection,
		)

		return mcp.NewToolResultText(text), nil
	}
}
