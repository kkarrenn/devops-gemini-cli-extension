// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	httpAddr  = flag.String("http", "", "if set, use streamable HTTP at this address, instead of stdin/stdout")
	pprofAddr = flag.String("pprof", "", "if set, host the pprof debugging server at this address")
)

func main() {
	flag.Parse()

	if *pprofAddr != "" {
		// For debugging memory leaks, add an endpoint to trigger a few garbage
		// collection cycles and ensure the /debug/pprof/heap endpoint only reports
		// reachable memory.
		http.DefaultServeMux.Handle("/gc", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			for range 3 {
				runtime.GC()
			}
			fmt.Fprintln(w, "GC'ed")
		}))
		go func() {
			// DefaultServeMux was mutated by the /debug/pprof import.
			http.ListenAndServe(*pprofAddr, http.DefaultServeMux)
		}()
	}

	opts := &mcp.ServerOptions{
		Instructions:      "Use this server!",
		CompletionHandler: complete, // support completions by setting this handler
	}

	server := mcp.NewServer(&mcp.Implementation{Name: "everything"}, opts)

	// Add tools that exercise different features of the protocol.
	mcp.AddTool(server, &mcp.Tool{Name: "greet", Description: "say hi"}, contentTool)
	mcp.AddTool(server, &mcp.Tool{Name: "greet (structured)"}, structuredTool) // returns structured output
	mcp.AddTool(server, &mcp.Tool{Name: "ping"}, pingingTool)                  // performs a ping
	mcp.AddTool(server, &mcp.Tool{Name: "log"}, loggingTool)                   // performs a log
	mcp.AddTool(server, &mcp.Tool{Name: "sample"}, samplingTool)               // performs sampling
	mcp.AddTool(server, &mcp.Tool{Name: "elicit"}, elicitingTool)              // performs elicitation

	// Add the new tools
	if err := addAllTools(context.Background(), server); err != nil {
		log.Fatalf("failed to add tools: %v", err)
	}

	// Add a basic prompt.
	server.AddPrompt(&mcp.Prompt{Name: "greet"}, prompt)

	// Serve over stdio, or streamable HTTP if -http is set.
	if *httpAddr != "" {
		handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
			return server
		}, nil)
		log.Printf("MCP handler listening at %s", *httpAddr)
		if *pprofAddr != "" {
			log.Printf("pprof listening at http://%s/debug/pprof", *pprofAddr)
		}
		log.Fatal(http.ListenAndServe(*httpAddr, handler))
	} else {
		t := &mcp.LoggingTransport{Transport: &mcp.StdioTransport{}, Writer: os.Stderr}
		if err := server.Run(context.Background(), t); err != nil {
			log.Printf("Server failed: %v", err)
		}
	}
}

func prompt(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	return &mcp.GetPromptResult{
		Description: "Hi prompt",
		Messages: []*mcp.PromptMessage{
			{
				Role:    "user",
				Content: &mcp.TextContent{Text: "Say hi to " + req.Params.Arguments["name"]},
			},
		},
	}, nil
}

type args struct {
	Name string `json:"name" jsonschema:"the name to say hi to"`
}

// contentTool is a tool that returns unstructured content.
//
// Since its output type is 'any', no output schema is created.
func contentTool(ctx context.Context, req *mcp.CallToolRequest, args args) (*mcp.CallToolResult, any, error) {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: "Hi " + args.Name},
		},
	}, nil, nil
}

type result struct {
	Message string `json:"message" jsonschema:"the message to convey"`
}

// structuredTool returns a structured result.
func structuredTool(ctx context.Context, req *mcp.CallToolRequest, args *args) (*mcp.CallToolResult, *result, error) {
	return nil, &result{Message: "Hi " + args.Name}, nil
}

func pingingTool(ctx context.Context, req *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
	if err := req.Session.Ping(ctx, nil); err != nil {
		return nil, nil, fmt.Errorf("ping failed")
	}
	return nil, nil, nil
}

func loggingTool(ctx context.Context, req *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
	if err := req.Session.Log(ctx, &mcp.LoggingMessageParams{
		Data:  "something happened!",
		Level: "error",
	}); err != nil {
		return nil, nil, fmt.Errorf("log failed")
	}
	return nil, nil, nil
}

func samplingTool(ctx context.Context, req *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
	res, err := req.Session.CreateMessage(ctx, new(mcp.CreateMessageParams))
	if err != nil {
		return nil, nil, fmt.Errorf("sampling failed: %v", err)
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			res.Content,
		},
	}, nil, nil
}

func elicitingTool(ctx context.Context, req *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
	res, err := req.Session.Elicit(ctx, &mcp.ElicitParams{
		Message: "provide a random string",
		RequestedSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"random": {Type: "string"},
			},
		},
	})
	if err != nil {
		return nil, nil, fmt.Errorf("eliciting failed: %v", err)
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: res.Content["random"].(string)},
		},
	}, nil, nil
}

func complete(ctx context.Context, req *mcp.CompleteRequest) (*mcp.CompleteResult, error) {
	return &mcp.CompleteResult{
		Completion: mcp.CompletionResultDetails{
			Total:  1,
			Values: []string{req.Params.Argument.Value + "x"},
		},
	}, nil
}
