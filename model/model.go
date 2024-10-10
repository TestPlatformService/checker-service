package model

type QuestionInfo struct {
	QuestionId string `json:"question_id"`
	IO          []InputOutput `json:"io"`
	TimeLimit   int32         `json:"timeLimit"`
	MemoryLimit int64         `json:"memoryLimit"`
}

type InputOutput struct {
	In  string `json:"in"`
	Out string `json:"out"`
}

type ApiRequest struct {
	Code        string        `json:"code"`
	Lang        string        `json:"lang"`
	TimeLimit   int32         `json:"timeLimit"`
	MemoryLimit int64         `json:"memoryLimit"`
	IO          []InputOutput `json:"io"`
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
