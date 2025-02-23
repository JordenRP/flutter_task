package models

import (
    "time"
    "todo-app/internal/db"
)

type Category struct {
    ID        uint      `json:"id"`
    Name      string    `json:"name"`
    UserID    uint      `json:"user_id"`
    CreatedAt time.Time `json:"created_at"`
}

func CreateCategory(name string, userID uint) (*Category, error) {
    var category Category
    err := db.DB.QueryRow(
        `INSERT INTO categories (name, user_id, created_at) 
         VALUES ($1, $2, NOW()) 
         RETURNING id, name, user_id, created_at`,
        name, userID,
    ).Scan(&category.ID, &category.Name, &category.UserID, &category.CreatedAt)

    if err != nil {
        return nil, err
    }
    return &category, nil
}

func GetUserCategories(userID uint) ([]Category, error) {
    rows, err := db.DB.Query(
        `SELECT id, name, user_id, created_at 
         FROM categories 
         WHERE user_id = $1 
         ORDER BY created_at DESC`,
        userID,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var categories []Category
    for rows.Next() {
        var category Category
        err := rows.Scan(&category.ID, &category.Name, &category.UserID, &category.CreatedAt)
        if err != nil {
            return nil, err
        }
        categories = append(categories, category)
    }
    return categories, nil
}

func DeleteCategory(id, userID uint) error {
    result, err := db.DB.Exec(
        "DELETE FROM categories WHERE id = $1 AND user_id = $2",
        id, userID,
    )
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return ErrNotFound
    }

    return nil
}

func GetTasksByCategory(categoryID, userID uint) ([]Task, error) {
    rows, err := db.DB.Query(
        `SELECT id, title, description, completed, user_id, due_date, priority, created_at, updated_at 
         FROM tasks 
         WHERE category_id = $1 AND user_id = $2 
         ORDER BY due_date ASC, priority DESC, created_at DESC`,
        categoryID, userID,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var tasks []Task
    for rows.Next() {
        var task Task
        err := rows.Scan(
            &task.ID, &task.Title, &task.Description, &task.Completed,
            &task.UserID, &task.DueDate, &task.Priority, &task.CreatedAt, &task.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        tasks = append(tasks, task)
    }
    return tasks, nil
}

func UpdateTaskCategory(taskID, categoryID, userID uint) error {
    result, err := db.DB.Exec(
        "UPDATE tasks SET category_id = $1 WHERE id = $2 AND user_id = $3",
        categoryID, taskID, userID,
    )
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return ErrNotFound
    }

    return nil
} 