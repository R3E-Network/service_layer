package models

import (
	"encoding/json"
	"time"
)

// Function represents a JavaScript function to be executed in the TEE
type Function struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	OwnerID    string    `json:"ownerId"`
	UserID     int       `json:"userId"`
	Code       string    `json:"code"`
	SecretRefs []string  `json:"secretRefs,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// ExecutionResult represents the outcome of a function execution
type ExecutionResult struct {
	ID            string        `json:"id"`
	ExecutionID   string        `json:"executionId"`
	FunctionID    string        `json:"functionId"`
	Status        string        `json:"status"`
	Result        interface{}   `json:"result"`
	Error         string        `json:"error,omitempty"`
	ExecutionTime time.Duration `json:"executionTime"`
	StartTime     time.Time     `json:"startTime"`
	EndTime       time.Time     `json:"endTime"`
	Timestamp     time.Time     `json:"timestamp"`
	GasUsed       float64       `json:"gasUsed,omitempty"`
	Logs          []string      `json:"logs,omitempty"`
}

// FunctionSecret represents a securely stored credential for functions
type FunctionSecret struct {
	ID        string    `json:"id"`
	OwnerID   string    `json:"ownerId"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

// Execution represents a function execution
type Execution struct {
	ID         int             `json:"id" db:"id"`
	FunctionID int             `json:"function_id" db:"function_id"`
	Status     string          `json:"status" db:"status"`
	StartTime  time.Time       `json:"start_time" db:"start_time"`
	EndTime    time.Time       `json:"end_time,omitempty" db:"end_time"`
	Duration   int             `json:"duration,omitempty" db:"duration"`
	Result     json.RawMessage `json:"result,omitempty" db:"result"`
	Error      string          `json:"error,omitempty" db:"error"`
	CreatedAt  time.Time       `json:"created_at" db:"created_at"`
	Logs       []ExecutionLog  `json:"logs,omitempty" db:"-"`
}

// ExecutionLog represents a log entry for a function execution
type ExecutionLog struct {
	ID          int       `json:"id" db:"id"`
	ExecutionID int       `json:"execution_id" db:"execution_id"`
	Timestamp   time.Time `json:"timestamp" db:"timestamp"`
	Level       string    `json:"level" db:"level"`
	Message     string    `json:"message" db:"message"`
}

// ExecutionRequest represents a request to execute a function
type ExecutionRequest struct {
	Params interface{} `json:"params"`
	Async  bool        `json:"async"`
}

// FunctionRepository defines methods for working with functions
type FunctionRepository interface {
	Create(function *Function) error
	GetByID(id int) (*Function, error)
	GetByUserIDAndName(userID int, name string) (*Function, error)
	List(userID int, offset, limit int) ([]*Function, error)
	Update(function *Function) error
	Delete(id int) error
	IncrementExecutionCount(id int) error
	UpdateLastExecution(id int, lastExecution time.Time) error
	GetSecrets(functionID int) ([]string, error)
	SetSecrets(functionID int, secrets []string) error
}

// ExecutionRepository defines methods for working with executions
type ExecutionRepository interface {
	Create(execution *Execution) error
	GetByID(id int) (*Execution, error)
	ListByFunctionID(functionID int, offset, limit int) ([]*Execution, error)
	Update(execution *Execution) error
	Delete(id int) error
	AddLog(log *ExecutionLog) error
	GetLogs(executionID int, offset, limit int) ([]*ExecutionLog, error)
}
