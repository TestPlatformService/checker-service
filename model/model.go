package model

type QuestionInfo struct {
	QuestionId  string        `json:"question_id"`
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

// Payload represents the inner structure of EventResponse
type Payload struct {
	Status       int    `json:"status"`
	Message      string `json:"message"`
	Test         int    `json:"test,omitempty"`
	Time         int    `json:"time,omitempty"`
	Memory       int    `json:"memory,omitempty"`
	CompileError string `json:"compile_error,omitempty"`
}

// EventResponse represents the response structure for an event
type EventResponse struct {
	EventId string  `json:"event_id"`
	Payload Payload `json:"payload"`
}

type Request struct {
	QuestionId     string `json:"question_id"`
	UserId         string `json:"user_id"`
	QuestionName   string `json:"question_name"`
	Status         string `json:"status"`
	Language       string `json:"language"`
	CompiledTime   string `json:"compiled_time"`
	CompiledMemory string `json:"compiled_memory"`
	Code           string `json:"code"`
	UserTaskId     string `json:"user_task_id"`
	SubmittedAt    string `json:"submitted_at"`
}
