package main

import (
	"bytes"
	"checker/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type InputOutput struct {
	In  string `json:"in"`
	Out string `json:"out"`
}

type RunRequest struct {
	Code        string        `json:"code"`
	Lang        string        `json:"lang"`
	TimeLimit   int32         `json:"timeLimit"`
	MemoryLimit int64         `json:"memoryLimit"`
	IO          []InputOutput `json:"io"`
}

type RunResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func sendRunRequest(runReq RunRequest) (*RunResponse, error) {
	// API URL
	apiURL := "https://capi.robocontest.uz/run"

	// Request body'ni JSON formatga o'tkazish
	reqBody, err := json.Marshal(runReq)
	if err != nil {
		return nil, err
	}

	// HTTP POST so'rovini yuborish
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Javobni o'qish
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Javobni JSON formatga o'girish
	var runResp RunResponse
	if err := json.Unmarshal(body, &runResp); err != nil {
		return nil, err
	}

	return &runResp, nil
}

func handleCheck(w http.ResponseWriter, r *http.Request) {
	// Request'dan ma'lumotlarni olish
	var req RunRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// API ga request yuborish
	runResp, err := sendRunRequest(req)
	if err != nil {
		http.Error(w, "Failed to send request", http.StatusInternalServerError)
		return
	}

	// SSE uchun javob sarlavhalarini sozlash
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Tekshiruv jarayonini o'rnatish
	for i := 1; i <= 3; i++ {
		time.Sleep(2 * time.Second) // Jarayonni simulatsiya qilish

		// Tekshiruv natijalarini SSE orqali yuborish
		result := fmt.Sprintf(`{"event":"status","payload":{"status":%d,"test":"%d","message":"%s"}}`, i, i, runResp.Message)
		fmt.Fprintf(w, "id: %s\n", time.Now().Format(time.RFC3339Nano))
		fmt.Fprintf(w, "data: %s\n\n", result)
		w.(http.Flusher).Flush()
	}

	// Tekshiruv yakunlanganda SSE stream'ini yopish
	fmt.Fprint(w, "event: close\ndata: end\n\n")
	w.(http.Flusher).Flush()
}

func main() {
	cfg := config.LoadConfig()
	http.HandleFunc("/check", handleCheck)
	log.Println("Checker Service is running on port 50054")
	http.ListenAndServe(cfg.CHECKER_SERVICE, nil)
}
