// Command server is the MCP server of the example plugin of Octopod: plugin-host runs it and
// calls its tools, with the workspace, the persona and the activation parameters in
// _meta.octopod.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// octopod is what plugin-host tells about each call.
type octopod struct {
	Workspace string            `json:"workspace"`
	Persona   string            `json:"persona"`
	Params    map[string]string `json:"params"`
}

func callerOf(req *mcp.CallToolRequest) octopod {
	var o octopod
	if raw, err := json.Marshal(req.Params.Meta["octopod"]); err == nil {
		_ = json.Unmarshal(raw, &o)
	}
	return o
}

func text(format string, args ...any) *mcp.CallToolResult {
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf(format, args...)}}}
}

func failure(format string, args ...any) *mcp.CallToolResult {
	result := text(format, args...)
	result.IsError = true
	return result
}

func main() {
	server := mcp.NewServer(&mcp.Implementation{Name: "octopod-plugin-example", Version: "1"}, nil)

	server.AddTool(&mcp.Tool{Name: "local_time", InputSchema: json.RawMessage(`{"type":"object","required":["timezone"],"properties":{"timezone":{"type":"string"}}}`)},
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args struct {
				Timezone string `json:"timezone"`
			}
			if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
				return failure("invalid arguments: %v", err), nil
			}
			location, err := time.LoadLocation(args.Timezone)
			if err != nil {
				return failure("unknown time zone %q", args.Timezone), nil
			}
			return text("It is %s in %s.", time.Now().In(location).Format("15:04 on Monday 2 January"), args.Timezone), nil
		})

	server.AddTool(&mcp.Tool{Name: "greet", InputSchema: json.RawMessage(`{"type":"object","required":["name"],"properties":{"name":{"type":"string"}}}`)},
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args struct {
				Name string `json:"name"`
			}
			if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
				return failure("invalid arguments: %v", err), nil
			}
			caller := callerOf(req)
			greeting := caller.Params["greeting"]
			if greeting == "" {
				greeting = "Hello"
			}
			return text("%s %s! (from the workspace %s)", greeting, args.Name, caller.Workspace), nil
		})

	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return server }, &mcp.StreamableHTTPOptions{Stateless: true, JSONResponse: true})
	log.Fatal((&http.Server{Addr: ":8080", Handler: handler, ReadHeaderTimeout: 10 * time.Second}).ListenAndServe())
}
