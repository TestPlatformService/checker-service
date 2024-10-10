package handler

import (
	"bytes"
	"checker/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func sendRunRequest(runReq model.ApiRequest, Log *slog.Logger) (*model.RunResponse, error) {
	apiURL := "https://capi.robocontest.uz/run"

	// Request body'ni JSON formatga o'tkazish
	reqBody, err := json.Marshal(runReq)
	if err != nil {
		Log.Error(fmt.Sprintf("Ma'lumotlarni jsonga o'girishda xatolik: %v", err))
		return nil, err
	}

	// HTTP POST so'rovini yuborish
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		Log.Error(fmt.Sprintf("POST so'rovini yuborishda xatolik: %v", err))
		return nil, err
	}
	defer resp.Body.Close()

	// Javobni o'qish
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Log.Error(fmt.Sprintf("Responseni o'qishda xatolik: %v", err))
		return nil, err
	}

	// Javobni JSON formatga o'girish
	var runResp model.RunResponse
	if err := json.Unmarshal(body, &runResp); err != nil {
		Log.Error(fmt.Sprintf("Ma'lumotlarni jsonga o'girishda xatolik: %v", err))
		return nil, err
	}

	return &runResp, nil
}

func (h *Handler) Check(c *gin.Context) {
	var req model.RunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Log.Error(fmt.Sprintf("Ma'lumotlarni olishda xatolik: %v", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	questionInfo, err := h.Service.QuestionInfo(req.QuestionId)
	if err != nil{
		h.Log.Error(fmt.Sprintf("Question ma'lumotlarini olishda xatolik: %v", err))
		c.JSON(http.StatusBadRequest, gin.H{"Message":"Question ma'lumotlarini olishda xatolik"})
		return 
	}

	// API ga request yuborish
	var apiReq = model.ApiRequest{
		Code: req.Code,
		Lang: req.Lang,
		MemoryLimit: questionInfo.MemoryLimit,
		TimeLimit: questionInfo.TimeLimit,
		IO: questionInfo.IO,
	}
	runResp, err := sendRunRequest(apiReq, h.Log)
	if err != nil {
		h.Log.Error(fmt.Sprintf("Robocontest api bilan bog'lanmadi: %v", err))
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
