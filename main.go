package main

import (
	"belajar-golang-chapter-48/helper"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", helper.HandleHome)
	http.HandleFunc("/ws", helper.HandleConnections)

	go helper.HandleMessages()

	log.Println("Server starting at :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
