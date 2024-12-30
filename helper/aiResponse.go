package helper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// ConversationTracker adalah struct untuk melacak riwayat percakapan.
type ConversationTracker struct {
	history    []map[string]string
	maxHistory int
}

// NewConversationTracker membuat instance baru dari ConversationTracker.
func NewConversationTracker(maxHistory int) *ConversationTracker {
	return &ConversationTracker{
		history:    []map[string]string{},
		maxHistory: maxHistory,
	}
}

// AddConversation menambahkan pasangan pertanyaan-jawaban ke riwayat percakapan.
func (ct *ConversationTracker) AddConversation(userInput, aiResponse string) {
	ct.history = append(ct.history, map[string]string{
		"question": userInput,
		"answer":   aiResponse,
	})

	// Jika jumlah percakapan melebihi batas, hapus yang paling lama
	if len(ct.history) > ct.maxHistory {
		ct.history = ct.history[1:]
	}
}

// GetHistory mengembalikan seluruh riwayat percakapan.
func (ct *ConversationTracker) GetHistory() []map[string]string {
	return ct.history
}

// ClearHistory menghapus seluruh riwayat percakapan.
func (ct *ConversationTracker) ClearHistory() {
	ct.history = []map[string]string{}
}

// UploadImage adalah fungsi untuk mengunggah gambar dan mendapatkan respons lengkap dari API.
func GetResponseAi(prompt string, listConversation []map[string]string) string {
	promptAwal := "Kamu adalah asisten virtual bernama Darmi, berikan jawaban singkat di bawah 200 karakter dan gunakan bahasa gaul bergaya bahasa anak muda Indonesia dengan kata aku dan kamu agar lebih terasa nyaman, dari pertanyaan dibawah:"
	if len(listConversation) > 0 {
		promptAwal = fmt.Sprintf("Kamu adalah asisten virtual bernama Darmi, berikan jawaban singkat di bawah 200 karakter dan gunakan bahasa gaul bergaya bahasa anak muda Indonesia dengan kata aku dan kamu agar lebih terasa nyaman, dengan menganalisa history percakapan berikut %v \nlalu beri jawaban dari pertanyaan dibawah:", listConversation)
	}

	// Encode prompt untuk URL
	encodedPrompt := url.QueryEscape(fmt.Sprintf("%s\n%s", promptAwal, prompt))
	apiURL := fmt.Sprintf("https://chatwithai.codesearch.workers.dev/?chat=%s", encodedPrompt)

	// Buat request GET
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return fmt.Sprintf("Gagal membuat request: %v", err)
	}

	// Eksekusi request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("Gagal mengirim request: %v", err)
	}
	defer resp.Body.Close()

	// Baca body response
	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("Gagal membaca response: %v", err)
	}

	// Parsing JSON response
	var jsonResponse map[string]string
	err = json.Unmarshal(resBody, &jsonResponse)
	if err != nil {
		return fmt.Sprintf("Gagal parsing response JSON: %v", err)
	}

	// Ambil data dari JSON response
	result, ok := jsonResponse["data"]
	if !ok || result == "" {
		return "Kurang tau nih, gagal dapet response"
	}

	return result
}
