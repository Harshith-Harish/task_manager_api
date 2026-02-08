package tests

import (
	"testing"

	"github.com/Harshith-Harish/task-manager-api/models"
)

func TestTaskValidation(t *testing.T) {
	tests := []struct {
		name    string
		task    models.Task
		wantErr bool
	}{
		{
			name: "valid task",
			task: models.Task{
				Title:    "Test Task",
				Status:   "pending",
				Priority: "medium",
			},
			wantErr: false,
		},
		{
			name: "empty title",
			task: models.Task{
				Title:    "",
				Status:   "pending",
				Priority: "medium",
			},
			wantErr: true,
		},
		{
			name: "invalid status",
			task: models.Task{
				Title:    "Test",
				Status:   "invalid",
				Priority: "medium",
			},
			wantErr: true,
		},
		{
			name: "invalid priority",
			task: models.Task{
				Title:    "Test",
				Status:   "pending",
				Priority: "urgent",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.task.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidStatus(t *testing.T) {
	validStatuses := []string{"pending", "in_progress", "completed"}

	for _, status := range validStatuses {
		task := models.Task{
			Title:    "Test",
			Status:   status,
			Priority: "medium",
		}
		if err := task.Validate(); err != nil {
			t.Errorf("Valid status %s should not error, got: %v", status, err)
		}
	}
}

func TestValidPriority(t *testing.T) {
	validPriorities := []string{"low", "medium", "high"}

	for _, priority := range validPriorities {
		task := models.Task{
			Title:    "Test",
			Status:   "pending",
			Priority: priority,
		}
		if err := task.Validate(); err != nil {
			t.Errorf("Valid priority %s should not error, got: %v", priority, err)
		}
	}
}