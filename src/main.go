package main

import (
	"embed"
	"io/fs"

	"templatenextjs_golangserver/src/server"
	"templatenextjs_golangserver/src/utils"
)

//go:embed all:out
var webuiFS embed.FS

func main() {
	utils.LoadEnv()

	utils.LogInfo("Starting Golang + Next.js server...")

	// Extract the sub-filesystem under "out"
	strippedFS, err := fs.Sub(webuiFS, "out")
	if err != nil {
		utils.LogWarn("Warning: failed to initialize embedded frontend FS: %v", err)
	} else {
		server.WebuiFS = strippedFS
	}

	server.RunServer()
}
