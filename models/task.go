package models

import (
	"time"
)

// Task represents a task in the system
type Task struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title" binding:"required,min=1,max=255"`
	Description string    `json:"description" db:"description"`
	Status      string    `json:"status" db:"status" binding:"required,oneof=pending in_progress completed"`
	Priority    string    `json:"priority" db:"priority" binding:"required,oneof=low medium high"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// TaskUpdate represents partial task updates
type TaskUpdate struct {
	Title       *string `json:"title" binding:"omitempty,min=1,max=255"`
	Description *string `json:"description"`
	Status      *string `json:"status" binding:"omitempty,oneof=pending in_progress completed"`
	Priority    *string `json:"priority" binding:"omitempty,oneof=low medium high"`
}

// TaskStats represents task statistics
type TaskStats struct {
	Total       int            `json:"total"`
	ByStatus    map[string]int `json:"by_status"`
	ByPriority  map[string]int `json:"by_priority"`
	CompletedToday int         `json:"completed_today"`
}

// Validate checks if task fields are valid
func (t *Task) Validate() error {
	if t.Title == "" {
		return ErrEmptyTitle
	}
	if len(t.Title) > 255 {
		return ErrTitleTooLong
	}
	if !isValidStatus(t.Status) {
		return ErrInvalidStatus
	}
	if !isValidPriority(t.Priority) {
		return ErrInvalidPriority
	}
	return nil
}

func isValidStatus(status string) bool {
	validStatuses := []string{"pending", "in_progress", "completed"}
	for _, v := range validStatuses {
		if v == status {
			return true
		}
	}
	return false
}

func isValidPriority(priority string) bool {
	validPriorities := []string{"low", "medium", "high"}
	for _, v := range validPriorities {
		if v == priority {
			return true
		}
	}
	return false
}