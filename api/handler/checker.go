package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type InputOutput struct {
	In  string      `json:"in"`
	Out interface{} `json:"out"`
}

type RunRequest struct {
	QuestionId string `json:"question_id"`
	Code       string `json:"code"`
	Lang       string `json:"lang"`
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

func (h *Handler) Check(c *gin.Context) {
	var req RunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// API ga request yuborish
	runResp, err := sendRunRequest(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
		return
	}

	// SSE uchun javob sarlavhalarini sozlash
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	// Tekshiruv jarayonini o'rnatish
	for i := 1; i <= 3; i++ {
		time.Sleep(2 * time.Second) // Jarayonni simulatsiya qilish

		// Tekshiruv natijalarini SSE orqali yuborish
		result := fmt.Sprintf(`{"event":"status","payload":{"status":%d,"test":"%d","message":"%s"}}`, i, i, runResp.Message)
		fmt.Fprintf(c.Writer, "id: %s\n", time.Now().Format(time.RFC3339Nano))
		fmt.Fprintf(c.Writer, "data: %s\n\n", result)
		c.Writer.(http.Flusher).Flush()
	}

	// Tekshiruv yakunlanganda SSE stream'ini yopish
	fmt.Fprint(c.Writer, "event: close\ndata: end\n\n")
	c.Writer.(http.Flusher).Flush()
}
