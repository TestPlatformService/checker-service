package handler

import (
	"bytes"
	"checker/model"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func sendRunRequest(runReq model.ApiRequest, Log *slog.Logger) ([]model.RunResponse, error) {
	apiURL := "https://capi.robocontest.uz/run"

	// Request body'ni JSON formatga o'tkazish
	Log.Info(fmt.Sprintf("Request: %v", runReq))
	reqBody, err := json.Marshal(runReq)
	if err != nil {
		Log.Error(fmt.Sprintf("Ma'lumotlarni jsonga o'girishda xatolik1: %v", err))
		return nil, err
	}

	// HTTP POST so'rovini yuborish
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		Log.Error(fmt.Sprintf("POST so'rovini yuborishda xatolik: %v", err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		Log.Error(fmt.Sprintf("POST so'rovining javobi xato: %v", resp.Status))
		return nil, fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}

	// SSE oqimini o'qish
	var runResponses []model.RunResponse
	decoder := json.NewDecoder(resp.Body)

	for {
		// SSE da har bir 'data' hodisasini o'qish
		var event map[string]json.RawMessage
		if err := decoder.Decode(&event); err == io.EOF {
			break // Barcha ma'lumotlar o'qildi
		} else if err != nil {
			Log.Error(fmt.Sprintf("SSE oqimini o'qishda xatolik: %v", err))
			return nil, err
		}

		// Agar 'data' kaliti mavjud bo'lsa, uni pars qilish
		if rawData, found := event["data"]; found {
			var runResp model.RunResponse
			if err := json.Unmarshal(rawData, &runResp); err != nil {
				Log.Error(fmt.Sprintf("Ma'lumotlarni jsonga o'girishda xatolik2: %v", err))
				return nil, err
			}
			runResponses = append(runResponses, runResp)
		}
	}

	return runResponses, nil
}

func (h *Handler) Check(c *gin.Context) {
	var req model.RunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Log.Error(fmt.Sprintf("Ma'lumotlarni olishda xatolik: %v", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	questionInfo, err := h.Service.QuestionInfo(req.QuestionId)
	if err != nil {
		h.Log.Error(fmt.Sprintf("Question ma'lumotlarini olishda xatolik: %v", err))
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Question ma'lumotlarini olishda xatolik"})
		return
	}

	// Prepare API request
	var apiReq = model.ApiRequest{
		Code:        req.Code,
		Lang:        req.Lang,
		MemoryLimit: questionInfo.MemoryLimit,
		TimeLimit:   questionInfo.TimeLimit,
		IO:          questionInfo.IO,
	}

	// SSE Headers
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no") // Disable buffering for SSE

	// Start the simulated checking process (or fetch from an actual external process)
	for i := 1; i <= 3; i++ {
		// Simulate an interval between steps
		time.Sleep(2 * time.Second)

		// Send request to external API and receive a slice of responses
		runResp, err := sendRunRequest(apiReq, h.Log)
		if err != nil {
			h.Log.Error(fmt.Sprintf("Robocontest API bilan bog'lanmadi: %v", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
			return
		}

		// Loop through each response in the slice and send SSE for each
		for _, response := range runResp {
			// Send an update to the client
			result := fmt.Sprintf(`{"event":"status","payload":{"status":"%s","message":"%s"}}`, response.Status, response.Message)
			fmt.Fprintf(c.Writer, "id: %s\n", time.Now().Format(time.RFC3339Nano))
			fmt.Fprintf(c.Writer, "data: %s\n\n", result)
			c.Writer.(http.Flusher).Flush()

			// If Status == "3", break the loop early (assuming "3" means completed)
			if response.Status == "3" {
				break
			}
		}
	}

	// Close the SSE stream when the process is done
	fmt.Fprint(c.Writer, "event: close\ndata: end\n\n")
	c.Writer.(http.Flusher).Flush()
}
