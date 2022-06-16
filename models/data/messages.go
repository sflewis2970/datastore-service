package data

import "database/sql"

type StatusCode int

// Status Response Message
type StatusResponse struct {
	Timestamp string     `json:"timestamp"`
	Status    StatusCode `json:"status"`
	Message   string     `json:"message,omitempty"`
	Warning   string     `json:"warning,omitempty"`
	Error     string     `json:"error,omitempty"`
}

// Question Request-Response Messages
type QuestionRequest struct {
	QuestionID string `json:"questionid"`
	Question   string `json:"question"`
	Category   string `json:"category"`
	Answer     string `json:"answer"`
}

type QuestionResponse struct {
	QuestionID      string `json:"questionid"`
	Question        string `json:"question"`
	Category        string `json:"category"`
	Answer          string `json:"answer"`
	Timestamp       string `json:"timestamp"`
	Action          string `json:"action"`
	RecordsAffected string `json:"recordsaffected"`
	Message         string `json:"message,omitempty"`
	Warning         string `json:"warning,omitempty"`
	Error           string `json:"error,omitempty"`
}

type QuestionTable struct {
	Question string `json:"question"`
	Category string `json:"category"`
	Answer   string `json:"answer"`
}

// Answer Request-Response Messages
type AnswerRequest struct {
	QuestionID string `json:"questionid"`
	Response   string `json:"response"`
}

type AnswerResponse struct {
	Question  string `json:"question"`
	Category  string `json:"category"`
	Answer    string `json:"answer"`
	Response  string `json:"response"`
	Timestamp string `json:"timestamp"`
	Correct   bool   `json:"correct"`
	Message   string `json:"message,omitempty"`
	Warning   string `json:"warning,omitempty"`
	Error     string `json:"error,omitempty"`
}

const RESULTS_DEFAULT int64 = 0

type IDBModel interface {
	Open(driverName string) (*sql.DB, error)
	Ping() error
	InsertQuestion(question QuestionRequest) (int64, error)
	GetQuestion(questionID string) (QuestionTable, error)
	UpdateQuestion(question QuestionRequest) (int64, error)
	DeleteQuestion(questionID string) (int64, error)
}
