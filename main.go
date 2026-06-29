// Package main is the entry point for the Glash CLI.
//
//	@title			Glash API
//	@version		1.0
//	@description	Glash is a terminal-based AI coding assistant. This API is served over a Unix socket (or Windows named pipe) and provides programmatic access to workspaces, sessions, agents, LSP, MCP, and more.
//	@contact.name	Charm
//	@contact.url	https://charm.sh
//	@license.name	MIT
//	@license.url	https://github.com/charmbracelet/glash/blob/main/LICENSE
//	@BasePath		/v1
package main

import (
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/charmbracelet/glash/internal/cmd"
	_ "github.com/charmbracelet/glash/internal/dns"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	if os.Getenv("GLASH_PROFILE") != "" {
		go func() {
			slog.Info("Serving pprof at localhost:6060")
			if httpErr := http.ListenAndServe("localhost:6060", nil); httpErr != nil {
				slog.Error("Failed to pprof listen", "error", httpErr)
			}
		}()
	}

	cmd.Execute()
}
