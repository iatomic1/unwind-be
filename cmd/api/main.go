// @title Unwind Api
// @version 1.0
// @description Api for Unwind
//
// @contact.name Al-Ameen Adeyemi
// @contact.url: https://github.com/adeyemialameen04
//
// @host localhost:8080
package main

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/adeyemialameen04/unwind-be/core/router"
	"github.com/adeyemialameen04/unwind-be/core/server"
	"github.com/adeyemialameen04/unwind-be/internal/config"
)

// GetRootDir returns the absolute path to the project root directory
func GetRootDir() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(b), "../..")
}

func main() {
	rootDir := GetRootDir()
	cfg, err := config.Load(rootDir)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	router.SetupRouter(srv)
	server.RunServer(srv)
}
