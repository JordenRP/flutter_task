package models

import (
    "todo-app/internal/db"
    "time"
)

type Priority int

const (
    Low Priority = iota
    Medium
    High
)

type Task struct {
    ID          uint      `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Completed   bool      `json:"completed"`
    UserID      uint      `json:"user_id"`
    DueDate     time.Time `json:"due_date"`
    Priority    Priority  `json:"priority"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

func CreateTask(title, description string, userID uint, dueDate time.Time, priority Priority) (*Task, error) {
    var task Task
    err := db.DB.QueryRow(
        `INSERT INTO tasks (title, description, completed, user_id, due_date, priority, created_at, updated_at) 
         VALUES ($1, $2, false, $3, $4, $5, NOW(), NOW()) 
         RETURNING id, title, description, completed, user_id, due_date, priority, created_at, updated_at`,
        title, description, userID, dueDate, priority,
    ).Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.UserID, &task.DueDate, &task.Priority, &task.CreatedAt, &task.UpdatedAt)

    if err != nil {
        return nil, err
    }
    return &task, nil
}

func GetUserTasks(userID uint) ([]Task, error) {
    rows, err := db.DB.Query(
        `SELECT id, title, description, completed, user_id, due_date, priority, created_at, updated_at 
         FROM tasks WHERE user_id = $1 
         ORDER BY due_date ASC, priority DESC, created_at DESC`,
        userID,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var tasks []Task
    for rows.Next() {
        var task Task
        err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.UserID, &task.DueDate, &task.Priority, &task.CreatedAt, &task.UpdatedAt)
        if err != nil {
            return nil, err
        }
        tasks = append(tasks, task)
    }
    return tasks, nil
}

func UpdateTask(id uint, title, description string, completed bool, dueDate time.Time, priority Priority) (*Task, error) {
    var task Task
    err := db.DB.QueryRow(
        `UPDATE tasks 
         SET title = $1, description = $2, completed = $3, due_date = $4, priority = $5, updated_at = NOW() 
         WHERE id = $6 
         RETURNING id, title, description, completed, user_id, due_date, priority, created_at, updated_at`,
        title, description, completed, dueDate, priority, id,
    ).Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.UserID, &task.DueDate, &task.Priority, &task.CreatedAt, &task.UpdatedAt)

    if err != nil {
        return nil, err
    }
    return &task, nil
}

func DeleteTask(id uint) error {
    _, err := db.DB.Exec("DELETE FROM tasks WHERE id = $1", id)
    return err
} 