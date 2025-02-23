package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "github.com/gorilla/mux"
    "todo-app/internal/models"
)

type CategoryHandler struct{}

type CreateCategoryRequest struct {
    Name string `json:"name"`
}

func NewCategoryHandler() *CategoryHandler {
    return &CategoryHandler{}
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
    userID := getUserIDFromToken(r)
    categories, err := models.GetUserCategories(userID)
    if err != nil {
        http.Error(w, "Could not get categories", http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(categories)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req CreateCategoryRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    userID := getUserIDFromToken(r)
    category, err := models.CreateCategory(req.Name, userID)
    if err != nil {
        http.Error(w, "Could not create category", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    categoryID, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        http.Error(w, "Invalid category ID", http.StatusBadRequest)
        return
    }

    userID := getUserIDFromToken(r)
    err = models.DeleteCategory(uint(categoryID), userID)
    if err != nil {
        http.Error(w, "Could not delete category", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func (h *CategoryHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    categoryID, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        http.Error(w, "Invalid category ID", http.StatusBadRequest)
        return
    }

    userID := getUserIDFromToken(r)
    tasks, err := models.GetTasksByCategory(uint(categoryID), userID)
    if err != nil {
        http.Error(w, "Could not get tasks", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(tasks)
}

func (h *CategoryHandler) UpdateTaskCategory(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    taskID, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        http.Error(w, "Invalid task ID", http.StatusBadRequest)
        return
    }

    var req struct {
        CategoryID uint `json:"category_id"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    userID := getUserIDFromToken(r)
    err = models.UpdateTaskCategory(uint(taskID), req.CategoryID, userID)
    if err != nil {
        http.Error(w, "Could not update task category", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
} 