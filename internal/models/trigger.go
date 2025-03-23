package models

import (
	"time"
)

// TriggerType defines the type of trigger
type TriggerType string

const (
	TriggerTypeSchedule    TriggerType = "schedule"
	TriggerTypePriceAlert  TriggerType = "price_alert"
	TriggerTypeBlockHeight TriggerType = "block_height"
	TriggerTypeTransaction TriggerType = "transaction"
)

// Trigger represents an automated function execution condition
type Trigger struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	UserID     string                 `json:"userId"`
	Type       TriggerType            `json:"type"`
	FunctionID string                 `json:"functionId"`
	Schedule   string                 `json:"schedule,omitempty"`   // CRON expression for schedule triggers
	Condition  string                 `json:"condition,omitempty"`  // Condition expression for other triggers
	Parameters map[string]interface{} `json:"parameters,omitempty"` // Parameters to pass to the function
	CreatedAt  time.Time              `json:"createdAt"`
	UpdatedAt  time.Time              `json:"updatedAt"`
	Status     string                 `json:"status"` // active, paused, error
}

// TriggerExecution represents a record of a triggered function execution
type TriggerExecution struct {
	ID           string    `json:"id"`
	TriggerID    string    `json:"triggerId"`
	FunctionID   string    `json:"functionId"`
	Status       string    `json:"status"` // success, error
	ExecutionID  string    `json:"executionId,omitempty"`
	ErrorMessage string    `json:"errorMessage,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
}

// TriggerRepository defines methods for working with triggers
type TriggerRepository interface {
	Create(trigger *Trigger) error
	GetByID(id int) (*Trigger, error)
	GetByUserIDAndName(userID int, name string) (*Trigger, error)
	List(userID int, offset, limit int) ([]*Trigger, error)
	ListActiveTriggers() ([]*Trigger, error)
	Update(trigger *Trigger) error
	UpdateStatus(id int, status string) error
	Delete(id int) error

	// Event related methods
	CreateEvent(event *TriggerExecution) error
	GetEventByID(id int) (*TriggerExecution, error)
	ListEventsByTriggerID(triggerID int, offset, limit int) ([]*TriggerExecution, error)
}
