package models

import (
    "time"
    "todo-app/internal/db"
)

type Notification struct {
    ID        uint      `json:"id"`
    UserID    uint      `json:"user_id"`
    TaskID    uint      `json:"task_id"`
    Message   string    `json:"message"`
    CreatedAt time.Time `json:"created_at"`
    Read      bool      `json:"read"`
}

func CreateNotification(userID, taskID uint, message string) error {
    _, err := db.DB.Exec(
        `INSERT INTO notifications (user_id, task_id, message, created_at, read) 
         VALUES ($1, $2, $3, NOW(), false)`,
        userID, taskID, message,
    )
    return err
}

func GetUserNotifications(userID uint) ([]Notification, error) {
    rows, err := db.DB.Query(
        `SELECT id, user_id, task_id, message, created_at, read 
         FROM notifications 
         WHERE user_id = $1 
         ORDER BY created_at DESC`,
        userID,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var notifications []Notification
    for rows.Next() {
        var n Notification
        err := rows.Scan(&n.ID, &n.UserID, &n.TaskID, &n.Message, &n.CreatedAt, &n.Read)
        if err != nil {
            return nil, err
        }
        notifications = append(notifications, n)
    }
    return notifications, nil
}

func MarkNotificationAsRead(id uint) error {
    _, err := db.DB.Exec(
        "UPDATE notifications SET read = true WHERE id = $1",
        id,
    )
    return err
}

func CheckDueTasks() error {
    _, err := db.DB.Exec(`
        INSERT INTO notifications (user_id, task_id, message, created_at, read)
        SELECT 
            user_id,
            id as task_id,
            CASE 
                WHEN due_date < NOW() THEN 'Задача просрочена: ' || title
                WHEN due_date < NOW() + INTERVAL '1 day' THEN 'Задача должна быть выполнена сегодня: ' || title
                WHEN due_date < NOW() + INTERVAL '3 days' THEN 'До срока выполнения задачи осталось менее 3 дней: ' || title
                ELSE 'Новая задача создана: ' || title
            END as message,
            NOW(),
            false
        FROM tasks
        WHERE 
            completed = false 
            AND (
                due_date < NOW() 
                OR due_date < NOW() + INTERVAL '1 day'
                OR due_date < NOW() + INTERVAL '3 days'
                OR id NOT IN (SELECT task_id FROM notifications)
            )
            AND NOT EXISTS (
                SELECT 1 FROM notifications n 
                WHERE n.task_id = tasks.id 
                AND n.created_at > NOW() - INTERVAL '1 hour'
            )
    `)
    return err
} 