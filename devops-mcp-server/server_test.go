package main

import (
	"context"
	"log"
	"testing"

	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestMCPServer(t *testing.T) {
	ctx := context.Background()

	// Create a client and server.
	// Wait for the client to initialize the session.
	client := mcp.NewClient(&mcp.Implementation{Name: "client", Version: "v0.0.1"}, nil)
	server := createServer()

	// Connect the server and client using in-memory transports.
	//
	// Connect the server first so that it's ready to receive initialization
	// messages from the client.
	t1, t2 := mcp.NewInMemoryTransports()
	serverSession, err := server.Connect(ctx, t1, nil)
	if err != nil {
		log.Fatal(err)
	}
	clientSession, err := client.Connect(ctx, t2, nil)
	if err != nil {
		log.Fatal(err)
	}

	tools, err := clientSession.ListTools(ctx, &mcp.ListToolsParams{})
	if err != nil {
		log.Fatalf("failed to list tools: %v", err)
	}
	log.Printf("Prompts: %v", len(tools.Tools))
	for _,tool := range tools.Tools {
		log.Printf("Tool %s: %s", tool.Name, tool.Title)
	}

	prompts, err := clientSession.ListPrompts(ctx, &mcp.ListPromptsParams{})
	if err != nil {
		log.Fatalf("failed to list prompts: %v", err)
	}
	log.Printf("Prompts: %v", len(prompts.Prompts))
	for _,prompt := range prompts.Prompts {
		log.Printf("Prompt %s: %s", prompt.Name, prompt.Title)
	}
	

	// Now shut down the session by closing the client, and waiting for the
	// server session to end.
	if err := clientSession.Close(); err != nil {
		log.Fatal(err)
	}
	if err := serverSession.Wait(); err != nil {
		log.Fatal(err)
	}
}
