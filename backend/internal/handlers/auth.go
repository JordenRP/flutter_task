package handlers

import (
    "encoding/json"
    "net/http"
    "github.com/golang-jwt/jwt/v5"
    "time"
    "todo-app/internal/models"
)

type AuthHandler struct {
    jwtSecret []byte
}

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type RegisterRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
    Name     string `json:"name"`
}

type AuthResponse struct {
    Token string `json:"token"`
}

func NewAuthHandler(secret string) *AuthHandler {
    return &AuthHandler{
        jwtSecret: []byte(secret),
    }
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    user, err := models.GetUserByEmail(req.Email)
    if err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    if !user.CheckPassword(req.Password) {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "email": user.Email,
        "exp":   time.Now().Add(time.Hour * 24).Unix(),
    })

    tokenString, err := token.SignedString(h.jwtSecret)
    if err != nil {
        http.Error(w, "Could not generate token", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(AuthResponse{Token: tokenString})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    user, err := models.CreateUser(req.Email, req.Password, req.Name)
    if err != nil {
        http.Error(w, "Could not create user", http.StatusInternalServerError)
        return
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "email": user.Email,
        "exp":   time.Now().Add(time.Hour * 24).Unix(),
    })

    tokenString, err := token.SignedString(h.jwtSecret)
    if err != nil {
        http.Error(w, "Could not generate token", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(AuthResponse{Token: tokenString})
} 