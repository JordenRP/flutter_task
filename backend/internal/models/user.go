package models

import (
    "todo-app/internal/db"
    "golang.org/x/crypto/bcrypt"
)

type User struct {
    ID       uint   `json:"id"`
    Email    string `json:"email"`
    Password string `json:"-"`
    Name     string `json:"name"`
}

func CreateUser(email, password, name string) (*User, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    var id uint
    err = db.DB.QueryRow(
        "INSERT INTO users (email, password, name) VALUES ($1, $2, $3) RETURNING id",
        email, string(hashedPassword), name,
    ).Scan(&id)
    if err != nil {
        return nil, err
    }

    return &User{
        ID:    id,
        Email: email,
        Name:  name,
    }, nil
}

func GetUserByEmail(email string) (*User, error) {
    var user User
    var hashedPassword string
    err := db.DB.QueryRow(
        "SELECT id, email, password, name FROM users WHERE email = $1",
        email,
    ).Scan(&user.ID, &user.Email, &hashedPassword, &user.Name)
    if err != nil {
        return nil, err
    }

    user.Password = hashedPassword
    return &user, nil
}

func (u *User) CheckPassword(password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
    return err == nil
} 