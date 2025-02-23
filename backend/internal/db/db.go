package db

import (
    "database/sql"
    _ "github.com/lib/pq"
    "fmt"
)

var DB *sql.DB

func InitDB(host, port, user, password, dbname string) error {
    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)
    
    var err error
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        return err
    }

    err = DB.Ping()
    if err != nil {
        return err
    }

    err = createTables()
    if err != nil {
        return err
    }

    err = migrateExistingData()
    if err != nil {
        return err
    }

    return nil
}

func createTables() error {
    queries := []string{
        `CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            email VARCHAR(255) UNIQUE NOT NULL,
            password VARCHAR(255) NOT NULL,
            name VARCHAR(255) NOT NULL
        )`,
        `CREATE TABLE IF NOT EXISTS tasks (
            id SERIAL PRIMARY KEY,
            title VARCHAR(255) NOT NULL,
            description TEXT,
            completed BOOLEAN DEFAULT FALSE,
            user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
            due_date TIMESTAMP NOT NULL DEFAULT NOW(),
            priority INTEGER NOT NULL DEFAULT 0,
            created_at TIMESTAMP NOT NULL,
            updated_at TIMESTAMP NOT NULL
        )`,
        `CREATE TABLE IF NOT EXISTS notifications (
            id SERIAL PRIMARY KEY,
            user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
            task_id INTEGER REFERENCES tasks(id) ON DELETE CASCADE,
            message TEXT NOT NULL,
            created_at TIMESTAMP NOT NULL,
            read BOOLEAN DEFAULT FALSE
        )`,
    }

    for _, query := range queries {
        _, err := DB.Exec(query)
        if err != nil {
            return err
        }
    }
    return nil
}

func migrateExistingData() error {
    // Проверяем существование колонок
    var hasDueDate, hasPriority bool
    err := DB.QueryRow(`
        SELECT EXISTS (
            SELECT 1 
            FROM information_schema.columns 
            WHERE table_name = 'tasks' AND column_name = 'due_date'
        ), EXISTS (
            SELECT 1 
            FROM information_schema.columns 
            WHERE table_name = 'tasks' AND column_name = 'priority'
        )
    `).Scan(&hasDueDate, &hasPriority)
    if err != nil {
        return err
    }

    // Добавляем недостающие колонки
    if !hasDueDate {
        _, err = DB.Exec(`ALTER TABLE tasks ADD COLUMN due_date TIMESTAMP NOT NULL DEFAULT NOW()`)
        if err != nil {
            return err
        }
    }

    if !hasPriority {
        _, err = DB.Exec(`ALTER TABLE tasks ADD COLUMN priority INTEGER NOT NULL DEFAULT 0`)
        if err != nil {
            return err
        }
    }

    return nil
} 