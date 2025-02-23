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

func GetCategory(id, userID uint) (*Category, error) {
    var category Category
    err := db.DB.QueryRow(
        `SELECT id, name, user_id, created_at 
         FROM categories 
         WHERE id = $1 AND user_id = $2`,
        id, userID,
    ).Scan(&category.ID, &category.Name, &category.UserID, &category.CreatedAt)

    if err != nil {
        return nil, err
    }
    return &category, nil
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
        `SELECT t.id, t.title, t.description, t.completed, t.user_id, t.category_id, t.due_date, t.priority, t.created_at, t.updated_at,
                c.id, c.name, c.user_id, c.created_at
         FROM tasks t
         LEFT JOIN categories c ON t.category_id = c.id
         WHERE t.category_id = $1 AND t.user_id = $2 
         ORDER BY t.due_date ASC, t.priority DESC, t.created_at DESC`,
        categoryID, userID,
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