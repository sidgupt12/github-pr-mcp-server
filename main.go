package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// A simplified MCP context structure
type Context struct {
	Type    string                 `json:"type"`
	Content map[string]interface{} `json:"content"`
	ID      string                 `json:"id,omitempty"`
}

// A simplified MCP message structure
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// MCP request structure
type MCPRequest struct {
	Messages []Message `json:"messages"`
	System   string    `json:"system,omitempty"`
	Context  []Context `json:"context,omitempty"`
}

// MCP response structure
type MCPResponse struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func main() {
	router := mux.NewRouter()

	// MCP endpoint
	router.HandleFunc("/mcp", mcpHandler).Methods("POST")

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}).Methods("GET")

	// Start server
	port := "8080"
	fmt.Printf("MCP Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func mcpHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the incoming request
	var req MCPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Log the received request
	reqJSON, _ := json.MarshalIndent(req, "", "  ")
	fmt.Printf("Received MCP request:\n%s\n", string(reqJSON))

	// Create a simple response
	// In a real MCP server, you'd process the context and messages
	lastMessage := ""
	if len(req.Messages) > 0 {
		lastMessage = req.Messages[len(req.Messages)-1].Content
	}

	contextInfo := "No context provided"
	if len(req.Context) > 0 {
		contextJSON, _ := json.MarshalIndent(req.Context, "", "  ")
		contextInfo = string(contextJSON)
	}

	response := MCPResponse{
		Role: "assistant",
		Content: fmt.Sprintf("Echo from MCP server at %s\n\nYou said: %s\n\nContext received: %s",
			time.Now().Format(time.RFC1123),
			lastMessage,
			contextInfo),
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
