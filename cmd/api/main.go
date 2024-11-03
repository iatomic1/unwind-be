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
	"context"
	"log"

	"github.com/adeyemialameen04/unwind-be/core/handlers/books"
	"github.com/adeyemialameen04/unwind-be/core/router"
	"github.com/adeyemialameen04/unwind-be/core/server"
	"github.com/adeyemialameen04/unwind-be/internal/config"
	"github.com/adeyemialameen04/unwind-be/internal/projectpath"
	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()
	cfg, err := config.Load(projectpath.Root)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	conn, err := pgx.Connect(context.Background(), cfg.DbUrl)
	if err != nil {
		log.Fatalf("Unable to connect to db: %v\n", err)
	}
	defer conn.Close(context.Background())

	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	books.GetB(ctx, conn)
	router.SetupRouter(srv)
	server.RunServer(srv)
}
