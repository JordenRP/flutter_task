package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "github.com/gorilla/mux"
    "todo-app/internal/models"
    "github.com/golang-jwt/jwt/v5"
)

type TaskHandler struct{}

type CreateTaskRequest struct {
    Title       string `json:"title"`
    Description string `json:"description"`
}

type UpdateTaskRequest struct {
    Title       string `json:"title"`
    Description string `json:"description"`
    Completed   bool   `json:"completed"`
}

func NewTaskHandler() *TaskHandler {
    return &TaskHandler{}
}

func getUserIDFromToken(r *http.Request) uint {
    claims := r.Context().Value("claims").(jwt.MapClaims)
    userID := uint(claims["user_id"].(float64))
    return userID
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req CreateTaskRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    userID := getUserIDFromToken(r)
    task, err := models.CreateTask(req.Title, req.Description, userID)
    if err != nil {
        http.Error(w, "Could not create task", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
    userID := getUserIDFromToken(r)
    tasks, err := models.GetUserTasks(userID)
    if err != nil {
        http.Error(w, "Could not get tasks", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    taskID, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        http.Error(w, "Invalid task ID", http.StatusBadRequest)
        return
    }

    var req UpdateTaskRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    task, err := models.UpdateTask(uint(taskID), req.Title, req.Description, req.Completed)
    if err != nil {
        http.Error(w, "Could not update task", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    taskID, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        http.Error(w, "Invalid task ID", http.StatusBadRequest)
        return
    }

    err = models.DeleteTask(uint(taskID))
    if err != nil {
        http.Error(w, "Could not delete task", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
} 