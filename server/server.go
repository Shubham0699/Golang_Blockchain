package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Shubham0699/go-mini-blockchain/block"
)

type Server struct {
	Blockchain *block.Blockchain
}

func NewServer(bc *block.Blockchain) *Server {
	return &Server{Blockchain: bc}
}

func (s *Server) Start(port string) {
	http.HandleFunc("/chain", s.handleGetChain)
	http.HandleFunc("/addblock", s.handleAddBlockQuery)  // GET way
	http.HandleFunc("/addblockjson", s.handleAddBlockPost) // POST way

	fmt.Println("🚀 Server running on port", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

// ---------------- GET /chain ----------------
func (s *Server) handleGetChain(w http.ResponseWriter, r *http.Request) {
	chain := s.Blockchain.GetBlocks() // ✅ FIXED: uses your alias
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chain)
}

// ---------------- GET /addblock?data=xxx ----------------
func (s *Server) handleAddBlockQuery(w http.ResponseWriter, r *http.Request) {
	data := r.URL.Query().Get("data")
	if data == "" {
		http.Error(w, "Missing data parameter", http.StatusBadRequest)
		return
	}
	s.Blockchain.AddBlock(data)
	fmt.Fprintf(w, "✅ Block added with data: %s", data)
}

// ---------------- POST /addblockjson ----------------
func (s *Server) handleAddBlockPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		Data string `json:"data"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil || body.Data == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	s.Blockchain.AddBlock(body.Data)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Block added successfully",
		"data":    body.Data,
	})
}
