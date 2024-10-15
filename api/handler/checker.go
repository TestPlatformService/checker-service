package handler

import (
	"bufio"
	"bytes"
	"checker/model"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func sendRunRequest(runReq model.ApiRequest, Log *slog.Logger) ([]model.Payload, error) {
	apiURL := "https://capi.robocontest.uz/run"

	// Serialize request to JSON
	Log.Info(fmt.Sprintf("Sending request: %v", runReq))
	reqBody, err := json.Marshal(runReq)
	if err != nil {
		Log.Error(fmt.Sprintf("Error marshalling to JSON: %v", err))
		return nil, err
	}

	// Send HTTP POST request
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		Log.Error(fmt.Sprintf("Error sending POST request: %v", err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		Log.Error(fmt.Sprintf("Received non-OK HTTP status: %v", resp.Status))
		return nil, fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}

	// SSE response processing using bufio.Scanner
	var runResponses []model.Payload
	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := scanner.Text()

		// SSE events often start with "data: " prefix
		if strings.HasPrefix(line, "data: ") {
			// Extract the JSON part after "data: "
			jsonData := strings.TrimPrefix(line, "data: ")

			// Try to unmarshal the JSON part into your EventResponse struct
			var eventResp model.EventResponse
			if err := json.Unmarshal([]byte(jsonData), &eventResp); err != nil {
				Log.Error(fmt.Sprintf("Error unmarshalling response: %v", err))
				return nil, err
			}

			// Log and append response based on event status
			Log.Info(fmt.Sprintf("Received status: %d, message: %v", eventResp.Payload.Status, eventResp.Payload.Message))

			// Check for compile error
			if eventResp.Payload.Status == 5 { // Compile error
				Log.Error(fmt.Sprintf("Compilation error: %s", eventResp.Payload.CompileError))
				return nil, fmt.Errorf("compilation error: %s", eventResp.Payload.CompileError)
			}

			runResponses = append(runResponses, eventResp.Payload)

			// Handle based on status
			if eventResp.Payload.Status == 9 { // Success
				Log.Info("Code execution completed successfully.")
				break
			}
		}
	}

	if err := scanner.Err(); err != nil {
		Log.Error(fmt.Sprintf("Error reading response body: %v", err))
		return nil, err
	}

	return runResponses, nil
}

// Check handles the incoming request, interacts with the external API, and sends SSE responses.
func (h *Handler) Check(c *gin.Context) {
	var req model.RunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Log.Error(fmt.Sprintf("Invalid request body: %v", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	questionInfo, err := h.Service.QuestionInfo(req.QuestionId)
	if err != nil {
		h.Log.Error(fmt.Sprintf("Error fetching question info: %v", err))
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Error fetching question info"})
		return
	}

	// Prepare API request
	apiReq := model.ApiRequest{
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

	// Send request to external API and get the response
	runResp, err := sendRunRequest(apiReq, h.Log)
	if err != nil {
		h.Log.Error(fmt.Sprintf("Error communicating with RoboContest API: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}

	// Send responses to client via SSE
	for _, response := range runResp {
		result := fmt.Sprintf(`{"event":"status","payload":{"status":%d,"message":"%s"}}`, response.Status, response.Message)
		fmt.Fprintf(c.Writer, "id: %s\n", time.Now().Format(time.RFC3339Nano))
		fmt.Fprintf(c.Writer, "data: %s\n\n", result)
		c.Writer.(http.Flusher).Flush()

		// If status is success (status == 9), stop processing further.
		if response.Status == 9 {
			break
		}
	}

	// Close the SSE stream
	fmt.Fprint(c.Writer, "event: close\ndata: end\n\n")
	c.Writer.(http.Flusher).Flush()
}
