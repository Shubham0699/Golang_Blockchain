package p2p

import (
	"log"
	"net/http"
	"sync"

	"github.com/Shubham0699/go-mini-blockchain/block"
	"github.com/gorilla/websocket"
)

type Node struct {
	Address    string
	Peers      map[string]*websocket.Conn
	Mutex      sync.Mutex
	Blockchain *block.Blockchain
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Create a new node
func NewNode(address string, bc *block.Blockchain) *Node {
	return &Node{
		Address:    address,
		Peers:      make(map[string]*websocket.Conn),
		Blockchain: bc,
	}
}

// Connect to a peer
func (n *Node) ConnectPeer(peerAddr string) {
	ws, _, err := websocket.DefaultDialer.Dial("ws://"+peerAddr+"/ws", nil)
	if err != nil {
		log.Println("Failed to connect to peer:", err)
		return
	}
	n.Mutex.Lock()
	n.Peers[peerAddr] = ws
	n.Mutex.Unlock()
	go n.ListenPeer(ws)
	log.Println("✅ Connected to peer:", peerAddr)
}

// Handle incoming peer connections
func (n *Node) PeerHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade websocket:", err)
		return
	}
	n.Mutex.Lock()
	n.Peers[ws.RemoteAddr().String()] = ws
	n.Mutex.Unlock()
	go n.ListenPeer(ws)
	log.Println("✅ New peer connected:", ws.RemoteAddr().String())
}

// Listen for messages from a peer
func (n *Node) ListenPeer(ws *websocket.Conn) {
	for {
		var incoming block.Block
		if err := ws.ReadJSON(&incoming); err != nil {
			log.Println("Error reading block from peer:", err)
			n.Mutex.Lock()
			delete(n.Peers, ws.RemoteAddr().String())
			n.Mutex.Unlock()
			return
		}
		// Validate and add block
		n.Blockchain.MineBlock(incoming.Transactions)
		log.Println("✅ Received block from peer and added to chain")
	}
}

// Broadcast a block to all peers
func (n *Node) BroadcastBlock(b *block.Block) {
	n.Mutex.Lock()
	defer n.Mutex.Unlock()
	for peer, ws := range n.Peers {
		if err := ws.WriteJSON(b); err != nil {
			log.Println("Failed to send block to peer", peer, err)
		}
	}
}

// Start WebSocket server
func (n *Node) StartServer() {
	http.HandleFunc("/ws", n.PeerHandler)
	log.Println("🌐 P2P Node listening at", n.Address)
	if err := http.ListenAndServe(n.Address, nil); err != nil {
		log.Fatal(err)
	}
}
