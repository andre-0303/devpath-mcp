package tools

import (
	"context"
	"fmt"

	"github.com/andre-0303/devpath-mcp/internal/service"
	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

func newMarkTopicCompleteTool() mcp.Tool {
	return mcp.NewTool("mark_topic_complete",
		mcp.WithDescription("Marks a topic as completed and saves progress to disk."),
		mcp.WithString("language",
			mcp.Required(),
			mcp.Description("Programming language"),
			mcp.Enum("golang"),
		),
		mcp.WithString("topic",
			mcp.Required(),
			mcp.Description("Name of the topic to mark as complete (case-insensitive)"),
		),
	)
}

func handleMarkTopicComplete(svc *service.Service) mcpserver.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		language, err := req.RequireString("language")
		if err != nil {
			return mcp.NewToolResultError("parameter 'language' is required"), nil
		}
		topic, err := req.RequireString("topic")
		if err != nil {
			return mcp.NewToolResultError("parameter 'topic' is required"), nil
		}

		alreadyDone := svc.IsTopicCompleted(language, topic)

		if err := svc.MarkTopicComplete(language, topic); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		done, total := svc.CompletedCountFor(language)
		pct := 0.0
		if total > 0 {
			pct = float64(done) / float64(total) * 100
		}

		prefix := "✓ Marked"
		if alreadyDone {
			prefix = "'%s' was already marked as complete for %s."
			// Reformat for already-done case.
			nextLine := ""
			if next, err2 := svc.GetNextTopic(language, topic); err2 == nil && next.NextTopic != nil {
				nextLine = fmt.Sprintf("\nNext up: %s [%s — %d min]", next.NextTopic.Name, next.NextTopic.Difficulty, next.NextTopic.EstimatedMinutes)
			}
			text := fmt.Sprintf(
				"'%s' was already marked as complete for %s.\n\nProgress: %d/%d topics completed (%.1f%%)%s",
				topic, language, done, total, pct, nextLine,
			)
			return mcp.NewToolResultText(text), nil
		}

		_ = prefix
		nextLine := ""
		if next, err2 := svc.GetNextTopic(language, topic); err2 == nil && next.NextTopic != nil {
			nextLine = fmt.Sprintf("\nNext up: %s [%s — %d min]", next.NextTopic.Name, next.NextTopic.Difficulty, next.NextTopic.EstimatedMinutes)
		}

		text := fmt.Sprintf(
			"✓ Marked '%s' as complete for %s!\n\nProgress: %d/%d topics completed (%.1f%%)%s",
			topic, language, done, total, pct, nextLine,
		)
		return mcp.NewToolResultText(text), nil
	}
}
