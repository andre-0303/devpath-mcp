package main

import (
	"fmt"
	"log"
	"os"

	mcpsetup "github.com/andre-0303/devpath-mcp/internal/mcp"
	"github.com/andre-0303/devpath-mcp/internal/models"
	"github.com/andre-0303/devpath-mcp/internal/service"
	"github.com/andre-0303/devpath-mcp/internal/storage"
	"github.com/andre-0303/devpath-mcp/internal/tools"
	"github.com/mark3labs/mcp-go/server"
)

const dataPath = "data/progress.json"

func main() {
	// CRITICAL: all diagnostic output must go to stderr.
	// Any stray bytes on stdout corrupt the STDIO MCP transport.
	log.SetOutput(os.Stderr)

	progress, err := storage.Load(dataPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "devpath-mcp: failed to load progress: %v\n", err)
		os.Exit(1)
	}

	saveFn := func(p *models.UserProgress) error {
		return storage.Save(dataPath, p)
	}

	svc := service.New(progress, saveFn)
	s := mcpsetup.NewServer()
	tools.RegisterAll(s, svc)

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "devpath-mcp: server error: %v\n", err)
		os.Exit(1)
	}
}
