package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "time"
    "log"
    "fmt"
    "github.com/gorilla/mux"
    "todo-app/internal/models"
    "github.com/golang-jwt/jwt/v5"
)

type TaskHandler struct{}

type CreateTaskRequest struct {
    Title       string `json:"title"`
    Description string `json:"description"`
    DueDate     string `json:"due_date"`
    Priority    int    `json:"priority"`
}

type UpdateTaskRequest struct {
    Title       string `json:"title"`
    Description string `json:"description"`
    Completed   bool   `json:"completed"`
    DueDate     string `json:"due_date"`
    Priority    int    `json:"priority"`
}

func NewTaskHandler() *TaskHandler {
    return &TaskHandler{}
}

func getUserIDFromToken(r *http.Request) uint {
    claims := r.Context().Value("claims").(jwt.MapClaims)
    userID := uint(claims["user_id"].(float64))
    return userID
}

func parseDate(dateStr string) (time.Time, error) {
    formats := []string{
        time.RFC3339,
        "2006-01-02T15:04:05.000",
        "2006-01-02T15:04:05Z",
        "2006-01-02T15:04:05",
        "2006-01-02",
    }

    for _, format := range formats {
        if t, err := time.Parse(format, dateStr); err == nil {
            return t, nil
        }
    }

    return time.Time{}, fmt.Errorf("unsupported date format: %s", dateStr)
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req CreateTaskRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        log.Printf("Error decoding request: %v", err)
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    log.Printf("Creating task: %+v", req)

    dueDate, err := parseDate(req.DueDate)
    if err != nil {
        log.Printf("Error parsing date: %v", err)
        http.Error(w, "Invalid date format", http.StatusBadRequest)
        return
    }

    userID := getUserIDFromToken(r)
    task, err := models.CreateTask(req.Title, req.Description, userID, dueDate, models.Priority(req.Priority))
    if err != nil {
        log.Printf("Error creating task: %v", err)
        http.Error(w, "Could not create task", http.StatusInternalServerError)
        return
    }

    log.Printf("Task created successfully: %+v", task)
    json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
    userID := getUserIDFromToken(r)
    tasks, err := models.GetUserTasks(userID)
    if err != nil {
        log.Printf("Error getting tasks: %v", err)
        http.Error(w, "Could not get tasks", http.StatusInternalServerError)
        return
    }

    log.Printf("Retrieved %d tasks", len(tasks))
    json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    taskID, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        log.Printf("Error parsing task ID: %v", err)
        http.Error(w, "Invalid task ID", http.StatusBadRequest)
        return
    }

    var req UpdateTaskRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        log.Printf("Error decoding request: %v", err)
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    log.Printf("Updating task %d: %+v", taskID, req)

    dueDate, err := parseDate(req.DueDate)
    if err != nil {
        log.Printf("Error parsing date: %v", err)
        http.Error(w, "Invalid date format", http.StatusBadRequest)
        return
    }

    task, err := models.UpdateTask(uint(taskID), req.Title, req.Description, req.Completed, dueDate, models.Priority(req.Priority))
    if err != nil {
        log.Printf("Error updating task: %v", err)
        http.Error(w, "Could not update task", http.StatusInternalServerError)
        return
    }

    log.Printf("Task updated successfully: %+v", task)
    json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    taskID, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        log.Printf("Error parsing task ID: %v", err)
        http.Error(w, "Invalid task ID", http.StatusBadRequest)
        return
    }

    log.Printf("Deleting task %d", taskID)

    err = models.DeleteTask(uint(taskID))
    if err != nil {
        log.Printf("Error deleting task: %v", err)
        http.Error(w, "Could not delete task", http.StatusInternalServerError)
        return
    }

    log.Printf("Task %d deleted successfully", taskID)
    w.WriteHeader(http.StatusOK)
} 