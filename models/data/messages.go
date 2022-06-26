package data

import "database/sql"

const NO_RESULTS_RETURNED_MSG string = "No results returned..."

const (
	RESULTS_DEFAULT          int64 = 0
	OPEN_ERROR_CODE          int64 = -1
	INSERT_ERROR_CODE        int64 = -2
	GET_ERROR_CODE           int64 = -3
	UPDATE_ERROR_CODE        int64 = -4
	DELETE_ERROR_CODE        int64 = -5
	ROWS_AFFECTED_ERROR_CODE int64 = -99
)

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

type IDBModel interface {
	Open(driverName string) (*sql.DB, error)
	Ping() error
	Insert(question QuestionRequest) (int64, error)
	Get(questionID string) (QuestionTable, error)
	Update(question QuestionRequest) (int64, error)
	Delete(questionID string) (int64, error)
}
