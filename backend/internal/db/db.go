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
        `CREATE TABLE IF NOT EXISTS categories (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
            created_at TIMESTAMP NOT NULL DEFAULT NOW()
        )`,
        `CREATE TABLE IF NOT EXISTS tasks (
            id SERIAL PRIMARY KEY,
            title VARCHAR(255) NOT NULL,
            description TEXT,
            completed BOOLEAN DEFAULT FALSE,
            user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
            category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
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
    // Проверяем существование таблицы categories
    var hasCategoriesTable bool
    err := DB.QueryRow(`
        SELECT EXISTS (
            SELECT 1 
            FROM information_schema.tables 
            WHERE table_name = 'categories'
        )
    `).Scan(&hasCategoriesTable)
    if err != nil {
        return err
    }

    // Если таблица categories не существует, создаем ее
    if !hasCategoriesTable {
        _, err = DB.Exec(`
            CREATE TABLE categories (
                id SERIAL PRIMARY KEY,
                name VARCHAR(255) NOT NULL,
                user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
                created_at TIMESTAMP NOT NULL DEFAULT NOW()
            )
        `)
        if err != nil {
            return err
        }
    }

    // Проверяем существование колонки category_id в таблице tasks
    var hasCategoryId bool
    err = DB.QueryRow(`
        SELECT EXISTS (
            SELECT 1 
            FROM information_schema.columns 
            WHERE table_name = 'tasks' AND column_name = 'category_id'
        )
    `).Scan(&hasCategoryId)
    if err != nil {
        return err
    }

    // Если колонка category_id не существует, добавляем ее
    if !hasCategoryId {
        _, err = DB.Exec(`
            ALTER TABLE tasks 
            ADD COLUMN category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL
        `)
        if err != nil {
            return err
        }
    }

    return nil
} 