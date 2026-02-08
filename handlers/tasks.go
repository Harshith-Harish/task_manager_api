package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/Harshith-Harish/task-manager-api/database"
	"github.com/Harshith-Harish/task-manager-api/models"
	"github.com/gin-gonic/gin"
)

// GetTasks retrieves all tasks with optional filtering
func GetTasks(c *gin.Context) {
	status := c.Query("status")
	priority := c.Query("priority")
	limit := c.DefaultQuery("limit", "100")

	query := `
		SELECT id, title, description, status, priority, created_at, updated_at 
		FROM tasks 
		WHERE 1=1
	`
	args := []interface{}{}
	argCount := 1

	if status != "" {
		query += " AND status = $" + strconv.Itoa(argCount)
		args = append(args, status)
		argCount++
	}

	if priority != "" {
		query += " AND priority = $" + strconv.Itoa(argCount)
		args = append(args, priority)
		argCount++
	}

	query += " ORDER BY created_at DESC LIMIT $" + strconv.Itoa(argCount)
	limitInt, _ := strconv.Atoi(limit)
	args = append(args, limitInt)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch tasks",
		})
		return
	}
	defer rows.Close()

	tasks := []models.Task{}
	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to scan task",
			})
			return
		}
		tasks = append(tasks, task)
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
		"count": len(tasks),
	})
}

// GetTask retrieves a single task by ID
func GetTask(c *gin.Context) {
	id := c.Param("id")

	var task models.Task
	query := `
		SELECT id, title, description, status, priority, created_at, updated_at 
		FROM tasks 
		WHERE id = $1
	`

	err := database.DB.QueryRow(query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch task",
		})
		return
	}

	c.JSON(http.StatusOK, task)
}

// CreateTask creates a new task
func CreateTask(c *gin.Context) {
	var task models.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Additional validation
	if err := task.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	query := `
		INSERT INTO tasks (title, description, status, priority, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	err := database.DB.QueryRow(
		query,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		now,
		now,
	).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create task",
		})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// UpdateTask updates an existing task
func UpdateTask(c *gin.Context) {
	id := c.Param("id")

	var update models.TaskUpdate
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	query := "UPDATE tasks SET updated_at = $1"
	args := []interface{}{time.Now()}
	argCount := 2

	if update.Title != nil {
		query += ", title = $" + strconv.Itoa(argCount)
		args = append(args, *update.Title)
		argCount++
	}

	if update.Description != nil {
		query += ", description = $" + strconv.Itoa(argCount)
		args = append(args, *update.Description)
		argCount++
	}

	if update.Status != nil {
		query += ", status = $" + strconv.Itoa(argCount)
		args = append(args, *update.Status)
		argCount++
	}

	if update.Priority != nil {
		query += ", priority = $" + strconv.Itoa(argCount)
		args = append(args, *update.Priority)
		argCount++
	}

	query += " WHERE id = $" + strconv.Itoa(argCount)
	query += " RETURNING id, title, description, status, priority, created_at, updated_at"
	args = append(args, id)

	var task models.Task
	err := database.DB.QueryRow(query, args...).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update task",
		})
		return
	}

	c.JSON(http.StatusOK, task)
}

// DeleteTask deletes a task by ID
func DeleteTask(c *gin.Context) {
	id := c.Param("id")

	result, err := database.DB.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete task",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task deleted successfully",
		"id":      id,
	})
}

// GetStats returns task statistics
func GetStats(c *gin.Context) {
	stats := models.TaskStats{
		ByStatus:   make(map[string]int),
		ByPriority: make(map[string]int),
	}

	// Total count
	err := database.DB.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&stats.Total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stats"})
		return
	}

	// By status
	rows, _ := database.DB.Query("SELECT status, COUNT(*) FROM tasks GROUP BY status")
	defer rows.Close()
	for rows.Next() {
		var status string
		var count int
		rows.Scan(&status, &count)
		stats.ByStatus[status] = count
	}

	// By priority
	rows, _ = database.DB.Query("SELECT priority, COUNT(*) FROM tasks GROUP BY priority")
	defer rows.Close()
	for rows.Next() {
		var priority string
		var count int
		rows.Scan(&priority, &count)
		stats.ByPriority[priority] = count
	}

	// Completed today
	database.DB.QueryRow(`
		SELECT COUNT(*) FROM tasks 
		WHERE status = 'completed' 
		AND DATE(updated_at) = CURRENT_DATE
	`).Scan(&stats.CompletedToday)

	c.JSON(http.StatusOK, stats)
}