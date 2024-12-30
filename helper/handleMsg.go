package helper

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

// Map untuk menyimpan koneksi klien
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

// Upgrader untuk koneksi WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Mengizinkan semua origin (pastikan aman untuk produksi)
	},
}

// Struktur pesan
type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

// HandleHome melayani halaman HTML
func HandleHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "client/index.html")
}

// HandleConnections menangani koneksi WebSocket
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket Upgrade Error: %v", err)
		return
	}
	defer ws.Close()

	// Tambahkan klien ke daftar
	clients[ws] = true
	log.Printf("New client connected: %v", ws.RemoteAddr())

	// Loop untuk membaca pesan dari klien
	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Read Error: %v", err)
			delete(clients, ws) // Hapus klien dari daftar
			break
		}

		// Kirim pesan ke channel broadcast
		broadcast <- msg
	}
}

// HandleMessages menangani pesan broadcast
func HandleMessages() {
	tracker := NewConversationTracker(10) // Misal buffer 10 riwayat
	for {
		msg := <-broadcast
		log.Printf("Message received: %+v", msg)

		// Kirim pesan ke semua klien
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Write Error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}

		// Generate respons dari AI Bot
		botResponse := generateBotResponse(tracker, msg.Message)
		if botResponse != "" {
			botMsg := Message{
				Username: "AI Bot",
				Message:  botResponse,
			}

			tracker.AddConversation(msg.Message, botResponse)
			for client := range clients {
				err := client.WriteJSON(botMsg)
				if err != nil {
					log.Printf("Bot Write Error: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}
}

// Generate respons dari bot
func generateBotResponse(tracker *ConversationTracker, message string) string {
	message = strings.ToLower(message)
	if strings.Contains(message, "darmi") {
		listConversation := tracker.GetHistory()
		return GetResponseAi(message, listConversation)
	}
	return ""
}
