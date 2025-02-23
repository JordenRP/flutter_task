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
    CategoryID  *uint     `json:"category_id"`
    DueDate     time.Time `json:"due_date"`
    Priority    Priority  `json:"priority"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    Category    *Category `json:"category,omitempty"`
}

func CreateTask(title, description string, userID uint, dueDate time.Time, priority Priority, categoryID *uint) (*Task, error) {
    var task Task
    err := db.DB.QueryRow(
        `INSERT INTO tasks (title, description, completed, user_id, category_id, due_date, priority, created_at, updated_at) 
         VALUES ($1, $2, false, $3, $4, $5, $6, NOW(), NOW()) 
         RETURNING id, title, description, completed, user_id, category_id, due_date, priority, created_at, updated_at`,
        title, description, userID, categoryID, dueDate, priority,
    ).Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.UserID, &task.CategoryID,
           &task.DueDate, &task.Priority, &task.CreatedAt, &task.UpdatedAt)

    if err != nil {
        return nil, err
    }

    if task.CategoryID != nil {
        category, err := GetCategory(*task.CategoryID, userID)
        if err == nil {
            task.Category = category
        }
    }

    return &task, nil
}

func GetUserTasks(userID uint) ([]Task, error) {
    rows, err := db.DB.Query(
        `SELECT t.id, t.title, t.description, t.completed, t.user_id, t.category_id, t.due_date, t.priority, t.created_at, t.updated_at,
                COALESCE(c.id, 0), COALESCE(c.name, ''), COALESCE(c.user_id, 0), COALESCE(c.created_at, NOW())
         FROM tasks t
         LEFT JOIN categories c ON t.category_id = c.id
         WHERE t.user_id = $1 
         ORDER BY t.due_date ASC, t.priority DESC, t.created_at DESC`,
        userID,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var tasks []Task
    for rows.Next() {
        var task Task
        var category Category
        var categoryID *uint
        err := rows.Scan(
            &task.ID, &task.Title, &task.Description, &task.Completed, &task.UserID, &categoryID,
            &task.DueDate, &task.Priority, &task.CreatedAt, &task.UpdatedAt,
            &category.ID, &category.Name, &category.UserID, &category.CreatedAt,
        )
        if err != nil {
            return nil, err
        }
        task.CategoryID = categoryID
        if categoryID != nil {
            task.Category = &category
        }
        tasks = append(tasks, task)
    }
    return tasks, nil
}

func UpdateTask(id uint, title, description string, completed bool, dueDate time.Time, priority Priority, categoryID *uint) (*Task, error) {
    var task Task
    err := db.DB.QueryRow(
        `UPDATE tasks 
         SET title = $1, description = $2, completed = $3, due_date = $4, priority = $5, category_id = $6, updated_at = NOW() 
         WHERE id = $7 
         RETURNING id, title, description, completed, user_id, category_id, due_date, priority, created_at, updated_at`,
        title, description, completed, dueDate, priority, categoryID, id,
    ).Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.UserID, &task.CategoryID,
           &task.DueDate, &task.Priority, &task.CreatedAt, &task.UpdatedAt)

    if err != nil {
        return nil, err
    }

    if task.CategoryID != nil {
        category, err := GetCategory(*task.CategoryID, task.UserID)
        if err == nil {
            task.Category = category
        }
    }

    return &task, nil
}

func DeleteTask(id uint) error {
    _, err := db.DB.Exec("DELETE FROM tasks WHERE id = $1", id)
    return err
} 